package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func singleValue(x int64) int64 {
	cnt := int64(1)
	for x > 3 {
		x = x/2 + 1
		cnt++
	}
	return cnt
}

func isSpecial(x int64) int64 {
	if x <= 2 {
		return 0
	}
	y := x - 1
	if y&(y-1) == 0 {
		return 1
	}
	return 0
}

func solveRef(input string) (string, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return "", err
	}
	var outputs []string
	for tc := 0; tc < t; tc++ {
		var n, q int
		fmt.Fscan(reader, &n, &q)
		prefSum := make([]int64, n+1)
		prefSpecial := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			var val int64
			fmt.Fscan(reader, &val)
			prefSum[i] = prefSum[i-1] + singleValue(val)
			prefSpecial[i] = prefSpecial[i-1] + isSpecial(val)
		}
		for qi := 0; qi < q; qi++ {
			var l, r int
			fmt.Fscan(reader, &l, &r)
			base := prefSum[r] - prefSum[l-1]
			cnt := prefSpecial[r] - prefSpecial[l-1]
			outputs = append(outputs, fmt.Sprintf("%d", base+cnt/2))
		}
	}
	return strings.Join(outputs, "\n"), nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstderr: %s", err, stderr.String())
	}
	if stderr.Len() > 0 {
		return "", fmt.Errorf("unexpected stderr output: %s", stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

type testCase struct {
	name  string
	input string
}

func makeCase(name string, cases []struct {
	n int
	q int
	a []int64
	l []int
	r []int
}) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(cases))
	for _, tc := range cases {
		fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.q)
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", v)
		}
		sb.WriteByte('\n')
		for i := 0; i < tc.q; i++ {
			fmt.Fprintf(&sb, "%d %d\n", tc.l[i], tc.r[i])
		}
	}
	return testCase{name: name, input: sb.String()}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	for idx := 0; idx < 30; idx++ {
		tcCount := rng.Intn(2) + 1
		cases := make([]struct {
			n int
			q int
			a []int64
			l []int
			r []int
		}, tcCount)
		for i := 0; i < tcCount; i++ {
			n := rng.Intn(5) + 1
			q := rng.Intn(5) + 1
			cases[i].n = n
			cases[i].q = q
			cases[i].a = make([]int64, n)
			for j := 0; j < n; j++ {
				cases[i].a[j] = int64(rng.Intn(20) + 2)
			}
			cases[i].l = make([]int, q)
			cases[i].r = make([]int, q)
			for j := 0; j < q; j++ {
				l := rng.Intn(n) + 1
				r := rng.Intn(n-l+1) + l
				cases[i].l[j] = l
				cases[i].r[j] = r
			}
		}
		tests = append(tests, makeCase(fmt.Sprintf("rand_%d", idx+1), cases))
	}
	return tests
}

func handcraftedTests() []testCase {
	cases := []struct {
		n int
		q int
		a []int64
		l []int
		r []int
	}{
		{
			n: 5,
			q: 5,
			a: []int64{4, 3, 2, 5, 6},
			l: []int{1, 1, 2, 4, 5},
			r: []int{1, 2, 3, 5, 5},
		},
		{
			n: 5,
			q: 5,
			a: []int64{10, 13, 14, 15, 9},
			l: []int{1, 2, 3, 4, 5},
			r: []int{1, 2, 5, 5, 5},
		},
	}
	return []testCase{makeCase("handcrafted", cases)}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := append(handcraftedTests(), randomTests()...)
	for idx, tc := range tests {
		expect, err := solveRef(tc.input)
		if err != nil {
			fmt.Printf("failed to compute reference for %s: %v\n", tc.name, err)
			os.Exit(1)
		}
		out, err := runCandidate(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d (%s) runtime error: %v\ninput:\n%s\n", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		expFields := strings.Fields(expect)
		gotFields := strings.Fields(out)
		if len(expFields) != len(gotFields) {
			fmt.Printf("test %d (%s) mismatch in answers count\ninput:\n%s\nexpect:\n%s\nactual:\n%s\n", idx+1, tc.name, tc.input, expect, out)
			os.Exit(1)
		}
		for i := range expFields {
			if expFields[i] != gotFields[i] {
				fmt.Printf("test %d (%s) mismatch at answer %d\ninput:\n%s\nexpect:\n%s\nactual:\n%s\n", idx+1, tc.name, i+1, tc.input, expect, out)
				os.Exit(1)
			}
		}
	}
	fmt.Println("all tests passed")
}
