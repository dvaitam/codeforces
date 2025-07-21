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

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveCase(n, m int, s string) int {
	cnt := [7]int{}
	for _, ch := range s {
		if ch >= 'A' && ch <= 'G' {
			cnt[ch-'A']++
		}
	}
	add := 0
	for i := 0; i < 7; i++ {
		if m > cnt[i] {
			add += m - cnt[i]
		}
	}
	return add
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(50) + 1
	m := rng.Intn(5) + 1
	b := make([]byte, n)
	for i := range b {
		b[i] = byte('A' + rng.Intn(7))
	}
	s := string(b)
	input := fmt.Sprintf("1\n%d %d\n%s\n", n, m, s)
	expect := fmt.Sprintf("%d", solveCase(n, m, s))
	return input, expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
