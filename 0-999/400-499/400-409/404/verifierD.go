package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const mod = 1000000007

type testCaseD struct{ s string }

func generateTestsD() []testCaseD {
	r := rand.New(rand.NewSource(4))
	tests := make([]testCaseD, 0, 100)
	for len(tests) < 100 {
		n := 1 + r.Intn(6)
		var b strings.Builder
		for i := 0; i < n; i++ {
			switch r.Intn(4) {
			case 0:
				b.WriteByte('*')
			case 1:
				b.WriteByte('?')
			case 2:
				b.WriteByte('0' + byte(r.Intn(3)))
			case 3:
				b.WriteByte('0' + byte(r.Intn(3)))
			}
		}
		tests = append(tests, testCaseD{s: b.String()})
	}
	return tests
}

func waysD(board []byte) int64 {
	n := len(board)
	for i := 0; i < n; i++ {
		if board[i] == '*' {
			continue
		}
		c := int(board[i] - '0')
		cnt := 0
		if i-1 >= 0 && board[i-1] == '*' {
			cnt++
		}
		if i+1 < n && board[i+1] == '*' {
			cnt++
		}
		if cnt != c {
			return 0
		}
	}
	return 1
}

func solveDRec(s []byte, idx int, cur []byte) int64 {
	if idx == len(s) {
		return waysD(cur) % mod
	}
	ch := s[idx]
	if ch == '?' {
		sum := int64(0)
		for _, c := range []byte{'*', '0', '1', '2'} {
			cur[idx] = c
			sum = (sum + solveDRec(s, idx+1, cur)) % mod
		}
		return sum
	}
	cur[idx] = ch
	return solveDRec(s, idx+1, cur)
}

func solveDString(s string) int64 {
	b := []byte(s)
	cur := make([]byte, len(b))
	return solveDRec(b, 0, cur)
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("run failed: %v\n%s", err, errb.String())
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTestsD()
	for i, t := range tests {
		out, err := runBinary(bin, t.s+"\n")
		if err != nil {
			fmt.Printf("test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		expect := solveDString(t.s)
		got := strings.TrimSpace(out)
		if fmt.Sprint(expect) != got {
			fmt.Printf("test %d failed: expected %d got %s\n", i+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
