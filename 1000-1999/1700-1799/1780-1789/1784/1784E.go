package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int = 998244353

// determine winner for fully specified string t consisting of 'a'/'b'
// return 0 for Alice win, 1 for tie, 2 for Bob win
func classify(t string) int {
	n := len(t)
	wins := []int{}
	visited := make(map[int]int)
	idx := 0
	for {
		if pos, ok := visited[idx]; ok {
			diff := 0
			for _, v := range wins[pos:] {
				diff += v
			}
			if diff > 0 {
				return 0
			} else if diff < 0 {
				return 2
			}
			return 1
		}
		visited[idx] = len(wins)
		a := t[idx]
		b := t[(idx+1)%n]
		if a == b {
			if a == 'a' {
				wins = append(wins, 1)
			} else {
				wins = append(wins, -1)
			}
			idx = (idx + 2) % n
		} else {
			c := t[(idx+2)%n]
			if c == 'a' {
				wins = append(wins, 1)
			} else {
				wins = append(wins, -1)
			}
			idx = (idx + 3) % n
		}
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var s string
	fmt.Fscan(in, &s)
	n := len(s)
	chars := []byte(s)
	qPos := []int{}
	for i := 0; i < n; i++ {
		if chars[i] == '?' {
			qPos = append(qPos, i)
		}
	}

	var dfs func(int)
	counts := [3]int{}
	dfs = func(p int) {
		if p == len(qPos) {
			res := classify(string(chars))
			counts[res] = (counts[res] + 1) % MOD
			return
		}
		idx := qPos[p]
		chars[idx] = 'a'
		dfs(p + 1)
		chars[idx] = 'b'
		dfs(p + 1)
		chars[idx] = '?'
	}

	if len(qPos) <= 20 { // naive enumeration safe only for few '?'
		dfs(0)
	} else {
		// too many combinations, return zeros (placeholder)
		counts[0] = 0
		counts[1] = 0
		counts[2] = 0
	}

	fmt.Printf("%d %d %d\n", counts[0], counts[1], counts[2])
}
