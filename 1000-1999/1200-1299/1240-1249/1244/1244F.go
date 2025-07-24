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

	var n int
	var k int64
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	var s string
	fmt.Fscan(reader, &s)

	b := []byte(s)
	stable := make([]bool, n)
	dist := make([]int64, n)
	for i := range dist {
		dist[i] = -1
	}
	anyStable := false
	for i := 0; i < n; i++ {
		prev := (i - 1 + n) % n
		next := (i + 1) % n
		if b[i] == b[prev] || b[i] == b[next] {
			stable[i] = true
			anyStable = true
		}
	}
	if !anyStable {
		if k%2 == 1 {
			for i := 0; i < n; i++ {
				if b[i] == 'W' {
					b[i] = 'B'
				} else {
					b[i] = 'W'
				}
			}
		}
		fmt.Fprintln(writer, string(b))
		return
	}

	// BFS from stable positions
	q := make([]int, 0)
	for i := 0; i < n; i++ {
		if stable[i] {
			dist[i] = 0
			q = append(q, i)
		}
	}
	head := 0
	for head < len(q) {
		cur := q[head]
		head++
		for _, nb := range []int{(cur - 1 + n) % n, (cur + 1) % n} {
			if dist[nb] == -1 {
				dist[nb] = dist[cur] + 1
				b[nb] = b[cur]
				q = append(q, nb)
			}
		}
	}

	orig := []byte(s)
	res := make([]byte, n)
	for i := 0; i < n; i++ {
		if dist[i] != -1 && dist[i] <= k {
			res[i] = b[i]
		} else {
			if k%2 == 0 {
				res[i] = orig[i]
			} else {
				if orig[i] == 'W' {
					res[i] = 'B'
				} else {
					res[i] = 'W'
				}
			}
		}
	}
	fmt.Fprintln(writer, string(res))
}
