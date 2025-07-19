package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

var (
	mod  int64
	n    int64
	a    []int64
	inv2 int64
	inv6 int64
)

func modadd(x, y int64) int64 {
	z := x + y
	if z >= mod {
		z -= mod
	}
	return z
}

func modsub(x, y int64) int64 {
	z := x - y
	if z < 0 {
		z += mod
	}
	return z
}

func modmul(x, y int64) int64 {
	return (x * y) % mod
}

func fpow(x, e int64) int64 {
	var res int64 = 1
	x %= mod
	for e > 0 {
		if e&1 != 0 {
			res = modmul(res, x)
		}
		x = modmul(x, x)
		e >>= 1
	}
	return res
}

func inv(x int64) int64 {
	return fpow(x, mod-2)
}

// sum of squares of arithmetic sequence length n, start x, diff d
func asSqrSum(x, d int64) int64 {
	// term1: n*x*x
	t1 := modmul(modmul(n, x), x)
	// term2: (n-1)*n*d*x
	t2 := modmul(modmul(modmul(n-1, n), d), x)
	// term3: d*d*(n-1)*n*(2*n-1)/6
	t3 := modmul(modmul(modmul(modmul(d, d), n-1), n), (2*n-1)%mod)
	t3 = modmul(t3, inv6)
	return modadd(modadd(t1, t2), t3)
}

// sum of cubes of arithmetic sequence length n, start x, diff d
func asCbcSum(x, d int64) int64 {
	// term1: n*x^3
	t1 := modmul(modmul(modmul(n, x), x), x)
	// term2: (n-1)^2 * n^2 * d^3 /4
	t2 := modmul(n-1, n-1)
	t2 = modmul(modmul(t2, n), n)
	t2 = modmul(modmul(modmul(t2, d), d), d)
	t2 = modmul(t2, inv2)
	t2 = modmul(t2, inv2)
	// term3: d^2 * x * (n-1)*n*(2*n-1)/2
	t3 := modmul(modmul(d, d), x)
	t3 = modmul(modmul(modmul(t3, n-1), n), (2*n-1)%mod)
	t3 = modmul(t3, inv2)
	// term4: 3*d*x^2*(n-1)*n/2
	t4 := modmul(modmul(d, x), x)
	t4 = modmul(modmul(modmul(t4, n-1), n), inv2)
	t4 = modmul(t4, 3)
	return modadd(modadd(modadd(t1, t2), t3), t4)
}

// build sequence and compare with sorted a
func check(x, d int64) bool {
	b := make([]int64, n)
	b[0] = x
	for i := int64(1); i < n; i++ {
		b[i] = modadd(b[i-1], d)
	}
	sort.Slice(b, func(i, j int) bool { return b[i] < b[j] })
	for i := range b {
		if b[i] != a[i] {
			return false
		}
	}
	return true
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	fmt.Fscan(reader, &mod, &n)
	a = make([]int64, n)
	for i := range a {
		fmt.Fscan(reader, &a[i])
	}
	if n == 1 {
		fmt.Fprintf(writer, "%d 0\n", a[0])
		return
	}
	sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
	// precompute sums
	var sum, sqrSum, cbcSum int64
	for _, v := range a {
		sum = modadd(sum, v)
		sqrSum = modadd(sqrSum, modmul(v, v))
		cbcSum = modadd(cbcSum, modmul(modmul(v, v), v))
	}
	inv2 = inv(2)
	inv6 = inv(6)
	invN := inv(n)
	// search d and x
	for j := int64(1); j < n; j++ {
		d := modsub(a[0], a[j])
		// x = (sum - n*(n-1)/2*d) * invN
		t := modmul(modmul(modmul(n, n-1), inv2), d)
		x := modmul(modsub(sum, t), invN)
		// check x in a
		idx := sort.Search(int(n), func(i int) bool { return a[i] >= x })
		if idx < int(n) && a[idx] == x {
			if asSqrSum(x, d) == sqrSum && asCbcSum(x, d) == cbcSum && check(x, d) {
				fmt.Fprintf(writer, "%d %d\n", x, d)
				return
			}
		}
	}
	fmt.Fprintln(writer, -1)
}
