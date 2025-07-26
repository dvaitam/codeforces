package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type tcE []string

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func solveE(types []string) string {
	soft, hard := 0, 0
	for _, t := range types {
		if t == "soft" {
			soft++
		} else {
			hard++
		}
	}
	for size := 1; ; size++ {
		total := size * size
		w := (total + 1) / 2
		b := total / 2
		if (soft <= w && hard <= b) || (soft <= b && hard <= w) {
			return fmt.Sprint(size)
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := []tcE{
		{"soft"},
		{"hard"},
		{"soft", "hard"},
		{"soft", "soft"},
		{"hard", "hard"},
		{"soft", "hard", "hard"},
		{"hard", "hard", "hard"},
		{"soft", "soft", "hard", "hard"},
		{"soft", "soft", "soft", "soft"},
		{"hard", "hard", "hard", "hard"},
		{"soft", "hard", "hard", "hard", "hard"},
		{"soft", "soft", "soft", "soft", "hard"},
		{"soft", "soft", "hard", "hard", "hard", "hard"},
		{"soft", "soft", "soft", "soft", "soft", "soft"},
		{"hard", "hard", "hard", "hard", "hard", "hard"},
		{"soft", "soft", "hard", "hard", "hard", "hard", "hard"},
		{"soft", "soft", "soft", "soft", "hard", "hard", "hard", "hard"},
		{"hard", "hard", "hard", "hard", "hard", "hard", "hard", "hard", "hard"},
		{"soft", "soft", "soft", "soft", "soft", "soft", "soft", "soft", "soft", "soft"},
		{"hard", "hard", "hard", "hard", "hard", "hard", "hard", "hard", "hard", "hard"},
	}

	for i, types := range tests {
		input := fmt.Sprintf("%d\n", len(types))
		for j, t := range types {
			name := fmt.Sprintf("c%d", j)
			input += fmt.Sprintf("%s %s\n", name, t)
		}
		want := solveE(types)
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
