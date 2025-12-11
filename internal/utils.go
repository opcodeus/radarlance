package internal

import (
	"bufio"
	"os"
	"strings"
	"time"

	"github.com/ditashi/jsbeautifier-go/jsbeautifier"
	"github.com/yosssi/gohtml"
)

func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines, scanner.Err()
}

func BeautifyJS(input string) (string, error) {
	opts := jsbeautifier.DefaultOptions()

	opts["indent_size"] = 2
	opts["space_in_empty_paren"] = true
	opts["jslint_happy"] = false
	opts["end_with_newline"] = true
	opts["brace_style"] = "collapse"

	result, err := jsbeautifier.Beautify(&input, opts)
	if err != nil {
		return "", err
	}
	return result, nil
}

func BeautifyHTML(input string) (string, error) {
    formatted := gohtml.Format(input)
    return formatted, nil
}

func Timestamp() string {
	return time.Now().Format("2006-01-02_15")
}