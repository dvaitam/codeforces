package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solve(input string) string {
	var n, k int
	fmt.Sscan(input, &n, &k)

	if k == 1 || k == 2 {
		dp := make([]int, n+1)
		dp[0] = 1
		for i := 1; ; i++ {
			w := i
			if k == 2 {
				w = i * (i + 1) / 2
			}
			if w > n {
				break
			}
			for j := w; j <= n; j++ {
				dp[j] = (dp[j] + dp[j-w]) % 998244353
			}
		}
		ans := 0
		for j := 1; j <= n; j++ {
			ans = (ans + dp[j]) % 998244353
		}
		return fmt.Sprintf("%d", ans)
	}

	C := make([][]int, 4050)
	for i := 0; i < 4050; i++ {
		C[i] = make([]int, 2050)
		C[i][0] = 1
		for j := 1; j <= i && j < 2050; j++ {
			C[i][j] = C[i-1][j-1] + C[i-1][j]
			if C[i][j] > n {
				C[i][j] = n + 1
			}
		}
	}

	ans := 0
	var dfs func(last_g, i, u_sum, current_cost int)
	dfs = func(last_g, i, u_sum, current_cost int) {
		new_u_sum := u_sum
		new_cost := current_cost
		for g := 1; g <= last_g; g++ {
			new_u_sum += i
			if new_u_sum+k-3 >= 4050 || k-2 >= 2050 {
				break
			}
			new_cost += C[new_u_sum+k-3][k-2]
			if new_cost > n {
				break
			}
			ans = (ans + 1) % 998244353
			dfs(g, i+1, new_u_sum, new_cost)
		}
	}

	dfs(n, 1, 0, 0)
	return fmt.Sprintf("%d", ans)
}

func runProg(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func genTest(rng *rand.Rand) string {
	n := rng.Intn(5) + 1
	k := rng.Intn(5) + 1
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%d %d\n", n, k)
	return buf.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/bin")
		os.Exit(1)
	}
	target := os.Args[1]

	rng := rand.New(rand.NewSource(42))
	for i := 0; i < 100; i++ {
		test := genTest(rng)
		expected := strings.TrimSpace(solve(test))
		got, err := runProg(target, test)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d execution error: %v\n", i+1, err)
			os.Exit(1)
		}
		if expected != got {
			fmt.Printf("test %d failed\ninput:\n%s\nexpected:%s\ngot:%s\n", i+1, test, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
