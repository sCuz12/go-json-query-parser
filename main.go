package main

import (
	"fmt"
)

func main() {
	// queryString := "select name where age > 26 OR age < 29"
	queryString := "select name where age > 26 and name=John"
	jsonData := `[{"name": "John", "age": 30, "city": "New York"}, {"name": "Jane", "age": 25, "city": "Chicago"}]`

	var query Query

	query.Parse(queryString)
	
	proccessedResp,_,err := query.ProcessQuery(jsonData)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(proccessedResp)
}