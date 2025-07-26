package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func bugEval(expr string) string {
	expr += "+"
	var result, cur int64
	sign := int64(1)
	for i := 0; i < len(expr); i++ {
		c := expr[i]
		if c == '+' || c == '-' {
			result += sign * cur
			cur = 0
		}
		if c == '-' {
			sign = -1
		}
		if c == '+' {
			sign = 1
		}
		cur = cur*10 + int64(int(c)-int('0'))
	}
	return fmt.Sprint(result)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := []string{
		"1+1",
		"10+20-5",
		"0-0+0",
		"255-255+255-255",
		"12+34+56",
		"100-50-25",
		"3-2+1-0",
		"9+8-7+6-5+4-3+2-1+0",
		"200-100+50-25+12-6+3-1",
		"1-2-3-4-5",
	}

	for i, expr := range tests {
		input := expr + "\n"
		want := bugEval(expr)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Printf("Test %d failed: expected %q, got %q\n", i+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
