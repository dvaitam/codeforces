package main

import (
	"bufio"
	"fmt"
	"os"
)

const N int = 200000
const B int = 450

var (
	presence   [N + 2]bool
	base       [N + 2]int64
	add        [N/B + 3]int64
	cnt        [N/B + 3]int64
	sumBase    [N/B + 3]int64
	sumContrib [N/B + 3]int64
	ans        int64
	d          int
)

func choose2(x int64) int64 {
	if x < 2 {
		return 0
	}
	return x * (x - 1) / 2
}

func rangeAdd(l, r int, delta int64) {
	if l > r {
		return
	}
	lb := l / B
	rb := r / B
	if lb == rb {
		end := r
		for i := l; i <= end; i++ {
			b := i / B
			if presence[i] {
				f := base[i] + add[b]
				if delta == 1 {
					ans += f
					sumContrib[b] += f
				} else {
					ans -= f - 1
					sumContrib[b] -= f - 1
				}
				sumBase[b] += delta
			}
			base[i] += delta
		}
		return
	}
	// left partial block
	le := (lb+1)*B - 1
	for i := l; i <= le; i++ {
		b := i / B
		if presence[i] {
			f := base[i] + add[b]
			if delta == 1 {
				ans += f
				sumContrib[b] += f
			} else {
				ans -= f - 1
				sumContrib[b] -= f - 1
			}
			sumBase[b] += delta
		}
		base[i] += delta
	}
	// full blocks
	for b := lb + 1; b < rb; b++ {
		if delta == 1 {
			diff := sumBase[b] + add[b]*cnt[b]
			ans += diff
			sumContrib[b] += diff
			add[b] += 1
		} else {
			diff := -(sumBase[b] + add[b]*cnt[b] - cnt[b])
			ans += diff
			sumContrib[b] += diff
			add[b] -= 1
		}
	}
	// right partial block
	rs := rb * B
	if rs < l {
		rs = l
	}
	for i := rs; i <= r; i++ {
		b := i / B
		if presence[i] {
			f := base[i] + add[b]
			if delta == 1 {
				ans += f
				sumContrib[b] += f
			} else {
				ans -= f - 1
				sumContrib[b] -= f - 1
			}
			sumBase[b] += delta
		}
		base[i] += delta
	}
}

func addPoint(x int) {
	b := x / B
	f := base[x] + add[b]
	presence[x] = true
	cnt[b]++
	sumBase[b] += base[x]
	val := choose2(f)
	sumContrib[b] += val
	ans += val
	r := x + d
	if r > N {
		r = N
	}
	rangeAdd(x+1, r, 1)
}

func removePoint(x int) {
	b := x / B
	f := base[x] + add[b]
	presence[x] = false
	cnt[b]--
	sumBase[b] -= base[x]
	val := choose2(f)
	sumContrib[b] -= val
	ans -= val
	r := x + d
	if r > N {
		r = N
	}
	rangeAdd(x+1, r, -1)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var q int
	fmt.Fscan(reader, &q, &d)
	for i := 0; i < q; i++ {
		var x int
		fmt.Fscan(reader, &x)
		if presence[x] {
			removePoint(x)
		} else {
			addPoint(x)
		}
		fmt.Fprintln(writer, ans)
	}
}
