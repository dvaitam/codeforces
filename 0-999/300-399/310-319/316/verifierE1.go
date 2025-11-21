package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const mod = 1000000000

func solveRef(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return ""
	}
	a := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	f := make([]int, n+1)
	if n >= 0 {
		f[0] = 1
	}
	if n >= 1 {
		f[1] = 1
	}
	for i := 2; i <= n; i++ {
		f[i] = f[i-1] + f[i-2]
		if f[i] >= mod {
			f[i] -= mod
		}
	}
	var outputs []string
	for q := 0; q < m; q++ {
		var t int
		fmt.Fscan(reader, &t)
		switch t {
		case 1:
			var x, v int
			fmt.Fscan(reader, &x, &v)
			a[x] = v % mod
		case 2:
			var l, r int
			fmt.Fscan(reader, &l, &r)
			sum := 0
			for idx := 0; idx <= r-l; idx++ {
				sum += f[idx] * a[l+idx]
				sum %= mod
			}
			outputs = append(outputs, fmt.Sprintf("%d", sum%mod))
		case 3:
			var l, r, d int
			fmt.Fscan(reader, &l, &r, &d)
			for i := l; i <= r; i++ {
				a[i] = (a[i] + d) % mod
			}
		}
	}
	return strings.Join(outputs, "\n")
}

type testCase struct {
	name  string
	input string
	// outputs compared after running exe
}

func makeCase(name string, n int, a []int, ops []string) testCase {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, len(ops)))
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	for _, op := range ops {
		sb.WriteString(op)
		sb.WriteByte('\n')
	}
	return testCase{name: name, input: sb.String()}
}

func randomArray(rng *rand.Rand, n int, maxVal int) []int {
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(maxVal)
	}
	return arr
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(316))
	var tests []testCase
	gen := func(prefix string, count, maxN, maxQ int) {
		for i := 0; i < count; i++ {
			n := rng.Intn(maxN) + 1
			q := rng.Intn(maxQ) + 1
			a := randomArray(rng, n, 1000)
			var ops []string
			for j := 0; j < q; j++ {
				t := rng.Intn(3) + 1
				switch t {
				case 1:
					x := rng.Intn(n) + 1
					v := rng.Intn(100000)
					ops = append(ops, fmt.Sprintf("1 %d %d", x, v))
				case 2:
					l := rng.Intn(n) + 1
					r := rng.Intn(n-l+1) + l
					ops = append(ops, fmt.Sprintf("2 %d %d", l, r))
				case 3:
					l := rng.Intn(n) + 1
					r := rng.Intn(n-l+1) + l
					d := rng.Intn(100000)
					ops = append(ops, fmt.Sprintf("3 %d %d %d", l, r, d))
				}
			}
			tests = append(tests, makeCase(fmt.Sprintf("%s_%d", prefix, i+1), n, a, ops))
		}
	}
	gen("small", 80, 5, 10)
	gen("medium", 80, 20, 40)
	gen("large", 40, 50, 60)
	return tests
}

func e1SpecificTests() []testCase {
	return []testCase{
		makeCase("single_update_query", 1, []int{5}, []string{"2 1 1", "1 1 3", "2 1 1"}),
		makeCase("all_type2", 3, []int{1, 2, 3}, []string{"2 1 3", "2 2 2", "2 1 1"}),
		makeCase("only_type1", 3, []int{0, 0, 0}, []string{"1 1 5", "1 2 6", "1 3 7"}),
		makeCase("type3_presence", 4, []int{1, 2, 3, 4}, []string{"3 1 4 5", "2 1 4"}),
		makeCase("mixed", 5, []int{1, 2, 3, 4, 5}, []string{
			"2 1 5",
			"1 3 10",
			"2 2 4",
			"3 2 5 3",
			"2 1 3",
		}),
	}
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := append(e1SpecificTests(), randomTests()...)
	for idx, tc := range tests {
		expect := solveRef(tc.input)
		got, err := runCandidate(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d (%s) runtime error: %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		if got != strings.TrimSpace(expect) {
			fmt.Printf("test %d (%s) failed\ninput:\n%s\nexpect:\n%s\nactual:\n%s\n", idx+1, tc.name, tc.input, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
