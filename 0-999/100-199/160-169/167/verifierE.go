package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func modPow(a, e, p int64) int64 {
	res := int64(1)
	a %= p
	for e > 0 {
		if e&1 != 0 {
			res = res * a % p
		}
		a = a * a % p
		e >>= 1
	}
	return res
}

func modInv(a, p int64) int64 {
	return modPow(a, p-2, p)
}

func solve(n, m int, p int64, edges [][2]int, queries [][]int) int64 {
	adj := make([][]int, n)
	indeg := make([]int, n)
	outdeg := make([]int, n)
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		indeg[v]++
		outdeg[u]++
	}
	var src, sink []int
	for i := 0; i < n; i++ {
		if indeg[i] == 0 {
			src = append(src, i)
		}
		if outdeg[i] == 0 {
			sink = append(sink, i)
		}
	}
	k := len(src)
	if k == 0 {
		return 1 % p
	}
	dq := make([]int, 0, n)
	indeg0 := make([]int, n)
	copy(indeg0, indeg)
	for i := 0; i < n; i++ {
		if indeg0[i] == 0 {
			dq = append(dq, i)
		}
	}
	topo := make([]int, 0, n)
	for idx := 0; idx < len(dq); idx++ {
		u := dq[idx]
		topo = append(topo, u)
		for _, v := range adj[u] {
			indeg0[v]--
			if indeg0[v] == 0 {
				dq = append(dq, v)
			}
		}
	}
	M := make([][]int64, k)
	for i := range M {
		M[i] = make([]int64, k)
	}
	dp := make([]int64, n)
	for i := 0; i < k; i++ {
		for j := 0; j < n; j++ {
			dp[j] = 0
		}
		dp[src[i]] = 1
		for _, u := range topo {
			if dp[u] != 0 {
				x := dp[u]
				for _, v := range adj[u] {
					dp[v] += x
					if dp[v] >= p {
						dp[v] -= p
					}
				}
			}
		}
		for j := 0; j < k; j++ {
			M[i][j] = dp[sink[j]]
		}
	}
	det := int64(1)
	sign := int64(1)
	for i := 0; i < k; i++ {
		piv := -1
		for r := i; r < k; r++ {
			if M[r][i] != 0 {
				piv = r
				break
			}
		}
		if piv == -1 {
			det = 0
			break
		}
		if piv != i {
			M[i], M[piv] = M[piv], M[i]
			sign = p - sign
		}
		inv := modInv(M[i][i], p)
		for r := i + 1; r < k; r++ {
			if M[r][i] != 0 {
				factor := M[r][i] * inv % p
				for c := i; c < k; c++ {
					M[r][c] = (M[r][c] - factor*M[i][c]) % p
					if M[r][c] < 0 {
						M[r][c] += p
					}
				}
			}
		}
		det = det * M[i][i] % p
	}
	det = det * sign % p
	if det < 0 {
		det += p
	}
	return det
}

func generateCase(rng *rand.Rand) (string, int64) {
	n := rng.Intn(4) + 2
	edges := make([][2]int, 0, n-1)
	for i := 0; i < n-1; i++ {
		edges = append(edges, [2]int{i, i + 1})
	}
	m := len(edges)
	p := int64(1000000007)
	ans := solve(n, m, p, edges, nil)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, p))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0]+1, e[1]+1))
	}
	return sb.String(), ans
}

func run(bin, input string) (int64, error) {
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
		return 0, fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	val, err := strconv.ParseInt(strings.TrimSpace(out.String()), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid integer output: %v", err)
	}
	return val, nil
}

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d mismatch: expected %d got %d\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
