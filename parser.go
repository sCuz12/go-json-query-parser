package parser

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/sCuz12/go-json-query-parser/sliceutils"
)

type Query struct {
	Fields     []string
	Conditions []string
	Operators  []string
	SelectAll  bool //flag to indicate user has specified to select all fields 
}

// Compile the regular expression once
var conditionRegex = regexp.MustCompile(`( in |[><=])`)


func (q *Query) Parse(query string) error {
	lowerQuery := strings.ToLower(query)

	parts := strings.Split(lowerQuery, " where ")

	fieldsPart := strings.TrimPrefix(parts[0], "select ")

	//check if asterics(*) and update flag
	if fieldsPart == "*" {
		q.SelectAll = true
	}else {
		fields := strings.Split(fieldsPart, ",") // get the fields for select
		//trim 
		for i,field := range fields {
			fields[i] = strings.TrimSpace(field)
			q.Fields = append(q.Fields, fields[i])
		}
	}

	if len(parts) > 1 {
		conditionsPart := parts[1]

		// Split conditions by " and " and " or "
		re := regexp.MustCompile(`\s+(and|or)\s+`)
		splitConditions := re.Split(conditionsPart, -1)
		q.Conditions = splitConditions

		//find all "and" and "or" operators
		operators := re.FindAllString(conditionsPart, -1)

		for _, op := range operators {
			q.Operators = append(q.Operators, strings.TrimSpace(op))
		}
	}
	return nil
}

func (q *Query) ProcessQuery(jsonData string) (string, int, error) {
	//unmarshal string to proccessable
	var data []map[string]interface{}

	result := []map[string]interface{}{}

	err := json.Unmarshal([]byte(jsonData), &data)

	if err != nil {
		return "", 0, fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	//If selected all append all fields 
	if q.SelectAll && len(data) > 0 {
		for key := range data[0] {
			q.Fields = append(q.Fields, key)
		}
	}

	// Iterate over the slice of items
	for _, item := range data {
		// check if apply to conditions
		if evaluateConditions(item, q.Conditions, q.Operators) {
			filteredItem := map[string]interface{}{}
			//iterate through available fields extracted from parser
			for _, field := range q.Fields {
				if val, ok := item[field]; ok {

					filteredItem[field] = val
				}
			}

			result = append(result, filteredItem)
		}
	}

	jsonDataResp, err := json.MarshalIndent(result, "", "  ")

	if err != nil {
		return "", 0, fmt.Errorf("error marshalling result: %w", err)
	}

	return string(jsonDataResp), len(result), nil
}

// evaluateConditions evaluates a list of conditions on a given data item and combines
// the results using logical operators ("and", "or").
func evaluateConditions(data map[string]interface{}, conditions []string, operators []string) bool {
	//Ex conditions : [age < 31 name=john]
	if len(conditions) == 0 {
		return true
	}

	//initially apply the first condition
	result := evaluateCondition(data, conditions[0])

	//loop to the rest of conditions
	for i := 1; i < len(conditions); i++ {
		operator := operators[i-1]
		conditionResult := evaluateCondition(data, conditions[i])

		if operator == "and" {
			result = result && conditionResult
		} else if operator == "or" {
			result = result || conditionResult
		}
	}
	return result
}


// evaluateCondition evaluates a single condition on a given data item.
func evaluateCondition(data map[string]interface{}, condition string) bool {
	parts := conditionRegex.Split(condition, -1)

	//if the parts are not 2 means no key , value something went wrong 
	if len(parts) != 2 {
		return false
	}

	key := strings.TrimSpace(parts[0]) //get key (ex: age)
	valueStr := strings.TrimSpace(parts[1]) //get value (ex: 18)

	operator := strings.TrimSpace(conditionRegex.FindString(condition)) //get operator (>,in,=)

	dataValue, exists := data[key]

	//return falses if not found
	if !exists {
		fmt.Println("Key not found in data:", key)
		return false
	}
	
	//check for in operator
	if operator == "in" {

		if !sliceutils.IsSlice(dataValue) {
			fmt.Println("Data value is not a slice")
			return false
		}
		
		//assuming the values is comma seperated trim parenthesis 
		values := strings.Split(strings.Trim(valueStr, "()"), ",")

		for i, v := range values {
			values[i] = strings.Trim(v, " '") // Remove leading and trailing spaces and single quotes
		}

		dataValues := reflect.ValueOf(dataValue)

		for i := 0; i < dataValues.Len(); i++ {
			elem := dataValues.Index(i).Interface()
			elemStr := fmt.Sprintf("%v", elem)
			fmt.Println("Checking element:", elemStr)
			for _, v := range values {
				fmt.Println(v)
				if elemStr == v {
					fmt.Println("Match found:", elemStr)
					return true
				}
			}
		}

		return false
	}

	conditionValue, err := strconv.ParseFloat(valueStr, 64)
	if err == nil {
		// Numeric comparison
		dataValueFloat, ok := dataValue.(float64)
		if !ok {
			return false
		}

		switch operator {
		case ">":
			return dataValueFloat > conditionValue
		case "<":
			return dataValueFloat < conditionValue
		case "=":
			return dataValueFloat == conditionValue
		default:
			return false
		}
	 
	} else {
		// String comparison
		dataValueStr, ok := dataValue.(string)

		if !ok {
			return false
		}

		switch operator {
		case "=":
			return strings.ToLower(dataValueStr) == valueStr
		default:
			return false
		}
	}
}
