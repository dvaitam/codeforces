package main

import (
	"bufio"
	"fmt"
	"os"
)

var bestWidth int
var bestA []int

func dfs(idx, lastA, rem, width int, arr []int) {
	if rem == 0 {
		if width < bestWidth {
			bestWidth = width
			bestA = append([]int(nil), arr...)
		}
		return
	}
	if width >= bestWidth {
		return
	}
	maxA := lastA
	if tmp := (rem - 1) / 2; tmp < maxA {
		maxA = tmp
	}
	for a := maxA; a >= 0; a-- {
		piece := 2*a + 1
		if piece > rem {
			continue
		}
		newW := width
		if idx+a > newW {
			newW = idx + a
		}
		arr = append(arr, a)
		dfs(idx+1, a, rem-piece, newW, arr)
		arr = arr[:len(arr)-1]
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)

	bestWidth = 101
	bestA = nil
	dfs(1, n, n, 0, []int{})
	if bestA == nil {
		fmt.Println(-1)
		return
	}
	m := bestWidth
	board := make([][]byte, m)
	for i := 0; i < m; i++ {
		row := make([]byte, m)
		for j := 0; j < m; j++ {
			row[j] = '.'
		}
		board[i] = row
	}
	for idx, a := range bestA {
		i := idx
		board[i][i] = 'o'
		for r := 1; r <= a && i+r < m; r++ {
			board[i][i+r] = 'o'
			board[i+r][i] = 'o'
		}
	}
	for i := 0; i < m/2; i++ {
		board[i], board[m-1-i] = board[m-1-i], board[i]
	}

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	fmt.Fprintln(out, m)
	for i := 0; i < m; i++ {
		out.Write(board[i])
		out.WriteByte('\n')
	}
}
