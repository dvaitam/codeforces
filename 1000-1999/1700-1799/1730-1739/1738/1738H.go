package main

import (
	"bufio"
	"fmt"
	"os"
)

func isPalindrome(b []byte) bool {
	i, j := 0, len(b)-1
	for i < j {
		if b[i] != b[j] {
			return false
		}
		i++
		j--
	}
	return true
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var q int
	fmt.Fscan(in, &q)

	queue := make([]byte, 0)
	for ; q > 0; q-- {
		var op string
		fmt.Fscan(in, &op)
		if op == "push" {
			var c string
			fmt.Fscan(in, &c)
			queue = append(queue, c[0])
		} else if op == "pop" {
			if len(queue) > 0 {
				queue = queue[1:]
			}
		}
		seen := make(map[string]struct{})
		for i := 0; i < len(queue); i++ {
			for j := i; j < len(queue); j++ {
				if isPalindrome(queue[i : j+1]) {
					seen[string(queue[i:j+1])] = struct{}{}
				}
			}
		}
		fmt.Fprintln(out, len(seen))
	}
}
