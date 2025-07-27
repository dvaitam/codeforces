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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveCase(a, b, c int64) string {
	sum := a + b + c
	if sum%9 != 0 {
		return "NO"
	}
	n := sum / 9
	if a < n || b < n || c < n {
		return "NO"
	}
	return "YES"
}

func deterministicCases() [][3]int64 {
	return [][3]int64{
		{1, 1, 1},
		{9, 9, 9},
		{3, 4, 2},
		{1, 2, 3},
		{100000000, 1, 1},
	}
}

func randomCase(rng *rand.Rand) (int64, int64, int64) {
	a := rng.Int63n(100000000) + 1
	b := rng.Int63n(100000000) + 1
	c := rng.Int63n(100000000) + 1
	return a, b, c
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := deterministicCases()
	for len(cases) < 100 {
		a, b, c := randomCase(rng)
		cases = append(cases, [3]int64{a, b, c})
	}
	for i, tc := range cases {
		in := fmt.Sprintf("1\n%d %d %d\n", tc[0], tc[1], tc[2])
		expect := solveCase(tc[0], tc[1], tc[2])
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%d %d %d\n", i+1, err, tc[0], tc[1], tc[2])
			os.Exit(1)
		}
		if strings.TrimSpace(out) != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", i+1, expect, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
