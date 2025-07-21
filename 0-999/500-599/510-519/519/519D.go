package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	weights := make([]int64, 26)
	for i := 0; i < 26; i++ {
		fmt.Fscan(reader, &weights[i])
	}
	var s string
	fmt.Fscan(reader, &s)
	n := len(s)
	// Compute prefix sums of weights
	prefix := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		prefix[i] = prefix[i-1] + weights[s[i-1]-'a']
	}
	// cnt holds for each letter a map from prefix sum to count of occurrences
	cnt := make([]map[int64]int64, 26)
	for i := 0; i < 26; i++ {
		cnt[i] = make(map[int64]int64)
	}
	var ans int64
	for i := 1; i <= n; i++ {
		c := s[i-1] - 'a'
		target := prefix[i-1]
		if v, ok := cnt[c][target]; ok {
			ans += v
		}
		// record this position as potential start
		cnt[c][prefix[i]]++
	}
	fmt.Println(ans)
}
