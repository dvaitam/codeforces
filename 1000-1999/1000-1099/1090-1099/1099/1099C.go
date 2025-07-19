package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var s string
	var k int
	fmt.Fscan(reader, &s)
	fmt.Fscan(reader, &k)
	n := len(s)
	c1, c2 := 0, 0
	for i := 0; i < n; i++ {
		if s[i] == '?' {
			c1++
		} else if s[i] == '*' {
			c2++
		}
	}
	// minimal length by removing each wildcard and preceding char
	ss := n - 2*(c1+c2)
	if k < ss {
		fmt.Println("Impossible")
		return
	}
	// maximal length without expansion (remove wildcards only)
	ss += c1 + c2
	if c2 == 0 && k > ss {
		fmt.Println("Impossible")
		return
	}
	// build answer
	if k <= ss {
		dif := ss - k
		var ans []byte
		for i := 0; i < n-1; i++ {
			ch := s[i]
			next := s[i+1]
			if ch == '?' || ch == '*' {
				continue
			}
			if next != '?' && next != '*' {
				ans = append(ans, ch)
			}
			if next == '?' || next == '*' {
				if dif > 0 {
					dif--
					continue
				}
				ans = append(ans, ch)
			}
		}
		last := s[n-1]
		if last != '?' && last != '*' {
			ans = append(ans, last)
		}
		fmt.Println(string(ans))
		return
	}
	// k > ss: need to expand using one '*'
	dif := k - ss
	var ans []byte
	for i := 0; i < n-1; i++ {
		ch := s[i]
		if ch == '?' {
			continue
		}
		if ch == '*' && c2 > 1 {
			c2--
			continue
		}
		if ch == '*' && c2 == 1 {
			c2--
			prev := s[i-1]
			for j := 0; j < dif; j++ {
				ans = append(ans, prev)
			}
			continue
		}
		ans = append(ans, ch)
	}
	last := s[n-1]
	if last != '?' && last != '*' {
		ans = append(ans, last)
	}
	if last == '*' && c2 == 1 {
		// expand at end
		c2--
		prev := s[n-2]
		for j := 0; j < dif; j++ {
			ans = append(ans, prev)
		}
	}
	fmt.Println(string(ans))
}
