package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

const eps = 1e-12

type node struct {
	prev int
	idx  int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, d int
	if _, err := fmt.Fscan(in, &n, &d); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	dpVal := make([]float64, 10)
	dpLen := make([]int, 10)
	stateID := make([]int, 10)
	for i := 0; i < 10; i++ {
		dpVal[i] = math.Inf(-1)
		dpLen[i] = 0
		stateID[i] = -1
	}
	nodes := []node{{prev: -1, idx: -1}}
	dpVal[1] = 0
	stateID[1] = 0

	for idx := 0; idx < n; idx++ {
		val := math.Log(float64(a[idx]))
		newVal := make([]float64, 10)
		newLen := make([]int, 10)
		newStateID := make([]int, 10)
		copy(newVal, dpVal)
		copy(newLen, dpLen)
		copy(newStateID, stateID)
		for r := 0; r < 10; r++ {
			if stateID[r] == -1 {
				continue
			}
			nr := (r * a[idx]) % 10
			cand := dpVal[r] + val
			cLen := dpLen[r] + 1
			update := false
			if cand > newVal[nr]+eps {
				update = true
			} else if math.Abs(cand-newVal[nr]) <= eps && cLen > newLen[nr] {
				update = true
			}
			if update {
				newVal[nr] = cand
				newLen[nr] = cLen
				newStateID[nr] = len(nodes)
				nodes = append(nodes, node{prev: stateID[r], idx: idx})
			}
		}
		dpVal = newVal
		dpLen = newLen
		stateID = newStateID
	}

	if stateID[d] <= 0 {
		fmt.Println(-1)
		return
	}

	resultIdx := []int{}
	cur := stateID[d]
	for cur != 0 {
		nd := nodes[cur]
		resultIdx = append(resultIdx, nd.idx)
		cur = nd.prev
	}
	for i, j := 0, len(resultIdx)-1; i < j; i, j = i+1, j-1 {
		resultIdx[i], resultIdx[j] = resultIdx[j], resultIdx[i]
	}

	out := bufio.NewWriter(os.Stdout)
	fmt.Fprintln(out, len(resultIdx))
	for i, idx := range resultIdx {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, a[idx])
	}
	fmt.Fprintln(out)
	out.Flush()
}
