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

const mod = 1000000007

var g = map[int]int64{
	1:  1,
	2:  0,
	3:  3,
	4:  0,
	5:  15,
	6:  0,
	7:  133,
	8:  0,
	9:  2025,
	10: 0,
	11: 37851,
	12: 0,
	13: 942073,
	14: 0,
	15: 31601835,
	16: 0,
}

func expected(n int) string {
	fact := int64(1)
	for i := 2; i <= n; i++ {
		fact = fact * int64(i) % mod
	}
	ans := fact * g[n] % mod
	return fmt.Sprintf("%d", ans)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(16) + 1
	input := fmt.Sprintf("%d\n", n)
	expect := expected(n)
	return input, expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// run for every n once
	for n := 1; n <= 16; n++ {
		in := fmt.Sprintf("%d\n", n)
		exp := expected(n)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "n=%d failed: %v\n", n, err)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "n=%d expected %s got %s\n", n, exp, out)
			os.Exit(1)
		}
	}

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
