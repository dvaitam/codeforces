package main

import (
	"bufio"
	"fmt"
	"os"
)

const INF int64 = 1 << 60

func factor(x int64) []int64 {
	if x < 2 {
		return nil
	}
	factors := []int64{}
	for d := int64(2); d*d <= x; d++ {
		if x%d == 0 {
			factors = append(factors, d)
			for x%d == 0 {
				x /= d
			}
		}
	}
	if x > 1 {
		factors = append(factors, x)
	}
	return factors
}

func cost(val, p, b int64) int64 {
	if val%p == 0 {
		return 0
	}
	if (val-1)%p == 0 || (val+1)%p == 0 {
		return b
	}
	return INF
}

func solve(arr []int64, a, b int64, p int64) int64 {
	n := len(arr)
	dp0 := int64(0)
	dp1 := INF
	dp2 := INF
	for i := 0; i < n; i++ {
		c := cost(arr[i], p, b)
		ndp0 := INF
		if dp0 < INF && c < INF {
			ndp0 = dp0 + c
		}
		ndp1 := INF
		m01 := dp0
		if dp1 < m01 {
			m01 = dp1
		}
		if m01 < INF {
			ndp1 = m01 + a
		}
		ndp2 := INF
		m12 := dp1
		if dp2 < m12 {
			m12 = dp2
		}
		if m12 < INF && c < INF {
			ndp2 = m12 + c
		}
		dp0, dp1, dp2 = ndp0, ndp1, ndp2
	}
	res := dp0
	if dp2 < res {
		res = dp2
	}
	if dp1 < res && dp1 != int64(n)*a {
		res = dp1
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	var a, b int64
	if _, err := fmt.Fscan(reader, &n, &a, &b); err != nil {
		return
	}
	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &arr[i])
	}
	candidates := make(map[int64]struct{})
	nums := []int64{arr[0] - 1, arr[0], arr[0] + 1, arr[n-1] - 1, arr[n-1], arr[n-1] + 1}
	for _, v := range nums {
		for _, p := range factor(abs64(v)) {
			candidates[p] = struct{}{}
		}
	}
	ans := INF
	for p := range candidates {
		c := solve(arr, a, b, p)
		if c < ans {
			ans = c
		}
	}
	fmt.Println(ans)
}

func abs64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}
