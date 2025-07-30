package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

var reader *bufio.Reader

func getLine() string {
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}

func solve() {
	n, _ := strconv.Atoi(getLine())
	parts := strings.Split(getLine(), " ")
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i], _ = strconv.Atoi(parts[i])
	}
	sort.Ints(a)

	for k := n; k >= 0; k-- {
		if canWin(k, n, a) {
			fmt.Println(k)
			return
		}
	}
}

func canWin(k int, n int, a []int) bool {
	if k == 0 {
		return true
	}

	b := make([]int, n)
	copy(b, a)

	for i := 1; i <= k; i++ {
		limit := k - i + 1

		// Alice's move: remove the largest element that is less than or equal to the limit.
		// This is the optimal strategy for Alice as it saves smaller elements for later stages
		// which have stricter (smaller) limits.
		pos := sort.Search(len(b), func(j int) bool { return b[j] > limit })
		if pos == 0 {
			// No element is <= limit, so Alice cannot make a move.
			return false
		}
		idxToRemove := pos - 1
		b = append(b[:idxToRemove], b[idxToRemove+1:]...)

		// Bob's move: add the limit to the smallest remaining element.
		// This is the optimal strategy for Bob as it maximizes the chance of making a small element
		// too large for Alice in future stages.
		if len(b) > 0 {
			newVal := b[0] + limit
			b = b[1:]

			// Insert the new value back into the sorted slice `b`.
			insertPos := sort.Search(len(b), func(j int) bool { return b[j] >= newVal })
			
			// A standard way to insert into a slice in Go.
			b = append(b, 0) // Make space at the end.
			copy(b[insertPos+1:], b[insertPos:]) // Shift elements to the right.
			b[insertPos] = newVal // Place the new value.
		}
	}

	return true
}

func main() {
	reader = bufio.NewReader(os.Stdin)
	t, _ := strconv.Atoi(getLine())
	for i := 0; i < t; i++ {
		solve()
	}
}
