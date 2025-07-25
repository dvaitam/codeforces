package main

import "fmt"

const mod int64 = 1000000007

func main() {
	var n int64
	if _, err := fmt.Scan(&n); err != nil {
		return
	}
	res := int64(1)
	for i := int64(2); i <= n; i++ {
		res = res * i % mod
	}
	fmt.Println(res)
}
