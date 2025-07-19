package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	sLine, err := reader.ReadString('\n')
	if err != nil {
		return
	}
	tLine, err := reader.ReadString('\n')
	if err != nil {
		return
	}
	s := strings.TrimSpace(sLine)
	t := strings.TrimSpace(tLine)
	if len(s) != len(t) {
		fmt.Println(-1)
		return
	}
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	const inf = 1000000000
	var d [26][26]int
	for i := 0; i < 26; i++ {
		for j := 0; j < 26; j++ {
			if i == j {
				d[i][j] = 0
			} else {
				d[i][j] = inf
			}
		}
	}
	for i := 0; i < n; i++ {
		var uStr, vStr string
		var l int
		if _, err := fmt.Fscan(reader, &uStr, &vStr, &l); err != nil {
			return
		}
		u := int(uStr[0] - 'a')
		v := int(vStr[0] - 'a')
		if l < d[u][v] {
			d[u][v] = l
		}
	}
	for k := 0; k < 26; k++ {
		for i := 0; i < 26; i++ {
			for j := 0; j < 26; j++ {
				if d[i][k] < inf && d[k][j] < inf {
					d[i][j] = min(d[i][j], d[i][k]+d[k][j])
				}
			}
		}
	}
	m := len(s)
	result := make([]byte, m)
	var ans int64
	for i := 0; i < m; i++ {
		si := int(s[i] - 'a')
		ti := int(t[i] - 'a')
		if d[si][ti] > d[ti][si] {
			si, ti = ti, si
		}
		best := d[si][ti]
		for j := 0; j < 26; j++ {
			if d[si][j] < inf && d[ti][j] < inf {
				if d[si][j]+d[ti][j] < best {
					best = d[si][j] + d[ti][j]
				}
			}
		}
		if best >= inf {
			fmt.Println(-1)
			return
		}
		var c int
		if best == d[si][ti] {
			c = ti
		} else {
			for j := 0; j < 26; j++ {
				if d[si][j] < inf && d[ti][j] < inf && d[si][j]+d[ti][j] == best {
					c = j
					break
				}
			}
		}
		result[i] = byte('a' + c)
		ans += int64(best)
	}
	fmt.Println(ans)
	fmt.Println(string(result))
}
