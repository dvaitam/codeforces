package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
)

var (
	n        int
	a, b     []*big.Int
	x        []int
	k        []int64
	children [][]int
	ok       = true
)

func dfs(v int) *big.Int {
	net := new(big.Int).Set(b[v])
	for _, u := range children[v] {
		net.Add(net, dfs(u))
	}
	net.Sub(net, a[v])
	if v == 1 {
		if net.Sign() < 0 {
			ok = false
		}
		return net
	}
	if net.Sign() >= 0 {
		return net
	}
	need := new(big.Int).Neg(net)
	req := new(big.Int).Mul(need, big.NewInt(k[v]))
	req.Neg(req)
	return req
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a = make([]*big.Int, n+1)
	b = make([]*big.Int, n+1)
	x = make([]int, n+1)
	k = make([]int64, n+1)
	children = make([][]int, n+1)

	for i := 1; i <= n; i++ {
		var val int64
		fmt.Fscan(in, &val)
		b[i] = big.NewInt(val)
	}
	for i := 1; i <= n; i++ {
		var val int64
		fmt.Fscan(in, &val)
		a[i] = big.NewInt(val)
	}
	for i := 2; i <= n; i++ {
		fmt.Fscan(in, &x[i], &k[i])
		children[x[i]] = append(children[x[i]], i)
	}

	dfs(1)

	if ok {
		fmt.Fprintln(out, "YES")
	} else {
		fmt.Fprintln(out, "NO")
	}
}
