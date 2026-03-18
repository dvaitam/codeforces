package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type edge struct{ u, v int }

const embeddedRefSource = `package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, k int
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}

	adj := make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}

	MOD := 1000000007

	dp := make([][]int, n+1)
	for i := 1; i <= n; i++ {
		dp[i] = make([]int, 2*k+1)
	}

	var dfs func(u, p int)
	dfs = func(u, p int) {
		dp[u][0] = 1
		if k+1 <= 2*k {
			dp[u][k+1] = 1
		}

		for _, v := range adj[u] {
			if v == p {
				continue
			}
			dfs(v, u)

			nextDp := make([]int, 2*k+1)
			for i := 0; i <= 2*k; i++ {
				if dp[u][i] == 0 {
					continue
				}
				for j := 0; j <= 2*k; j++ {
					if dp[v][j] == 0 {
						continue
					}

					cU := i
					if i > k {
						cU = 1000000000
					}
					dU := -1000000000
					if i > k {
						dU = i - k - 1
					}

					cV := j + 1
					if j > k {
						cV = 1000000000
					}
					dV := -1000000000
					if j > k {
						dV = j - k
					}

					newC := cU
					if cV < newC {
						newC = cV
					}

					newD := -1000000000
					if dU != -1000000000 && dU+cV > k {
						if dU > newD {
							newD = dU
						}
					}
					if dV != -1000000000 && dV+cU > k {
						if dV > newD {
							newD = dV
						}
					}

					var newState int
					if newD != -1000000000 {
						newState = newD + k + 1
					} else {
						newState = newC
					}

					if newState <= 2*k {
						nextDp[newState] = (nextDp[newState] + (dp[u][i]*dp[v][j])%MOD) % MOD
					}
				}
			}
			for i := 0; i <= 2*k; i++ {
				dp[u][i] = nextDp[i]
			}
		}
	}

	dfs(1, 0)

	ans := 0
	for i := 0; i <= k; i++ {
		ans = (ans + dp[1][i]) % MOD
	}
	fmt.Println(ans)
}
`

func buildRef() (string, error) {
	ref := "./refE.bin"
	tmpGo := "refE_src.go"
	if err := os.WriteFile(tmpGo, []byte(embeddedRefSource), 0644); err != nil {
		return "", fmt.Errorf("write embedded source: %v", err)
	}
	defer os.Remove(tmpGo)
	cmd := exec.Command("go", "build", "-o", ref, tmpGo)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference: %v: %s", err, string(out))
	}
	return ref, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func genTestE() (int, int, []edge) {
	n := rand.Intn(10) + 1
	k := rand.Intn(min(20, n-1) + 1)
	edges := make([]edge, 0, n-1)
	for i := 1; i < n; i++ {
		p := rand.Intn(i)
		edges = append(edges, edge{p, i})
	}
	return n, k, edges
}

func runBinary(path string, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: verifierE /path/to/binary")
		os.Exit(1)
	}
	rand.Seed(time.Now().UnixNano())
	path := os.Args[1]

	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	for i := 0; i < 100; i++ {
		n, k, edges := genTestE()
		var b strings.Builder
		fmt.Fprintf(&b, "%d %d\n", n, k)
		for _, e := range edges {
			fmt.Fprintf(&b, "%d %d\n", e.u+1, e.v+1)
		}
		input := b.String()
		expStr, err := runBinary(ref, input)
		if err != nil {
			fmt.Printf("test %d: reference error: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		gotStr, err := runBinary(path, input)
		if err != nil {
			fmt.Printf("test %d: execution error: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(gotStr) != strings.TrimSpace(expStr) {
			fmt.Printf("test %d failed\ninput:\n%sexpected:%s\ngot:%s\n", i+1, input, expStr, gotStr)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
