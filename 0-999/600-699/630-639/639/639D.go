package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, k int
	var b, c int64
	if _, err := fmt.Fscan(in, &n, &k, &b, &c); err != nil {
		return
	}
	t := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &t[i])
	}

	// Using blogs to add 5 contribution may or may not be better than 5 comments
	per5 := b
	if per5 > 5*c {
		per5 = 5 * c
	}

	const inf int64 = 1<<63 - 1
	ans := inf

	for r := int64(0); r < 5; r++ {
		type pair struct{ tp, base int64 }
		arr := make([]pair, n)
		for i := 0; i < n; i++ {
			rem := ((t[i] % 5) + 5) % 5
			diff := (r - rem + 5) % 5
			tp := t[i] + diff
			w := tp / 5
			base := diff*c - w*per5
			arr[i] = pair{tp: tp, base: base}
		}

		sort.Slice(arr, func(i, j int) bool { return arr[i].tp < arr[j].tp })

		// max-heap for k smallest base values
		heap := make([]int64, 0, k)
		var sum int64
		push := func(x int64) {
			heap = append(heap, x)
			i := len(heap) - 1
			for i > 0 {
				p := (i - 1) / 2
				if heap[p] >= heap[i] {
					break
				}
				heap[p], heap[i] = heap[i], heap[p]
				i = p
			}
			sum += x
			if len(heap) > k {
				// pop max
				sum -= heap[0]
				heap[0] = heap[len(heap)-1]
				heap = heap[:len(heap)-1]
				i = 0
				for {
					l := 2*i + 1
					if l >= len(heap) {
						break
					}
					r := l + 1
					if r < len(heap) && heap[r] > heap[l] {
						l = r
					}
					if heap[i] >= heap[l] {
						break
					}
					heap[i], heap[l] = heap[l], heap[i]
					i = l
				}
			}
		}

		i := 0
		for i < n {
			curW := arr[i].tp / 5
			for i < n && arr[i].tp/5 == curW {
				push(arr[i].base)
				i++
			}
			if len(heap) == k {
				cost := sum + int64(k)*curW*per5
				if cost < ans {
					ans = cost
				}
			}
		}
	}

	fmt.Println(ans)
}
