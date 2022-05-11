# KVStore

I read a [great blog post](http://justinjaffray.com/durability-and-redo-logging/) on durability, and I wanted to try it out for myself!

The goal here is to understand how `fsync` impacts performance.

## High-Level Architecture

For this project, I have a nested key/value store. This is stored in-memory via a recursive struct. Nodes can both have values and children.

Operations are persisted to disk via a Write-Ahead Log in the format `WRITE\t<keypath>\t<data>`.

## Usage

```
# in the server
$ ./kvstore server "test.dat"
INFO[0000] starting server

# in the client
$ ./kvstore client put "test/val" "this is a test val"  
success

$ ./kvstore client get "test/val"                     
got data:
this is a test val
```

This write is persisted in `test.dat` as follows:

```
WRITE	test/val	this is a test val
```

## First Pass

For my first pass, each operation results in a write to the log and an `fsync` so that we can guarantee durability. One run of the test takes about 20ms:

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

## Second Pass

`fsync` is very expensive, so hopefully, if we can reduce the number of `fsync`, we can improve our throughput.

Intuitively, if we do 1/x the fsyncs, we should see 1/x the write latency.

To achieve this, I vectorized the writing to the WAL. Writes are put in a queue, and the queue is flushed either after it is full or after a timeout. If our intuition is correct, a queue size of X should have a total latency of 1/X since we're doing 1/X the fsyncs.

This intuition turns out to be true!

Vectorized execution with queue size of 10 - **10x improvement**

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

Vectorized execution with queue size of 50 - **50x improvement**

```
$ go test -bench=. -count 5
goos: darwin
goarch: amd64
pkg: github.com/pawalt/kvstore
cpu: VirtualApple @ 2.50GHz
BenchmarkKeyWriting-10              2559            39565 ns/op
BenchmarkKeyWriting-10              3038            406752 ns/op
BenchmarkKeyWriting-10              2992            398902 ns/op
BenchmarkKeyWriting-10              2697            411237 ns/op
BenchmarkKeyWriting-10              3001            409003 ns/op
PASS
ok      github.com/pawalt/kvstore       6.411s
```
