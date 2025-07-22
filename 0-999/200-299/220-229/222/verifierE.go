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

type test struct {
	input    string
	expected string
}

func charToIndex(c byte) int {
	if c >= 'a' && c <= 'z' {
		return int(c - 'a')
	}
	return int(c-'A') + 26
}

func indexToChar(i int) byte {
	if i < 26 {
		return byte('a' + i)
	}
	return byte('A' + i - 26)
}

func mul(a, b [][]int64, m int) [][]int64 {
	c := make([][]int64, m)
	for i := 0; i < m; i++ {
		c[i] = make([]int64, m)
		for k := 0; k < m; k++ {
			if a[i][k] == 0 {
				continue
			}
			aik := a[i][k]
			for j := 0; j < m; j++ {
				c[i][j] = (c[i][j] + aik*b[k][j]) % mod
			}
		}
	}
	return c
}

func matPow(a [][]int64, e int64, m int) [][]int64 {
	res := make([][]int64, m)
	for i := 0; i < m; i++ {
		res[i] = make([]int64, m)
		res[i][i] = 1
	}
	for e > 0 {
		if e&1 == 1 {
			res = mul(res, a, m)
		}
		a = mul(a, a, m)
		e >>= 1
	}
	return res
}

func solve(input string) string {
	reader := strings.NewReader(input)
	var n int64
	var m, k int
	if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
		return ""
	}
	adj := make([][]int64, m)
	for i := 0; i < m; i++ {
		adj[i] = make([]int64, m)
		for j := 0; j < m; j++ {
			adj[i][j] = 1
		}
	}
	for i := 0; i < k; i++ {
		var s string
		fmt.Fscan(reader, &s)
		if len(s) != 2 {
			continue
		}
		u := charToIndex(s[0])
		v := charToIndex(s[1])
		if u < m && v < m {
			adj[u][v] = 0
		}
	}
	var result int64
	if n == 1 {
		result = int64(m) % mod
	} else {
		p := matPow(adj, n-1, m)
		var sum int64
		for i := 0; i < m; i++ {
			for j := 0; j < m; j++ {
				sum = (sum + p[i][j]) % mod
			}
		}
		result = sum
	}
	return fmt.Sprintf("%d", result)
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(46))
	var tests []test
	fixed := []string{
		"1 1 0\n",
		"2 2 1\nab\n",
	}
	for _, f := range fixed {
		tests = append(tests, test{f, solve(f)})
	}
	for len(tests) < 100 {
		n := int64(rng.Intn(5) + 1)
		m := rng.Intn(4) + 1
		maxPairs := m * m
		k := rng.Intn(maxPairs + 1)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, k))
		used := map[[2]int]bool{}
		for i := 0; i < k; i++ {
			for {
				u := rng.Intn(m)
				v := rng.Intn(m)
				if !used[[2]int{u, v}] {
					used[[2]int{u, v}] = true
					sb.WriteByte(indexToChar(u))
					sb.WriteByte(indexToChar(v))
					sb.WriteByte('\n')
					break
				}
			}
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
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
