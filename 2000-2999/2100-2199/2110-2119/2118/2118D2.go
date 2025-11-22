package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type light struct {
	pos int64
	d   int64
}

func mod(x, k int64) int64 {
	x %= k
	if x < 0 {
		x += k
	}
	return x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		var k int64
		fmt.Fscan(in, &n, &k)

		p := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &p[i])
		}
		d := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &d[i])
		}

		keyS := make([]int64, n) // (p + d) % k
		keyT := make([]int64, n) // (p - d) % k
		for i := 0; i < n; i++ {
			keyS[i] = mod(p[i]+d[i], k)
			keyT[i] = mod(p[i]-d[i], k)
		}

		// Build maps from residue to indices (in increasing position order).
		sMap := make(map[int64][]int)
		tMap := make(map[int64][]int)
		for i := 0; i < n; i++ {
			sMap[keyS[i]] = append(sMap[keyS[i]], i)
			tMap[keyT[i]] = append(tMap[keyT[i]], i)
		}

		// Transitions: Ltrans goes to previous light with same keyS (moving left after a turn),
		// Rtrans goes to next light with same keyT (moving right after a turn).
		Ltrans := make([]int, n)
		for i := range Ltrans {
			Ltrans[i] = -1
		}
		for _, lst := range sMap {
			for i := 1; i < len(lst); i++ {
				Ltrans[lst[i]] = lst[i-1]
			}
		}

		Rtrans := make([]int, n)
		for i := range Rtrans {
			Rtrans[i] = -1
		}
		for _, lst := range tMap {
			for i := 0; i+1 < len(lst); i++ {
				Rtrans[lst[i]] = lst[i+1]
			}
		}

		// State graph has 2*n nodes: 0..n-1 represent state (at light i, moving right),
		// n..2n-1 represent state (at light i, moving left).
		state := make([]int8, 2*n) // 0=unvisited,1=escapes,2=visiting,3=cycle

		nextNode := func(v int) int {
			if v < n {
				nxt := Rtrans[v]
				if nxt == -1 {
					return -1
				}
				return n + nxt // switch to moving left at the next light
			}
			idx := v - n
			nxt := Ltrans[idx]
			if nxt == -1 {
				return -1
			}
			return nxt // switch to moving right at the previous light
		}

		for v := 0; v < 2*n; v++ {
			if state[v] != 0 {
				continue
			}
			stack := []int{}
			cur := v
			for {
				if cur == -1 {
					// escape
					for len(stack) > 0 {
						last := stack[len(stack)-1]
						stack = stack[:len(stack)-1]
						state[last] = 1
					}
					break
				}
				if state[cur] == 1 {
					for len(stack) > 0 {
						last := stack[len(stack)-1]
						stack = stack[:len(stack)-1]
						state[last] = 1
					}
					break
				}
				if state[cur] == 2 {
					// detected a cycle; all nodes in stack are part of cycle or lead into it
					for len(stack) > 0 {
						last := stack[len(stack)-1]
						stack = stack[:len(stack)-1]
						state[last] = 3
						if last == cur {
							break
						}
					}
					// remaining nodes (if any) also lead into cycle
					for len(stack) > 0 {
						last := stack[len(stack)-1]
						stack = stack[:len(stack)-1]
						state[last] = 3
					}
					break
				}
				if state[cur] == 3 {
					for len(stack) > 0 {
						last := stack[len(stack)-1]
						stack = stack[:len(stack)-1]
						state[last] = 3
					}
					break
				}
				// unvisited
				state[cur] = 2
				stack = append(stack, cur)
				cur = nextNode(cur)
			}
		}

		// Prepare position slices for binary search per residue of keyT.
		posByT := make(map[int64][]int64)
		idxByT := make(map[int64][]int)
		for i, resid := range keyT {
			posByT[resid] = append(posByT[resid], p[i])
			idxByT[resid] = append(idxByT[resid], i)
		}

		var q int
		fmt.Fscan(in, &q)
		for ; q > 0; q-- {
			var a int64
			fmt.Fscan(in, &a)
			r := mod(a, k)
			posList, ok := posByT[r]
			if !ok {
				fmt.Fprintln(out, "YES")
				continue
			}
			idx := sort.Search(len(posList), func(i int) bool { return posList[i] >= a })
			if idx == len(posList) {
				fmt.Fprintln(out, "YES")
				continue
			}
			lightIdx := idxByT[r][idx]
			// We arrive at this light from the left; state is moving left after turning here.
			if state[n+lightIdx] == 1 {
				fmt.Fprintln(out, "YES")
			} else {
				fmt.Fprintln(out, "NO")
			}
		}
	}
}
