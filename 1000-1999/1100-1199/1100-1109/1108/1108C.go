package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	var s string
	if _, err := fmt.Fscan(reader, &s); err != nil {
		return
	}
	perms := []string{"BGR", "BRG", "GBR", "GRB", "RBG", "RGB"}
	best := 0
	minChanges := n + 1
	for i, p := range perms {
		cnt := 0
		for j := 0; j < n; j++ {
			if s[j] != p[j%3] {
				cnt++
			}
		}
		if cnt < minChanges {
			minChanges = cnt
			best = i
		}
	}
	result := make([]byte, n)
	p := perms[best]
	for i := 0; i < n; i++ {
		result[i] = p[i%3]
	}
	// output
	fmt.Println(minChanges)
	fmt.Println(string(result))
}
