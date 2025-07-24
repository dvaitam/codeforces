package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var s string
	fmt.Fscan(in, &s)
	n := len(s)

	// Precompute nearest occurrences of each letter to the left and right
	left := make([][26]int, n)
	right := make([][26]int, n)

	for c := 0; c < 26; c++ {
		prev := -1
		for i := 0; i < n; i++ {
			left[i][c] = prev
			if int(s[i]-'a') == c {
				prev = i
			}
		}
		next := -1
		for i := n - 1; i >= 0; i-- {
			right[i][c] = next
			if int(s[i]-'a') == c {
				next = i
			}
		}
	}

	var total int64
	for start := 0; start < n; start++ {
		dist := make([]int, n)
		for i := range dist {
			dist[i] = -1
		}
		queue := make([]int, 0, n)
		head := 0
		dist[start] = 0
		queue = append(queue, start)
		for head < len(queue) {
			v := queue[head]
			head++
			d := dist[v] + 1
			for c := 0; c < 26; c++ {
				if l := left[v][c]; l != -1 && dist[l] == -1 {
					dist[l] = d
					queue = append(queue, l)
				}
				if r := right[v][c]; r != -1 && dist[r] == -1 {
					dist[r] = d
					queue = append(queue, r)
				}
			}
		}
		for _, v := range dist {
			total += int64(v)
		}
	}

	fmt.Println(total)
}
