package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

const MOD int64 = 1000000007

type testCase struct {
	name  string
	input string
}

func solveRef(input string) (string, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	var T int
	if _, err := fmt.Fscan(reader, &T); err != nil {
		return "", err
	}
	results := make([]int64, T)
	for tc := 0; tc < T; tc++ {
		var n, k int
		fmt.Fscan(reader, &n, &k)
		g := make([][]int, n)
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(reader, &u, &v)
			u--
			v--
			g[u] = append(g[u], v)
			g[v] = append(g[v], u)
		}
		if k == 2 {
			results[tc] = int64(n) * int64(n-1) / 2 % MOD
			continue
		}
		ans := int64(0)
		for r := 0; r < n; r++ {
			m := len(g[r])
			if m < k {
				continue
			}
			dist := make([]int, n)
			branch := make([]int, n)
			for i := 0; i < n; i++ {
				dist[i] = -1
				branch[i] = -1
			}
			queue := []int{}
			dist[r] = 0
			for idx, v := range g[r] {
				dist[v] = 1
				branch[v] = idx
				queue = append(queue, v)
			}
			for front := 0; front < len(queue); front++ {
				u := queue[front]
				for _, w := range g[u] {
					if dist[w] == -1 {
						dist[w] = dist[u] + 1
						branch[w] = branch[u]
						queue = append(queue, w)
					}
				}
			}
			maxD := 0
			cnt := make([][]int, m)
			for i := range cnt {
				cnt[i] = make([]int, n+1)
			}
			for v := 0; v < n; v++ {
				if v == r || dist[v] == -1 {
					continue
				}
				b := branch[v]
				d := dist[v]
				cnt[b][d]++
				if d > maxD {
					maxD = d
				}
			}
			for d := 1; d <= maxD; d++ {
				dp := make([]int64, k+1)
				dp[0] = 1
				for i := 0; i < m; i++ {
					c := cnt[i][d]
					if c == 0 {
						continue
					}
					for j := k; j >= 1; j-- {
						dp[j] = (dp[j] + dp[j-1]*int64(c)) % MOD
					}
				}
				ans = (ans + dp[k]) % MOD
			}
		}
		results[tc] = ans % MOD
	}
	var sb strings.Builder
	for i, val := range results {
		if i > 0 {
			sb.WriteByte('\n')
		}
		fmt.Fprintf(&sb, "%d", val)
	}
	return sb.String(), nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstderr: %s", err, stderr.String())
	}
	if stderr.Len() > 0 {
		return "", fmt.Errorf("unexpected stderr output: %s", stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func parseOutput(out string, expectedLines int) ([]string, error) {
	lines := strings.Fields(out)
	if len(lines) != expectedLines {
		return nil, fmt.Errorf("expected %d outputs, got %d", expectedLines, len(lines))
	}
	return lines, nil
}

func makeCase(name string, trees []graph) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(trees))
	for _, g := range trees {
		fmt.Fprintf(&sb, "%d %d\n", g.n, g.k)
		for _, e := range g.edges {
			fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
		}
	}
	return testCase{name: name, input: sb.String()}
}

type graph struct {
	n     int
	k     int
	edges [][2]int
}

func randomTree(rng *rand.Rand, n int) [][2]int {
	edges := make([][2]int, 0, n-1)
	for i := 2; i <= n; i++ {
		j := rng.Intn(i-1) + 1
		edges = append(edges, [2]int{i, j})
	}
	return edges
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	for tcase := 0; tcase < 30; tcase++ {
		tcCount := rng.Intn(3) + 1
		trees := make([]graph, tcCount)
		for i := 0; i < tcCount; i++ {
			n := rng.Intn(7) + 2
			k := rng.Intn(n-1) + 2
			trees[i] = graph{
				n:     n,
				k:     k,
				edges: randomTree(rng, n),
			}
		}
		tests = append(tests, makeCase(fmt.Sprintf("rand_%d", tcase+1), trees))
	}
	return tests
}

func handcraftedTests() []testCase {
	return []testCase{
		makeCase("simple_line", []graph{
			{
				n: 4, k: 2,
				edges: [][2]int{{1, 2}, {2, 3}, {3, 4}},
			},
		}),
		makeCase("star", []graph{
			{
				n: 5, k: 3,
				edges: [][2]int{{1, 2}, {1, 3}, {1, 4}, {1, 5}},
			},
		}),
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := append(handcraftedTests(), randomTests()...)
	for idx, tc := range tests {
		expect, err := solveRef(tc.input)
		if err != nil {
			fmt.Printf("failed to compute reference for %s: %v\n", tc.name, err)
			os.Exit(1)
		}
		out, err := runCandidate(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d (%s) runtime error: %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		expectLines := strings.Fields(expect)
		gotLines, err := parseOutput(out, len(expectLines))
		if err != nil {
			fmt.Printf("test %d (%s) invalid output: %v\ninput:\n%s\noutput:\n%s\n", idx+1, tc.name, err, tc.input, out)
			os.Exit(1)
		}
		for i := range expectLines {
			if expectLines[i] != gotLines[i] {
				fmt.Printf("test %d (%s) mismatch at case %d\ninput:\n%s\nexpect:%s\nactual:%s\n", idx+1, tc.name, i+1, tc.input, expect, out)
				os.Exit(1)
			}
		}
	}
	fmt.Println("all tests passed")
}
