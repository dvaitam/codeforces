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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func getExpectedLength(s string) int {
	cnt := make(map[rune]int)
	for _, ch := range s {
		cnt[ch]++
	}
	horiz := min(cnt['L'], cnt['R'])
	vert := min(cnt['U'], cnt['D'])
	if horiz == 0 && vert == 0 {
		return 0
	}
	if horiz == 0 {
		return 2
	}
	if vert == 0 {
		return 2
	}
	return (horiz + vert) * 2
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

type Point struct {
	x, y int
}

func verifyPath(s string, path string) error {
	sCnt := make(map[rune]int)
	for _, c := range s {
		sCnt[c]++
	}
	pCnt := make(map[rune]int)
	for _, c := range path {
		pCnt[c]++
	}
	for k, v := range pCnt {
		if v > sCnt[k] {
			return fmt.Errorf("used more %c than available", k)
		}
	}

	visited := make(map[Point]bool)
	curr := Point{0, 0}
	visited[curr] = true

	for i, c := range path {
		switch c {
		case 'L':
			curr.x--
		case 'R':
			curr.x++
		case 'U':
			curr.y++
		case 'D':
			curr.y--
		default:
			return fmt.Errorf("invalid char %c", c)
		}

		if i == len(path)-1 {
			if curr.x != 0 || curr.y != 0 {
				return fmt.Errorf("path did not end at (0,0)")
			}
		} else {
			if visited[curr] {
				return fmt.Errorf("visited (%d,%d) twice", curr.x, curr.y)
			}
			visited[curr] = true
		}
	}
	return nil
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

	out, err := runBinary(bin, input.String())
	if err != nil {
		fmt.Fprintf(os.Stderr, "runtime error: %v\n", err)
		os.Exit(1)
	}
	scanner := bufio.NewScanner(strings.NewReader(out))
	scanner.Split(bufio.ScanWords)

	for i, tc := range tests {
		expLen := getExpectedLength(tc.s)

		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, "missing output for test %d\n", i+1)
			os.Exit(1)
		}
		token := scanner.Text()
		length, err := strconv.Atoi(token)
		if err != nil {
			fmt.Fprintf(os.Stderr, "non-integer output on test %d: %q\n", i+1, token)
			os.Exit(1)
		}

		if length != expLen {
			fmt.Fprintf(os.Stderr, "wrong length on test %d: expected %d, got %d\n", i+1, expLen, length)
			os.Exit(1)
		}

		if length > 0 {
			if !scanner.Scan() {
				fmt.Fprintf(os.Stderr, "missing path on test %d\n", i+1)
				os.Exit(1)
			}
			path := scanner.Text()
			if len(path) != length {
				fmt.Fprintf(os.Stderr, "path length mismatch on test %d: declared %d, actual %d\n", i+1, length, len(path))
				os.Exit(1)
			}
			if err := verifyPath(tc.s, path); err != nil {
				fmt.Fprintf(os.Stderr, "invalid path on test %d: %v\n", i+1, err)
				os.Exit(1)
			}
		}
	}
	if scanner.Scan() {
		fmt.Fprintln(os.Stderr, "extra output")
		os.Exit(1)
	}
	fmt.Println("Accepted")
}