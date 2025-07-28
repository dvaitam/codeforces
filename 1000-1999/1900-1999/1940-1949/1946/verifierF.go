package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type test struct {
	input    string
	expected string
}

func solveOne(a []int, l, r int) int64 {
	n := r - l + 1
	dp := make([]int64, n)
	var ans int64
	for i := l; i <= r; i++ {
		dp[i-l] = 1
		for j := l; j < i; j++ {
			if a[i]%a[j] == 0 {
				dp[i-l] += dp[j-l]
			}
		}
		ans += dp[i-l]
	}
	return ans
}

func solve(input string) string {
	reader := strings.NewReader(input)
	var t int
	fmt.Fscan(reader, &t)
	var out strings.Builder
	for ; t > 0; t-- {
		var n, q int
		fmt.Fscan(reader, &n, &q)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		for ; q > 0; q-- {
			var l, r int
			fmt.Fscan(reader, &l, &r)
			l--
			r--
			ans := solveOne(a, l, r)
			out.WriteString(fmt.Sprintf("%d\n", ans))
		}
	}
	return strings.TrimSpace(out.String())
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(51))
	var tests []test
	fixed := []string{
		"1\n1 1\n1\n1 1\n",
	}
	for _, f := range fixed {
		tests = append(tests, test{f, solve(f)})
	}
	for len(tests) < 100 {
		n := rng.Intn(5) + 1
		q := rng.Intn(3) + 1
		a := make([]int, n)
		for i := 0; i < n; i++ {
			a[i] = rng.Intn(6) + 1
		}
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d %d\n", n, q))
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", a[i])
		}
		sb.WriteByte('\n')
		for j := 0; j < q; j++ {
			l := rng.Intn(n) + 1
			r := rng.Intn(n-l+1) + l
			fmt.Fprintf(&sb, "%d %d\n", l, r)
		}
		inp := sb.String()
		tests = append(tests, test{inp, solve(inp)})
	}
	return tests
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		got, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Printf("Runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(t.expected) {
			fmt.Printf("Wrong answer on test %d\nInput:\n%sExpected:%s\nGot:%s\n", i+1, t.input, t.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
