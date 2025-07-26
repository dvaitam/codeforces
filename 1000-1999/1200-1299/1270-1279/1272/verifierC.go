package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testC struct {
	n       int
	k       int
	s       string
	allowed []byte
}

func genTestsC() []testC {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testC, 100)
	letters := []byte("abcdefghijklmnopqrstuvwxyz")
	for i := range tests {
		n := r.Intn(100) + 1
		b := make([]byte, n)
		for j := range b {
			b[j] = letters[r.Intn(26)]
		}
		k := r.Intn(26) + 1
		perm := r.Perm(26)
		allowed := make([]byte, k)
		for j := 0; j < k; j++ {
			allowed[j] = letters[perm[j]]
		}
		tests[i] = testC{n: n, k: k, s: string(b), allowed: allowed}
	}
	return tests
}

func solveC(tc testC) int64 {
	allowed := make(map[byte]bool)
	for _, ch := range tc.allowed {
		allowed[ch] = true
	}
	var cur int64
	var ans int64
	for i := 0; i < tc.n; i++ {
		if allowed[tc.s[i]] {
			cur++
		} else {
			ans += cur * (cur + 1) / 2
			cur = 0
		}
	}
	ans += cur * (cur + 1) / 2
	return ans
}

func runBinary(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", bin)
	} else {
		cmd = exec.CommandContext(ctx, bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("timeout")
	}
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsC()
	for i, tc := range tests {
		var input strings.Builder
		fmt.Fprintf(&input, "%d %d\n", tc.n, tc.k)
		fmt.Fprintln(&input, tc.s)
		for j := 0; j < tc.k; j++ {
			if j > 0 {
				input.WriteByte(' ')
			}
			input.WriteByte(tc.allowed[j])
		}
		input.WriteByte('\n')
		expected := fmt.Sprintf("%d", solveC(tc))
		out, err := runBinary(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != expected {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected:%s\ngot:%s\n", i+1, input.String(), expected, out)
			os.Exit(1)
		}
	}
	fmt.Println("Accepted")
}
