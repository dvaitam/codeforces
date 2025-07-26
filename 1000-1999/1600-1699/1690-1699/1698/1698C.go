package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		arr := make([]int, n)
		mp := make(map[int]bool)
		posCount, negCount, zeroCount := 0, 0, 0
		posVals := []int{}
		negVals := []int{}
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
			mp[arr[i]] = true
			if arr[i] > 0 {
				posCount++
				if len(posVals) < 2 {
					posVals = append(posVals, arr[i])
				}
			} else if arr[i] < 0 {
				negCount++
				if len(negVals) < 2 {
					negVals = append(negVals, arr[i])
				}
			} else {
				zeroCount++
			}
		}
		if posCount > 2 || negCount > 2 {
			fmt.Fprintln(out, "NO")
			continue
		}
		vals := []int{}
		vals = append(vals, posVals...)
		vals = append(vals, negVals...)
		for i := 0; i < zeroCount && i < 2; i++ {
			vals = append(vals, 0)
		}
		ok := true
		l := len(vals)
		for i := 0; i < l && ok; i++ {
			for j := i + 1; j < l && ok; j++ {
				for k := j + 1; k < l; k++ {
					s := vals[i] + vals[j] + vals[k]
					if !mp[s] {
						ok = false
						break
					}
				}
			}
		}
		if ok {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
