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

const MOD = 1000000007

func generateSeq(pairs int, rng *rand.Rand) string {
	seq := make([]byte, 0, 2*pairs)
	open := 0
	for len(seq) < 2*pairs {
		remaining := 2*pairs - len(seq)
		if open == 0 {
			seq = append(seq, '(')
			open++
		} else if open == remaining {
			seq = append(seq, ')')
			open--
		} else if rng.Intn(2) == 0 {
			seq = append(seq, '(')
			open++
		} else {
			seq = append(seq, ')')
			open--
		}
	}
	return string(seq)
}

func matchBrackets(s string) []int {
	n := len(s)
	match := make([]int, n)
	stack := []int{}
	for i, c := range s {
		if c == '(' {
			stack = append(stack, i)
		} else {
			l := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			match[i] = l
			match[l] = i
		}
	}
	return match
}

func countColorings(s string) int {
	n := len(s)
	match := matchBrackets(s)
	colors := make([]int, n)
	var dfs func(int, int) int
	dfs = func(idx int, prev int) int {
		if idx == n {
			return 1
		}
		if colors[idx] != 0 {
			col := colors[idx]
			if prev > 0 && col == prev {
				return 0
			}
			return dfs(idx+1, col)
		}
		if s[idx] == ')' {
			return dfs(idx+1, prev)
		}
		j := match[idx]
		total := 0
		for _, col := range []int{1, 2} {
			if prev > 0 && col == prev {
				continue
			}
			colors[idx] = col
			total = (total + dfs(idx+1, col)) % MOD
			colors[idx] = 0
		}
		for _, col := range []int{1, 2} {
			colors[j] = col
			total = (total + dfs(idx+1, prev)) % MOD
			colors[j] = 0
		}
		return total % MOD
	}
	return dfs(0, 0) % MOD
}

func runCase(bin string, s string, expect int) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(s + "\n")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if got%MOD != expect%MOD {
		return fmt.Errorf("expected %d got %d", expect%MOD, got%MOD)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 100; i++ {
		pairs := rng.Intn(4) + 1 // up to 5 pairs => length up to 10
		s := generateSeq(pairs, rng)
		expect := countColorings(s)
		if err := runCase(bin, s, expect); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput: %s\n", i+1, err, s)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
