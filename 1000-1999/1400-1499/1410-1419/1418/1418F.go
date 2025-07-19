package main

import (
	"bufio"
	"fmt"
	"os"
)

const N = 200100
const B = 420

var mindiv [N + 1]int
var pr []int

func getdiv(n int) []int {
	if n == 1 {
		return []int{1}
	}
	p := mindiv[n]
	deg := 0
	m := n
	for m%p == 0 {
		deg++
		m /= p
	}
	ret := getdiv(m)
	dp := p
	k := len(ret)
	for i := 0; i < deg; i++ {
		for j := 0; j < k; j++ {
			ret = append(ret, ret[j]*dp)
		}
		dp *= p
	}
	return ret
}

func query(ql, qr, l, r int) (int, int) {
	for ql <= qr && ql < B {
		if (r/ql)*ql >= l {
			return ql, (r / ql) * ql
		}
		ql++
	}
	for qr >= ql {
		fall := r / qr
		if fall*qr >= l {
			return qr, fall * qr
		}
		fall++
		qr = r / fall
	}
	return -1, -1
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	// sieve
	for i := 2; i <= N; i++ {
		if mindiv[i] == 0 {
			mindiv[i] = i
			pr = append(pr, i)
		}
		for _, p := range pr {
			if p > mindiv[i] || i*p > N {
				break
			}
			mindiv[i*p] = p
		}
	}

	var n, m int
	var l, r int64
	fmt.Fscan(reader, &n, &m, &l, &r)
	for x1 := 1; x1 <= n; x1++ {
		if x1 == n {
			fmt.Fprintln(writer, -1)
			continue
		}
		ry1 := min(int(r/int64(x1)), m)
		ly1 := min(int((l+int64(x1)-1)/int64(x1)), m+1)
		if ly1 > ry1 {
			fmt.Fprintln(writer, -1)
			continue
		}
		found := false
		for _, d := range getdiv(x1) {
			lz1 := x1/d + 1
			rz1 := n / d
			i, j := query(lz1, rz1, ly1, ry1)
			if i != -1 {
				z1 := i * d
				y2 := (x1 / d) * (j / i)
				fmt.Fprintf(writer, "%d %d %d %d\n", x1, j, z1, y2)
				found = true
				break
			}
		}
		if !found {
			fmt.Fprintln(writer, -1)
		}
	}
}
