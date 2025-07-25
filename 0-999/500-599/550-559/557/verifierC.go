package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type group struct {
	length int
	costs  []int
}

func sumSmallest(freq []int, k int) int {
	if k <= 0 {
		return 0
	}
	sum := 0
	for cost := 1; cost < len(freq) && k > 0; cost++ {
		if freq[cost] > 0 {
			take := freq[cost]
			if take > k {
				take = k
			}
			sum += take * cost
			k -= take
		}
	}
	return sum
}

func solveC(n int, legs []struct{ l, d int }) int {
	sort.Slice(legs, func(i, j int) bool { return legs[i].l < legs[j].l })
	groups := make([]group, 0)
	for i := 0; i < n; {
		j := i
		for j < n && legs[j].l == legs[i].l {
			j++
		}
		g := group{length: legs[i].l, costs: make([]int, j-i)}
		for k := i; k < j; k++ {
			g.costs[k-i] = legs[k].d
		}
		groups = append(groups, g)
		i = j
	}
	m := len(groups)
	suffixCost := make([]int, m+1)
	for i := m - 1; i >= 0; i-- {
		sum := 0
		for _, v := range groups[i].costs {
			sum += v
		}
		suffixCost[i] = suffixCost[i+1] + sum
	}
	freq := make([]int, 201)
	prefixCost := 0
	countShorter := 0
	ans := int(^uint(0) >> 1)
	for idx, g := range groups {
		sort.Ints(g.costs)
		c := len(g.costs)
		prefixGroup := make([]int, c+1)
		for i := 0; i < c; i++ {
			prefixGroup[i+1] = prefixGroup[i] + g.costs[i]
		}
		groupTotal := prefixGroup[c]
		costLonger := suffixCost[idx+1]
		for x := 1; x <= c; x++ {
			keepOthers := x - 1
			if keepOthers > countShorter {
				keepOthers = countShorter
			}
			costKeepOthers := sumSmallest(freq, keepOthers)
			costRemoveShorter := prefixCost - costKeepOthers
			costRemoveGroup := groupTotal - prefixGroup[x]
			total := costLonger + costRemoveShorter + costRemoveGroup
			if total < ans {
				ans = total
			}
		}
		for _, v := range g.costs {
			freq[v]++
		}
		prefixCost += groupTotal
		countShorter += c
	}
	if ans < 0 {
		ans = 0
	}
	return ans
}

func genCase() (string, int) {
	n := rand.Intn(10) + 1
	legs := make([]struct{ l, d int }, n)
	for i := 0; i < n; i++ {
		legs[i].l = rand.Intn(10) + 1
	}
	for i := 0; i < n; i++ {
		legs[i].d = rand.Intn(200) + 1
	}
	res := solveC(n, append([]struct{ l, d int }{}, legs...))
	input := fmt.Sprintf("%d\n", n)
	for i, v := range legs {
		if i > 0 {
			input += " "
		}
		input += fmt.Sprintf("%d", v.l)
	}
	input += "\n"
	for i, v := range legs {
		if i > 0 {
			input += " "
		}
		input += fmt.Sprintf("%d", v.d)
	}
	input += "\n"
	return input, res
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	rand.Seed(1)
	for i := 0; i < 100; i++ {
		input, expected := genCase()
		cmd := exec.Command(binary)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		cmd.Stdout = &out
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		gotStr := strings.TrimSpace(out.String())
		var got int
		if _, err := fmt.Sscan(gotStr, &got); err != nil {
			fmt.Fprintf(os.Stderr, "test %d: failed to parse output\n", i+1)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "test %d failed\ninput:\n%sexpected: %d\ngot: %s\n", i+1, input, expected, gotStr)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
