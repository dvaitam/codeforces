package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	var k int64
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &arr[i])
	}
	sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })

	l, r := 0, n-1
	for l < r && k > 0 {
		for l < r && arr[l] == arr[l+1] {
			l++
		}
		for l < r && arr[r] == arr[r-1] {
			r--
		}
		if l >= r {
			break
		}
		cntL := int64(l + 1)
		cntR := int64(n - r)
		if cntL <= cntR {
			diff := arr[l+1] - arr[l]
			need := diff * cntL
			if need <= k {
				k -= need
				arr[l] += diff
				l++
			} else {
				arr[l] += k / cntL
				k = 0
			}
		} else {
			diff := arr[r] - arr[r-1]
			need := diff * cntR
			if need <= k {
				k -= need
				arr[r] -= diff
				r--
			} else {
				arr[r] -= k / cntR
				k = 0
			}
		}
	}
	if l >= r {
		fmt.Fprintln(writer, 0)
	} else {
		fmt.Fprintln(writer, arr[r]-arr[l])
	}
}
