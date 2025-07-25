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

func expectedResult(k2, k3, k5, k6 int64) string {
	min := func(a, b int64) int64 {
		if a < b {
			return a
		}
		return b
	}
	x := min(k2, min(k5, k6))
	k2 -= x
	y := min(k2, k3)
	res := x*256 + y*32
	return fmt.Sprint(res)
}

func generateCase(rng *rand.Rand) (string, string) {
	k2 := rng.Int63n(50)
	k3 := rng.Int63n(50)
	k5 := rng.Int63n(50)
	k6 := rng.Int63n(50)
	input := fmt.Sprintf("%d %d %d %d\n", k2, k3, k5, k6)
	exp := expectedResult(k2, k3, k5, k6)
	return input, exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
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
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
