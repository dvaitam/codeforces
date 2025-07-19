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

	var n, k int
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	// nxt[i] stores the next character index for transitions from state i
	nxt := make([]int, k)
	now := 0
	// use a byte slice for output
	res := make([]byte, 0, n)
	res = append(res, 'a')
	for i := 2; i <= n; i++ {
		ch := byte(nxt[now] + 'a')
		res = append(res, ch)
		p := nxt[now]
		nxt[now]++
		if nxt[now] >= k {
			for j := now + 1; j < k; j++ {
				if nxt[j] < now+1 {
					nxt[j] = now + 1
				}
			}
		}
		now = p
		if nxt[k-1] >= k {
			// reset
			for j := 0; j < k; j++ {
				nxt[j] = 0
			}
		}
	}
	writer.Write(res)
	writer.WriteByte('\n')
}
