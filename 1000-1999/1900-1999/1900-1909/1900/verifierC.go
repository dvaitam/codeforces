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

func solveOne(n int, s string, l, r []int) int {
	type node struct{ idx, cost int }
	stack := []node{{1, 0}}
	ans := int(1e9)
	for len(stack) > 0 {
		cur := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		i := cur.idx
		if l[i] == 0 && r[i] == 0 {
			if cur.cost < ans {
				ans = cur.cost
			}
		}
		if l[i] != 0 {
			c := cur.cost
			if s[i-1] != 'L' {
				c++
			}
			stack = append(stack, node{l[i], c})
		}
		if r[i] != 0 {
			c := cur.cost
			if s[i-1] != 'R' {
				c++
			}
			stack = append(stack, node{r[i], c})
		}
	}
	return ans
}

func generateTree(rng *rand.Rand, n int) ([]int, []int) {
	l := make([]int, n+1)
	r := make([]int, n+1)
	parents := make([]int, n+1)
	for i := 2; i <= n; i++ {
		for {
			p := rng.Intn(i-1) + 1
			side := rng.Intn(2)
			if side == 0 && l[p] == 0 {
				l[p] = i
				parents[i] = p
				break
			} else if side == 1 && r[p] == 0 {
				r[p] = i
				parents[i] = p
				break
			}
		}
	}
	return l, r
}

func generateCase(rng *rand.Rand) testCase {
	t := rng.Intn(3) + 1
	var in strings.Builder
	var out strings.Builder
	in.WriteString(fmt.Sprintf("%d\n", t))
	for i := 0; i < t; i++ {
		n := rng.Intn(10) + 1
		bytesStr := make([]byte, n)
		for j := 0; j < n; j++ {
			pick := rng.Intn(3)
			if pick == 0 {
				bytesStr[j] = 'U'
			} else if pick == 1 {
				bytesStr[j] = 'L'
			} else {
				bytesStr[j] = 'R'
			}
		}
		s := string(bytesStr)
		l, r := generateTree(rng, n)
		in.WriteString(fmt.Sprintf("%d\n%s\n", n, s))
		for j := 1; j <= n; j++ {
			in.WriteString(fmt.Sprintf("%d %d\n", l[j], r[j]))
		}
		out.WriteString(fmt.Sprintf("%d\n", solveOne(n, s, l, r)))
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	// small deterministic case
	in := "1\n1\nU\n0 0\n"
	out := fmt.Sprintf("%d\n", solveOne(1, "U", []int{0, 0}, []int{0, 0}))
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
