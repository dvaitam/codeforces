package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, a, b, k int
	if _, err := fmt.Fscan(reader, &n, &a, &b, &k); err != nil {
		return
	}
	freq := make(map[int]int, n)
	for i := 0; i < n; i++ {
		var r int
		fmt.Fscan(reader, &r)
		freq[r]++
	}

	unique := make([]int, 0, len(freq))
	for key := range freq {
		unique = append(unique, key)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(unique)))

	ans := 0
	for _, r := range unique {
		if r%k != 0 {
			continue
		}
		s := r / k
		cntR := freq[r]
		cntS, ok := freq[s]
		if !ok || cntR == 0 || cntS == 0 {
			continue
		}
		t := cntR / b
		if v := cntS / a; v < t {
			t = v
		}
		if t > 0 {
			freq[r] -= t * b
			freq[s] -= t * a
			ans += t
		}
	}

	fmt.Fprintln(writer, ans)
}
