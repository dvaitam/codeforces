package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

// extend computes possible states after appending x.
// each state is (incLast, decLast).
func extend(states [][2]int, x int) [][2]int {
	next := make([][2]int, 0, len(states)*2)
	for _, st := range states {
		if x > st[0] {
			next = append(next, [2]int{x, st[1]})
		}
		if x < st[1] {
			next = append(next, [2]int{st[0], x})
		}
	}
	res := make([][2]int, 0, 2)
	for _, st := range next {
		dominated := false
		for i := 0; i < len(res); {
			ot := res[i]
			if st[0] >= ot[0] && st[1] <= ot[1] {
				dominated = true
				break
			}
			if st[0] <= ot[0] && st[1] >= ot[1] {
				res[i] = res[len(res)-1]
				res = res[:len(res)-1]
			} else {
				i++
			}
		}
		if !dominated {
			res = append(res, st)
		}
	}
	return res
}

func isDecinc(arr []int) bool {
	states := [][2]int{{math.MinInt64, math.MaxInt64}}
	for _, x := range arr {
		states = extend(states, x)
		if len(states) == 0 {
			return false
		}
	}
	return true
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	p := make([]int, n)
	for i := range p {
		fmt.Fscan(in, &p[i])
	}

	var ans int64
	for l := 0; l < n; l++ {
		states := [][2]int{{math.MinInt64, math.MaxInt64}}
		for r := l; r < n; r++ {
			states = extend(states, p[r])
			if len(states) == 0 {
				break
			}
			ans++
		}
	}
	fmt.Fprintln(out, ans)
}
