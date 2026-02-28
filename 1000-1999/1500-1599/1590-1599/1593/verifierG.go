package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func generateCase(rng *rand.Rand) (string, [][2]int) {
	n := rng.Intn(10) + 2
	chars := []byte{'(', ')', '[', ']'}
	s := make([]byte, n)
	for i := range s {
		s[i] = chars[rng.Intn(len(chars))]
	}
	q := rng.Intn(5) + 1
	queries := make([][2]int, q)
	for i := 0; i < q; i++ {
		l := rng.Intn(n-1) + 1
		maxLen := n - l + 1
		length := 2 * (rng.Intn(maxLen/2) + 1)
		r := l + length - 1
		queries[i] = [2]int{l, r}
	}
	return string(s), queries
}

func runProg(exe, input string) (string, error) {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func minCost(sub string) int {
	const inf = int(1e9)
	type state struct {
		pos   int
		stack string
	}
	memo := make(map[state]int)

	var solve func(pos int, stack string) int
	solve = func(pos int, stack string) int {
		st := state{pos: pos, stack: stack}
		if v, ok := memo[st]; ok {
			return v
		}
		if pos == len(sub) {
			if len(stack) == 0 {
				return 0
			}
			return inf
		}

		options := [][2]int{}
		if sub[pos] == '(' || sub[pos] == ')' {
			options = append(options, [2]int{0, 0})
		} else {
			options = append(options, [2]int{1, 0})
			options = append(options, [2]int{0, 1})
		}

		best := inf
		for _, op := range options {
			typ := byte(op[0])
			cost := op[1]

			openCost := solve(pos+1, stack+string('0'+typ))
			if openCost+cost < best {
				best = openCost + cost
			}

			if len(stack) > 0 && stack[len(stack)-1] == '0'+typ {
				closeCost := solve(pos+1, stack[:len(stack)-1])
				if closeCost+cost < best {
					best = closeCost + cost
				}
			}
		}
		memo[st] = best
		return best
	}

	return solve(0, "")
}

func expectedOutput(s string, queries [][2]int) string {
	ans := make([]string, len(queries))
	for i, q := range queries {
		ans[i] = strconv.Itoa(minCost(s[q[0]-1 : q[1]]))
	}
	return strings.Join(ans, "\n")
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		s, queries := generateCase(rng)
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(s)
		sb.WriteByte('\n')
		fmt.Fprintf(&sb, "%d\n", len(queries))
		for _, q := range queries {
			fmt.Fprintf(&sb, "%d %d\n", q[0], q[1])
		}
		input := sb.String()
		exp := expectedOutput(s, queries)
		got, err := runProg(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d mismatch\nexpected:%s\n got:%s\ninput:%s", i+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
