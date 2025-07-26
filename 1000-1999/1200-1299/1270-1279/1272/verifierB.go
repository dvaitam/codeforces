package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type testB struct{ s string }

func genTestsB() []testB {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testB, 100)
	letters := []byte{'L', 'R', 'U', 'D'}
	for i := range tests {
		n := r.Intn(50) + 1
		b := make([]byte, n)
		for j := range b {
			b[j] = letters[r.Intn(len(letters))]
		}
		tests[i] = testB{s: string(b)}
	}
	return tests
}

func solveB(t testB) (int, string) {
	cnt := make(map[rune]int)
	for _, ch := range t.s {
		cnt[ch]++
	}
	horiz := min(cnt['L'], cnt['R'])
	vert := min(cnt['U'], cnt['D'])
	if horiz == 0 && vert == 0 {
		return 0, ""
	}
	if horiz == 0 {
		return 2, "UD"
	}
	if vert == 0 {
		return 2, "LR"
	}
	total := (horiz + vert) * 2
	var sb strings.Builder
	for i := 0; i < vert; i++ {
		sb.WriteByte('U')
	}
	for i := 0; i < horiz; i++ {
		sb.WriteByte('R')
	}
	for i := 0; i < vert; i++ {
		sb.WriteByte('D')
	}
	for i := 0; i < horiz; i++ {
		sb.WriteByte('L')
	}
	return total, sb.String()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
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
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsB()
	var input strings.Builder
	fmt.Fprintf(&input, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintln(&input, tc.s)
	}
	outputs := make([]struct {
		len  int
		path string
	}, len(tests))
	for i, tc := range tests {
		outputs[i].len, outputs[i].path = solveB(tc)
	}
	out, err := runBinary(bin, input.String())
	if err != nil {
		fmt.Fprintf(os.Stderr, "runtime error: %v\n", err)
		os.Exit(1)
	}
	scanner := bufio.NewScanner(strings.NewReader(out))
	for i, exp := range outputs {
		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, "wrong output format on test %d\n", i+1)
			os.Exit(1)
		}
		length, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
		if err != nil {
			fmt.Fprintf(os.Stderr, "non-integer output on test %d\n", i+1)
			os.Exit(1)
		}
		var path string
		if length > 0 {
			if !scanner.Scan() {
				fmt.Fprintf(os.Stderr, "missing path on test %d\n", i+1)
				os.Exit(1)
			}
			path = strings.TrimSpace(scanner.Text())
		}
		if length != exp.len || path != exp.path {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\n", i+1)
			os.Exit(1)
		}
	}
	if scanner.Scan() {
		fmt.Fprintln(os.Stderr, "extra output")
		os.Exit(1)
	}
	fmt.Println("Accepted")
}
