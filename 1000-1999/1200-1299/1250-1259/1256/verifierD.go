package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solveCase(n int, k int64, s string) string {
	zeroPos := []int{}
	for i := 0; i < n; i++ {
		if s[i] == '0' {
			zeroPos = append(zeroPos, i)
		}
	}
	res := make([]byte, n)
	for i := range res {
		res[i] = '1'
	}
	pos := 0
	for _, idx := range zeroPos {
		dist := idx - pos
		shift := int64(dist)
		if shift > k {
			shift = k
		}
		final := idx - int(shift)
		res[final] = '0'
		k -= shift
		pos++
	}
	return string(res)
}

func generate() (string, string) {
	const T = 100
	rand.Seed(4)
	var in strings.Builder
	var out strings.Builder
	fmt.Fprintf(&in, "%d\n", T)
	for i := 0; i < T; i++ {
		n := rand.Intn(20) + 1
		maxK := int64(n*(n-1)/2 + 1)
		k := rand.Int63n(maxK)
		bytes := make([]byte, n)
		for j := 0; j < n; j++ {
			if rand.Intn(2) == 0 {
				bytes[j] = '0'
			} else {
				bytes[j] = '1'
			}
		}
		s := string(bytes)
		fmt.Fprintf(&in, "%d %d\n", n, k)
		fmt.Fprintf(&in, "%s\n", s)
		fmt.Fprintln(&out, solveCase(n, k, s))
	}
	return in.String(), out.String()
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, buf.String())
	}
	return strings.TrimSpace(buf.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	in, exp := generate()
	out, err := runCandidate(bin, in)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if strings.TrimSpace(out) != strings.TrimSpace(exp) {
		fmt.Fprintln(os.Stderr, "wrong answer")
		fmt.Fprintln(os.Stderr, "expected:\n"+exp)
		fmt.Fprintln(os.Stderr, "got:\n"+out)
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
