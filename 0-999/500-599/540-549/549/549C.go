package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)

	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}

	odd := 0
	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(in, &x)
		if x%2 == 1 {
			odd++
		}
	}
	even := n - odd
	moves := n - k

	var res string
	if moves == 0 {
		if odd%2 == 1 {
			res = "Stannis"
		} else {
			res = "Daenerys"
		}
	} else if moves%2 == 1 { // Stannis moves last
		d := moves / 2
		if odd <= d || (k%2 == 0 && even <= d) {
			res = "Daenerys"
		} else {
			res = "Stannis"
		}
	} else { // Daenerys moves last
		s := moves / 2
		if k%2 == 1 && even <= s {
			res = "Stannis"
		} else {
			res = "Daenerys"
		}
	}

	fmt.Println(res)
}
