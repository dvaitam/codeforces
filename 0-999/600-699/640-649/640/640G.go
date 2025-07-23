package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var name string
	if _, err := fmt.Fscan(reader, &name); err != nil {
		return
	}
	var value string
	if _, err := fmt.Fscan(reader, &value); err != nil {
		return
	}

	prefix := 'i'
	if strings.Contains(value, ".") {
		prefix = 'f'
	}

	runes := []rune(name)
	if len(runes) > 0 {
		runes[0] = unicode.ToUpper(runes[0])
	}
	fmt.Fprintf(writer, "%c%s\n", prefix, string(runes))
}
