package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var k int
	if _, err := fmt.Fscan(reader, &k); err != nil {
		return
	}
	var n string
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	sum := 0
	diffs := make([]int, len(n))
	for i := 0; i < len(n); i++ {
		d := int(n[i] - '0')
		sum += d
		diffs[i] = 9 - d
	}
	if sum >= k {
		fmt.Println(0)
		return
	}
	need := k - sum
	sort.Slice(diffs, func(i, j int) bool {
		return diffs[i] > diffs[j]
	})
	changes := 0
	for _, inc := range diffs {
		need -= inc
		changes++
		if need <= 0 {
			break
		}
	}
	fmt.Println(changes)
}
