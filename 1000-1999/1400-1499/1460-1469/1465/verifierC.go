package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type rook struct {
	x, y int
}

type testCaseC struct {
	n     int
	rooks []rook
}

func genTestsC() []testCaseC {
	rand.Seed(44)
	tests := make([]testCaseC, 100)
	for i := range tests {
		n := rand.Intn(20) + 2
		m := rand.Intn(n-1) + 1
		rows := rand.Perm(n)[:m]
		cols := rand.Perm(n)[:m]
		rs := make([]rook, m)
		for j := 0; j < m; j++ {
			rs[j] = rook{rows[j] + 1, cols[j] + 1}
		}
		tests[i] = testCaseC{n: n, rooks: rs}
	}
	// add sample-like cases
	tests = append(tests, testCaseC{n: 3, rooks: []rook{{2, 3}}})
	tests = append(tests, testCaseC{n: 3, rooks: []rook{{2, 1}, {1, 2}}})
	tests = append(tests, testCaseC{n: 5, rooks: []rook{{2, 3}, {3, 1}, {1, 2}}})
	tests = append(tests, testCaseC{n: 5, rooks: []rook{{4, 5}, {5, 1}, {2, 2}, {3, 3}}})
	return tests
}

func solveC(tc testCaseC) int {
	mp := make(map[int]int)
	notDiag := 0
	for _, r := range tc.rooks {
		if r.x != r.y {
			mp[r.x] = r.y
			notDiag++
		}
	}
	vis := make(map[int]int)
	cycles := 0
	for x := range mp {
		if vis[x] != 0 {
			continue
		}
		cur := x
		for vis[cur] == 0 {
			vis[cur] = x
			next, ok := mp[cur]
			if !ok {
				break
			}
			cur = next
		}
		if vis[cur] == x {
			cycles++
		}
	}
	return notDiag + cycles
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	cmd.Env = append(os.Environ(), "GOMAXPROCS=1")
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("runtime error: %v", err)
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsC()
	for i, tc := range tests {
		var sb strings.Builder
		fmt.Fprintf(&sb, "1\n%d %d\n", tc.n, len(tc.rooks))
		for _, r := range tc.rooks {
			fmt.Fprintf(&sb, "%d %d\n", r.x, r.y)
		}
		input := sb.String()
		exp := solveC(tc)
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: %v\noutput:\n%s", i+1, err, got)
			os.Exit(1)
		}
		got = strings.TrimSpace(got)
		if got != fmt.Sprintf("%d", exp) {
			fmt.Printf("test %d failed: expected %d got %s\ninput:\n%s", i+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
