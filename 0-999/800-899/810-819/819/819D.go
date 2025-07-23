package main

import (
	"bufio"
	"fmt"
	"os"
)

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func exgcd(a, b int64) (int64, int64, int64) {
	if b == 0 {
		return a, 1, 0
	}
	g, x1, y1 := exgcd(b, a%b)
	x, y := y1, x1-(a/b)*y1
	return g, x, y
}

func modInv(a, mod int64) int64 {
	g, x, _ := exgcd(a, mod)
	if g != 1 {
		return 0
	}
	x %= mod
	if x < 0 {
		x += mod
	}
	return x
}

func floorSum(n, m, a, b int64) int64 {
	var ans int64
	for {
		if a >= m {
			ans += (n - 1) * n * (a / m) / 2
			a %= m
		}
		if b >= m {
			ans += n * (b / m)
			b %= m
		}
		yMax := a*n + b
		if yMax < m {
			break
		}
		n = yMax / m
		b = yMax % m
		m, a = a, m
	}
	return ans
}

func countPrefix(n, m, d, R int64) int64 {
	if R == 0 || n == 0 {
		return 0
	}
	c1 := floorSum(n, m, d, 0)
	c2 := floorSum(n, m, d, m-R)
	return n + c1 - c2
}

func countInterval(l, r, m, d, R int64) int64 {
	if r < l || R == 0 {
		return 0
	}
	return countPrefix(r+1, m, d, R) - countPrefix(l, m, d, R)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int64
	var n int
	fmt.Fscan(in, &T, &n)
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	prefix := make([]int64, n)
	var sum int64
	for i := 0; i < n; i++ {
		prefix[i] = sum
		sum += a[i]
	}
	A := sum
	g := gcd(A, T)
	A1 := A / g
	T1 := T / g
	invT1 := modInv(T1%A1, A1)

	type pair struct {
		t  int64
		id int
	}
	groups := make(map[int64][]pair)
	for i, s := range prefix {
		r := s % g
		s2 := (s - r) / g
		tVal := (s2 * invT1) % A1
		groups[r] = append(groups[r], pair{tVal, i})
	}

	Q := T / g
	base := Q / A1
	R := Q % A1
	ans := make([]int64, n)
	invStep := modInv(invT1, A1)

	for _, arr := range groups {
		// sort by t
		for i := 0; i < len(arr); i++ {
			for j := i + 1; j < len(arr); j++ {
				if arr[j].t < arr[i].t {
					arr[i], arr[j] = arr[j], arr[i]
				}
			}
		}
		m := len(arr)
		tVals := make([]int64, m)
		ids := make([]int, m)
		for i := 0; i < m; i++ {
			tVals[i] = arr[i].t
			ids[i] = arr[i].id
		}
		// full cycle counts
		cnt := make([]int64, m)
		cnt[0] += tVals[0] + 1
		for i := 1; i < m; i++ {
			cnt[i] += tVals[i] - tVals[i-1]
		}
		cnt[0] += A1 - 1 - tVals[m-1]
		for i := 0; i < m; i++ {
			ans[ids[i]] += cnt[i] * base
		}
		// remainder part
		for i := 0; i < m; i++ {
			l := int64(0)
			if i > 0 {
				l = tVals[i-1] + 1
			}
			r := tVals[i]
			ans[ids[i]] += countInterval(l, r, A1, invStep, R)
		}
		l := tVals[m-1] + 1
		r := A1 - 1
		ans[ids[0]] += countInterval(l, r, A1, invStep, R)
	}

	for i := 0; i < n; i++ {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, ans[i])
	}
	fmt.Fprintln(out)
}
