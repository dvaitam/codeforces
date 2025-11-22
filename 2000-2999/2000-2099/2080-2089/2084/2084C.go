package main

import (
	"bufio"
	"fmt"
	"os"
)

type pair struct {
	x, y int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		b := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
		}

		posMap := make(map[pair][]int)
		for i := 0; i < n; i++ {
			key := pair{a[i], b[i]}
			posMap[key] = append(posMap[key], i)
		}

		oddSame := 0
		ok := true
		visited := make(map[pair]bool)
		for k, list := range posMap {
			if visited[k] {
				continue
			}
			if k.x == k.y {
				if len(list)%2 == 1 {
					oddSame++
				}
				visited[k] = true
				continue
			}
			rev := pair{k.y, k.x}
			l1 := len(list)
			l2 := len(posMap[rev])
			if l1 != l2 {
				ok = false
				break
			}
			visited[k] = true
			visited[rev] = true
		}
		if !ok || oddSame > 1 || (oddSame == 1 && n%2 == 0) {
			fmt.Fprintln(out, -1)
			continue
		}

		res := make([]int, n)
		l, r := 0, n-1
		var mid int = -1

		for k, list := range posMap {
			if k.x == k.y {
				continue
			}
			if k.x > k.y {
				continue // handle only once when x<y or first appearance
			}
			rev := pair{k.y, k.x}
			list2 := posMap[rev]
			for i := 0; i < len(list); i++ {
				res[l] = list[i]
				res[r] = list2[i]
				l++
				r--
			}
		}

		for k, list := range posMap {
			if k.x != k.y {
				continue
			}
			for i := 0; i+1 < len(list); i += 2 {
				res[l] = list[i]
				res[r] = list[i+1]
				l++
				r--
			}
			if len(list)%2 == 1 {
				mid = list[len(list)-1]
			}
		}

		if mid != -1 {
			res[l] = mid
			l++
		}

		// res now contains desired order of original indices
		arr := make([]int, n)
		posOf := make([]int, n)
		for i := 0; i < n; i++ {
			arr[i] = i
			posOf[i] = i
		}
		ops := make([][2]int, 0)

		for i := 0; i < n; i++ {
			need := res[i]
			if arr[i] == need {
				continue
			}
			j := posOf[need]
			arr[i], arr[j] = arr[j], arr[i]
			posOf[arr[i]] = i
			posOf[arr[j]] = j
			ops = append(ops, [2]int{i + 1, j + 1})
		}

		fmt.Fprintln(out, len(ops))
		for _, op := range ops {
			fmt.Fprintln(out, op[0], op[1])
		}
	}
}
