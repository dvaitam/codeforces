package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const MOD int64 = 1000000007

type testCase struct {
	input  string
	expect string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()
	for i, tc := range tests {
		got, err := run(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\ninput:\n%s\n", i+1, err, tc.input)
			os.Exit(1)
		}
		got = strings.TrimSpace(got)
		if err := checkOutput(got, tc.expect); err != nil {
			fmt.Printf("test %d failed: %v\ninput:\n%s\nexpected:\n%s\nactual:\n%s\n", i+1, err, tc.input, tc.expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}

func checkOutput(out, expect string) error {
	if out == "" {
		return fmt.Errorf("empty output")
	}
	gotVal, err := strconv.ParseInt(out, 10, 64)
	if err != nil {
		return fmt.Errorf("output is not a valid integer: %v", err)
	}
	expVal, _ := strconv.ParseInt(expect, 10, 64)
	if (gotVal%MOD+MOD)%MOD != (expVal%MOD+MOD)%MOD {
		return fmt.Errorf("expected %s but got %s", expect, out)
	}
	return nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func genTests() []testCase {
	rand.Seed(42)
	var tests []testCase
	tests = append(tests, makeTest(1, 1, [][2]int{}, 1, 0))
	tests = append(tests, makeTest(2, 3, [][2]int{{1, 2}}, 2, 1))
	for t := 0; t < 200; t++ {
		n := rand.Intn(8) + 1
		edges := randomTree(n)
		m := rand.Intn(10) + 1
		k := rand.Intn(m) + 1
		x := rand.Intn(4) + 1
		tests = append(tests, makeTest(n, m, edges, k, x))
	}
	return tests
}

func randomTree(n int) [][2]int {
	if n <= 1 {
		return nil
	}
	edges := make([][2]int, 0, n-1)
	for v := 2; v <= n; v++ {
		u := rand.Intn(v-1) + 1
		edges = append(edges, [2]int{u, v})
	}
	return edges
}

func makeTest(n, m int, edges [][2]int, k, x int) testCase {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	sb.WriteString(fmt.Sprintf("%d %d\n", k, x))
	input := sb.String()
	expect := solveRef(input)
	return testCase{
		input:  input,
		expect: expect,
	}
}

func solveRef(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	var n int
	var m int64
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return ""
	}
	adj := make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	var k int64
	var x int
	fmt.Fscan(reader, &k, &x)

	dp := make([][][]int64, n+1)
	for i := range dp {
		dp[i] = make([][]int64, 3)
		for j := 0; j < 3; j++ {
			dp[i][j] = make([]int64, x+1)
		}
	}

	var dfs func(int, int)
	dfs = func(u, p int) {
		dp[u][0][0] = (k - 1) % MOD
		if dp[u][0][0] < 0 {
			dp[u][0][0] += MOD
		}
		if x >= 1 {
			dp[u][1][1] = 1
		}
		dp[u][2][0] = (m - k) % MOD
		if dp[u][2][0] < 0 {
			dp[u][2][0] += MOD
		}

		for _, v := range adj[u] {
			if v == p {
				continue
			}
			dfs(v, u)
			new0 := make([]int64, x+1)
			new1 := make([]int64, x+1)
			new2 := make([]int64, x+1)
			for i := 0; i <= x; i++ {
				if dp[u][0][i] != 0 {
					for j := 0; i+j <= x; j++ {
						if dp[v][0][j] != 0 {
							new0[i+j] = (new0[i+j] + dp[u][0][i]*dp[v][0][j]) % MOD
						}
						if dp[v][1][j] != 0 {
							new0[i+j] = (new0[i+j] + dp[u][0][i]*dp[v][1][j]) % MOD
						}
						if dp[v][2][j] != 0 {
							new0[i+j] = (new0[i+j] + dp[u][0][i]*dp[v][2][j]) % MOD
						}
					}
				}
				if dp[u][1][i] != 0 {
					for j := 0; i+j <= x; j++ {
						if dp[v][0][j] != 0 {
							new1[i+j] = (new1[i+j] + dp[u][1][i]*dp[v][0][j]) % MOD
						}
					}
				}
				if dp[u][2][i] != 0 {
					for j := 0; i+j <= x; j++ {
						if dp[v][0][j] != 0 {
							new2[i+j] = (new2[i+j] + dp[u][2][i]*dp[v][0][j]) % MOD
						}
						if dp[v][2][j] != 0 {
							new2[i+j] = (new2[i+j] + dp[u][2][i]*dp[v][2][j]) % MOD
						}
					}
				}
			}
			dp[u][0] = new0
			dp[u][1] = new1
			dp[u][2] = new2
		}
	}

	if n > 0 {
		dfs(1, 0)
	}

	var ans int64
	for s := 0; s < 3; s++ {
		for c := 0; c <= x; c++ {
			ans = (ans + dp[1][s][c]) % MOD
		}
	}
	return fmt.Sprintf("%d", (ans%MOD+MOD)%MOD)
}
