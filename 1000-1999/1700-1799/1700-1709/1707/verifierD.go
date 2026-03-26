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

// ---------- embedded reference solver for 1707D ----------

func referenceSolve(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)

	var n int
	var p int64
	fmt.Fscan(reader, &n, &p)

	adj := make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}

	children := make([][]int, n+1)
	var buildTree func(u, parent int)
	buildTree = func(u, parent int) {
		for _, v := range adj[u] {
			if v != parent {
				children[u] = append(children[u], v)
				buildTree(v, u)
			}
		}
	}
	buildTree(1, 0)

	dp := make([][]int64, n+1)
	F := make([][]int64, n+1)
	for i := 1; i <= n; i++ {
		dp[i] = make([]int64, n+1)
		F[i] = make([]int64, n+1)
	}

	var postOrder func(u int)
	postOrder = func(u int) {
		for _, v := range children[u] {
			postOrder(v)
		}

		if u == 1 {
			return
		}

		d := len(children[u])
		if d == 0 {
			for i := 1; i <= n; i++ {
				dp[u][i] = 1
				F[u][i] = (F[u][i-1] + dp[u][i]) % p
			}
			return
		}

		S := make([]int64, d)
		pref := make([]int64, d)
		suff := make([]int64, d)

		for i := 1; i <= n; i++ {
			for j := 0; j < d; j++ {
				v := children[u][j]
				val := F[v][i]
				if j == 0 {
					pref[j] = val
				} else {
					pref[j] = (pref[j-1] * val) % p
				}
			}
			for j := d - 1; j >= 0; j-- {
				v := children[u][j]
				val := F[v][i]
				if j == d-1 {
					suff[j] = val
				} else {
					suff[j] = (suff[j+1] * val) % p
				}
			}

			ans := pref[d-1]
			for j := 0; j < d; j++ {
				v := children[u][j]
				term := (dp[v][i] * S[j]) % p
				ans = (ans + term) % p
			}
			dp[u][i] = ans
			F[u][i] = (F[u][i-1] + ans) % p

			for j := 0; j < d; j++ {
				var left, right int64 = 1, 1
				if j > 0 {
					left = pref[j-1]
				}
				if j < d-1 {
					right = suff[j+1]
				}
				E_v_i := (left * right) % p
				S[j] = (S[j] + E_v_i) % p
			}
		}
	}

	postOrder(1)

	g := make([]int64, n+1)
	for x := 1; x <= n; x++ {
		var ways int64 = 1
		for _, v := range children[1] {
			ways = (ways * F[v][x]) % p
		}
		g[x] = ways
	}

	C := make([][]int64, n+1)
	for i := 0; i <= n; i++ {
		C[i] = make([]int64, i+1)
		C[i][0] = 1
		for j := 1; j < i; j++ {
			C[i][j] = (C[i-1][j-1] + C[i-1][j]) % p
		}
		C[i][i] = 1
	}

	for k := 1; k <= n-1; k++ {
		var exact int64 = 0
		for j := 1; j <= k; j++ {
			term := (C[k][j] * g[j]) % p
			if (k-j)%2 == 1 {
				exact = (exact - term + p) % p
			} else {
				exact = (exact + term) % p
			}
		}
		if k == n-1 {
			fmt.Fprint(writer, exact, "\n")
		} else {
			fmt.Fprint(writer, exact, " ")
		}
	}

	writer.Flush()
	return strings.TrimSpace(buf.String())
}

// ---------- test generator ----------

func genCase(r *rand.Rand) string {
	n := r.Intn(5) + 2
	p := int64(1000000007)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, p))
	for i := 2; i <= n; i++ {
		u := i
		v := r.Intn(i-1) + 1
		sb.WriteString(fmt.Sprintf("%d %d\n", u, v))
	}
	return sb.String()
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		input := genCase(rng)
		expect := referenceSolve(input)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, err, input)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All 100 tests passed")
}
