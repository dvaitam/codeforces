package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	rdr := bufio.NewReader(os.Stdin)
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	var tc int
	if _, err := fmt.Fscan(rdr, &tc); err != nil {
		return
	}
	for tc > 0 {
		tc--
		var L, v, l, r int64
		fmt.Fscan(rdr, &L, &v)
		fmt.Fscan(rdr, &l, &r)
		total := L / v
		inRange := r/v - ((l - 1) / v)
		fmt.Fprintln(w, total-inRange)
	}
}
