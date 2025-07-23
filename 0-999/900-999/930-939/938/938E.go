package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const mod int64 = 1000000007

func modPow(a, e, m int64) int64 {
	res := int64(1)
	a %= m
	for e > 0 {
		if e&1 == 1 {
			res = res * a % m
		}
		a = a * a % m
		e >>= 1
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &arr[i])
	}

	// factorial of n modulo mod
	fact := int64(1)
	for i := 2; i <= n; i++ {
		fact = fact * int64(i) % mod
	}

	// prepare sorted copy for frequency and greater counts
	b := append([]int64(nil), arr...)
	sort.Slice(b, func(i, j int) bool { return b[i] < b[j] })

	// compress values
	vals := make([]int64, 0)
	freq := make([]int64, 0)
	for i := 0; i < n; {
		j := i
		for j < n && b[j] == b[i] {
			j++
		}
		vals = append(vals, b[i])
		freq = append(freq, int64(j-i))
		i = j
	}

	m := len(vals)
	greater := make([]int64, m)
	var sumGreater int64
	for i := m - 1; i >= 0; i-- {
		greater[i] = sumGreater
		sumGreater += freq[i]
	}

	var total int64
	for i := 0; i < m; i++ {
		if greater[i] == 0 {
			continue
		}
		inv := modPow(greater[i]+1, mod-2, mod)
		contrib := freq[i] % mod
		contrib = contrib * (vals[i] % mod) % mod
		contrib = contrib * inv % mod
		total = (total + contrib) % mod
	}

	ans := total * fact % mod
	fmt.Fprintln(writer, ans)
}
