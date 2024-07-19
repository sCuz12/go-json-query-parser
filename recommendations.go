package parser

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/sCuz12/go-json-query-parser/sliceutils"
)

// returns query recommendations
func (q *Query) GenerateRecommendations(jsonData string) ([]string, error) {
	var data []map[string]interface{}
	err := json.Unmarshal([]byte(jsonData), &data)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var recommendationsQuery []string

	randomizeSample := sliceutils.RandomSubset(data, 6).([]map[string]interface{})

	for _, item := range randomizeSample{
		randomNumber := rand.Intn(2)
		for key, value := range item {
			if randomNumber == 1 {
				recommendationsQuery = append(recommendationsQuery, fmt.Sprintf("select %v where %v=%v", key, key, value))
			} else {
				recommendationsQuery = append(recommendationsQuery, fmt.Sprintf("select * where %v=%v", key, value))
			}
		}
	}

	return sliceutils.RandomSubset(recommendationsQuery,6).([]string), nil
}
