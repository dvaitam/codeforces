package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	scores := make([]int, 3)
	for i := 0; i < 3; i++ {
		if _, err := fmt.Fscan(in, &scores[i]); err != nil {
			return
		}
	}

	sorted := append([]int(nil), scores...)
	sort.Ints(sorted)

	if sorted[2]-sorted[0] >= 10 {
		fmt.Println("check again")
		return
	}

	fmt.Printf("final %d\n", sorted[1])
}

