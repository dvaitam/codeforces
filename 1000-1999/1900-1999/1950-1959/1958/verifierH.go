package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
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

func solveCase(n int, hp int64, basic []int64, comps [][]int) int {
	spells := make([]Spell, 0)
	spells = append(spells, Spell{})
	for _, v := range basic {
		spells = append(spells, Spell{basic: true, val: v})
	}
	for _, seq := range comps {
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
		return -1
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
		return res
	}
	return -1
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rand.Seed(49)
	const t = 100
	cases := make([]struct {
		n     int
		hp    int64
		basic []int64
		comps [][]int
	}, t)
	for i := 0; i < t; i++ {
		n := rand.Intn(3) + 1 // small n
		hp := int64(rand.Intn(50) + 1)
		basic := make([]int64, n)
		for j := 0; j < n; j++ {
			basic[j] = int64(rand.Intn(21)) - 10
		}
		m := rand.Intn(3)
		comps := make([][]int, m)
		for j := 0; j < m; j++ {
			s := rand.Intn(n+m-1) + 1
			seq := make([]int, s)
			for k := 0; k < s; k++ {
				seq[k] = rand.Intn(n+j) + 1
			}
			comps[j] = seq
		}
		cases[i] = struct {
			n     int
			hp    int64
			basic []int64
			comps [][]int
		}{n, hp, basic, comps}
	}

	var input strings.Builder
	fmt.Fprintln(&input, t)
	for _, c := range cases {
		fmt.Fprintf(&input, "%d %d\n", c.n, c.hp)
		for i, v := range c.basic {
			if i > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprint(&input, v)
		}
		input.WriteByte('\n')
		fmt.Fprintln(&input, len(c.comps))
		for _, seq := range c.comps {
			fmt.Fprintf(&input, "%d", len(seq))
			for _, x := range seq {
				fmt.Fprintf(&input, " %d", x)
			}
			input.WriteByte('\n')
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input.String())
	outBytes, err := cmd.Output()
	if err != nil {
		fmt.Println("error running binary:", err)
		os.Exit(1)
	}
	outs := strings.Fields(strings.TrimSpace(string(outBytes)))
	if len(outs) != t {
		fmt.Printf("expected %d lines, got %d\n", t, len(outs))
		os.Exit(1)
	}
	for i, s := range outs {
		var got int
		fmt.Sscan(s, &got)
		want := solveCase(cases[i].n, cases[i].hp, cases[i].basic, cases[i].comps)
		if got != want {
			fmt.Printf("mismatch on case %d expected %d got %d\n", i+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
