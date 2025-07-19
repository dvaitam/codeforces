package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, sx, sy int
	if _, err := fmt.Fscan(reader, &n, &sx, &sy); err != nil {
		return
	}
	x1, x2, y1, y2 := 0, 0, 0, 0
	for i := 0; i < n; i++ {
		var x, y int
		fmt.Fscan(reader, &x, &y)
		if x < sx {
			x1++
		}
		if x > sx {
			x2++
		}
		if y < sy {
			y1++
		}
		if y > sy {
			y2++
		}
	}
	// Determine best direction
	ans := x1
	ansx, ansy := sx, sy
	if x2 > ans {
		ans = x2
	}
	if y1 > ans {
		ans = y1
	}
	if y2 > ans {
		ans = y2
	}
	if ans == x1 {
		ansx = sx - 1
	} else if ans == x2 {
		ansx = sx + 1
	} else if ans == y1 {
		ansy = sy - 1
	} else if ans == y2 {
		ansy = sy + 1
	}
	fmt.Println(ans)
	fmt.Println(ansx, ansy)
}
