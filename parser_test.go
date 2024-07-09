package main

import (
	"io"
	"os"
	"testing"
)

type TestQuery struct {
	QueryString string
	Expected int
}

func TestProcessQuery(t *testing.T) {
	jsonData := `[
		{"name": "John", "age": 30, "city": "New York"},
		{"name": "Jane", "age": 25, "city": "Chicago"},
		{"name": "Alice Johnson", "age": 40, "city": "San Francisco"}
	]`

	tests := [] TestQuery{ 
		{"select name where age > 25", 2},
		{"select name where age < 30", 1},
		{"select name where age = 40", 1},
	}

	for _, test := range tests {
		var query Query
		query.Parse(test.QueryString)

		_,total,err := query.ProcessQuery(jsonData)
		
		if err != nil {
			t.Fatalf("ProcessQuery() returned error: %v", err)
		}

		if total != test.Expected{
			t.Errorf("ProcessQuery() returned %d results; want %d", total, test.Expected)
		}
	}
}

func TestReadFromJson(t *testing.T) {
	file , err := os.Open("./testdata/data.json")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	jsonData,err := io.ReadAll(file)

	if err != nil {
        t.Fatalf("Failed to read file: %v", err)
    }

	tests := []TestQuery {
		{"select name where age > 30", 18},
		{"select name where name = Hannah Purple", 2},
		{"select name where name = Giorkos nana", 0},
	}


	for _, test := range tests {
		var query Query
		query.Parse(test.QueryString)

		_,total,err := query.ProcessQuery(string(jsonData))
		
		if err != nil {
			t.Fatalf("ProcessQuery() returned error: %v", err)
		}

		if total != test.Expected{
			t.Errorf("ProcessQuery() returned %d results; want %d", total, test.Expected)
		}
	}
}

func TestMultipleConditions(t *testing.T) {
	file , err := os.Open("./testdata/data.json") //OPEN FILE
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	jsonData,err := io.ReadAll(file)

	if err != nil {
        t.Fatalf("Failed to read file: %v", err)
    }

	tests := []TestQuery{
		{
			QueryString: "select name where age < 31  and name=John",
			Expected: 1,
		},
		{
			QueryString: "select name where age > 100 and name=John" ,
			Expected: 0,
		},
		{
			QueryString: "select name where age > 200 OR name=John" ,
			Expected: 4,
		},
		{
			QueryString: "select name where age > 210 OR age = 201" ,
			Expected: 1,
		},
		{
			QueryString: "select name where age < 31  and name=John AND age>300",
			Expected: 0,
		},
	}

	for _,test := range(tests) {
		var query Query 

		//parse the query
		query.Parse(test.QueryString)
		_,total,err := query.ProcessQuery(string(jsonData))
		
		if err != nil {
			t.Fatalf("Failed to process query %v",err)
		}

		if total != test.Expected{
			t.Errorf("ProcessQuery() returned %d results; want %d", total, test.Expected)
		}


	}

}

func TestParseQuery(t *testing.T) {

	tests:= []TestQuery{
		{"select name where age > 25 and city = new york", 1},
		{"select name where age < 30 or city = chicago", 2},
	}

	for _,test := range tests {
		var query Query 

		err := query.Parse(test.QueryString)

		if err != nil {
			t.Fatalf("Error to parse the query ")
		}
	}
}