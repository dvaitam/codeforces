package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var s string
	if _, err := fmt.Fscan(reader, &s); err != nil {
		return
	}
	pos := strings.IndexByte(s, '.')
	if pos == -1 {
		pos = len(s)
	}
	digits := strings.ReplaceAll(s, ".", "")
	first := 0
	for first < len(digits) && digits[first] == '0' {
		first++
	}
	last := len(digits) - 1
	for last >= 0 && digits[last] == '0' {
		last--
	}
	digits = digits[first : last+1]
	exponent := pos - (first + 1)
	var sb strings.Builder
	sb.WriteByte(digits[0])
	if len(digits) > 1 {
		sb.WriteByte('.')
		sb.WriteString(digits[1:])
	}
	if exponent != 0 {
		sb.WriteByte('E')
		sb.WriteString(strconv.Itoa(exponent))
	}
	fmt.Println(sb.String())
}
