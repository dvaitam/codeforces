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

	var n, p int
	if _, err := fmt.Fscan(reader, &n, &p); err != nil {
		return
	}

	a := make([]int64, n)
	for i := 0; i < n; i++ {
		var x int64
		fmt.Fscan(reader, &x)
		a[i] = x
	}

	sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })

	// compute minimal starting value L
	L := int64(1)
	for i := 0; i < n; i++ {
		v := a[i] - int64(i)
		if v > L {
			L = v
		}
	}

	maxIdx := 2 * n
	prefix := make([]int, maxIdx)
	j := 0
	for idx := 0; idx < maxIdx; idx++ {
		val := L + int64(idx)
		for j < n && a[j] <= val {
			j++
		}
		prefix[idx] = j
	}

	positions := make([][]int, p)
	for idx := 0; idx < maxIdx; idx++ {
		r := (prefix[idx] - idx) % p
		if r < 0 {
			r += p
		}
		positions[r] = append(positions[r], idx)
	}

	result := make([]int64, 0, n)
	for k := 0; k < n; k++ {
		residue := (p - (k % p)) % p
		arr := positions[residue]
		// binary search for first index >= k
		i := sort.SearchInts(arr, k)
		if i < len(arr) && arr[i] <= k+n-1 {
			continue
		}
		result = append(result, L+int64(k))
	}

	fmt.Fprintln(writer, len(result))
	for i, v := range result {
		if i > 0 {
			fmt.Fprint(writer, " ")
		}
		fmt.Fprint(writer, v)
	}
	if len(result) > 0 {
		fmt.Fprintln(writer)
	}
}
