package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const maxN = 400000
const sieveLimit = 6000000

var primePrefix []int64

func init() {
	sieve := make([]bool, sieveLimit+1)
	primePrefix = make([]int64, maxN+1)
	count := 0
	sum := int64(0)
	for i := 2; i <= sieveLimit && count < maxN; i++ {
		if !sieve[i] {
			sum += int64(i)
			count++
			primePrefix[count] = sum
			if count == maxN {
				break
			}
			for j := i * 2; j <= sieveLimit; j += i {
				sieve[j] = true
			}
		}
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
		}
		sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
		prefix := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			prefix[i] = prefix[i-1] + arr[i-1]
		}
		low, high := 0, n
		for low < high {
			mid := (low + high + 1) >> 1
			sumLargest := prefix[n] - prefix[n-mid]
			if sumLargest >= primePrefix[mid] {
				low = mid
			} else {
				high = mid - 1
			}
		}
		fmt.Fprintln(out, n-low)
	}
}
