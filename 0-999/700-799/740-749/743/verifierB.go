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

func solveCaseB(n, k int64) int64 {
	for {
		mid := int64(1) << (n - 1)
		if k == mid {
			return n
		} else if k < mid {
			n--
		} else {
			k -= mid
			n--
		}
	}
}

func generateCaseB(rng *rand.Rand) (string, string) {
	n := rng.Int63n(50) + 1
	var maxK int64
	if n >= 62 {
		maxK = int64(1)<<61 - 1
	} else {
		maxK = (int64(1) << n) - 1
	}
	k := rng.Int63n(maxK) + 1
	input := fmt.Sprintf("%d %d\n", n, k)
	expect := fmt.Sprintf("%d", solveCaseB(n, k))
	return input, expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCaseB(rng)
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
