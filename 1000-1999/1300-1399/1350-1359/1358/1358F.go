package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	A := make([]int64, n)
	B := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &A[i])
	}
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &B[i])
	}

	if n == 1 {
		if A[0] == B[0] {
			fmt.Fprintln(out, "SMALL")
			fmt.Fprintln(out, 0)
			fmt.Fprintln(out, "")
		} else {
			fmt.Fprintln(out, "IMPOSSIBLE")
		}
		return
	}

	operations := make([]byte, 0)
	pCount := 0

	copyB := make([]int64, n)
	copy(copyB, B)

	for {
		same := true
		for i := 0; i < n; i++ {
			if copyB[i] != A[i] {
				same = false
				break
			}
		}
		if same {
			break
		}
		revsame := true
		for i := 0; i < n; i++ {
			if copyB[i] != A[n-1-i] {
				revsame = false
				break
			}
		}
		if revsame {
			operations = append(operations, 'R')
			for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
				copyB[i], copyB[j] = copyB[j], copyB[i]
			}
			break
		}
		inc := true
		dec := true
		for i := 1; i < n; i++ {
			if copyB[i] <= copyB[i-1] {
				inc = false
			}
			if copyB[i] >= copyB[i-1] {
				dec = false
			}
		}
		if inc {
			for i := n - 1; i > 0; i-- {
				copyB[i] -= copyB[i-1]
				if copyB[i] <= 0 {
					fmt.Fprintln(out, "IMPOSSIBLE")
					return
				}
			}
			operations = append(operations, 'P')
			pCount++
		} else if dec {
			for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
				copyB[i], copyB[j] = copyB[j], copyB[i]
			}
			operations = append(operations, 'R')
		} else {
			fmt.Fprintln(out, "IMPOSSIBLE")
			return
		}
		if len(operations) > 500000 {
			break
		}
	}

	if pCount > 200000 {
		fmt.Fprintln(out, "BIG")
		fmt.Fprintln(out, pCount)
		return
	}
	fmt.Fprintln(out, "SMALL")
	fmt.Fprintln(out, len(operations))
	for i := len(operations) - 1; i >= 0; i-- {
		out.WriteByte(operations[i])
	}
	fmt.Fprintln(out)
}
