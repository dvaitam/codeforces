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
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(reader, &n, &m)
		tot := n + m + 1
		a := make([]int, tot)
		b := make([]int, tot)
		for i := 0; i < tot; i++ {
			fmt.Fscan(reader, &a[i])
		}
		for i := 0; i < tot; i++ {
			fmt.Fscan(reader, &b[i])
		}

		role := make([]bool, tot) // true programmer
		pref := make([]bool, tot) // true prefers programmer
		forced := make([]bool, tot)
		skill := make([]int, tot)

		P := n
		T := m
		var totalSkill int64
		for i := 0; i < tot; i++ {
			if a[i] > b[i] {
				pref[i] = true
			}
			if pref[i] {
				if P > 0 {
					role[i] = true
					skill[i] = a[i]
					P--
				} else {
					role[i] = false
					skill[i] = b[i]
					forced[i] = true
					T--
				}
			} else {
				if T > 0 {
					role[i] = false
					skill[i] = b[i]
					T--
				} else {
					role[i] = true
					skill[i] = a[i]
					forced[i] = true
					P--
				}
			}
			totalSkill += int64(skill[i])
		}

		nForcedProg := make([]int, tot)
		nForcedTest := make([]int, tot)
		nextP := -1
		nextT := -1
		for i := tot - 1; i >= 0; i-- {
			nForcedProg[i] = nextP
			nForcedTest[i] = nextT
			if pref[i] && forced[i] && !role[i] {
				nextP = i
			}
			if !pref[i] && forced[i] && role[i] {
				nextT = i
			}
		}

		for i := 0; i < tot; i++ {
			ans := totalSkill - int64(skill[i])
			if role[i] {
				j := nForcedProg[i]
				if j != -1 {
					ans += int64(a[j] - b[j])
				}
			} else {
				j := nForcedTest[i]
				if j != -1 {
					ans += int64(b[j] - a[j])
				}
			}
			if i > 0 {
				writer.WriteByte(' ')
			}
			fmt.Fprint(writer, ans)
		}
		writer.WriteByte('\n')
	}
}
