package main

import "fmt"

func C(n, k int64) int64 {
	if k < 0 || k > n {
		return 0
	}
	if k > n-k {
		k = n - k
	}
	res := int64(1)
	for i := int64(1); i <= k; i++ {
		res = res * (n - k + i) / i
	}
	return res
}

func main() {
	var n int64
	if _, err := fmt.Scan(&n); err != nil {
		return
	}
	a := C(n+4, 5)
	b := C(n+2, 3)
	fmt.Println(a * b)
}
