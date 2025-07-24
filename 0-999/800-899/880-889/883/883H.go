package main

import (
	"bufio"
	"fmt"
	"os"
)

func idx(b byte) int {
	if b >= 'a' && b <= 'z' {
		return int(b - 'a')
	}
	if b >= 'A' && b <= 'Z' {
		return 26 + int(b-'A')
	}
	return 52 + int(b-'0')
}

func ch(i int) byte {
	if i < 26 {
		return byte('a' + i)
	}
	if i < 52 {
		return byte('A' + i - 26)
	}
	return byte('0' + i - 52)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	var s string
	fmt.Fscan(in, &s)

	freq := make([]int, 62)
	for i := 0; i < n; i++ {
		freq[idx(s[i])]++
	}

	oddCnt := 0
	for _, v := range freq {
		if v%2 == 1 {
			oddCnt++
		}
	}

	k := 0
	length := 0
	for cand := 1; cand <= n; cand++ {
		if n%cand != 0 {
			continue
		}
		l := n / cand
		if l%2 == 0 {
			if oddCnt == 0 {
				k = cand
				length = l
				break
			}
		} else {
			if cand >= oddCnt {
				k = cand
				length = l
				break
			}
		}
	}

	halfLen := length / 2
	centers := make([]byte, 0, k)
	for i := 0; i < 62; i++ {
		if freq[i]%2 == 1 {
			centers = append(centers, ch(i))
			freq[i]--
		}
	}

	for len(centers) < k {
		for i := 0; i < 62 && len(centers) < k; i++ {
			if freq[i] >= 2 {
				freq[i] -= 2
				centers = append(centers, ch(i))
			}
		}
	}

	pairs := make([]byte, 0, (n-k)/2)
	for i := 0; i < 62; i++ {
		for freq[i] >= 2 {
			pairs = append(pairs, ch(i))
			freq[i] -= 2
		}
	}

	res := make([]string, k)
	for i := 0; i < k; i++ {
		half := pairs[i*halfLen : (i+1)*halfLen]
		pal := make([]byte, 0, length)
		pal = append(pal, half...)
		if length%2 == 1 {
			pal = append(pal, centers[i])
		}
		for j := len(half) - 1; j >= 0; j-- {
			pal = append(pal, half[j])
		}
		res[i] = string(pal)
	}

	fmt.Fprintln(out, k)
	for i := 0; i < k; i++ {
		if i > 0 {
			out.WriteByte(' ')
		}
		out.WriteString(res[i])
	}
	out.WriteByte('\n')
}
