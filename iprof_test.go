package iprof

import (
	"math"
	"testify/assert"
	"testing"
	"time"
)

func TestSimplejson(t *testing.T) {
	now := time.Now()

	SetWindow("test", 5000)
	for i := 0; i < 80000; i += 4 {
		Log("test", time.Duration(i)*time.Millisecond, now)
	}

	avg, perc := Stat("test")
	assert.Equal(t, 69998, avg)
	assert.Equal(t, 79004, math.Ceil(perc(95)))
	assert.Equal(t, 79804, math.Ceil(perc(99)))
	assert.Equal(t, 79996, perc(100))

	ps := Stats()
	p, ok := ps["test"]
	if assert.True(t, ok) {
		assert.Equal(t, 69998, p.Average)
		assert.Equal(t, 79004, math.Ceil(p.Percentile(95)))
		assert.Equal(t, 79804, math.Ceil(p.Percentile(99)))
		assert.Equal(t, 79996, p.Percentile(100))
	}
}
