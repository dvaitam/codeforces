package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

// expectedA constructs the required array or returns "-1\n" if impossible.
func expectedA(x, y, n int) string {
	b := make([]int, n)
	b[0] = x
	b[n-1] = y
	for i := 0; i < n/2; i++ {
		b[i], b[n-1-i] = b[n-1-i], b[i]
	}
	if n > 1 {
		b[1] = b[0] - 1
		for i := 2; i < n-1; i++ {
			b[i] = b[i-1] - i
		}
	}
	for i := 0; i < n/2; i++ {
		b[i], b[n-1-i] = b[n-1-i], b[i]
	}
	ok := true
	if n > 1 {
		last := b[1] - b[0]
		for i := 2; i < n; i++ {
			diff := b[i] - b[i-1]
			if last <= diff || b[i] == b[i-1] {
				ok = false
				break
			}
			last = diff
		}
	}
	if !ok {
		return "-1\n"
	}
	var sb strings.Builder
	for i, v := range b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func generateCase(rng *rand.Rand) (string, string) {
	x := rng.Intn(998) + 1
	y := rng.Intn(1000-x) + x + 1
	n := rng.Intn(10) + 3
	input := fmt.Sprintf("1\n%d %d %d\n", x, y, n)
	expect := expectedA(x, y, n)
	return input, expect
}

func runCase(bin, input, expected string) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	expStr := strings.TrimSpace(expected)
	if outStr != expStr {
		return fmt.Errorf("expected %q got %q", expStr, outStr)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
