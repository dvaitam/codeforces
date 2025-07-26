package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func reverse(n int64) int64 {
	var r int64
	for n > 0 {
		r = r*10 + n%10
		n /= 10
	}
	return r
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	rand.Seed(42)
	bin := os.Args[1]
	for i := 1; i <= 100; i++ {
		a := rand.Int63n(1e9 + 1)
		b := rand.Int63n(1e9 + 1)
		input := fmt.Sprintf("%d %d\n", a, b)
		expected := fmt.Sprintf("%d", a+reverse(b))
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
