package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		var n, p1, p2, p3, t1, t2 int
		if _, err := fmt.Fscan(reader, &n, &p1, &p2, &p3, &t1, &t2); err != nil {
			break
		}
		var lPrev, rPrev int
		var pc int64
		// first lecture
		fmt.Fscan(reader, &lPrev, &rPrev)
		pc += int64(rPrev-lPrev) * int64(p1)
		// subsequent lectures
		for i := 1; i < n; i++ {
			var l, r int
			fmt.Fscan(reader, &l, &r)
			// lecture duration cost
			pc += int64(r-l) * int64(p1)
			// break duration
			x := l - rPrev
			if x <= t1 {
				pc += int64(x) * int64(p1)
			} else {
				pc += int64(t1) * int64(p1)
				x -= t1
				if x <= t2 {
					pc += int64(x) * int64(p2)
				} else {
					pc += int64(t2) * int64(p2)
					x -= t2
					pc += int64(x) * int64(p3)
				}
			}
			rPrev = r
		}
		fmt.Println(pc)
	}
}
