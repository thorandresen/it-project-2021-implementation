package main

import (
	"log"
	"time"
)

func benchmark() {

}

func performBenchmark(db DatabaseRequester) {
	start := time.Now()

	db.commenceDatabase()

	for i := 0; i < 100; i++ {
		db.storeIdentity("hi", "hi")
	}

	for i := 0; i < 100; i++ {
		db.initiatePuf(i)
	}

	for i := 0; i < 2000; i++ {
		db.updateOwner()
	}

	benchmarkTime := time.Since(start)
	log.Printf("BenchmarkTime VERIFICATION: %s", benchmarkTime)
}
