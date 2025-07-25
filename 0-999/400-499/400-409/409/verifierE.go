package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func run(binary string, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, binary)
	cmd.Stdin = strings.NewReader(input)
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
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	idx := 0
	for p := 1; p <= 10; p++ {
		for q := 1; q <= 10; q++ {
			input := strconv.FormatFloat(float64(p)/float64(q), 'f', 6, 64) + "\n"
			want := fmt.Sprintf("%d %d", p, q)
			out, err := run(bin, input)
			if err != nil {
				fmt.Printf("test %d error: %v\n", idx+1, err)
				os.Exit(1)
			}
			if out != want {
				fmt.Printf("test %d failed: expected %q got %q\n", idx+1, want, out)
				os.Exit(1)
			}
			idx++
		}
	}
	fmt.Println("All tests passed")
}
