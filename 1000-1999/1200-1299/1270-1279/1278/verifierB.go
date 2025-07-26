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
		return out.String() + errBuf.String(), fmt.Errorf("runtime error: %v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

func solveB(a, b int64) int64 {
	if a < b {
		a, b = b, a
	}
	diff := a - b
	if diff == 0 {
		return 0
	}
	var k int64
	for {
		k++
		sum := k * (k + 1) / 2
		if sum >= diff && (sum-diff)%2 == 0 {
			return k
		}
	}
}

func generateCase(r *rand.Rand) (string, string) {
	a := r.Int63n(1_000_000_000) + 1
	b := r.Int63n(1_000_000_000) + 1
	expect := fmt.Sprintf("%d", solveB(a, b))
	input := fmt.Sprintf("1\n%d %d\n", a, b)
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
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		out = strings.TrimSpace(out)
		if val, err := strconv.ParseInt(out, 10, 64); err != nil || fmt.Sprintf("%d", val) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
