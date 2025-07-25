package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		d := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &d[i])
		}

		left := make([]int, n)
		right := make([]int, n)
		alive := make([]bool, n)
		for i := 0; i < n; i++ {
			if i == 0 {
				left[i] = -1
			} else {
				left[i] = i - 1
			}
			if i == n-1 {
				right[i] = -1
			} else {
				right[i] = i + 1
			}
			alive[i] = true
		}

		inQueue := make([]bool, n)
		queue := make([]int, 0)
		for i := 0; i < n; i++ {
			dmg := int64(0)
			if left[i] != -1 {
				dmg += a[left[i]]
			}
			if right[i] != -1 {
				dmg += a[right[i]]
			}
			if dmg > d[i] {
				queue = append(queue, i)
				inQueue[i] = true
			}
		}

		ans := make([]int, 0, n)
		for len(queue) > 0 {
			ans = append(ans, len(queue))
			nextCheck := make(map[int]struct{})
			// first remove all monsters in queue
			for _, idx := range queue {
				if !alive[idx] {
					continue
				}
				alive[idx] = false
				inQueue[idx] = false
				L := left[idx]
				R := right[idx]
				if L != -1 {
					right[L] = R
				}
				if R != -1 {
					left[R] = L
				}
				if L != -1 && alive[L] && !inQueue[L] {
					nextCheck[L] = struct{}{}
				}
				if R != -1 && alive[R] && !inQueue[R] {
					nextCheck[R] = struct{}{}
				}
			}
			nextQueue := make([]int, 0)
			for idx := range nextCheck {
				dmg := int64(0)
				if left[idx] != -1 {
					dmg += a[left[idx]]
				}
				if right[idx] != -1 {
					dmg += a[right[idx]]
				}
				if dmg > d[idx] {
					nextQueue = append(nextQueue, idx)
					inQueue[idx] = true
				}
			}
			queue = nextQueue
		}
		for len(ans) < n {
			ans = append(ans, 0)
		}
		for i, v := range ans {
			if i > 0 {
				writer.WriteByte(' ')
			}
			fmt.Fprint(writer, v)
		}
		fmt.Fprintln(writer)
	}
}
