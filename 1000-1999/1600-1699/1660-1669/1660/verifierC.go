package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func buildBinary(src, tag string) (string, error) {
	if strings.HasSuffix(src, ".go") {
		out := filepath.Join(os.TempDir(), tag)
		cmd := exec.Command("go", "build", "-o", out, src)
		if outb, err := cmd.CombinedOutput(); err != nil {
			return "", fmt.Errorf("build %s: %v\n%s", src, err, string(outb))
		}
		return out, nil
	}
	return src, nil
}

func runCase(exe, input, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func solveCase(s string) int {
	dp := make([]int, 26)
	for i := 0; i < 26; i++ {
		dp[i] = math.MinInt32
	}
	dpNo := 0
	for _, ch := range s {
		idx := int(ch - 'a')
		newDp := make([]int, 26)
		copy(newDp, dp)
		newDpNo := dpNo
		if newDp[idx] < dpNo {
			newDp[idx] = dpNo
		}
		for c := 0; c < 26; c++ {
			if newDp[idx] < dp[c] {
				newDp[idx] = dp[c]
			}
		}
		if dp[idx] != math.MinInt32 && dp[idx]+2 > newDpNo {
			newDpNo = dp[idx] + 2
		}
		dp = newDp
		dpNo = newDpNo
	}
	return len(s) - dpNo
}

func randString(r *rand.Rand, n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = byte(r.Intn(26) + 'a')
	}
	return string(b)
}

func generateCase(r *rand.Rand) (string, string) {
	n := r.Intn(20) + 1
	s := randString(r, n)
	input := fmt.Sprintf("1\n%s\n", s)
	expect := fmt.Sprintf("%d\n", solveCase(s))
	return input, expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		return
	}
	candSrc := os.Args[1]
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	refSrc := filepath.Join(dir, "1660C.go")

	cand, err := buildBinary(candSrc, "candC.bin")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	ref, err := buildBinary(refSrc, "refC.bin")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(ref, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "reference failed: %v\n", err)
			os.Exit(1)
		}
		if err := runCase(cand, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
