package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var s string
	if _, err := fmt.Fscan(in, &s); err != nil {
		return
	}
	cnt := make([]int, 26)
	for _, ch := range s {
		cnt[ch-'a']++
	}
	i, j := 0, 25
	for i < j {
		for i < j && cnt[i]%2 == 0 {
			i++
		}
		for i < j && cnt[j]%2 == 0 {
			j--
		}
		if i >= j {
			break
		}
		cnt[i]++
		cnt[j]--
		i++
		j--
	}
	first := make([]byte, 0, len(s)/2)
	var mid byte
	for k := 0; k < 26; k++ {
		for cnt[k] >= 2 {
			first = append(first, byte('a'+k))
			cnt[k] -= 2
		}
		if cnt[k] == 1 {
			mid = byte('a' + k)
		}
	}
	res := make([]byte, 0, len(s))
	res = append(res, first...)
	if mid != 0 {
		res = append(res, mid)
	}
	for i := len(first) - 1; i >= 0; i-- {
		res = append(res, first[i])
	}
	fmt.Println(string(res))
}
