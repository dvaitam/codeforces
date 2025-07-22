package main

import (
	"bufio"
	"fmt"
	"os"
)

type group struct {
	ch  byte
	len int
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	s, _ := reader.ReadString('\n')
	if len(s) > 0 && s[len(s)-1] == '\n' {
		s = s[:len(s)-1]
	}
	if len(s) == 0 {
		fmt.Println(0)
		return
	}
	// Build list of consecutive character groups
	groups := make([]group, 0, len(s))
	cur := group{ch: s[0], len: 1}
	for i := 1; i < len(s); i++ {
		if s[i] == cur.ch {
			cur.len++
		} else {
			groups = append(groups, cur)
			cur = group{ch: s[i], len: 1}
		}
	}
	groups = append(groups, cur)

	ops := 0
	for len(groups) > 1 {
		newGroups := make([]group, 0, len(groups))
		m := len(groups)
		for i, g := range groups {
			dec := 1
			if i > 0 && i < m-1 {
				dec = 2
			}
			g.len -= dec
			if g.len > 0 {
				if len(newGroups) > 0 && newGroups[len(newGroups)-1].ch == g.ch {
					newGroups[len(newGroups)-1].len += g.len
				} else {
					newGroups = append(newGroups, g)
				}
			}
		}
		groups = newGroups
		ops++
	}
	fmt.Println(ops)
}
