package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, k, A, B int64
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	fmt.Fscan(reader, &k, &A, &B)

	var cost int64
	if k == 1 {
		cost = A * (n - 1)
		fmt.Println(cost)
		return
	}

	for n > 1 {
		if n < k {
			cost += A * (n - 1)
			break
		}
		if n%k != 0 {
			mod := n % k
			cost += A * mod
			n -= mod
		} else {
			newN := n / k
			diff := n - newN
			subCost := diff * A
			if subCost < B {
				cost += subCost
			} else {
				cost += B
			}
			n = newN
		}
	}
	fmt.Println(cost)
}
