package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	fmt.Fscan(reader, &n, &m)
	arr := make([]int, n)
	freq := make([]int, m+1)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &arr[i])
		if arr[i] <= m {
			freq[arr[i]]++
		}
	}

	// multiple[x] = number of elements divisible by x
	multiple := make([]int, m+1)
	for i := 1; i <= m; i++ {
		for j := i; j <= m; j += i {
			multiple[i] += freq[j]
		}
	}

	// count[l] = number of elements that divide l
	count := make([]int, m+1)
	for v := 1; v <= m; v++ {
		if freq[v] == 0 {
			continue
		}
		for j := v; j <= m; j += v {
			count[j] += freq[v]
		}
	}

	// sieve for smallest prime factor
	spf := make([]int, m+1)
	for i := 2; i <= m; i++ {
		if spf[i] == 0 {
			for j := i; j <= m; j += i {
				if spf[j] == 0 {
					spf[j] = i
				}
			}
		}
	}

	bestL, bestCount := 1, 0
	for l := 1; l <= m; l++ {
		if count[l] <= bestCount {
			continue
		}
		temp := l
		possible := true
		for temp > 1 {
			p := spf[temp]
			pow := 1
			for temp%p == 0 {
				temp /= p
				pow *= p
			}
			if multiple[pow] == 0 {
				possible = false
				break
			}
		}
		if possible {
			bestL = l
			bestCount = count[l]
		}
	}

	fmt.Fprintln(writer, bestL, bestCount)
	if bestCount == 0 {
		fmt.Fprintln(writer)
		return
	}
	first := true
	for i, val := range arr {
		if bestL%val == 0 {
			if !first {
				writer.WriteByte(' ')
			}
			first = false
			fmt.Fprint(writer, i+1)
		}
	}
	writer.WriteByte('\n')
}
