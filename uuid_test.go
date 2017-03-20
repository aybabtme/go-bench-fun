package main

import (
	"testing"

	"golang.org/x/crypto/bcrypt"

	"runtime"

	"sync"

	"github.com/pborman/uuid"
)

func BenchmarkUUID(b *testing.B)         { benchmarkUUID(b, b.N) }
func BenchmarkParallelUUID(b *testing.B) { parallelBench(b, runtime.NumCPU()/2, b.N, benchmarkUUID) }

func benchmarkUUID(b *testing.B, n int) {
	dontOptimizeMeAway := true
	for i := 0; i < n; i++ {
		uuid := uuid.NewRandom()
		dontOptimizeMeAway = dontOptimizeMeAway && uuid != nil
	}
	if !dontOptimizeMeAway {
		b.Fatal("need to use the value to avoid the compiler thinking it can be optimized away")
	}
}

func BenchmarkBcrypt(b *testing.B)         { benchmarkBcrypt(b, b.N) }
func BenchmarkParallelBcrypt(b *testing.B) { parallelBench(b, runtime.NumCPU()/2, b.N, benchmarkBcrypt) }

func benchmarkBcrypt(b *testing.B, n int) {
	dontOptimizeMeAway := true
	uuid := uuid.NewRandom()
	code := []byte(uuid[0:8])
	for i := 0; i < n; i++ {
		ecode, _ := bcrypt.GenerateFromPassword(code, bcrypt.DefaultCost)
		dontOptimizeMeAway = dontOptimizeMeAway && ecode != nil
	}
	if !dontOptimizeMeAway {
		b.Fatal("need to use the value to avoid the compiler thinking it can be optimized away")
	}
}

func parallelBench(b *testing.B, para, repeats int, benchFn func(b *testing.B, n int)) {
	var (
		scheduled sync.WaitGroup
		finished  sync.WaitGroup
	)
	start := make(chan struct{})
	for i := 0; i < para; i++ {
		scheduled.Add(1)
		finished.Add(1)
		go func() {
			scheduled.Done()
			<-start
			benchFn(b, repeats/para)
			finished.Done()
		}()
	}

	scheduled.Wait() // wait until all our goroutines have been scheduled
	b.ResetTimer()   // reset the benchmark
	close(start)     // start the parallel benchmarks at the same time
	finished.Wait()  // wait until they're finished
}
