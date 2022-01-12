package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

func performBenchmark(db DatabaseRequester) {
	fmt.Println("COMMENCING BENCHMARK OF " + db.getDatabaseType())
	totalBenchTime := time.Now()

	fmt.Println("COMMENCING DATABASE")

	db.commenceDatabase()

	var storeIdentityTime []int64

	fmt.Println("STORING IDENTITIES")
	var siBar Bar
	siBar.NewOption(0, 100)
	for i := 0; i < 100; i++ {
		start := time.Now()

		db.storeIdentity("hi", "hi")

		elapsed := time.Since(start).Milliseconds()
		storeIdentityTime = append(storeIdentityTime, elapsed)
		siBar.Play(int64(i))
	}
	siBar.Finish()

	var initPufTime []int64

	fmt.Println("INITIATING PUFS")
	var ipBar Bar
	siBar.NewOption(0, 100)

	for i := 0; i < 100; i++ {
		start := time.Now()

		db.initiatePuf(i)

		elapsed := time.Since(start).Milliseconds()
		initPufTime = append(initPufTime, elapsed)
		ipBar.Play(int64(i))
	}
	siBar.Finish()

	var updateOwnerTime []int64

	fmt.Println("UPDATE OWNERS")

	var uoBar Bar
	uoBar.NewOption(0, 2000)
	for i := 0; i < 2000; i++ {
		start := time.Now()

		db.updateOwner()

		elapsed := time.Since(start).Milliseconds()
		updateOwnerTime = append(updateOwnerTime, elapsed)
		uoBar.Play(int64(i))
	}
	uoBar.Finish()

	benchmarkTime := time.Since(totalBenchTime).Milliseconds()
	var benchmark []int64
	benchmark = append(benchmark, benchmarkTime)

	fmt.Println("LOGGING RESULTS TO FILE")

	log2file(db.getDatabaseType()+"_store_id.log", "STORE_IDENTITY", storeIdentityTime)
	log2file(db.getDatabaseType()+"init_puf.log", "INITATEP_PUF", initPufTime)
	log2file(db.getDatabaseType()+"update_owner.log", "UPDATE_OWNER", updateOwnerTime)
	log2file(db.getDatabaseType()+"total_benchmark.log", "TOTAL_TIME", benchmark)
}

func log2file(filename string, call string, values []int64) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)

	for i := 0; i < len(values); i++ {
		var toMS float64 = float64(values[i]) / 1000000
		log.Printf("%s done in: %f ms", call, toMS)
	}

}
