package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, p int
	var s string
	fmt.Fscan(reader, &n, &p)
	fmt.Fscan(reader, &s)
	a := []byte(s)
	for i := n - 1; i >= 0; i-- {
		for c := a[i] + 1; c < byte('a'+p); c++ {
			if i >= 1 && a[i-1] == c {
				continue
			}
			if i >= 2 && a[i-2] == c {
				continue
			}
			a[i] = c
			ok := true
			for j := i + 1; j < n; j++ {
				found := false
				for d := byte('a'); d < byte('a'+p); d++ {
					if j >= 1 && a[j-1] == d {
						continue
					}
					if j >= 2 && a[j-2] == d {
						continue
					}
					a[j] = d
					found = true
					break
				}
				if !found {
					ok = false
					break
				}
			}
			if ok {
				fmt.Println(string(a))
				return
			}
		}
	}
	fmt.Println("NO")
}
