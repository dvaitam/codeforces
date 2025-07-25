package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func run(bin string, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("timeout")
	}
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	userBin := os.Args[1]
	refBin := "./773B_ref"
	if err := exec.Command("go", "build", "-o", refBin, "773B.go").Run(); err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference solution:", err)
		os.Exit(1)
	}
	rand.Seed(42)
	for t := 0; t < 100; t++ {
		n := rand.Intn(5) + 2
		lines := make([][]int, n)
		for i := 0; i < n; i++ {
			lines[i] = make([]int, 5)
			solved := false
			for j := 0; j < 5; j++ {
				if rand.Float64() < 0.2 {
					lines[i][j] = -1
				} else {
					lines[i][j] = rand.Intn(120)
					solved = true
				}
			}
			if !solved {
				idx := rand.Intn(5)
				lines[i][idx] = rand.Intn(120)
			}
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			for j := 0; j < 5; j++ {
				if j > 0 {
					sb.WriteByte(' ')
				}
				sb.WriteString(fmt.Sprintf("%d", lines[i][j]))
			}
			sb.WriteByte('\n')
		}
		input := sb.String()
		expect, err := run(refBin, input)
		if err != nil {
			fmt.Fprintln(os.Stderr, "reference failed on test", t+1, ":", err)
			os.Exit(1)
		}
		got, err := run(userBin, input)
		if err != nil {
			fmt.Fprintln(os.Stderr, "program failed on test", t+1, ":", err)
			os.Exit(1)
		}
		if expect != strings.TrimSpace(got) {
			fmt.Fprintf(os.Stderr, "mismatch on test %d: expected %s got %s\n", t+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
