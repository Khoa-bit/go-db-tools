package tool

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/gofiber/fiber/v2/log"
)

const (
	contextLines       = 3
	maxTrackTraceDepth = 10
)

func Assert(condition bool, message string, data ...any) {
	if condition {
		return
	}

	dataBlock := GetDataBlock(data)

	sourceContextRows := make([]string, 0, 10)
	skip := 1
	for ; skip <= maxTrackTraceDepth; skip++ {
		pc, file, line, ok := runtime.Caller(skip)
		if !ok {
			break
		}

		fn := runtime.FuncForPC(pc)
		if strings.HasPrefix(fn.Name(), "runtime.") {
			break
		}

		sourceContext := getSourceContext(file, line)
		sourceContextRows = append(sourceContextRows, fmt.Sprintf("%s:%d %s()\n%v", FormatBold(file), line, fn.Name(), sourceContext))
	}
	if skip > maxTrackTraceDepth {
		sourceContextRows = append(sourceContextRows, "(...)")
	}

	log.Fatalf(`Assertion failed
Message: %v
%v
%v
`, message, dataBlock, strings.Join(sourceContextRows, "\n"))
}

// getSourceContext reads the source file and returns lines around the failure.
func getSourceContext(file string, line int) string {
	f, err := os.Open(file)
	if err != nil {
		return ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	start := max(1, line-contextLines)
	end := line + contextLines

	lines := make([]string, 0, end-start+1)
	for lineNumber := 1; scanner.Scan() && lineNumber <= end; lineNumber++ {
		if lineNumber < start {
			continue
		}

		if lineNumber == line {
			lines = append(lines, FormatRed(fmt.Sprintf("%4d | %s", lineNumber, scanner.Text())))
		} else {
			text := scanner.Text()
			if strings.Contains(text, "//") {
				text = FormatBrightBlack(text)
			}
			lines = append(lines, fmt.Sprintf("%s %s", FormatBrightBlack(fmt.Sprintf("%4d |", lineNumber)), text))
		}
	}

	return strings.Join(lines, "\n")
}

func GetDataBlock(data []any) string {
	dataBlock := ""
	if len(data)%2 != 0 {
		log.Fatalf("Assert: data must be a list of key-value pairs: %v\n%v", data, string(debug.Stack()))
	} else if len(data) > 0 {
		dataRows := make([]string, 0, len(data)/2)
		for i := 0; i < len(data); i += 2 {
			dataRows = append(dataRows, fmt.Sprintf("    %v: %v", data[i], data[i+1]))
		}
		dataBlock = fmt.Sprintf("Data:\n%v\n", strings.Join(dataRows, "\n"))
	}
	return dataBlock
}

func GetSimpleStackTrace(depth int) string {
	if depth == 0 {
		return ""
	}
	depth = min(maxTrackTraceDepth, depth)

	sourceContextRows := make([]string, 0, 10)
	skip := 2
	for ; skip <= depth; skip++ {
		pc, file, line, ok := runtime.Caller(skip)
		if !ok {
			break
		}

		fn := runtime.FuncForPC(pc)
		if strings.HasPrefix(fn.Name(), "runtime.") {
			break
		}

		sourceContextRows = append(sourceContextRows, fmt.Sprintf("%s:%d %s()", FormatBold(file), line, fn.Name()))
	}
	if skip > depth {
		sourceContextRows = append(sourceContextRows, "(...)")
	}

	return strings.Join(sourceContextRows, "\n")
}
