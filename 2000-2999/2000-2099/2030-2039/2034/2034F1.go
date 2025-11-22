package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const MOD int64 = 998244353
const MAXN = 400000

var fact []int64
var invFact []int64

type event struct {
	t    int // step index (number of draws)
	x    int // number of red gems drawn by step t
	base int // base satchel value after t draws ignoring any doublings
	prob int64
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

func initFact() {
	fact = make([]int64, MAXN+1)
	invFact = make([]int64, MAXN+1)
	fact[0] = 1
	for i := 1; i <= MAXN; i++ {
		fact[i] = fact[i-1] * int64(i) % MOD
	}
	invFact[MAXN] = modPow(fact[MAXN], MOD-2)
	for i := MAXN; i >= 1; i-- {
		invFact[i-1] = invFact[i] * int64(i) % MOD
	}
}

func comb(n, k int) int64 {
	if k < 0 || k > n {
		return 0
	}
	return fact[n] * invFact[k] % MOD * invFact[n-k] % MOD
}

// hyper returns probability of drawing needR reds in len draws
// from a pool with r reds and b blues.
func hyper(r, b, length, needR int) int64 {
	if needR < 0 || needR > r || length-needR < 0 || length-needR > b {
		return 0
	}
	num := comb(r, needR) * comb(b, length-needR) % MOD
	den := comb(r+b, length)
	return num * modPow(den, MOD-2) % MOD
}

func main() {
	in := bufio.NewReader(os.Stdin)
	initFact()

	var T int
	fmt.Fscan(in, &T)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	for ; T > 0; T-- {
		var n, m, k int
		fmt.Fscan(in, &n, &m, &k)

		events := make([]event, 0, k)
		for i := 0; i < k; i++ {
			var r, b int
			fmt.Fscan(in, &r, &b)
			t := n + m - (r + b)
			x := n - r
			base := t + x
			events = append(events, event{t: t, x: x, base: base})
		}

		sort.Slice(events, func(i, j int) bool {
			return events[i].t < events[j].t
		})

		for i := range events {
			ev := &events[i]
			ev.prob = hyper(n, m, ev.t, ev.x)
		}

		A := make([]int64, len(events))
		for i, ev := range events {
			val := int64(ev.base%int(MOD)) * ev.prob % MOD
			for j := 0; j < i; j++ {
				prev := events[j]
				remR := n - prev.x
				remB := m - (prev.t - prev.x)
				d := ev.t - prev.t
				need := ev.x - prev.x
				pCond := hyper(remR, remB, d, need)
				val = (val + A[j]*pCond) % MOD
			}
			A[i] = val
		}

		ans := (int64(2*n+m) % MOD)
		for _, v := range A {
			ans = (ans + v) % MOD
		}
		if ans < 0 {
			ans += MOD
		}
		fmt.Fprintln(writer, ans)
	}
}
