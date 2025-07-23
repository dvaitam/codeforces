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
	temps := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &temps[i])
	}

	if n <= 1 {
		// Though constraints guarantee n>=2, handle generically
		fmt.Println(temps[n-1])
		return
	}
	diff := temps[1] - temps[0]
	isAP := true
	for i := 2; i < n; i++ {
		if temps[i]-temps[i-1] != diff {
			isAP = false
			break
		}
	}
	if isAP {
		fmt.Println(temps[n-1] + diff)
	} else {
		fmt.Println(temps[n-1])
	}
}
