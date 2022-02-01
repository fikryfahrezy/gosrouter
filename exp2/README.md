## Benchmark

```
goos: windows
goarch: amd64
pkg: github.com/fikryfahrezy/gosrouter/v2
cpu: Intel(R) Core(TM) i5-8265U CPU @ 1.60GHz
BenchmarkDynamicRoute
BenchmarkDynamicRoute-8           322195              4135 ns/op            2865 B/op         47 allocs/op
BenchmarkDynamicRoute5
BenchmarkDynamicRoute5-8           19870             59143 ns/op           45552 B/op        585 allocs/op
BenchmarkDynamicRoute20
BenchmarkDynamicRoute20-8           1474            706236 ns/op          554424 B/op       7257 allocs/op
```