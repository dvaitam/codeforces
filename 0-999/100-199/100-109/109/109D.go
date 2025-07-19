package main

import (
 	"bufio"
 	"fmt"
 	"os"
 	"sort"
)

// pair holds a value and its original index
type pair struct {
 	val int
 	idx int
}

// lucky returns true if x's digits are all 4 or 7 (matching C++ behavior: lucky(0)==true)
func lucky(x int) bool {
 	for x > 0 {
 		d := x % 10
 		if d != 4 && d != 7 {
 			return false
 		}
 		x /= 10
 	}
 	return true
}

func main() {
 	in := bufio.NewReader(os.Stdin)
 	out := bufio.NewWriter(os.Stdout)
 	defer out.Flush()

 	var N int
 	fmt.Fscan(in, &N)

 	ar := make([]pair, N+1)
 	for i := 1; i <= N; i++ {
 		var tmp int
 		fmt.Fscan(in, &tmp)
 		ar[i] = pair{val: tmp, idx: i}
 	}

 	// sort ar[1..N] by value then index
 	sub := ar[1:]
 	sort.Slice(sub, func(i, j int) bool {
 		if sub[i].val != sub[j].val {
 			return sub[i].val < sub[j].val
 		}
 		return sub[i].idx < sub[j].idx
 	})

 	// find a lucky pivot (original index)
 	var luck int
 	for i := 1; i <= N; i++ {
 		if lucky(ar[i].val) {
 			luck = ar[i].idx
 			break
 		}
 	}

 	// record moves
 	X := make([]int, 0, 2*N)
 	Y := make([]int, 0, 2*N)
 	vis := make([]bool, N+1)
 	mv := func(a, b int) {
 		X = append(X, a)
 		Y = append(Y, b)
 	}

 	if luck == 0 {
 		// no lucky pivot: check if already sorted
 		ok := true
 		for i := 1; i+1 <= N; i++ {
 			if ar[i].idx > ar[i+1].idx {
 				ok = false
 				break
 			}
 		}
 		if !ok {
 			fmt.Fprintln(out, -1)
 		} else {
 			fmt.Fprintln(out, 0)
 		}
 		return
 	}

 	// first cycle with pivot
 	cur := luck
 	for ar[cur].idx != luck {
 		vis[cur] = true
 		mv(cur, ar[cur].idx)
 		cur = ar[cur].idx
 	}
 	vis[cur] = true
 	luck = cur

 	// process remaining cycles
 	for i := 1; i <= N; i++ {
 		if !vis[i] {
 			mv(luck, i)
 			cur = i
 			for ar[cur].idx != i {
 				vis[cur] = true
 				nxt := ar[cur].idx
 				mv(cur, nxt)
 				cur = nxt
 			}
 			vis[cur] = true
 			mv(luck, cur)
 		}
 	}

 	// output moves
 	fmt.Fprintln(out, len(X))
 	for i := range X {
 		fmt.Fprintln(out, X[i], Y[i])
 	}
}
