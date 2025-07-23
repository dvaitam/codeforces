package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var k int
	if _, err := fmt.Fscan(reader, &k); err != nil {
		return
	}
	if k%2 == 0 {
		fmt.Println("NO")
		return
	}

	n := 4*k - 2
	m := 2*((k-1)+(k-1)*(k-1)+((k-1)/2)) + 1

	fmt.Println("YES")
	fmt.Printf("%d %d\n", n, m)

	goFunc := func(start, limit, n int) {
		for i := start; i < limit; i += 2 {
			fmt.Printf("%d %d\n", i, i+1)
			for j := 0; j < n-1; j++ {
				fmt.Printf("%d %d\n", i, limit+j)
				fmt.Printf("%d %d\n", i+1, limit+j)
			}
		}
		for i := limit; i < limit+n-1; i++ {
			fmt.Printf("%d %d\n", i, limit+n-1)
		}
	}

	goFunc(1, k, k)
	goFunc(2*k, 3*k-1, k)
	fmt.Printf("%d %d\n", 2*k-1, 4*k-2)
}
