package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func expected(n, d int, positions []int) int {
	gems := map[int]int{}
	maxPos := 0
	for _, p := range positions {
		gems[p]++
		if p > maxPos {
			maxPos = p
		}
	}
	memo := make(map[[2]int]int)
	var dfs func(pos, jump int) int
	dfs = func(pos, jump int) int {
		key := [2]int{pos, jump}
		if v, ok := memo[key]; ok {
			return v
		}
		best := 0
		for delta := -1; delta <= 1; delta++ {
			nj := jump + delta
			if nj <= 0 {
				continue
			}
			np := pos + nj
			if np > maxPos {
				continue
			}
			val := gems[np] + dfs(np, nj)
			if val > best {
				best = val
			}
		}
		memo[key] = best
		return best
	}
	return gems[d] + dfs(d, d)
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for tc := 0; tc < 100; tc++ {
		n := rng.Intn(8) + 1
		d := rng.Intn(10) + 1
		positions := make([]int, n)
		for i := 0; i < n; i++ {
			positions[i] = d + rng.Intn(20)
		}
		sort.Ints(positions)
		input := fmt.Sprintf("%d %d\n", n, d)
		for _, p := range positions {
			input += fmt.Sprintf("%d\n", p)
		}
		exp := expected(n, d, positions)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", tc+1, err, input)
			os.Exit(1)
		}
		if out != fmt.Sprintf("%d", exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\ninput:\n%s", tc+1, exp, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
