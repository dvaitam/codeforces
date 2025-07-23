package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1000000007

func modPow(a, b int64) int64 {
	a %= mod
	res := int64(1)
	for b > 0 {
		if b&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		b >>= 1
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var p, k int64
	if _, err := fmt.Fscan(reader, &p, &k); err != nil {
		return
	}
	if k == 0 {
		fmt.Println(modPow(p, p-1))
		return
	}
	if k == 1 {
		fmt.Println(modPow(p, p))
		return
	}
	visited := make([]bool, p)
	var cycles int64
	for i := int64(1); i < p; i++ {
		if !visited[i] {
			cycles++
			j := i
			for !visited[j] {
				visited[j] = true
				j = (j * k) % p
			}
		}
	}
	fmt.Println(modPow(p, cycles))
}
