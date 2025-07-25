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

func simulate(n int, s string, ops string, p int) string {
	pair := make([]int, n+1)
	stack := make([]int, 0, n/2)
	for i := 1; i <= n; i++ {
		if s[i-1] == '(' {
			stack = append(stack, i)
		} else {
			j := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			pair[i] = j
			pair[j] = i
		}
	}
	next := make([]int, n+2)
	prev := make([]int, n+2)
	for i := 1; i <= n; i++ {
		next[i] = i + 1
		prev[i] = i - 1
	}
	next[n] = 0
	prev[1] = 0
	cur := p
	for i := 0; i < len(ops); i++ {
		switch ops[i] {
		case 'L':
			cur = prev[cur]
		case 'R':
			cur = next[cur]
		case 'D':
			l := cur
			r := pair[cur]
			if l > r {
				l, r = r, l
			}
			lp := prev[l]
			rn := next[r]
			if lp != 0 {
				next[lp] = rn
			}
			if rn != 0 {
				prev[rn] = lp
			}
			if rn != 0 {
				cur = rn
			} else {
				cur = lp
			}
		}
	}
	for prev[cur] != 0 {
		cur = prev[cur]
	}
	var out strings.Builder
	for i := cur; i != 0; i = next[i] {
		out.WriteByte(s[i-1])
	}
	return out.String()
}

func randomBalanced(n int, rng *rand.Rand) string {
	open := n / 2
	close := n / 2
	var sb strings.Builder
	stack := 0
	for open > 0 || close > 0 {
		if open == 0 {
			sb.WriteByte(')')
			close--
			stack--
		} else if stack == 0 {
			sb.WriteByte('(')
			open--
			stack++
		} else {
			if rng.Intn(open+close) < open {
				sb.WriteByte('(')
				open--
				stack++
			} else {
				sb.WriteByte(')')
				close--
				stack--
			}
		}
	}
	return sb.String()
}

func generateCase(rng *rand.Rand) (string, string) {
	n := (rng.Intn(10) + 1) * 2
	s := randomBalanced(n, rng)
	m := rng.Intn(10) + 1
	p := rng.Intn(n) + 1
	pair := make([]int, n+1)
	stack := make([]int, 0, n/2)
	for i := 1; i <= n; i++ {
		if s[i-1] == '(' {
			stack = append(stack, i)
		} else {
			j := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			pair[i] = j
			pair[j] = i
		}
	}
	next := make([]int, n+2)
	prev := make([]int, n+2)
	for i := 1; i <= n; i++ {
		next[i] = i + 1
		prev[i] = i - 1
	}
	next[n] = 0
	prev[1] = 0
	cur := p
	var ops []byte
	length := n
	for len(ops) < m {
		choices := []byte{'L', 'R', 'D'}
		op := choices[rng.Intn(len(choices))]
		valid := true
		switch op {
		case 'L':
			if prev[cur] == 0 {
				valid = false
			}
		case 'R':
			if next[cur] == 0 {
				valid = false
			}
		case 'D':
			l := cur
			r := pair[cur]
			if l > r {
				l, r = r, l
			}
			if length-(r-l+1) == 0 {
				valid = false
			}
		}
		if !valid {
			continue
		}
		ops = append(ops, op)
		switch op {
		case 'L':
			cur = prev[cur]
		case 'R':
			cur = next[cur]
		case 'D':
			l := cur
			r := pair[cur]
			if l > r {
				l, r = r, l
			}
			lp := prev[l]
			rn := next[r]
			if lp != 0 {
				next[lp] = rn
			}
			if rn != 0 {
				prev[rn] = lp
			}
			length -= r - l + 1
			if rn != 0 {
				cur = rn
			} else {
				cur = lp
			}
		}
	}
	input := fmt.Sprintf("%d %d %d\n%s\n%s\n", n, m, p, s, string(ops))
	expected := simulate(n, s, string(ops), p)
	return input, expected
}

func runProg(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("%v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, exp := generateCase(rng)
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\noutput:\n%s", i+1, err, got)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed:\nexpected: %s\n got: %s\n", i+1, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
