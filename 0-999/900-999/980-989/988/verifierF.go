package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type segment struct{ l, r int }
type umbrella struct {
	x int
	p int64
}

type testCase struct {
	a    int
	rain []segment
	umb  []umbrella
}

func solve(tc testCase) string {
	a := tc.a
	rainArr := make([]bool, a)
	for _, seg := range tc.rain {
		for x := seg.l; x < seg.r; x++ {
			rainArr[x] = true
		}
	}
	umbIndex := make([]int, a+1)
	weights := []int64{0}
	for _, u := range tc.umb {
		if u.x < 0 || u.x > a {
			continue
		}
		if umbIndex[u.x] == 0 {
			weights = append(weights, u.p)
			umbIndex[u.x] = len(weights) - 1
		} else if u.p < weights[umbIndex[u.x]] {
			weights[umbIndex[u.x]] = u.p
		}
	}
	U := len(weights) - 1
	const INF int64 = 1 << 60
	dp := make([][]int64, a+1)
	for i := 0; i <= a; i++ {
		dp[i] = make([]int64, U+1)
		for j := 0; j <= U; j++ {
			dp[i][j] = INF
		}
	}
	dp[0][0] = 0
	if umbIndex[0] != 0 {
		dp[0][umbIndex[0]] = 0
	}

	for pos := 0; pos < a; pos++ {
		for j := 0; j <= U; j++ {
			val := dp[pos][j]
			if val == INF {
				continue
			}
			candidates := []int{}
			if j != 0 {
				candidates = append(candidates, j)
			}
			if umbIndex[pos] != 0 && umbIndex[pos] != j {
				candidates = append(candidates, umbIndex[pos])
			}
			if !rainArr[pos] {
				candidates = append(candidates, 0)
			}
			if len(candidates) == 0 && rainArr[pos] {
				continue
			}
			used := make(map[int]bool)
			for _, k := range candidates {
				if used[k] {
					continue
				}
				used[k] = true
				if k == 0 && rainArr[pos] {
					continue
				}
				cost := val + weights[k]
				if cost < dp[pos+1][k] {
					dp[pos+1][k] = cost
				}
				if umbIndex[pos+1] != 0 {
					idx := umbIndex[pos+1]
					if cost < dp[pos+1][idx] {
						dp[pos+1][idx] = cost
					}
				}
				if cost < dp[pos+1][0] {
					dp[pos+1][0] = cost
				}
			}
		}
	}
	ans := INF
	for j := 0; j <= U; j++ {
		if dp[a][j] < ans {
			ans = dp[a][j]
		}
	}
	if ans == INF {
		return "-1\n"
	}
	return fmt.Sprintf("%d\n", ans)
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", tc.a, len(tc.rain), len(tc.umb))
	for _, seg := range tc.rain {
		fmt.Fprintf(&sb, "%d %d\n", seg.l, seg.r)
	}
	for _, u := range tc.umb {
		fmt.Fprintf(&sb, "%d %d\n", u.x, u.p)
	}
	return sb.String()
}

func runCase(bin string, tc testCase) error {
	input := buildInput(tc)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	expected := strings.TrimSpace(solve(tc))
	got := strings.TrimSpace(out.String())
	if expected != got {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	var cases []testCase
	cases = append(cases, testCase{a: 1})
	rng := rand.New(rand.NewSource(5))
	for i := 0; i < 100; i++ {
		a := rng.Intn(10) + 1
		n := rng.Intn(a/2 + 1)
		m := rng.Intn(a + 1)
		var rain []segment
		pos := 0
		for j := 0; j < n; j++ {
			l := rng.Intn(a - pos)
			r := l + rng.Intn(a-l) + 1
			if r > a {
				r = a
			}
			rain = append(rain, segment{l: l + pos, r: r + pos})
			pos += r - l
			if pos >= a {
				break
			}
		}
		var umb []umbrella
		for j := 0; j < m; j++ {
			x := rng.Intn(a + 1)
			p := int64(rng.Intn(10) + 1)
			umb = append(umb, umbrella{x: x, p: p})
		}
		cases = append(cases, testCase{a: a, rain: rain, umb: umb})
	}

	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\ninput:\n%s", i+1, err, buildInput(tc))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
