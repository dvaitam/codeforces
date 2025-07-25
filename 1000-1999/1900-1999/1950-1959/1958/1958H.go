package main

import (
	"bufio"
	"fmt"
	"os"
)

const INF int64 = 1e18

type Spell struct {
	basic bool
	val   int64
	seq   []int
}

func clamp(x int64) int64 {
	if x > INF {
		return INF
	}
	if x < -INF {
		return -INF
	}
	return x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		var hp int64
		fmt.Fscan(in, &n, &hp)
		spells := make([]Spell, 0)
		spells = append(spells, Spell{})
		for i := 1; i <= n; i++ {
			var v int64
			fmt.Fscan(in, &v)
			spells = append(spells, Spell{basic: true, val: v})
		}
		var m int
		fmt.Fscan(in, &m)
		for i := 0; i < m; i++ {
			var s int
			fmt.Fscan(in, &s)
			seq := make([]int, s)
			for j := 0; j < s; j++ {
				fmt.Fscan(in, &seq[j])
			}
			spells = append(spells, Spell{seq: seq})
		}

		total := make([]int64, len(spells))
		minPref := make([]int64, len(spells))

		for i := 1; i < len(spells); i++ {
			sp := spells[i]
			if sp.basic {
				total[i] = sp.val
				if sp.val < 0 {
					minPref[i] = sp.val
				} else {
					minPref[i] = 0
				}
			} else {
				run := int64(0)
				mp := int64(0)
				for _, sub := range sp.seq {
					if run+minPref[sub] < mp {
						mp = run + minPref[sub]
					}
					run = clamp(run + total[sub])
				}
				total[i] = run
				minPref[i] = clamp(mp)
			}
		}

		last := len(spells) - 1
		if hp+minPref[last] > 0 {
			fmt.Fprintln(out, -1)
			continue
		}

		var kill func(id int, cur int64) int
		kill = func(id int, cur int64) int {
			sp := spells[id]
			if sp.basic {
				if cur+sp.val <= 0 {
					return id
				}
				return -1
			}
			r := cur
			for _, sub := range sp.seq {
				if r+minPref[sub] <= 0 {
					return kill(sub, r)
				}
				r = clamp(r + total[sub])
			}
			return -1
		}

		res := kill(last, hp)
		if res > 0 && res <= n {
			fmt.Fprintln(out, res)
		} else {
			fmt.Fprintln(out, -1)
		}
	}
}
