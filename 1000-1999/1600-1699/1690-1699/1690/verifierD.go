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

func run(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", bin)
	} else {
		cmd = exec.CommandContext(ctx, bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func solveCase(n, k int, s string) int {
	cnt := 0
	for i := 0; i < k; i++ {
		if s[i] == 'W' {
			cnt++
		}
	}
	best := cnt
	for i := k; i < n; i++ {
		if s[i-k] == 'W' {
			cnt--
		}
		if s[i] == 'W' {
			cnt++
		}
		if cnt < best {
			best = cnt
		}
	}
	return best
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(45))
	letters := []byte{'W', 'B'}
	for tc := 0; tc < 100; tc++ {
		n := rng.Intn(20) + 1
		k := rng.Intn(n) + 1
		b := make([]byte, n)
		for i := range b {
			b[i] = letters[rng.Intn(2)]
		}
		s := string(b)
		var sb strings.Builder
		fmt.Fprintf(&sb, "1\n%d %d\n%s\n", n, k, s)
		input := sb.String()
		expected := solveCase(n, k, s)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n", tc+1, err)
			os.Exit(1)
		}
		var got int
		fmt.Sscanf(strings.TrimSpace(out), "%d", &got)
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\ninput:\n%s", tc+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
