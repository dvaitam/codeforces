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

const MOD int = 998244353

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func countWays(d int, lanterns []int, points []int) int {
	L := len(lanterns)
	powers := make([]int, L)
	ways := 0
	var dfs func(int)
	dfs = func(idx int) {
		if idx == L {
			for _, p := range points {
				ok := false
				for i := 0; i < L; i++ {
					if abs(lanterns[i]-p) <= powers[i] {
						ok = true
						break
					}
				}
				if !ok {
					return
				}
			}
			ways = (ways + 1) % MOD
			return
		}
		for pw := 0; pw <= d; pw++ {
			powers[idx] = pw
			dfs(idx + 1)
		}
	}
	dfs(0)
	return ways
}

func solveCase(d int, lanterns []int, points []int, queries []int) []int {
	res := make([]int, len(queries))
	for i, f := range queries {
		ls := append(lanterns, f)
		res[i] = countWays(d, ls, points)
	}
	return res
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(7))
	tests := []testCase{}
	for len(tests) < 100 {
		d := rng.Intn(3) + 4 // 4..6
		n := rng.Intn(2) + 1 // 1..2
		m := rng.Intn(2) + 1 // 1..2
		used := map[int]bool{}
		lanterns := make([]int, n)
		for i := 0; i < n; i++ {
			for {
				x := rng.Intn(d-1) + 1
				if !used[x] {
					used[x] = true
					lanterns[i] = x
					break
				}
			}
		}
		points := make([]int, m)
		for i := 0; i < m; i++ {
			for {
				x := rng.Intn(d-1) + 1
				if !used[x] {
					used[x] = true
					points[i] = x
					break
				}
			}
		}
		q := rng.Intn(2) + 1
		queries := make([]int, q)
		for i := 0; i < q; i++ {
			for {
				x := rng.Intn(d-1) + 1
				if !used[x] {
					used[x] = true
					queries[i] = x
					break
				}
			}
		}
		res := solveCase(d, lanterns, points, queries)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d\n", d, n, m))
		for i, v := range lanterns {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		for i, v := range points {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		sb.WriteString(fmt.Sprintf("%d\n", q))
		for i, v := range queries {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		expSb := strings.Builder{}
		for i, v := range res {
			if i > 0 {
				expSb.WriteByte('\n')
			}
			expSb.WriteString(fmt.Sprint(v))
		}
		tests = append(tests, testCase{input: sb.String(), expected: expSb.String()})
	}
	return tests
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	_ = rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := generateTests()
	for i, tc := range tests {
		got, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
		if got != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, tc.expected, got, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
