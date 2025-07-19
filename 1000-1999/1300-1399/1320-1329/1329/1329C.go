package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	A     []int
	G     int
	bound int
	res   []int
)

// find a descendant of n where both children are zero, following the larger child
func chkDel(n int) int {
	for {
		n1 := n * 2
		n2 := n*2 + 1
		if A[n1] == 0 && A[n2] == 0 {
			break
		}
		if A[n1] > A[n2] {
			n = n1
		} else {
			n = n2
		}
	}
	return n
}

// delete at n: promote children until leaf and set to zero
func doDel(n int) {
	for {
		n1 := n * 2
		n2 := n*2 + 1
		if A[n1] == 0 && A[n2] == 0 {
			A[n] = 0
			break
		}
		t := n1
		if A[n2] > A[n1] {
			t = n2
		}
		A[n] = A[t]
		n = t
	}
}

// recursively perform deletions to ensure all nodes in [bound, ...) removed
func doit(n int) {
	if n >= bound {
		return
	}
	for {
		k := chkDel(n)
		if k < bound {
			break
		}
		doDel(n)
		res = append(res, n)
	}
	doit(n * 2)
	doit(n*2 + 1)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)
	for T > 0 {
		T--
		var H int
		fmt.Fscan(reader, &H, &G)
		bound = 1 << G
		// allocate tree array of size enough for H levels
		size := 1 << (H + 1)
		A = make([]int, size)
		// read values for nodes 1 to 2^H - 1
		for i := 1; i < (1 << H); i++ {
			fmt.Fscan(reader, &A[i])
		}
		// prepare result
		res = res[:0]
		// perform deletions
		doit(1)
		// compute sum of remaining nodes up to level G-1
		var sum int64
		for i := 1; i < bound; i++ {
			sum += int64(A[i])
		}
		// output
		fmt.Fprintln(writer, sum)
		if len(res) > 0 {
			for i, v := range res {
				if i > 0 {
					writer.WriteByte(' ')
				}
				fmt.Fprint(writer, v)
			}
		}
		writer.WriteByte('\n')
	}
}
