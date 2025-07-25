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
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var s string
		fmt.Fscan(in, &s)
		even := make([]byte, 0, len(s))
		odd := make([]byte, 0, len(s))
		for i := 0; i < len(s); i++ {
			if (s[i]-'0')%2 == 0 {
				even = append(even, s[i])
			} else {
				odd = append(odd, s[i])
			}
		}
		res := make([]byte, 0, len(s))
		i, j := 0, 0
		for i < len(even) && j < len(odd) {
			if even[i] < odd[j] {
				res = append(res, even[i])
				i++
			} else {
				res = append(res, odd[j])
				j++
			}
		}
		if i < len(even) {
			res = append(res, even[i:]...)
		}
		if j < len(odd) {
			res = append(res, odd[j:]...)
		}
		fmt.Fprintln(out, string(res))
	}
}
