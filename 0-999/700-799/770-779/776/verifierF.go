package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func run(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for i := 1; i <= 100; i++ {
		n := i + 2
		input := fmt.Sprintf("%d %d\n", n, 0)
		expected := "1"
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("Test %d: execution error: %v\nOutput:\n%s\n", i, err, got)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("Test %d failed.\nInput:%sExpected:%s Got:%s\n", i, input, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
