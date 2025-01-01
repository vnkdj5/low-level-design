package main

import (
	"fmt"

	inmemorydb "github.com/vnkdj5/low-level-design/in-memory-db"
)

func main() {
	db := inmemorydb.NewInMemoryDB()

	// Create a table
	schema := map[string]string{
		"name": "string",
		"age":  "int",
		"city": "string",
	}
	db.CreateTable("users", schema)

	// Insert records into the table
	db.Insert("users", "1", inmemorydb.Record{"name": "Alice", "age": 30, "city": "Pune"})
	db.Insert("users", "2", inmemorydb.Record{"name": "Bob", "age": 25, "city": "Mumbai"})
	db.Insert("users", "3", inmemorydb.Record{"name": "Charlie", "age": 30, "city": "Pune"})

	// Query records
	results, _ := db.Select("users", "name", "city", "Pune")
	fmt.Println(results) // Output: [Alice Charlie]

	// Create index on age
	db.CreateIndex("users", "age")

	// Query with index
	results, _ = db.Select("users", "name", "age", 25)
	fmt.Println(results)

	// Query: Select name and age where city = 'Pune' AND age > 28
	conditions := []inmemorydb.Condition{
		{Attribute: "city", Operator: "=", Value: "Pune"},
		{Attribute: "age", Operator: ">", Value: 28},
	}
	results1, err := db.SelectWithConditions("users", []string{"name", "age"}, conditions, "AND")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Query: Select name and age where city = 'Pune' AND age > 28")
	fmt.Println(results1)

}
