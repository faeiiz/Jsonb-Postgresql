package main

import (
	"Jsonb/benchmark"
	"Jsonb/internal/db"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
)

func init() {
	start := time.Now()

	if err := godotenv.Load(); err != nil {
		log.Fatal("❌ Error loading .env file")
	}
	fmt.Println("✅ .env loaded")

	pgStart := time.Now()
	if err := db.ConnectPostgres(); err != nil {
		log.Fatalf("❌ Postgres connection failed: %v", err)
	}
	fmt.Println("✅ Connected to PostgreSQL in", time.Since(pgStart))

	mongoStart := time.Now()
	if err := db.ConnectMongo(); err != nil {
		log.Fatalf("❌ MongoDB connection failed: %v", err)
	}
	fmt.Println("✅ Connected to MongoDB in", time.Since(mongoStart))

	fmt.Println("🔗 Total DB connection setup time:", time.Since(start))
}

func main() {
	// Parse command line flags
	runs := flag.Int("runs", 100, "Number of benchmark runs to execute")
	flag.Parse()

	// // ✅ 1. CREATE: Insert single record
	// record := data.DummyData{
	// 	ID:   999, // Unique JSON ID
	// 	Name: "CRUD Test User",
	// 	Age:  42,
	// 	Address: map[string]string{
	// 		"city": "TestCity",
	// 		"zip":  "123456",
	// 	},
	// 	Tags: []string{"test", "crud", "jsonb"},
	// }

	// if err := db.CreatePostgres(record); err != nil {
	// 	log.Fatal("❌ Failed to create record:", err)
	// }
	// fmt.Println("✅ Created single record with data.id = 999")

	// // ✅ 2. READ: Find by JSON field data.id
	// readRecord, err := db.FindPostgresByID(999)
	// if err != nil {
	// 	log.Fatal("❌ Failed to find record by JSON id:", err)
	// }

	// // Convert datatypes.JSON to map for pretty printing
	// var jsonMap map[string]interface{}
	// if err := json.Unmarshal(readRecord.Data, &jsonMap); err != nil {
	// 	log.Fatal("❌ Failed to unmarshal JSON data:", err)
	// }
	// fmt.Printf("🔍 Read record with data.id = 999: %+v\n", jsonMap)

	// // ✅ 3. UPDATE: Change city for records with age >= 40
	// if err := db.UpdatePostgresCity(40, "UpdatedCity"); err != nil {
	// 	log.Fatal("❌ Failed to update record(s):", err)
	// }
	// fmt.Println("✅ Updated city to 'UpdatedCity' for records where age >= 40")

	// // ✅ 4. DELETE: Delete records where age < 50
	// if err := db.DeletePostgresByAge(50); err != nil {
	// 	log.Fatal("❌ Failed to delete record(s):", err)
	// }
	// fmt.Println("🗑️ Deleted record(s) where age < 50")

	// fmt.Println("🏁 Full CRUD test completed!")
	
	// Run benchmark tests
	fmt.Printf("🚀 Running benchmark %d time(s)...\n", *runs)
	for i := 1; i <= *runs; i++ {
		fmt.Printf("📊 Benchmark run %d/%d\n", i, *runs)
		benchmark.RunBenchmark()
	}
	fmt.Printf("🏁 Completed %d benchmark run(s)!\n", *runs)
}
