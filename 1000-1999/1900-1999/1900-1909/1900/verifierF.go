package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCase struct {
	input    string
	expected string
}

func f(arr []int) int {
	x := append([]int(nil), arr...)
	opMin := true
	for len(x) > 1 {
		var y []int
		if opMin {
			for i := range x {
				if i == 0 {
					if len(x) == 1 || x[i] < x[i+1] {
						y = append(y, x[i])
					}
				} else if i == len(x)-1 {
					if x[i] < x[i-1] {
						y = append(y, x[i])
					}
				} else if x[i] < x[i-1] && x[i] < x[i+1] {
					y = append(y, x[i])
				}
			}
		} else {
			for i := range x {
				if i == 0 {
					if len(x) == 1 || x[i] > x[i+1] {
						y = append(y, x[i])
					}
				} else if i == len(x)-1 {
					if x[i] > x[i-1] {
						y = append(y, x[i])
					}
				} else if x[i] > x[i-1] && x[i] > x[i+1] {
					y = append(y, x[i])
				}
			}
		}
		if len(y) == 0 {
			// shouldn't happen for permutations, but guard
			y = append(y, x[0])
		}
		x = y
		opMin = !opMin
	}
	return x[0]
}

func solveOne(n, q int, a []int, queries [][2]int) []int {
	res := make([]int, q)
	for i, qu := range queries {
		l, r := qu[0], qu[1]
		arr := append([]int(nil), a[l:r]...)
		res[i] = f(arr)
	}
	return res
}

func generateCase(rng *rand.Rand) testCase {
	n := rng.Intn(6) + 1
	q := rng.Intn(5) + 1
	perm := rand.Perm(n)
	for i := range perm {
		perm[i]++
	}
	queries := make([][2]int, q)
	var in strings.Builder
	var out strings.Builder
	in.WriteString(fmt.Sprintf("%d %d\n", n, q))
	for i := 0; i < n; i++ {
		if i > 0 {
			in.WriteByte(' ')
		}
		in.WriteString(fmt.Sprintf("%d", perm[i]))
	}
	in.WriteByte('\n')
	for i := 0; i < q; i++ {
		l := rng.Intn(n) + 1
		r := rng.Intn(n-l+1) + l
		queries[i] = [2]int{l - 1, r} // zero based for solve
		in.WriteString(fmt.Sprintf("%d %d\n", l, r))
	}
	ans := solveOne(n, q, perm, queries)
	for _, v := range ans {
		out.WriteString(fmt.Sprintf("%d\n", v))
	}
	return testCase{input: in.String(), expected: out.String()}
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(tc.expected)
	if got != exp {
		return fmt.Errorf("expected:\n%s\ngot:\n%s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	// trivial case
	in := "1 1\n1\n1 1\n"
	ans := solveOne(1, 1, []int{1}, [][2]int{{0, 1}})
	out := fmt.Sprintf("%d\n", ans[0])
	cases := []testCase{{input: in, expected: out}}
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
