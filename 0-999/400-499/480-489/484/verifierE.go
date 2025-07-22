package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type test struct {
	in  string
	out string
}

func solve(h []int, queries [][3]int) []int {
	res := make([]int, len(queries))
	for qi, q := range queries {
		l, r, w := q[0], q[1], q[2]
		best := 0
		for i := l; i+w <= r+1; i++ {
			mn := math.MaxInt32
			for j := i; j < i+w; j++ {
				if h[j] < mn {
					mn = h[j]
				}
			}
			if mn > best {
				best = mn
			}
		}
		res[qi] = best
	}
	return res
}

func generateTests() []test {
	rnd := rand.New(rand.NewSource(5))
	tests := make([]test, 0, 100)
	for len(tests) < 100 {
		n := rnd.Intn(10) + 1
		h := make([]int, n)
		var in bytes.Buffer
		fmt.Fprintln(&in, n)
		for i := 0; i < n; i++ {
			v := rnd.Intn(20) + 1
			h[i] = v
			if i > 0 {
				fmt.Fprint(&in, " ")
			}
			fmt.Fprint(&in, v)
		}
		in.WriteByte('\n')
		m := rnd.Intn(3) + 1
		fmt.Fprintln(&in, m)
		qs := make([][3]int, m)
		for j := 0; j < m; j++ {
			l := rnd.Intn(n)
			r := l + rnd.Intn(n-l)
			w := rnd.Intn(r-l+1) + 1
			fmt.Fprintf(&in, "%d %d %d\n", l+1, r+1, w)
			qs[j] = [3]int{l, r, w}
		}
		res := solve(h, qs)
		var out bytes.Buffer
		for _, v := range res {
			fmt.Fprintln(&out, v)
		}
		tests = append(tests, test{in: in.String(), out: out.String()})
	}
	return tests
}

func runBinary(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierE.go <binary>")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		got, err := runBinary(bin, t.in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		g := strings.TrimSpace(got)
		e := strings.TrimSpace(t.out)
		if g != e {
			fmt.Fprintf(os.Stderr, "Test %d failed\nInput:\n%sExpected:\n%sGot:\n%s\n", i+1, t.in, e, g)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
