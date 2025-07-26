package main

import (
	"bufio"
	"fmt"
	"os"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		var a, b, c int64
		fmt.Fscan(in, &n, &a, &b, &c)
		var s string
		fmt.Fscan(in, &s)
		cnt00 := 0
		cnt11 := 0
		for i := 0; i+1 < n; i++ {
			if s[i] == '0' && s[i+1] == '0' {
				cnt00++
			}
			if s[i] == '1' && s[i+1] == '1' {
				cnt11++
			}
		}
		singles := 0
		i := 0
		for i < n {
			if s[i] == '0' {
				j := i
				for j < n && s[j] == '0' {
					j++
				}
				if j-i == 1 && i > 0 && j < n {
					singles++
				}
				i = j
			} else {
				i++
			}
		}
		ans := int64(min(cnt00, cnt11)) * (a + b)
		if cnt00 > cnt11 {
			ans += a
		} else if cnt11 > cnt00 {
			ans += b
		}
		if b > c {
			ans += (b - c) * int64(singles)
		}
		fmt.Fprintln(out, ans)
	}
}
