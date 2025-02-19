package additional

import (
	"fmt"
	"go-db-tools/tool"
	"regexp"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2/log"
)

var FindParamsRegex = regexp.MustCompile(`\?(\d+)`)
var FindINClauseRegex = regexp.MustCompile(`([\w._]+ IN )\(\?(\d+)\)`)

const initStepSize = 10
const sqlite3VariablesNumberLimit = 32766 // "The NNN value must be between 1 and the sqlite3_limit() parameter SQLITE_LIMIT_VARIABLE_NUMBER (default value: 32766)." (https://www.sqlite.org/c3ref/bind_blob.html)
const TimeFormat = "2006-01-02T15:04:05.999Z"

func buildINQuery(query string, rawInputs []any) (string, []any) {
	inputs := make([]any, 0, len(rawInputs))
	arrayInputs := make([][]any, 0, len(rawInputs))
	foundArray := false
	for _, v := range rawInputs {
		switch slice := v.(type) {
		case int, int32, int64, float32, float64, []byte, string, bool:
			tool.Assert(!foundArray, "non-array inputs must always go before any array inputs", "rawInputs", rawInputs)
			inputs = append(inputs, v)
		case []int:
			inner := make([]any, len(slice))
			for i, val := range slice {
				inner[i] = val
			}
			foundArray = true
			arrayInputs = append(arrayInputs, inner)
		case []int64:
			inner := make([]any, len(slice))
			for i, val := range slice {
				inner[i] = val
			}
			foundArray = true
			arrayInputs = append(arrayInputs, inner)
		case []string:
			inner := make([]any, len(slice))
			for i, val := range slice {
				inner[i] = val
			}
			foundArray = true
			arrayInputs = append(arrayInputs, inner)
		default:
			tool.Assert(false, "Unknown type", "type", fmt.Sprintf("%T", v))
		}
	}

	lowerBound := len(FindParamsRegex.FindAllIndex([]byte(query), -1)) - len(arrayInputs) + 1
	tool.Assert(lowerBound >= 0, "Not enough parameters in query. Can only find index parameters '?NNN', where 'NNN' is a positive integer", "lowerBound", lowerBound, "len(arrayInputs)", len(arrayInputs), "FindParamsRegex", FindParamsRegex)

	// As to not prepare too many query, we will prepare the input to with a few slots initially then double everytime it's not enough
	for i := 0; i < len(arrayInputs); i++ {
		goalLen := 0
		step := initStepSize
		for goalLen < len(arrayInputs[i]) {
			goalLen += step
			step *= 2
		}
		newInput := make([]any, goalLen)
		copy(newInput, arrayInputs[i])
		arrayInputs[i] = newInput
	}

	newQuery := query
	subMatches := FindINClauseRegex.FindAllSubmatch([]byte(query), -1)
	counter := -1
	for _, match := range subMatches {
		replace := string(match[0])
		prefix := string(match[1])
		num := string(match[2])
		paramNum64, err := strconv.ParseInt(num, 10, 64)
		tool.Assert(err == nil, "failed to parse int", "err", err, "replace", replace, "num", num)
		paramNum := int(paramNum64)
		if counter == -1 {
			counter = paramNum
		}

		tool.Assert(paramNum >= lowerBound, "All IN parameters must be at the end", "matched", replace, "paramNum", paramNum, "lowerBound", lowerBound)
		inputList := arrayInputs[paramNum-lowerBound]
		placeholders := make([]string, len(inputList))
		for i, _ := range inputList {
			placeholders[i] = fmt.Sprintf("?%d", counter+i)
		}
		counter += len(inputList)
		newQuery = strings.ReplaceAll(newQuery, replace, fmt.Sprintf("%s(%s)", prefix, strings.Join(placeholders, ", ")))
	}
	if counter >= sqlite3VariablesNumberLimit-10000 {
		tool.Assert(counter < sqlite3VariablesNumberLimit, "Too many variables in query", "counter", counter, "limit", sqlite3VariablesNumberLimit)
		log.Warnf("Approaching the limit of variables in query: %d/%d", counter, sqlite3VariablesNumberLimit)
	}

	lenArgs := 0
	for i := 0; i < len(arrayInputs); i++ {
		lenArgs += len(arrayInputs[i])
	}
	args := make([]any, 0, len(inputs)+lenArgs)
	copy(args, inputs)
	for _, input := range arrayInputs {
		args = append(args, input...)
	}

	return newQuery, args
}
