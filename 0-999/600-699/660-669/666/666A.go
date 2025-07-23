package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var s string
	if _, err := fmt.Fscan(reader, &s); err != nil {
		return
	}
	n := len(s)
	dp2 := make([]bool, n+1)
	dp3 := make([]bool, n+1)

	for i := n - 2; i >= 5; i-- {
		if i+2 <= n {
			if i+2 == n {
				dp2[i] = true
			} else {
				if dp2[i+2] && s[i:i+2] != s[i+2:i+4] {
					dp2[i] = true
				}
				if !dp2[i] && dp3[i+2] && s[i:i+2] != s[i+2:i+5] {
					dp2[i] = true
				}
			}
		}
		if i+3 <= n {
			if i+3 == n {
				dp3[i] = true
			} else {
				if dp2[i+3] && s[i:i+3] != s[i+3:i+5] {
					dp3[i] = true
				}
				if !dp3[i] && dp3[i+3] && s[i:i+3] != s[i+3:i+6] {
					dp3[i] = true
				}
			}
		}
	}

	set := make(map[string]struct{})
	for i := 5; i < n; i++ {
		if dp2[i] {
			set[s[i:i+2]] = struct{}{}
		}
		if dp3[i] {
			set[s[i:i+3]] = struct{}{}
		}
	}

	ans := make([]string, 0, len(set))
	for k := range set {
		ans = append(ans, k)
	}
	sort.Strings(ans)
	fmt.Fprintln(writer, len(ans))
	for _, v := range ans {
		fmt.Fprintln(writer, v)
	}
}
