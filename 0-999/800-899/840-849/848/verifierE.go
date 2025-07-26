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

func runSolution(bin, input string) (string, error) {
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
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func dist(a, b, n int) int {
	d := a - b
	if d < 0 {
		d = -d
	}
	if d > 2*n-d {
		d = 2*n - d
	}
	return d
}

func beauty(match []int, n int) int {
	N := 2 * n
	opp := make([]bool, N)
	countOpp := 0
	for i := 0; i < N; i++ {
		if match[i] == (i+n)%N {
			opp[i] = true
			countOpp++
		}
	}
	if countOpp == 0 {
		return 0
	}
	removed := make([]bool, N)
	for i := 0; i < N; i++ {
		if opp[i] {
			removed[i] = true
			removed[(i+n)%N] = true
		}
	}
	var segs []int
	cur := 0
	for i := 0; i < N; i++ {
		if removed[i] {
			if cur > 0 {
				segs = append(segs, cur)
				cur = 0
			} else if len(segs) == 0 {
				segs = append(segs, 0)
			}
			continue
		}
		cur++
	}
	if cur > 0 {
		segs = append(segs, cur)
	}
	if len(segs) == 0 {
		return 1
	}
	prod := 1
	for _, v := range segs {
		if v == 0 {
			prod *= 0
		} else {
			prod *= v
		}
	}
	return prod
}

func brute(n int) int {
	N := 2 * n
	used := make([]bool, N)
	match := make([]int, N)
	for i := 0; i < N; i++ {
		match[i] = -1
	}
	var res int
	var dfs func(int)
	dfs = func(i int) {
		for i < N && used[i] {
			i++
		}
		if i == N {
			res += beauty(match, n)
			return
		}
		for j := i + 1; j < N; j++ {
			if !used[j] && (dist(i, j, n) <= 2 || dist(i, j, n) == n) {
				i2 := (i + n) % N
				j2 := (j + n) % N
				if used[i2] || used[j2] || !(dist(i2, j2, n) <= 2 || dist(i2, j2, n) == n) {
					continue
				}
				used[i], used[j], used[i2], used[j2] = true, true, true, true
				match[i], match[j] = j, i
				match[i2], match[j2] = j2, i2
				dfs(i + 1)
				used[i], used[j], used[i2], used[j2] = false, false, false, false
				match[i], match[j], match[i2], match[j2] = -1, -1, -1, -1
			}
		}
	}
	dfs(0)
	return res
}

func solveE(n int) int {
	if n > 7 {
		return 0
	}
	return brute(n) % 998244353
}

func generateCaseE(rng *rand.Rand) string {
	n := rng.Intn(6) + 1
	return fmt.Sprintf("%d\n", n)
}

func verifyE(input, output string) error {
	var n int
	if _, err := fmt.Sscan(strings.TrimSpace(input), &n); err != nil {
		return fmt.Errorf("bad input: %v", err)
	}
	expected := solveE(n)
	var got int
	if _, err := fmt.Sscan(strings.TrimSpace(output), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if got != expected {
		return fmt.Errorf("expected %d got %d", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []string{"1\n", "2\n", "3\n"}
	for len(cases) < 100 {
		cases = append(cases, generateCaseE(rng))
	}
	for i, tc := range cases {
		out, err := runSolution(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
		if err := verifyE(tc, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%soutput:\n%s", i+1, err, tc, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
