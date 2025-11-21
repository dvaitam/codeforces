package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var g, c, l int
	if _, err := fmt.Fscan(in, &g, &c, &l); err != nil {
		return
	}
	scores := []int{g, c, l}
	sort.Ints(scores)
	if scores[2]-scores[0] >= 10 {
		fmt.Println("check again")
		return
	}
	fmt.Printf("final %d\n", scores[1])
}
