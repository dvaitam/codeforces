package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	sum := 0
	const inf = int(1e9)
	minPosOdd := inf
	maxNegOdd := -inf
	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(reader, &x)
		if x > 0 {
			sum += x
		}
		if x%2 != 0 {
			if x > 0 {
				if x < minPosOdd {
					minPosOdd = x
				}
			} else {
				if x > maxNegOdd {
					maxNegOdd = x
				}
			}
		}
	}
	if sum%2 == 1 {
		fmt.Println(sum)
		return
	}
	best := -inf
	if minPosOdd != inf {
		if s := sum - minPosOdd; s%2 != 0 && s > best {
			best = s
		}
	}
	if maxNegOdd != -inf {
		if s := sum + maxNegOdd; s%2 != 0 && s > best {
			best = s
		}
	}
	fmt.Println(best)
}
