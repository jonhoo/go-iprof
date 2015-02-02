package iprof

import (
	"math"
	"sort"
	"time"
)

const defaultWsize uint = 5000

type reading struct {
	duration time.Duration
	end      time.Time
}

type nreading struct {
	reading
	section string
}

var profs chan nreading
var stats map[string][]reading
var wsize map[string]uint

func init() {
	stats = make(map[string][]reading)
	wsize = make(map[string]uint)
	profs = make(chan nreading)
	go func() {
		for r := range profs {
			l := uint(len(stats[r.section]))
			w, ok := wsize[r.section]
			if !ok {
				w = defaultWsize
			}

			if l == w {
				stats[r.section] = append(stats[r.section][1:], r.reading)
			} else {
				stats[r.section] = append(stats[r.section], r.reading)
			}
		}
	}()
}

func Start(section string) func() {
	start := time.Now()
	return func() {
		end := time.Now()
		profs <- nreading{reading{end.Sub(start), end}, section}
	}
}

func Log(section string, duration time.Duration, end time.Time) {
	profs <- nreading{reading{duration, end}, section}
}

func SetWindow(section string, window uint) {
	wsize[section] = window
}

type durationSlice []float64

func (s durationSlice) Len() int {
	return len(s)
}
func (s durationSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s durationSlice) Less(i, j int) bool {
	return s[i] < s[j]
}

func Stat(section string) (average float64, percentile func(float64) float64) {
	total := float64(0)
	vals := make(durationSlice, len(stats[section]))
	for i, r := range stats[section] {
		v := r.duration.Seconds() * 1000
		vals[i] = v
		total += v
	}

	length := uint(len(vals))
	average = total / float64(length)

	sort.Sort(vals)
	percentile = func(perc float64) float64 {
		n := (perc / 100.0) * float64(length+1)
		k := uint(math.Floor(n))
		if k == 0 {
			return vals[0]
		}
		if k >= length {
			return vals[length-1]
		}

		d := n - float64(k)
		return vals[k] + d*(vals[k+1]-vals[k])
	}

	return
}
