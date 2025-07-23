package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var nA, nB int
	if _, err := fmt.Fscan(in, &nA, &nB); err != nil {
		return
	}
	var k, m int
	if _, err := fmt.Fscan(in, &k, &m); err != nil {
		return
	}
	A := make([]int, nA)
	for i := 0; i < nA; i++ {
		fmt.Fscan(in, &A[i])
	}
	B := make([]int, nB)
	for i := 0; i < nB; i++ {
		fmt.Fscan(in, &B[i])
	}
	if A[k-1] < B[nB-m] {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}
}
