go-iprof provides higher-level, instrumented profiling for Go
applications. This is useful for applications that want to obtain
information about where an application is spending its time, without
resorting to the low-level data provided by pprof. Instrumenting the
code for profiling also allows more insight into time spent in external
code, which pprof does not currently handle well.
