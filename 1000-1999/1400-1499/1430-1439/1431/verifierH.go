package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type testCase struct {
	name  string
	input string
}

func solveRef(input string) (int, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	var n, m int
	var k int64
	if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
		return 0, err
	}
	a := make([]int64, n+1)
	b := make([]int64, m+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	for j := 1; j <= m; j++ {
		fmt.Fscan(reader, &b[j])
	}
	c := make([][]int64, n+1)
	for i := 1; i <= n; i++ {
		c[i] = make([]int64, m+1)
		for j := 1; j <= m; j++ {
			fmt.Fscan(reader, &c[i][j])
		}
	}
	origG := make([][]int64, n+1)
	idxX := make([][]int, n+1)
	idxY := make([][]int, n+1)
	for i := 0; i <= n; i++ {
		origG[i] = make([]int64, m+1)
		idxX[i] = make([]int, m+1)
		idxY[i] = make([]int, m+1)
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			v := c[i][j]
			x, y := i, j
			if origG[i-1][j] > v {
				v = origG[i-1][j]
				x = idxX[i-1][j]
				y = idxY[i-1][j]
			}
			if origG[i][j-1] > v {
				v = origG[i][j-1]
				x = idxX[i][j-1]
				y = idxY[i][j-1]
			}
			origG[i][j] = v
			idxX[i][j] = x
			idxY[i][j] = y
		}
	}
	A0, B0 := 1, 1
	for A0 < n && a[A0+1] == 0 {
		A0++
	}
	for B0 < m && b[B0+1] == 0 {
		B0++
	}
	type state struct{ A, B int }
	var states []state
	var S int64
	A, B := A0, B0
	runs0 := 0
	for A < n || B < m {
		nextT := int64(math.MaxInt64 / 4)
		if A < n && a[A+1] < nextT {
			nextT = a[A+1]
		}
		if B < m && b[B+1] < nextT {
			nextT = b[B+1]
		}
		g := origG[A][B]
		need := nextT - S
		t := 0
		if need > 0 {
			t = int((need + g - 1) / g)
		}
		runs0 += t
		S += int64(t) * g
		states = append(states, state{A, B})
		for A < n && a[A+1] <= S {
			A++
		}
		for B < m && b[B+1] <= S {
			B++
		}
	}
	ans := runs0
	candidates := make(map[[2]int]struct{})
	for _, st := range states {
		x := idxX[st.A][st.B]
		y := idxY[st.A][st.B]
		candidates[[2]int{x, y}] = struct{}{}
	}
	if k == 0 {
		return ans, nil
	}
	for key := range candidates {
		x, y := key[0], key[1]
		var S2 int64
		A2, B2 := A0, B0
		runs := 0
		for A2 < n || B2 < m {
			nextT := int64(math.MaxInt64 / 4)
			if A2 < n && a[A2+1] < nextT {
				nextT = a[A2+1]
			}
			if B2 < m && b[B2+1] < nextT {
				nextT = b[B2+1]
			}
			g := origG[A2][B2]
			if A2 >= x && B2 >= y {
				if c[x][y]+k > g {
					g = c[x][y] + k
				}
			}
			need := nextT - S2
			t := 0
			if need > 0 {
				t = int((need + g - 1) / g)
			}
			runs += t
			if runs >= ans {
				break
			}
			S2 += int64(t) * g
			for A2 < n && a[A2+1] <= S2 {
				A2++
			}
			for B2 < m && b[B2+1] <= S2 {
				B2++
			}
		}
		if runs < ans {
			ans = runs
		}
	}
	return ans, nil
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

func parseIntOutput(out string) (int, error) {
	if out == "" {
		return 0, fmt.Errorf("empty output")
	}
	var val int
	if _, err := fmt.Sscan(out, &val); err != nil {
		return 0, fmt.Errorf("failed to parse integer: %v", err)
	}
	return val, nil
}

func makeCase(name string, n, m int, k int64, a []int64, b []int64, c [][]int64) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, m, k)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", a[i])
	}
	sb.WriteByte('\n')
	for j := 0; j < m; j++ {
		if j > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", b[j])
	}
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", c[i][j])
		}
		sb.WriteByte('\n')
	}
	return testCase{name: name, input: sb.String()}
}

func handcraftedTests() []testCase {
	a := []int64{0, 5, 10}
	b := []int64{0, 3}
	c := [][]int64{
		{1, 4},
		{5, 2},
		{3, 6},
	}
	test1 := makeCase("hand1", 3, 2, 2, a, b, c)

	a2 := []int64{0, 2}
	b2 := []int64{0, 1, 3}
	c2 := [][]int64{
		{2, 1, 3},
		{4, 5, 2},
	}
	test2 := makeCase("hand2", 2, 3, 1, a2, b2, c2)

	return []testCase{test1, test2}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	for idx := 0; idx < 40; idx++ {
		n := rng.Intn(4) + 1
		m := rng.Intn(4) + 1
		k := rng.Int63n(5)
		a := make([]int64, n)
		b := make([]int64, m)
		for i := 0; i < n; i++ {
			if i == 0 {
				a[i] = 0
			} else {
				a[i] = a[i-1] + int64(rng.Intn(5))
			}
		}
		for j := 0; j < m; j++ {
			if j == 0 {
				b[j] = 0
			} else {
				b[j] = b[j-1] + int64(rng.Intn(5))
			}
		}
		c := make([][]int64, n)
		for i := 0; i < n; i++ {
			c[i] = make([]int64, m)
			for j := 0; j < m; j++ {
				c[i][j] = int64(rng.Intn(5) + 1)
			}
		}
		tests = append(tests, makeCase(fmt.Sprintf("rand_%d", idx+1), n, m, k, a, b, c))
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierH.go /path/to/binary")
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
			fmt.Printf("test %d (%s) runtime error: %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		got, err := parseIntOutput(out)
		if err != nil {
			fmt.Printf("test %d (%s) invalid output: %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, out)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("test %d (%s) mismatch\ninput:\n%s\nexpect:%d\nactual:%d\n", idx+1, tc.name, tc.input, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
