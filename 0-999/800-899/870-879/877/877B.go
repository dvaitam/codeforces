package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var s string
	if _, err := fmt.Fscan(reader, &s); err != nil {
		return
	}

	dp1, dp2, dp3 := 0, 0, 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == 'a' {
			ndp3 := dp3
			if dp3+1 > ndp3 {
				ndp3 = dp3 + 1
			}
			if dp2+1 > ndp3 {
				ndp3 = dp2 + 1
			}
			ndp2 := dp2
			if dp1 > ndp2 {
				ndp2 = dp1
			}
			dp1 = dp1 + 1
			dp2 = ndp2
			dp3 = ndp3
		} else { // 'b'
			ndp3 := dp3
			ndp2 := dp2
			if dp2+1 > ndp2 {
				ndp2 = dp2 + 1
			}
			if dp1+1 > ndp2 {
				ndp2 = dp1 + 1
			}
			dp2 = ndp2
			dp3 = ndp3
		}
	}

	ans := dp1
	if dp2 > ans {
		ans = dp2
	}
	if dp3 > ans {
		ans = dp3
	}
	fmt.Fprintln(writer, ans)
}
