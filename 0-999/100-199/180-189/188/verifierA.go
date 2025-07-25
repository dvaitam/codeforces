package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func hexagonal(n int64) int64 {
	return 2*n*n - n
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for i := int64(1); i <= 100; i++ {
		input := fmt.Sprintf("%d\n", i)
		expected := fmt.Sprintf("%d", hexagonal(i))
		out, err := run(bin, input)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n", i, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != expected {
			fmt.Printf("wrong answer on test %d: expected %s got %s\n", i, expected, strings.TrimSpace(out))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
