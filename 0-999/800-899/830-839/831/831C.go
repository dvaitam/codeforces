package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var k, n int
	if _, err := fmt.Fscan(in, &k, &n); err != nil {
		return
	}
	a := make([]int, k)
	for i := 0; i < k; i++ {
		fmt.Fscan(in, &a[i])
	}
	b := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &b[i])
	}

	pref := make([]int, k)
	sum := 0
	for i := 0; i < k; i++ {
		sum += a[i]
		pref[i] = sum
	}

	prefSet := make(map[int]struct{})
	for _, val := range pref {
		prefSet[val] = struct{}{}
	}

	candidateSet := make(map[int]struct{})
	for _, val := range pref {
		candidateSet[b[0]-val] = struct{}{}
	}

	cnt := 0
	for x := range candidateSet {
		ok := true
		for _, bj := range b {
			if _, ok2 := prefSet[bj-x]; !ok2 {
				ok = false
				break
			}
		}
		if ok {
			cnt++
		}
	}

	fmt.Println(cnt)
}
