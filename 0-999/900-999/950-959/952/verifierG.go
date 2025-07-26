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

type row struct{ a, b, c byte }

func gen(s string) string {
	rows := []row{{'.', '.', '.'}}
	rows = append(rows, row{'.', 'X', 'X'})
	cur := 0
	for i := 0; i < len(s); i++ {
		target := int(s[i])
		diff := (cur - target + 256) % 256
		for j := 0; j < diff; j++ {
			rows = append(rows, row{'.', '.', '.'})
			rows = append(rows, row{'.', 'X', '.'})
		}
		rows = append(rows, row{'.', 'X', '.'})
		cur = target
	}
	rows = append(rows, row{'.', '.', '.'})

	var b strings.Builder
	for i, r := range rows {
		b.WriteByte(r.a)
		b.WriteByte(r.b)
		b.WriteByte(r.c)
		if i+1 < len(rows) {
			b.WriteByte('\n')
		}
	}
	return b.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := []string{
		"A",
		"Hello",
		"World!",
		"Test123",
		"GoLang",
		"Contest",
		"1234567890",
		"!@#$%^&*",
		"abcdefg",
		"XYZxyz",
	}

	for i, s := range tests {
		input := s + "\n"
		want := gen(s)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Printf("Test %d failed: expected:\n%s\nGot:\n%s\n", i+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
