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
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		n := (i % 10) + 1
		ySum := 0.0
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for j := 0; j < n; j++ {
			x := j
			y := i + j
			ySum += float64(y)
			sb.WriteString(fmt.Sprintf("%d %d\n", x, y))
		}
		input := sb.String()
		want := fmt.Sprintf("%.10f", 5+ySum/float64(n))
		out, err := run(bin, input)
		if err != nil {
			fmt.Printf("test %d error: %v\n", i+1, err)
			os.Exit(1)
		}
		if out != want {
			fmt.Printf("test %d failed: expected %q got %q\n", i+1, want, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
