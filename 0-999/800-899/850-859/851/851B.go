package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var ax, ay, bx, by, cx, cy int64
	if _, err := fmt.Fscan(in, &ax, &ay, &bx, &by, &cx, &cy); err != nil {
		return
	}

	abx := bx - ax
	aby := by - ay
	bcx := cx - bx
	bcy := cy - by

	// squared lengths of AB and BC
	d1 := abx*abx + aby*aby
	d2 := bcx*bcx + bcy*bcy

	// cross product to check collinearity
	cross := abx*bcy - aby*bcx

	if d1 == d2 && cross != 0 {
		fmt.Fprintln(out, "Yes")
	} else {
		fmt.Fprintln(out, "No")
	}
}
