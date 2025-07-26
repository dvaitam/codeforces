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
	n int
	s int
	k int
	r []int
	c string
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func expected(tc testCase) int {
	const INF = int(1e9)
	n, s, k := tc.n, tc.s-1, tc.k
	dp := make([][]int, n)
	for i := range dp {
		dp[i] = make([]int, k)
		for j := range dp[i] {
			dp[i][j] = INF
		}
	}
	best := INF
	for i := 0; i < n; i++ {
		cost := abs(i - s)
		if tc.r[i] >= k {
			best = min(best, cost)
		} else {
			if tc.r[i] < k {
				dp[i][tc.r[i]] = cost
			}
		}
	}
	for rating := 1; rating <= 50; rating++ {
		for i := 0; i < n; i++ {
			if tc.r[i] != rating {
				continue
			}
			for sum := 0; sum < k; sum++ {
				curr := dp[i][sum]
				if curr >= INF {
					continue
				}
				for j := 0; j < n; j++ {
					if tc.r[j] > tc.r[i] && tc.c[i] != tc.c[j] {
						s2 := sum + tc.r[j]
						cost2 := curr + abs(j-i)
						if s2 >= k {
							best = min(best, cost2)
						} else if cost2 < dp[j][s2] {
							dp[j][s2] = cost2
						}
					}
				}
			}
		}
	}
	if best >= INF {
		return -1
	}
	return best
}

func generateCase(rng *rand.Rand) testCase {
	n := rng.Intn(7) + 1
	s := rng.Intn(n) + 1
	r := make([]int, n)
	total := 0
	for i := range r {
		r[i] = rng.Intn(50) + 1
		total += r[i]
	}
	k := rng.Intn(total+50) + 1
	if k > 2000 {
		k = 2000
	}
	var sb strings.Builder
	colors := []byte{'R', 'G', 'B'}
	for i := 0; i < n; i++ {
		sb.WriteByte(colors[rng.Intn(3)])
	}
	return testCase{n: n, s: s, k: k, r: r, c: sb.String()}
}

func runCase(bin string, tc testCase) error {
	var b strings.Builder
	fmt.Fprintf(&b, "%d %d %d\n", tc.n, tc.s, tc.k)
	for i, v := range tc.r {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(&b, "%d", v)
	}
	b.WriteByte('\n')
	b.WriteString(tc.c)
	b.WriteByte('\n')

	got, err := run(bin, b.String())
	if err != nil {
		return err
	}
	exp := fmt.Sprintf("%d", expected(tc))
	if strings.TrimSpace(got) != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	cases := make([]testCase, 0, 100)
	// deterministic simple case
	cases = append(cases, testCase{n: 1, s: 1, k: 1, r: []int{1}, c: "R"})
	for len(cases) < 100 {
		cases = append(cases, generateCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
