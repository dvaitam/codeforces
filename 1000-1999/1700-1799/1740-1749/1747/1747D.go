package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	a := make([]int, n)
	for i := range a {
		fmt.Fscan(in, &a[i])
	}

	px := make([]int, n+1)
	nz := make([]int, n+1)
	even := map[int][]int{}
	odd := map[int][]int{}
	for i := 0; i <= n; i++ {
		if i > 0 {
			px[i] = px[i-1] ^ a[i-1]
			if a[i-1] != 0 {
				nz[i] = nz[i-1] + 1
			} else {
				nz[i] = nz[i-1]
			}
		}
		if i%2 == 0 {
			even[px[i]] = append(even[px[i]], i)
		} else {
			odd[px[i]] = append(odd[px[i]], i)
		}
	}

	for ; q > 0; q-- {
		var l, r int
		fmt.Fscan(in, &l, &r)
		xor := px[r] ^ px[l-1]
		if xor != 0 {
			fmt.Fprintln(out, -1)
			continue
		}
		if nz[r]-nz[l-1] == 0 {
			fmt.Fprintln(out, 0)
			continue
		}
		if (r-l+1)%2 == 1 {
			fmt.Fprintln(out, 1)
			continue
		}
		if a[l-1] == 0 || a[r-1] == 0 {
			fmt.Fprintln(out, 1)
			continue
		}
		// check for an odd length prefix with xor 0
		par := (l - 1) % 2
		var arr []int
		if par == 0 {
			arr = odd[px[l-1]]
		} else {
			arr = even[px[l-1]]
		}
		i := sort.SearchInts(arr, l)
		if i < len(arr) && arr[i] <= r {
			fmt.Fprintln(out, 2)
			continue
		}
		par = r % 2
		if par == 0 {
			arr = odd[px[r]]
		} else {
			arr = even[px[r]]
		}
		i = sort.SearchInts(arr, l-1)
		if i < len(arr) && arr[i] < r {
			fmt.Fprintln(out, 2)
			continue
		}
		fmt.Fprintln(out, -1)
	}
}
