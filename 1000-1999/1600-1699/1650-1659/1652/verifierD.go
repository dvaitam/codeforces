package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

const mod int64 = 998244353

var spf []int

func sieve(n int) {
	spf = make([]int, n+1)
	for i := 2; i <= n; i++ {
		if spf[i] == 0 {
			spf[i] = i
			if int64(i)*int64(i) <= int64(n) {
				for j := i * i; j <= n; j += i {
					if spf[j] == 0 {
						spf[j] = i
					}
				}
			}
		}
	}
}

func modPow(a, e int64) int64 {
	res := int64(1)
	a %= mod
	for e > 0 {
		if e&1 != 0 {
			res = res * a % mod
		}
		a = a * a % mod
		e >>= 1
	}
	return res
}

func modInv(a int64) int64 {
	return modPow(a, mod-2)
}

type Edge struct {
	to int
	x  int
	y  int
}

func factorUpdate(x int, sign int, cnt, minCnt map[int]int, updateMin bool) {
	for x > 1 {
		p := spf[x]
		c := 0
		for x%p == 0 {
			x /= p
			c++
		}
		cnt[p] += sign * c
		if updateMin && cnt[p] < minCnt[p] {
			minCnt[p] = cnt[p]
		}
	}
}

func dfs(u, p int, adj [][]Edge, val []int64, inv []int64, cnt, minCnt map[int]int) {
	for _, e := range adj[u] {
		if e.to == p {
			continue
		}
		factorUpdate(e.y, 1, cnt, minCnt, true)
		factorUpdate(e.x, -1, cnt, minCnt, true)
		val[e.to] = val[u] * int64(e.y) % mod * inv[e.x] % mod
		dfs(e.to, u, adj, val, inv, cnt, minCnt)
		factorUpdate(e.y, -1, cnt, minCnt, false)
		factorUpdate(e.x, 1, cnt, minCnt, false)
	}
}

func runReferenceLogic() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}

	const MAX = 200000
	sieve(MAX)
	inv := make([]int64, MAX+1)
	for i := 1; i <= MAX; i++ {
		inv[i] = modInv(int64(i))
	}

	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		adj := make([][]Edge, n+1)
		for i := 0; i < n-1; i++ {
			var u, v, x, y int
			fmt.Fscan(in, &u, &v, &x, &y)
			adj[u] = append(adj[u], Edge{v, x, y})
			adj[v] = append(adj[v], Edge{u, y, x})
		}
		val := make([]int64, n+1)
		val[1] = 1
		cnt := make(map[int]int)
		minCnt := make(map[int]int)
		dfs(1, 0, adj, val, inv, cnt, minCnt)

		base := int64(1)
		for p, e := range minCnt {
			if e < 0 {
				base = base * modPow(int64(p), int64(-e)) % mod
			}
		}
		ans := int64(0)
		for i := 1; i <= n; i++ {
			ans = (ans + val[i]*base) % mod
		}
		fmt.Fprintln(out, ans)
	}
}

func runBinary(bin, input string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin, args...)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("timeout")
	}
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateTests() []string {
	r := rand.New(rand.NewSource(4))
	tests := make([]string, 0, 100)
	for t := 0; t < 100; t++ {
		n := r.Intn(5) + 2
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 2; i <= n; i++ {
			u := r.Intn(i-1) + 1
			x := r.Intn(5) + 1
			y := r.Intn(5) + 1
			sb.WriteString(fmt.Sprintf("%d %d %d %d\n", u, i, x, y))
		}
		tests = append(tests, sb.String())
	}
	return tests
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--ref" {
		runReferenceLogic()
		return
	}
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		return
	}
	
	// Build ourselves as the reference binary
	refBin := "./verifier_ref_build"
	if err := exec.Command("go", "build", "-o", refBin, "verifierD.go").Run(); err != nil {
		fmt.Printf("failed to build reference from verifier source: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)
	
cand := os.Args[1]
tests := generateTests()
	
	for i, tc := range tests {
		exp, eerr := runBinary(refBin, tc, "--ref")
		got, gerr := runBinary(cand, tc)
		if eerr != nil {
			fmt.Printf("official solution failed on test %d: %v\n", i+1, eerr)
			os.Exit(1)
		}
		if gerr != nil {
			fmt.Printf("candidate failed on test %d: %v\n", i+1, gerr)
			os.Exit(1)
		}
		if exp != got {
			fmt.Printf("wrong answer on test %d\ninput:\n%sexpected: %s\ngot: %s\n", i+1, tc, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
