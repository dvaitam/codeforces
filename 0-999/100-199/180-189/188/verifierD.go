package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func expectedPattern(n int) string {
	var sb strings.Builder
	for i := 1; i <= n; i++ {
		sb.WriteString(strings.Repeat("*", i))
		if i < n {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func run(binary, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, binary)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for i := 1; i <= 100; i++ {
		n := ((i - 1) % 50) + 1 // values 1..50 repeated twice
		input := fmt.Sprintf("%d\n", n)
		expected := expectedPattern(n)
		out, err := run(bin, input)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n", i, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(expected) {
			fmt.Printf("wrong answer on test %d: expected\n%s\nbut got\n%s\n", i, expected, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
