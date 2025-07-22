package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &arr[i])
	}

	expected := make([]int, n)
	copy(expected, arr)
	sort.Slice(expected, func(i, j int) bool { return expected[i] > expected[j] })

	a := make([]int, n)
	copy(a, arr)
	output := make([]int, 0, n)

	for len(a) > 0 {
		changed := true
		for changed {
			changed = false
			for i := 0; i < len(a)-1; i++ {
				if a[i]-a[i+1] >= 2 {
					a[i]--
					a[i+1]++
					changed = true
				} else if a[i+1]-a[i] >= 2 {
					a[i+1]--
					a[i]++
					changed = true
				}
			}
		}
		idx := 0
		for i := 1; i < len(a); i++ {
			if a[i] > a[idx] {
				idx = i
			}
		}
		output = append(output, a[idx])
		a = append(a[:idx], a[idx+1:]...)
	}

	for i := 0; i < n; i++ {
		if output[i] != expected[i] {
			fmt.Println("NO")
			return
		}
	}
	fmt.Println("YES")
}
