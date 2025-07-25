package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
)

const MOD int64 = 1000000007

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, p int
	if _, err := fmt.Fscan(in, &n, &p); err != nil {
		return
	}
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &arr[i])
	}

	filtered := make([]int, 0, n)
	for _, v := range arr {
		if bits.Len(uint(v)) <= p {
			filtered = append(filtered, v)
		}
	}
	sort.Ints(filtered)

	base := make(map[int]struct{})
	for _, x := range filtered {
		y := x
		keep := true
		for y > 0 {
			if _, ok := base[y]; ok {
				keep = false
				break
			}
			if y%2 == 1 {
				y /= 2
			} else if y%4 == 0 {
				y /= 4
			} else {
				break
			}
		}
		if keep {
			base[x] = struct{}{}
		}
	}

	fib := make([]int64, p+5)
	fib[1] = 1
	for i := 2; i < len(fib); i++ {
		fib[i] = (fib[i-1] + fib[i-2]) % MOD
	}

	ans := int64(0)
	for x := range base {
		l := bits.Len(uint(x))
		if l > p {
			continue
		}
		idx := p - l + 3
		val := fib[idx] - 1
		if val < 0 {
			val += MOD
		}
		ans += val
		ans %= MOD
	}
	fmt.Fprintln(out, ans)
}
