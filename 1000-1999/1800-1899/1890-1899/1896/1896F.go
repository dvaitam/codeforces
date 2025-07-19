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

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for t > 0 {
		t--
		var n int
		fmt.Fscan(reader, &n)
		n *= 2
		var s string
		fmt.Fscan(reader, &s)
		// count ones
		cnt := 0
		for i := 0; i < n; i++ {
			if s[i] == '1' {
				cnt++
			}
		}
		if s[0] != s[n-1] || cnt%2 != 0 {
			fmt.Fprintln(writer, -1)
			continue
		}
		// compute flips
		flips := make([][]int, 2)
		arr := []byte(s)
		for i := 1; i < n-2; i++ {
			if arr[i] == '1' {
				flips[i%2] = append(flips[i%2], i)
				arr[i] = '0'
				if arr[i+1] == '0' {
					arr[i+1] = '1'
				} else {
					arr[i+1] = '0'
				}
			}
		}
		// build base sequence of parentheses
		half := n / 2
		base := make([]rune, 0, n)
		base = append(base, '(')
		for i := 0; i < half-1; i++ {
			base = append(base, '(', ')')
		}
		base = append(base, ')')
		// collect answers
		ans := make([]string, 0, 3)
		if s[0] == '1' {
			ans = append(ans, string(base))
		}
		for p := 0; p < 2; p++ {
			tmp := make([]rune, len(base))
			copy(tmp, base)
			for _, pos := range flips[p] {
				// swap positions pos and pos+1
				tmp[pos], tmp[pos+1] = tmp[pos+1], tmp[pos]
			}
			ans = append(ans, string(tmp))
		}
		// output
		fmt.Fprintln(writer, len(ans))
		for _, seq := range ans {
			fmt.Fprintln(writer, seq)
		}
	}
}
