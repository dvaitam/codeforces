package main

import (
	"bufio"
	"fmt"
	"os"
)

// DSU-based helper to skip already filled positions
var nextPos []int

func getNext(x int) int {
	if nextPos[x] != x {
		nextPos[x] = getNext(nextPos[x])
	}
	return nextPos[x]
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	type Info struct {
		str string
		pos []int
	}
	infos := make([]Info, n)
	maxLen := 0
	for i := 0; i < n; i++ {
		var t string
		var k int
		fmt.Fscan(in, &t, &k)
		p := make([]int, k)
		for j := 0; j < k; j++ {
			fmt.Fscan(in, &p[j])
			end := p[j] + len(t) - 1
			if end > maxLen {
				maxLen = end
			}
		}
		infos[i] = Info{t, p}
	}

	res := make([]byte, maxLen+1) // 1-based indexing
	for i := 1; i <= maxLen; i++ {
		res[i] = 'a'
	}

	nextPos = make([]int, maxLen+2)
	for i := 0; i <= maxLen+1; i++ {
		nextPos[i] = i
	}

	for _, info := range infos {
		l := len(info.str)
		for _, st := range info.pos {
			cur := getNext(st)
			for cur < st+l {
				res[cur] = info.str[cur-st]
				nextPos[cur] = getNext(cur + 1)
				cur = getNext(cur)
			}
		}
	}

	out.Write(res[1:])
}
