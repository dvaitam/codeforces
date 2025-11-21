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

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		arr := make([]int, n)
		freq := make(map[int]int)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
			freq[arr[i]]++
		}
		sort.Ints(arr)

		dupVals := make([]int, 0)
		for val, count := range freq {
			if count >= 2 {
				dupVals = append(dupVals, val)
			}
		}
		pairs := make([]int, 0)
		for val, count := range freq {
			for k := 0; k < count/2; k++ {
				pairs = append(pairs, val)
			}
		}
		sort.Ints(pairs)

		if len(pairs) >= 2 {
			a := pairs[0]
			b := pairs[1]
			fmt.Fprintln(out, a, a, b, b)
			continue
		}

		if len(pairs) == 0 {
			fmt.Fprintln(out, -1)
			continue
		}

		x := pairs[0]
		if freq[x] < 2 {
			fmt.Fprintln(out, -1)
			continue
		}

		rest := make([]int, 0, n-2)
		removed := 0
		for _, val := range arr {
			if val == x && removed < 2 {
				removed++
				continue
			}
			rest = append(rest, val)
		}

		found := false
		i := 0
		for j := 1; j < len(rest); j++ {
			for i < j && rest[j]-rest[i] >= 2*x {
				i++
			}
			if i < j && rest[j]-rest[i] < 2*x {
				fmt.Fprintln(out, x, x, rest[i], rest[j])
				found = true
				break
			}
		}

		if !found {
			fmt.Fprintln(out, -1)
		}
	}
}
