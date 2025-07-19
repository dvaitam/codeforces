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

	var n, k int
	fmt.Fscan(reader, &n, &k)
	var s string
	fmt.Fscan(reader, &s)
	buf := []byte(s)

	// Compute maximum possible distance
	maxDist := 0
	for i := 0; i < n; i++ {
		c := buf[i]
		t1 := int('z' - c)
		t2 := int(c - 'a')
		if t1 > t2 {
			maxDist += t1
		} else {
			maxDist += t2
		}
	}
	if maxDist < k {
		fmt.Fprint(writer, "-1")
		return
	}

	// Build result
	for i := 0; i < n && k > 0; i++ {
		c := buf[i]
		t1 := int('z' - c)
		t2 := int(c - 'a')
		if t1 >= k {
			buf[i] = c + byte(k)
			k = 0
			break
		} else if t2 >= k {
			buf[i] = c - byte(k)
			k = 0
			break
		} else {
			if t1 > t2 {
				buf[i] = 'z'
				k -= t1
			} else {
				buf[i] = 'a'
				k -= t2
			}
		}
	}

	fmt.Fprint(writer, string(buf))
}
