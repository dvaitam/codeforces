package main

import (
	"bufio"
	"fmt"
	"os"
)

func f(pIndex int, s []rune, m [][]int) bool {
	n := len(s)
	for i := 0; i < n; i++ {
		if int(s[i]-'a') == pIndex {
			j := i
			for j < n && s[j] == s[i] {
				j++
			}
			if i == 0 || j == n {
				i = j - 1
				continue
			}
			q1 := int(s[i-1] - 'a')
			q2 := int(s[j] - 'a')
			if m[q1][q2] != 1 {
				return false
			}
			i = j - 1
		} else {
			if i != 0 && int(s[i-1]-'a') != pIndex {
				q1 := int(s[i] - 'a')
				q2 := int(s[i-1] - 'a')
				if m[q1][q2] != 1 {
					return false
				}
			}
			if i+1 != n && int(s[i+1]-'a') != pIndex {
				q1 := int(s[i] - 'a')
				q2 := int(s[i+1] - 'a')
				if m[q1][q2] != 1 {
					return false
				}
			}
		}
	}
	return true
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, p int
	fmt.Fscan(reader, &n, &p)
	var str string
	fmt.Fscan(reader, &str)
	m := make([][]int, p)
	for i := 0; i < p; i++ {
		m[i] = make([]int, p)
		for j := 0; j < p; j++ {
			fmt.Fscan(reader, &m[i][j])
		}
	}
	tmp := []rune(str)
	used := make([]bool, p)
	for i := 0; i < p; i++ {
		for j := 0; j < p; j++ {
			if !used[j] && f(j, tmp, m) {
				used[j] = true
				// remove all occurrences of letter j
				newTmp := tmp[:0]
				for _, ch := range tmp {
					if int(ch-'a') != j {
						newTmp = append(newTmp, ch)
					}
				}
				tmp = newTmp
			}
		}
	}
	fmt.Fprint(writer, len(tmp))
}
