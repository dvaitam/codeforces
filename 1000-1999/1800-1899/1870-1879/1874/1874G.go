package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type State struct {
	hp  int
	dmg int
	val int64
}

var (
	n, m   int
	typ    []int
	pa, pb []int // for type1
	px, py []int // for type2/3
	pw     []int // for type4
	g      [][]int

	dpExit  [][]State
	visExit []bool
	propDP  []int64
	visProp []bool
)

func prune(states []State) []State {
	mp := make(map[[2]int]int64)
	for _, s := range states {
		key := [2]int{s.hp, s.dmg}
		if v, ok := mp[key]; !ok || s.val > v {
			mp[key] = s.val
		}
	}
	arr := make([]State, 0, len(mp))
	for k, v := range mp {
		arr = append(arr, State{hp: k[0], dmg: k[1], val: v})
	}
	sort.Slice(arr, func(i, j int) bool {
		if arr[i].hp != arr[j].hp {
			return arr[i].hp < arr[j].hp
		}
		if arr[i].dmg != arr[j].dmg {
			return arr[i].dmg < arr[j].dmg
		}
		return arr[i].val > arr[j].val
	})
	res := make([]State, 0, len(arr))
	for _, s := range arr {
		dominated := false
		for i := 0; i < len(res); {
			r := res[i]
			if r.hp >= s.hp && r.dmg >= s.dmg && r.val >= s.val {
				dominated = true
				break
			}
			if s.hp >= r.hp && s.dmg >= r.dmg && s.val >= r.val {
				res = append(res[:i], res[i+1:]...)
			} else {
				i++
			}
		}
		if !dominated {
			res = append(res, s)
		}
	}
	return res
}

func solveExit(v int) []State {
	if visExit[v] {
		return dpExit[v]
	}
	visExit[v] = true
	states := make([]State, 0)
	if len(g[v]) == 0 {
		states = append(states, State{})
	} else {
		for _, to := range g[v] {
			states = append(states, solveExit(to)...)
		}
		states = prune(states)
	}
	res := make([]State, 0, len(states))
	for _, s := range states {
		ns := s
		switch typ[v] {
		case 1:
			ns.val += int64(pa[v]) * int64(pb[v])
		case 2:
			ns.hp += px[v]
		case 3:
			ns.dmg += py[v]
		case 4:
			ns.val += int64(pw[v])
		}
		res = append(res, ns)
	}
	dpExit[v] = prune(res)
	return dpExit[v]
}

func startStates(v int) []State {
	states := make([]State, 0)
	if len(g[v]) == 0 {
		states = append(states, State{})
	} else {
		for _, to := range g[v] {
			states = append(states, solveExit(to)...)
		}
		states = prune(states)
	}
	res := make([]State, 0, len(states))
	for _, s := range states {
		ns := s
		switch typ[v] {
		case 2:
			ns.hp += px[v]
		case 3:
			ns.dmg += py[v]
		case 4:
			ns.val += int64(pw[v])
		}
		res = append(res, ns)
	}
	return prune(res)
}

func propValue(v int) int64 {
	if visProp[v] {
		return propDP[v]
	}
	visProp[v] = true
	best := int64(-1 << 60)
	if len(g[v]) == 0 {
		best = 0
	} else {
		for _, to := range g[v] {
			val := propValue(to)
			if val > best {
				best = val
			}
		}
	}
	add := int64(0)
	switch typ[v] {
	case 1:
		add = int64(pa[v]) * int64(pb[v])
	case 4:
		add = int64(pw[v])
	}
	propDP[v] = best + add
	return propDP[v]
}

func main() {
	in := bufio.NewReader(os.Stdin)
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	typ = make([]int, n+1)
	pa = make([]int, n+1)
	pb = make([]int, n+1)
	px = make([]int, n+1)
	py = make([]int, n+1)
	pw = make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &typ[i])
		switch typ[i] {
		case 1:
			fmt.Fscan(in, &pa[i], &pb[i])
		case 2:
			fmt.Fscan(in, &px[i])
		case 3:
			fmt.Fscan(in, &py[i])
		case 4:
			fmt.Fscan(in, &pw[i])
		}
	}
	g = make([][]int, n+1)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		g[u] = append(g[u], v)
	}

	dpExit = make([][]State, n+1)
	visExit = make([]bool, n+1)
	propDP = make([]int64, n+1)
	visProp = make([]bool, n+1)

	for i := n; i >= 1; i-- {
		if !visExit[i] {
			solveExit(i)
		}
	}

	ans := int64(0)
	hasCard := false
	for i := 1; i <= n; i++ {
		if typ[i] == 1 {
			hasCard = true
			st := startStates(i)
			a := int64(pa[i])
			b := int64(pb[i])
			for _, s := range st {
				hp := a + int64(s.hp)
				dmg := b + int64(s.dmg)
				val := s.val
				cand := hp*dmg*1_000_000_000 + val
				if cand > ans {
					ans = cand
				}
			}
		}
	}

	if !hasCard {
		ans = propValue(1)
	} else {
		pv := propValue(1)
		if pv > ans {
			ans = pv
		}
	}

	fmt.Println(ans)
}
