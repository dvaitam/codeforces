package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func allZero(a []int) bool {
	for _, v := range a {
		if v != 0 {
			return false
		}
	}
	return true
}

func applyShift(src []int, shift int, dst []int) bool {
	n := len(src)
	zero := true
	for i := 0; i < n; i++ {
		v := src[i] ^ src[(i+shift)%n]
		dst[i] = v
		if v != 0 {
			zero = false
		}
	}
	return zero
}

func minimalSteps(a []int) int {
	if allZero(a) {
		return 0
	}
	n := len(a)
	cur := append([]int(nil), a...)
	tmp := make([]int, n)
	steps := 0
	for p := bits.Len(uint(n)) - 1; p >= 0; p-- {
		shift := 1 << p
		if applyShift(cur, shift, tmp) {
			continue
		}
		steps += shift
		cur, tmp = tmp, cur
	}
	return steps + 1
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &arr[i])
	}
	ans := minimalSteps(arr)
	fmt.Fprintln(out, ans)
}
