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

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(in, &n, &k)
		var s string
		fmt.Fscan(in, &s)

		cnt0, cnt1, cnt2 := 0, 0, 0
		for _, ch := range s {
			switch ch {
			case '0':
				cnt0++
			case '1':
				cnt1++
			default:
				cnt2++
			}
		}

		res := make([]byte, n)
		topForced := cnt0
		topMax := cnt0 + cnt2
		bottomMax := cnt1 + cnt2

		for i := 1; i <= n; i++ {
			alwaysStay := i > topMax && i <= n-bottomMax
			if alwaysStay {
				res[i-1] = '+'
				continue
			}

			left := i - (n - bottomMax)
			if left < 0 {
				left = 0
			}
			right := i - 1 - topForced
			if right > cnt2 {
				right = cnt2
			}
			if right >= left {
				res[i-1] = '?'
			} else {
				res[i-1] = '-'
			}
		}

		fmt.Fprintln(out, string(res))
	}
}
