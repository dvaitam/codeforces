package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	freq := make(map[int]int)
	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(in, &x)
		freq[x]++
	}
	for _, c := range freq {
		if c%2 == 1 {
			fmt.Println("Conan")
			return
		}
	}
	fmt.Println("Agasa")
}
