package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func precompute() []string {
	res := make([]string, 0, 63)
	for i := 0; i <= 60; i++ {
		res = append(res, fmt.Sprintf("%d", 1<<uint(i)))
	}
	return res
}

func lcsPrefix(a, b string) int {
	j := 0
	for i := 0; i < len(a) && j < len(b); i++ {
		if a[i] == b[j] {
			j++
		}
	}
	return j
}

func solve(n string, powers []string) int {
	minOps := len(n) + 1e9
	for _, p := range powers {
		l := lcsPrefix(n, p)
		ops := len(n) - l + len(p) - l
		if ops < minOps {
			minOps = ops
		}
	}
	return minOps
}

type testCase struct {
	n string
}

func genTests() []testCase {
	r := rand.New(rand.NewSource(4))
	tests := make([]testCase, 100)
	for i := range tests {
		l := r.Intn(15) + 1
		var sb strings.Builder
		if r.Intn(2) == 0 {
			sb.WriteByte(byte('1' + byte(r.Intn(9))))
		} else {
			sb.WriteByte('1')
		}
		for j := 1; j < l; j++ {
			sb.WriteByte(byte('0' + r.Intn(10)))
		}
		tests[i].n = sb.String()
	}
	return tests
}

func run(bin, input string) (string, error) {
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	powers := precompute()
	tests := genTests()
	for i, tc := range tests {
		input := fmt.Sprintf("1\n%s\n", tc.n)
		want := fmt.Sprintf("%d", solve(tc.n, powers))
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			return
		}
		if strings.TrimSpace(got) != want {
			fmt.Printf("case %d failed: expected %s got %s\n", i+1, want, got)
			fmt.Printf("input:\n%s", input)
			return
		}
	}
	fmt.Println("All tests passed")
}
