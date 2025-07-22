package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type testCase struct {
	in  string
	out string
}

func countLuckyDigits(x int) int {
	c := 0
	for x > 0 {
		d := x % 10
		if d == 4 || d == 7 {
			c++
		}
		x /= 10
	}
	return c
}

func solve(m int) int64 {
	digits := make([]int, m+1)
	for i := 1; i <= m; i++ {
		digits[i] = countLuckyDigits(i)
	}
	nums := make([]int, m)
	for i := 0; i < m; i++ {
		nums[i] = i + 1
	}
	used := make([]bool, m)
	var perm [7]int
	var cnt int64
	var dfs func(pos int)
	dfs = func(pos int) {
		if pos == 7 {
			sum := 0
			for i := 1; i < 7; i++ {
				sum += digits[perm[i]]
			}
			if digits[perm[0]] > sum {
				cnt++
			}
			return
		}
		for i := 0; i < m; i++ {
			if !used[i] {
				used[i] = true
				perm[pos] = nums[i]
				dfs(pos + 1)
				used[i] = false
			}
		}
	}
	dfs(0)
	return cnt % 1000000007
}

func generateTests() []testCase {
	var tests []testCase
	m := 7
	for len(tests) < 100 {
		out := solve(m)
		tests = append(tests, testCase{in: fmt.Sprintf("%d\n", m), out: fmt.Sprint(out)})
		m++
		if m > 9 {
			m = 7
		}
	}
	return tests
}

func runTest(bin string, tc testCase) (string, error) {
	cmd := exec.Command(bin)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", err
	}
	go func() {
		defer stdin.Close()
		stdin.Write([]byte(tc.in))
	}()
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	if strings.HasSuffix(bin, ".go") {
		tmp, err := ioutil.TempFile("", "solbin*")
		if err != nil {
			fmt.Println("cannot create temp file:", err)
			os.Exit(1)
		}
		tmp.Close()
		exec.Command("go", "build", "-o", tmp.Name(), bin).Run()
		bin = tmp.Name()
		defer os.Remove(bin)
	}
	tests := generateTests()
	for i, tc := range tests {
		got, err := runTest(bin, tc)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != tc.out {
			fmt.Printf("wrong answer on test %d\ninput: %sexpected: %s\ngot: %s\n", i+1, tc.in, tc.out, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
