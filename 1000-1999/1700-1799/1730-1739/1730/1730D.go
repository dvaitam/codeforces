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
	for ; t > 0; t-- {
		if solve(reader, writer) != nil {
			return
		}
	}
}

func solve(reader *bufio.Reader, writer *bufio.Writer) error {
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return err
	}
	var s1, s2 string
	fmt.Fscan(reader, &s1)
	fmt.Fscan(reader, &s2)

	counts := make(map[[2]byte]int)
	b1 := []byte(s1)
	b2 := []byte(s2)
	for i := 0; i < n; i++ {
		a := b1[i]
		b := b2[n-1-i]
		if a > b {
			a, b = b, a
		}
		counts[[2]byte{a, b}]++
	}

	oddSame := 0
	ok := true
	for pair, c := range counts {
		if pair[0] != pair[1] {
			if c%2 == 1 {
				ok = false
				break
			}
		} else {
			if c%2 == 1 {
				oddSame++
			}
		}
	}
	if ok {
		if n%2 == 0 {
			if oddSame > 0 {
				ok = false
			}
		} else {
			if oddSame != 1 {
				ok = false
			}
		}
	}

	if ok {
		fmt.Fprintln(writer, "YES")
	} else {
		fmt.Fprintln(writer, "NO")
	}
	return nil
}
