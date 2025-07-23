package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var s string
	if _, err := fmt.Fscan(reader, &s); err != nil {
		return
	}
	ePos := strings.IndexByte(s, 'e')
	if ePos == -1 {
		fmt.Println(s)
		return
	}
	beforeE := s[:ePos]
	afterE := s[ePos+1:]
	var exp int
	fmt.Sscanf(afterE, "%d", &exp)

	dotPos := strings.IndexByte(beforeE, '.')
	if dotPos == -1 {
		dotPos = len(beforeE)
	}
	intPart := beforeE[:dotPos]
	fracPart := ""
	if dotPos < len(beforeE) {
		fracPart = beforeE[dotPos+1:]
	}
	digits := intPart + fracPart
	decIndex := len(intPart) + exp

	if decIndex >= len(digits) {
		out := digits + strings.Repeat("0", decIndex-len(digits))
		out = strings.TrimLeft(out, "0")
		if out == "" {
			out = "0"
		}
		fmt.Println(out)
		return
	}

	iPart := digits[:decIndex]
	fPart := digits[decIndex:]
	iPart = strings.TrimLeft(iPart, "0")
	if iPart == "" {
		iPart = "0"
	}
	fPart = strings.TrimRight(fPart, "0")
	if fPart == "" {
		fmt.Println(iPart)
	} else {
		fmt.Printf("%s.%s\n", iPart, fPart)
	}
}
