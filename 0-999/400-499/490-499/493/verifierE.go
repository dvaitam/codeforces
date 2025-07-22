package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
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

func max(x, y int64) int64 {
	if x > y {
		return x
	}
	return y
}

func solveE(r *bufio.Reader) string {
	var t, a, b int64
	if _, err := fmt.Fscan(r, &t, &a, &b); err != nil {
		return ""
	}
	if t == 1 && a == 1 {
		if b == 1 {
			return "inf"
		}
		return "0"
	}
	C := int64(0)
	powA := int64(1)
	aa := a
	for aa > 0 {
		d := aa % t
		C += d * powA
		aa /= t
		if powA > (1<<62)/max(a, 1) {
			powA = 0
		} else {
			powA *= a
		}
	}
	denom := t*a - 1
	if denom <= 0 || b < C || (b-C)%denom != 0 {
		return "0"
	}
	return "1"
}

func generateCaseE(rng *rand.Rand) string {
	t := rng.Int63n(5) + 1
	a := rng.Int63n(1000) + 1
	b := rng.Int63n(1000) + 1
	return fmt.Sprintf("%d %d %d\n", t, a, b)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCaseE(rng)
		expect := solveE(bufio.NewReader(strings.NewReader(tc)))
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
