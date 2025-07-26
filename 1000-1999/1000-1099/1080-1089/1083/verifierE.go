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
)

type test struct {
	input    string
	expected string
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

type square struct{ x, y, w int64 }

func solveE(input string) string {
	rdr := bufio.NewReader(strings.NewReader(input))
	var n int
	fmt.Fscan(rdr, &n)
	s := make([]square, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(rdr, &s[i].x, &s[i].y, &s[i].w)
	}
	sort.Slice(s, func(i, j int) bool { return s[i].y > s[j].y })
	f := make([]int64, n)
	q := make([]int, 0, n)
	var ans int64
	eps := 1e-12
	getSlope := func(u, v int) float64 {
		dx := float64(s[v].x - s[u].x)
		if math.Abs(dx) < eps {
			return math.Inf(1)
		}
		return float64(f[v]-f[u]) / dx
	}
	l := 0
	for i := 0; i < n; i++ {
		for len(q)-l >= 2 && getSlope(q[l], q[l+1]) >= float64(s[i].y) {
			l++
		}
		fi := s[i].x*s[i].y - s[i].w
		if len(q)-l > 0 {
			j := q[l]
			val := f[j] + (s[i].x-s[j].x)*s[i].y - s[i].w
			if val > fi {
				fi = val
			}
		}
		f[i] = fi
		if fi > ans {
			ans = fi
		}
		for len(q)-l >= 2 && getSlope(q[len(q)-2], q[len(q)-1]) <= getSlope(q[len(q)-1], i) {
			q = q[:len(q)-1]
		}
		q = append(q, i)
	}
	return fmt.Sprintf("%d", ans)
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(46))
	var tests []test
	for len(tests) < 100 {
		n := rng.Intn(6) + 1
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", n)
		squares := make([]square, n)
		for i := 0; i < n; i++ {
			squares[i].x = int64(rng.Intn(10))
			squares[i].y = int64(rng.Intn(10))
			squares[i].w = int64(rng.Intn(10))
			fmt.Fprintf(&sb, "%d %d %d\n", squares[i].x, squares[i].y, squares[i].w)
		}
		input := sb.String()
		expected := solveE(input)
		tests = append(tests, test{input, expected})
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		out, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(t.expected) {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected:%s\ngot:%s\n", i+1, t.input, t.expected, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
