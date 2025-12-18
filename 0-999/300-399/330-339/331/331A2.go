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
	fmt.Fscan(in, &n)
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	maxSum := int64(-1e18) // Initialize with a very small number
	var bestL, bestR int = -1, -1

	// Iterate over all pairs (l, r)
	for l := 0; l < n; l++ {
		for r := l + 1; r < n; r++ {
			if a[l] == a[r] {
				currentSum := a[l] + a[r]
				// Add positive numbers between l and r
				for k := l + 1; k < r; k++ {
					if a[k] > 0 {
						currentSum += a[k]
					}
				}
				if currentSum > maxSum {
					maxSum = currentSum
					bestL, bestR = l, r
				}
			}
		}
	}

    if bestL == -1 {
        // This case should theoretically not happen based on problem constraints (n >= 2)
        // unless no two elements are equal.
        // But the problem implies a solution always exists?
        // "there must be at least two trees... esthetic appeal ... must be the same"
        // If no such pair exists, the problem is unsolvable.
        // But for the purpose of this reference on random inputs, we might need to handle it.
        // However, the random generator in verifier generates:
        // arr[0], arr[1] = common, common
        // So a solution always exists.
        return 
    }

	remove := []int{}
	// Remove everything before bestL
	for i := 0; i < bestL; i++ {
		remove = append(remove, i+1)
	}
	// Remove negatives between bestL and bestR
	for i := bestL + 1; i < bestR; i++ {
		if a[i] < 0 {
			remove = append(remove, i+1)
		}
	}
	// Remove everything after bestR
	for i := bestR + 1; i < n; i++ {
		remove = append(remove, i+1)
	}

	fmt.Fprintf(out, "%d %d\n", maxSum, len(remove))
	for i, x := range remove {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, x)
	}
	if len(remove) > 0 {
		fmt.Fprintln(out)
	}
}