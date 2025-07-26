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

const MOD int64 = 1000000007

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func phi(n int) int {
	result := n
	p := 2
	for p*p <= n {
		if n%p == 0 {
			for n%p == 0 {
				n /= p
			}
			result -= result / p
		}
		p++
	}
	if n > 1 {
		result -= result / n
	}
	return result
}

func modPow(a, e int64) int64 {
	res := int64(1)
	for e > 0 {
		if e&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		e >>= 1
	}
	return res
}

func distance(adj [][]int, u, v int) int {
	q := []int{u}
	dist := make([]int, len(adj))
	for i := range dist {
		dist[i] = -1
	}
	dist[u] = 0
	for len(q) > 0 {
		x := q[0]
		q = q[1:]
		if x == v {
			return dist[x]
		}
		for _, to := range adj[x] {
			if dist[to] == -1 {
				dist[to] = dist[x] + 1
				q = append(q, to)
			}
		}
	}
	return 0
}

func expectedE(n int, a []int, edges [][2]int) int64 {
	adj := make([][]int, n)
	for _, e := range edges {
		x, y := e[0]-1, e[1]-1
		adj[x] = append(adj[x], y)
		adj[y] = append(adj[y], x)
	}
	sum := int64(0)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			d := distance(adj, i, j)
			val := phi(a[i] * a[j])
			sum += int64(val * d)
		}
	}
	numerator := sum * 2 % MOD
	denom := int64(n) * int64(n-1) % MOD
	return numerator * modPow(denom, MOD-2) % MOD
}

func generateCaseE(rng *rand.Rand) (int, []int, [][2]int) {
	n := rng.Intn(5) + 2
	perm := rng.Perm(n)
	a := make([]int, n)
	for i, v := range perm {
		a[i] = v + 1
	}
	edges := make([][2]int, n-1)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		edges[i-2] = [2]int{i, p}
	}
	return n, a, edges
}

func runCase(bin string, n int, a []int, edges [][2]int) error {
	var input strings.Builder
	input.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range a {
		if i > 0 {
			input.WriteByte(' ')
		}
		input.WriteString(fmt.Sprint(v))
	}
	input.WriteByte('\n')
	for _, e := range edges {
		input.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int64
	if _, err := fmt.Fscan(bytes.NewReader(out.Bytes()), &got); err != nil {
		return fmt.Errorf("invalid output: %v", err)
	}
	expect := expectedE(n, a, edges)
	if got != expect {
		return fmt.Errorf("expected %d got %d", expect, got)
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
	for i := 0; i < 100; i++ {
		n, a, edges := generateCaseE(rng)
		if err := runCase(bin, n, a, edges); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
