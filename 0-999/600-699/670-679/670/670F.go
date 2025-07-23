package main

import (
	"bufio"
	"fmt"
	"os"
)

func numDigits(x int) int {
	if x == 0 {
		return 1
	}
	c := 0
	for x > 0 {
		c++
		x /= 10
	}
	return c
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var s, t string
	if _, err := fmt.Fscan(in, &s); err != nil {
		return
	}
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	n := len(s)
	length := 1
	for d := 1; d <= 7; d++ {
		l := n - d
		if l >= 1 && numDigits(l) == d {
			length = l
			break
		}
	}
	cnt := make([]int, 10)
	for i := 0; i < len(s); i++ {
		cnt[s[i]-'0']++
	}
	strL := fmt.Sprintf("%d", length)
	for i := 0; i < len(strL); i++ {
		cnt[strL[i]-'0']--
	}
	cntT := make([]int, 10)
	for i := 0; i < len(t); i++ {
		cntT[t[i]-'0']++
	}
	rem := make([]int, 10)
	for i := 0; i < 10; i++ {
		rem[i] = cnt[i] - cntT[i]
	}
	build := func(arr []int) string {
		first := int(t[0] - '0')
		var less, equal, greater []byte
		for d := 0; d < 10; d++ {
			for i := 0; i < arr[d]; i++ {
				if d < first {
					less = append(less, byte('0'+d))
				} else if d == first {
					equal = append(equal, byte('0'+d))
				} else {
					greater = append(greater, byte('0'+d))
				}
			}
		}
		cand1 := string(append(append(append([]byte{}, less...), []byte(t)...), append(equal, greater...)...))
		cand2 := string(append(append(append([]byte{}, less...), equal...), append([]byte(t), greater...)...))
		if cand1 < cand2 {
			return cand1
		}
		return cand2
	}

	ans := ""
	have := false
	if t[0] != '0' || length == 1 {
		var rest []byte
		for d := 0; d < 10; d++ {
			for i := 0; i < rem[d]; i++ {
				rest = append(rest, byte('0'+d))
			}
		}
		cand := t + string(rest)
		ans = cand
		have = true
	}
	for d := 1; d <= 9; d++ {
		if rem[d] > 0 {
			rem[d]--
			cand := string('0'+byte(d)) + build(rem)
			rem[d]++
			if !have || cand < ans {
				ans = cand
				have = true
			}
			break
		}
	}
	if !have {
		ans = t
	}
	fmt.Println(ans)
}
