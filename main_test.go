package main

import (
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"testing"

	"github.com/pawalt/kvstore/pkg/server"
)

var num = 1000
var CONCURRENCY = 20
var KEYSPACE_SIZE = 100

func BenchmarkKeyWriting(b *testing.B) {
	tmpfile, err := ioutil.TempFile("", "testing")
	if err != nil {
		b.Fatal(err)
	}
	defer func() {
		tmpfile.Close()
		os.Remove(tmpfile.Name())
	}()

	srv, err := server.New(tmpfile.Name())
	if err != nil {
		b.Fatal(err)
	}
	go srv.HandleWrites()

	var wg sync.WaitGroup
	wg.Add(CONCURRENCY)

	b.ResetTimer()

	// i apologize to anyone reading this. there's gotta be a better way right?
	for i := 0; i < CONCURRENCY; i++ {
		go func() {
			for j := 0; j < b.N/CONCURRENCY; j++ {
				srv.Put(
					[]string{
						strconv.Itoa(rand.Intn(KEYSPACE_SIZE)),
						strconv.Itoa(rand.Intn(KEYSPACE_SIZE)),
					},
					[]byte("goober"))
			}

			wg.Done()
		}()
	}

	wg.Wait()
}
