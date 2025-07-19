package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func solve(in *bufio.Reader, out *bufio.Writer) {
	var n int64
	var x int64
	fmt.Fscan(in, &n, &x)
	n--
	ans := int64(0)
	temp := int64(0)
	for ; n > 0; n-- {
		var y int64
		fmt.Fscan(in, &y)
		if y == 1 && x != 1 {
			ans = -1
		}
		if ans == -1 {
			x = y
			continue
		}
		val := math.Log(float64(x)) / math.Log(float64(y))
		val = math.Log2(val)
		inc := int64(math.Ceil(val))
		if inc+temp < 0 {
			temp = 0
		} else {
			temp += inc
			if temp < 0 {
				temp = 0
			}
		}
		ans += temp
		x = y
	}
	fmt.Fprintln(out, ans)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		solve(in, out)
	}
}
