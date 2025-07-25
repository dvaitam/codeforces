package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type pair struct {
	bal int
	idx int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	if _, err := fmt.Fscan(in, &s); err != nil {
		return
	}
	n := len(s)
	arr := make([]pair, n)
	bal := 0
	for i := 0; i < n; i++ {
		arr[i] = pair{bal, i}
		if s[i] == '(' {
			bal++
		} else {
			bal--
		}
	}
	sort.Slice(arr, func(i, j int) bool {
		if arr[i].bal == arr[j].bal {
			return arr[i].idx > arr[j].idx
		}
		return arr[i].bal < arr[j].bal
	})
	res := make([]byte, n)
	for i := 0; i < n; i++ {
		res[i] = s[arr[i].idx]
	}
	fmt.Fprintln(out, string(res))
}
