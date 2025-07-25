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

func run(binary string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, binary)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("timeout")
		}
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		out, err := run(bin)
		if err != nil {
			fmt.Printf("test %d error: %v\n", i+1, err)
			os.Exit(1)
		}
		if out != "Brainfuck" {
			fmt.Printf("test %d failed: expected %q got %q\n", i+1, "Brainfuck", out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
