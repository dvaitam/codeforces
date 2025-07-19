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
	a := make([]int, n)
	total := 0
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
		total += a[i]
	}
	leader := a[0]
	// if leader alone has more than half
	if leader*2 > total {
		fmt.Println(1)
		fmt.Println(1)
		return
	}
	// try to add supporters
	supporters := make([]int, 0, n)
	curr := leader
	for i := 1; i < n; i++ {
		if a[i]*2 <= leader {
			curr += a[i]
			supporters = append(supporters, i+1)
			if curr*2 > total {
				// print group size and indices
				fmt.Println(len(supporters) + 1)
				// print leader index
				fmt.Print(1)
				for _, idx := range supporters {
					fmt.Print(" ", idx)
				}
				fmt.Println()
				return
			}
		}
	}
	// no valid group
	fmt.Println(0)
}
