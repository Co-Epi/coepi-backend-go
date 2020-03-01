
# Quick Benchmark


```
$ go test -bench=.
goos: darwin
goarch: amd64
pkg: github.com/Co-Epi/coepi-backend-go/chaintest
BenchmarkForwardMiMC-4   	    1000	   2240659 ns/op
BenchmarkSHA256-4        	    2000	    987269 ns/op
BenchmarkSHA1-4          	    2000	    707570 ns/op
BenchmarkBlake2b-4       	    2000	    710009 ns/op
BenchmarkAES-4           	   30000	     57148 ns/op
PASS
ok  	github.com/Co-Epi/coepi-backend-go/chaintest	9.831s
```
