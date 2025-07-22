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
	soft, hard := 0, 0
	for i := 0; i < n; i++ {
		var name, t string
		fmt.Fscan(in, &name, &t)
		if t == "soft" {
			soft++
		} else {
			hard++
		}
	}
	for size := 1; ; size++ {
		total := size * size
		w := (total + 1) / 2
		b := total / 2
		if (soft <= w && hard <= b) || (soft <= b && hard <= w) {
			fmt.Print(size)
			return
		}
	}
}
