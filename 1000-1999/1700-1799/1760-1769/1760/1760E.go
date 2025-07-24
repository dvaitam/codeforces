package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
		}

		prefixOnes := make([]int, n+1)
		for i := 1; i <= n; i++ {
			prefixOnes[i] = prefixOnes[i-1]
			if arr[i-1] == 1 {
				prefixOnes[i]++
			}
		}

		suffixZeros := make([]int, n+2)
		for i := n; i >= 1; i-- {
			suffixZeros[i] = suffixZeros[i+1]
			if arr[i-1] == 0 {
				suffixZeros[i]++
			}
		}

		var inv int64
		for i := 1; i <= n; i++ {
			if arr[i-1] == 0 {
				inv += int64(prefixOnes[i-1])
			}
		}

		best := inv
		for i := 1; i <= n; i++ {
			if arr[i-1] == 0 {
				newInv := inv - int64(prefixOnes[i-1]) + int64(suffixZeros[i+1])
				if newInv > best {
					best = newInv
				}
			} else {
				newInv := inv + int64(prefixOnes[i-1]) - int64(suffixZeros[i+1])
				if newInv > best {
					best = newInv
				}
			}
		}

		fmt.Fprintln(writer, best)
	}
}
