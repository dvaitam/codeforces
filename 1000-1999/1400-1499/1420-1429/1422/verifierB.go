package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type test struct {
	input    string
	expected string
}

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func solve(input string) string {
	r := strings.NewReader(input)
	var t int
	fmt.Fscan(r, &t)
	var out strings.Builder
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(r, &n, &m)
		a := make([][]int64, n)
		for i := 0; i < n; i++ {
			a[i] = make([]int64, m)
			for j := 0; j < m; j++ {
				fmt.Fscan(r, &a[i][j])
			}
		}
		var res int64
		for i := 0; i <= (n-1)/2; i++ {
			for j := 0; j <= (m-1)/2; j++ {
				i2 := n - 1 - i
				j2 := m - 1 - j
				vals := make([]int64, 0, 4)
				vals = append(vals, a[i][j])
				if i2 != i {
					vals = append(vals, a[i2][j])
				}
				if j2 != j {
					vals = append(vals, a[i][j2])
				}
				if i2 != i && j2 != j {
					vals = append(vals, a[i2][j2])
				}
				sort.Slice(vals, func(x, y int) bool { return vals[x] < vals[y] })
				median := vals[len(vals)/2]
				for _, v := range vals {
					res += abs(v - median)
				}
			}
		}
		fmt.Fprintf(&out, "%d\n", res)
	}
	return strings.TrimRight(out.String(), "\n")
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(43))
	tests := []test{}
	fixed := []string{
		"1\n1 1\n5\n",
		"1\n2 2\n1 2\n3 4\n",
	}
	for _, f := range fixed {
		tests = append(tests, test{f, solve(f)})
	}
	for len(tests) < 100 {
		n := rng.Intn(4) + 1
		m := rng.Intn(4) + 1
		var sb strings.Builder
		sb.WriteString("1\n")
		fmt.Fprintf(&sb, "%d %d\n", n, m)
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				val := rng.Intn(11) - 5
				if j > 0 {
					sb.WriteByte(' ')
				}
				fmt.Fprintf(&sb, "%d", val)
			}
			sb.WriteByte('\n')
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
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
