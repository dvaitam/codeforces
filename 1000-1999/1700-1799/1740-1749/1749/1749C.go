package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// check if Alice can win with given k
func canWin(a []int, k int) bool {
	arr := make([]int, len(a))
	copy(arr, a)
	sort.Ints(arr)
	for i := 0; i < k; i++ {
		need := k - i
		// largest index where arr[idx] <= need
		idx := sort.Search(len(arr), func(j int) bool { return arr[j] > need }) - 1
		if idx < 0 {
			return false
		}
		// remove arr[idx]
		arr = append(arr[:idx], arr[idx+1:]...)
		if len(arr) > 0 {
			arr[len(arr)-1] += need
			sort.Ints(arr)
		}
	}
	return true
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		ans := 0
		for k := n; k >= 0; k-- {
			if canWin(a, k) {
				ans = k
				break
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
