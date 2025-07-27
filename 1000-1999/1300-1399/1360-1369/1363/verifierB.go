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

func solveCase(s string) string {
	n := len(s)
	prefix0 := make([]int, n+1)
	prefix1 := make([]int, n+1)
	for i := 0; i < n; i++ {
		prefix0[i+1] = prefix0[i]
		prefix1[i+1] = prefix1[i]
		if s[i] == '0' {
			prefix0[i+1]++
		} else {
			prefix1[i+1]++
		}
	}
	total0 := prefix0[n]
	total1 := prefix1[n]
	minOps := n
	for i := 0; i <= n; i++ {
		ops01 := prefix1[i] + (total0 - prefix0[i])
		if ops01 < minOps {
			minOps = ops01
		}
		ops10 := prefix0[i] + (total1 - prefix1[i])
		if ops10 < minOps {
			minOps = ops10
		}
	}
	return fmt.Sprintf("%d", minOps)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(15) + 1
	var sb strings.Builder
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			sb.WriteByte('0')
		} else {
			sb.WriteByte('1')
		}
	}
	s := sb.String()
	input := fmt.Sprintf("1\n%s\n", s)
	expect := solveCase(s)
	return input, expect
}

func runCase(bin, in, expect string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	var errBuf bytes.Buffer
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expect {
		return fmt.Errorf("expected %q got %q", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, expect := generateCase(rng)
		if err := runCase(bin, in, expect); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
