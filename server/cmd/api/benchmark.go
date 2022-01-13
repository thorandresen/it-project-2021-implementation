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

	var storeIdentityTime []time.Duration

	fmt.Println("STORING IDENTITIES")
	var siBar Bar
	siBar.NewOption(0, 5000)
	for i := 0; i < 5000; i++ {
		start := time.Now()

		db.storeIdentity("hi", "hi")

		elapsed := time.Since(start)
		storeIdentityTime = append(storeIdentityTime, elapsed)
		siBar.Play(int64(i + 1))
	}
	siBar.Finish()

	var initPufTime []time.Duration

	fmt.Println("INITIATING PUFS")
	var ipBar Bar
	ipBar.NewOption(0, 50)

	for i := 0; i < 50; i++ {
		start := time.Now()

		db.initiatePuf(i)

		elapsed := time.Since(start)
		initPufTime = append(initPufTime, elapsed)
		ipBar.Play(int64(i + 1))
	}
	siBar.Finish()

	var updateOwnerTime []time.Duration

	fmt.Println("UPDATE OWNERS")

	var uoBar Bar
	uoBar.NewOption(0, 10000)
	for i := 0; i < 10000; i++ {
		start := time.Now()

		db.updateOwner()

		elapsed := time.Since(start)
		updateOwnerTime = append(updateOwnerTime, elapsed)
		uoBar.Play(int64(i + 1))
	}
	uoBar.Finish()

	benchmarkTime := time.Since(totalBenchTime)
	var benchmark []time.Duration
	benchmark = append(benchmark, benchmarkTime)

	fmt.Println("LOGGING RESULTS TO FILE")
	os.Mkdir("log", os.ModePerm)
	log2file("log/"+db.getDatabaseType()+"_store_id.log", "STORE_IDENTITY", storeIdentityTime)
	log2file("log/"+db.getDatabaseType()+"_init_puf.log", "INITATEP_PUF", initPufTime)
	log2file("log/"+db.getDatabaseType()+"_update_owner.log", "UPDATE_OWNER", updateOwnerTime)
	log2file("log/"+db.getDatabaseType()+"_total_benchmark.log", "TOTAL_TIME", benchmark)
}

func log2file(filename string, call string, values []time.Duration) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)

	for i := 0; i < len(values); i++ {
		log.Printf("%s done in: %s", call, values[i])
	}

}
