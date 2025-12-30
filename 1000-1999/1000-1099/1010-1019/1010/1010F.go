package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

const (
	Mod = 998244353
	G   = 3
)

type Poly []int

func qpow(a, b int) int {
	ans := 1
	a %= Mod
	for b > 0 {
		if b&1 == 1 {
			ans = int((int64(ans) * int64(a)) % Mod)
		}
		a = int((int64(a) * int64(a)) % Mod)
		b >>= 1
	}
	return ans
}

func inv(n int) int {
	return qpow(n, Mod-2)
}

var rev []int

func initRev(length int) {
	if length == len(rev) {
		return
	}
	rev = make([]int, length)
	for i := 0; i < length; i++ {
		rev[i] = (rev[i>>1] >> 1) | ((i & 1) * (length >> 1))
	}
}

func ntt(a Poly, length int, type_ int) {
	initRev(length)
	for i := 0; i < length; i++ {
		if i < rev[i] {
			a[i], a[rev[i]] = a[rev[i]], a[i]
		}
	}
	for mid := 1; mid < length; mid <<= 1 {
		var wn int
		if type_ == 1 {
			wn = qpow(G, (Mod-1)/(mid<<1))
		} else {
			wn = qpow(inv(G), (Mod-1)/(mid<<1))
		}
		for j := 0; j < length; j += (mid << 1) {
			w := 1
			for k := 0; k < mid; k++ {
				x, y := a[j+k], int((int64(w)*int64(a[j+k+mid]))%Mod)
				a[j+k] = (x + y) % Mod
				a[j+k+mid] = (x - y + Mod) % Mod
				w = int((int64(w) * int64(wn)) % Mod)
			}
		}
	}
	if type_ == -1 {
		iv := inv(length)
		for i := 0; i < length; i++ {
			a[i] = int((int64(a[i]) * int64(iv)) % Mod)
		}
	}
}

func polyMul(a, b Poly) Poly {
	n, m := len(a), len(b)
	l := 1
	for l < n+m-1 {
		l <<= 1
	}
	A, B := make(Poly, l), make(Poly, l)
	copy(A, a)
	copy(B, b)
	ntt(A, l, 1)
	ntt(B, l, 1)
	for i := 0; i < l; i++ {
		A[i] = int((int64(A[i]) * int64(B[i])) % Mod)
	}
	ntt(A, l, -1)
	return A[:n+m-1]
}

func polyAdd(a, b Poly) Poly {
	if len(a) < len(b) {
		a, b = b, a
	}
	res := make(Poly, len(a))
	copy(res, a)
	for i := 0; i < len(b); i++ {
		res[i] = (res[i] + b[i]) % Mod
	}
	return res
}

func shiftOne(v Poly) Poly {
	res := make(Poly, len(v)+1)
	copy(res[1:], v)
	return res
}

func addOne(v Poly) Poly {
	res := make(Poly, len(v))
	copy(res, v)
	res[0] = (res[0] + 1) % Mod
	return res
}

func decOne(v Poly) Poly {
	res := make(Poly, len(v))
	copy(res, v)
	res[0] = (res[0] + Mod - 1) % Mod
	return res
}

var (
	adj    [][]int
	son    []int
	sz     []int
	parent []int
	dp     []Poly
)

func dfs(x, fa int) {
	sz[x] = 1
	parent[x] = fa
	for _, y := range adj[x] {
		if y != fa {
			dfs(y, x)
			sz[x] += sz[y]
			if son[x] == 0 || sz[y] > sz[son[x]] {
				son[x] = y
			}
		}
	}
}

func solve(l, r int, q []Poly) (Poly, Poly) {
	if l == r {
		return addOne(q[l]), q[l]
	}
	mid := (l + r) >> 1
	lFi, lSe := solve(l, mid, q)
	rFi, rSe := solve(mid+1, r, q)
	term1 := polyMul(rSe, decOne(lFi))
	fi := polyAdd(term1, rFi)
	se := polyMul(lSe, rSe)
	return fi, se
}

func work(x int) {
	if son[x] == 0 {
		dp[x] = Poly{1, 1}
		return
	}
	for p := x; p != 0; p = son[p] {
		for _, y := range adj[p] {
			if y != parent[p] && y != son[p] {
				work(y)
			}
		}
	}
	chainQ := make([]Poly, 0)
	for p := x; p != 0; p = son[p] {
		var lightRes Poly = nil
		for _, y := range adj[p] {
			if y != parent[p] && y != son[p] {
				if lightRes == nil {
					lightRes = shiftOne(dp[y])
				} else {
					lightRes = polyMul(lightRes, shiftOne(dp[y]))
				}
			}
		}
		if lightRes == nil {
			chainQ = append(chainQ, decOne(Poly{1, 1}))
		} else {
			chainQ = append(chainQ, lightRes)
		}
	}
	for i, j := 0, len(chainQ)-1; i < j; i, j = i+1, j-1 {
		chainQ[i], chainQ[j] = chainQ[j], chainQ[i]
	}
	resFi, _ := solve(0, len(chainQ)-1, chainQ)
	dp[x] = resFi
}

func main() {
	sc := bufio.NewScanner(os.Stdin)
	sc.Split(bufio.ScanWords)
	next := func() string { sc.Scan(); return sc.Text() }
	nextInt := func() int { i, _ := strconv.Atoi(next()); return i }

	n := nextInt()
	X64, _ := strconv.ParseInt(next(), 10, 64)
	X := int(X64 % int64(Mod))

	adj = make([][]int, n+1)
	sz = make([]int, n+1)
	son = make([]int, n+1)
	parent = make([]int, n+1)
	dp = make([]Poly, n+1)

	for i := 0; i < n-1; i++ {
		u, v := nextInt(), nextInt()
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}

	dfs(1, 0)
	work(1)

	ans := 0
	coef := 1
	finalPoly := dp[1]
	for i := 1; i < len(finalPoly); i++ {
		term := (int64(coef) * int64(finalPoly[i])) % int64(Mod)
		ans = int((int64(ans) + term) % int64(Mod))
		num := (int64(X) + int64(i)) % int64(Mod)
		coef = int((int64(coef) * num % int64(Mod) * int64(inv(i))) % int64(Mod))
	}
	fmt.Println(ans)
}
