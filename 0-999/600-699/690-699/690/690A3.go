package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var T int
	if _, err := fmt.Fscan(reader, &T); err != nil {
		return
	}
	for T > 0 {
		T--
		var n, r int64
		fmt.Fscan(reader, &n, &r)
		var s int64
		for i := int64(1); i < n; i++ {
			var x int64
			fmt.Fscan(reader, &x)
			s += x
		}
		smod := s % n
		rmod := r % n
		delta := (rmod - smod + n) % n
		var ans int64
		if delta == 0 {
			ans = n
		} else {
			ans = delta
		}
		fmt.Println(ans)
	}
}
