before vectorized execution

```bash
$ go test -bench=. -count 5
goos: darwin
goarch: amd64
pkg: github.com/pawalt/kvstore
cpu: VirtualApple @ 2.50GHz
BenchmarkKeyWriting-10               100          19942306 ns/op
BenchmarkKeyWriting-10               100          19615203 ns/op
BenchmarkKeyWriting-10               100          23982191 ns/op
BenchmarkKeyWriting-10               100          19845203 ns/op
BenchmarkKeyWriting-10               100          19713636 ns/op
PASS
ok      github.com/pawalt/kvstore       10.723s
```

vectorized execution with queue size of 10 - **10x improvement**

```
$ go test -bench=. -count 5
goos: darwin
goarch: amd64
pkg: github.com/pawalt/kvstore
cpu: VirtualApple @ 2.50GHz
BenchmarkKeyWriting-10               592           1939244 ns/op
BenchmarkKeyWriting-10               592           1978630 ns/op
BenchmarkKeyWriting-10               590           1984245 ns/op
BenchmarkKeyWriting-10               614           1876579 ns/op
BenchmarkKeyWriting-10               619           1909888 ns/op
PASS
ok      github.com/pawalt/kvstore       7.231s
```

vectorized execution with queue size of 100 - **degraded performance**

i've honestly got no idea why increasing the write queue size would incur a performance penalty. shit makes 0 sense to me.

```
$ go test -bench=. -count 5
goos: darwin
goarch: amd64
pkg: github.com/pawalt/kvstore
cpu: VirtualApple @ 2.50GHz
BenchmarkKeyWriting-10               100          13944411 ns/op
BenchmarkKeyWriting-10               100          13615975 ns/op
BenchmarkKeyWriting-10               100          13591150 ns/op
BenchmarkKeyWriting-10               100          13576833 ns/op
BenchmarkKeyWriting-10               100          13579472 ns/op
PASS
ok      github.com/pawalt/kvstore       7.243s
```