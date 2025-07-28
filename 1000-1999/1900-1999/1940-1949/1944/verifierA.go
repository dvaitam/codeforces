package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func solveCase(n, k int) string {
	ans := n
	for m := 1; m <= n; m++ {
		if m*(n-m) <= k {
			ans = m
			break
		}
	}
	return strconv.Itoa(ans)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(100) + 1
	maxK := n * (n - 1) / 2
	k := rng.Intn(maxK + 1)
	input := fmt.Sprintf("1\n%d %d\n", n, k)
	expect := solveCase(n, k)
	return input, expect
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	// deterministic edge cases
	cases := []struct {
		in  string
		exp string
	}{
		{"1\n1 0\n", "1"},
		{"1\n2 1\n", "1"},
		{"1\n3 0\n", "3"},
	}
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		cases = append(cases, struct{ in, exp string }{in, exp})
	}
	for i, tc := range cases {
		out, err := runCandidate(bin, tc.in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(tc.exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, tc.exp, out, tc.in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
