package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var phase string
	if _, err := fmt.Fscan(in, &phase); err != nil {
		return
	}

	if phase == "first" {
		handleFirstRun(in, out)
	} else if phase == "second" {
		handleSecondRun(in, out)
	}
}

func handleFirstRun(in *bufio.Reader, out *bufio.Writer) {
	var n int
	fmt.Fscan(in, &n)
	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &arr[i])
	}

	var builder strings.Builder
	builder.Grow((n + 1) * 8) // rough upper bound to minimize reallocs

	writeEncodedNumber(&builder, int64(n))
	for i := 0; i < n; i++ {
		builder.WriteByte('z')
		writeEncodedNumber(&builder, arr[i])
	}

	fmt.Fprintln(out, builder.String())
}

func handleSecondRun(in *bufio.Reader, out *bufio.Writer) {
	var s string
	fmt.Fscan(in, &s)
	if len(s) == 0 {
		fmt.Fprintln(out, 0)
		return
	}

	parts := strings.Split(s, "z")
	numbers := make([]int64, 0, len(parts))
	for _, part := range parts {
		if part == "" {
			continue
		}
		numbers = append(numbers, decodeNumber(part))
	}

	if len(numbers) == 0 {
		fmt.Fprintln(out, 0)
		return
	}

	n := int(numbers[0])
	fmt.Fprint(out, n)
	for i := 0; i < n && i+1 < len(numbers); i++ {
		fmt.Fprintf(out, " %d", numbers[i+1])
	}
	fmt.Fprintln(out)
}

func writeEncodedNumber(builder *strings.Builder, value int64) {
	if value == 0 {
		builder.WriteByte('a')
		return
	}

	var digits [32]byte
	pos := len(digits)
	for value > 0 {
		pos--
		digits[pos] = byte('a' + (value % 25))
		value /= 25
	}
	builder.Write(digits[pos:])
}

func decodeNumber(token string) int64 {
	var value int64
	for i := 0; i < len(token); i++ {
		value = value*25 + int64(token[i]-'a')
	}
	return value
}
