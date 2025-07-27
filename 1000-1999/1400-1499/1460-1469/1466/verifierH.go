package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const mod int64 = 1000000007

func allPerms(n int) [][]int {
	res := [][]int{}
	p := make([]int, n)
	for i := 0; i < n; i++ {
		p[i] = i
	}
	var gen func(int)
	gen = func(i int) {
		if i == n {
			cp := make([]int, n)
			copy(cp, p)
			res = append(res, cp)
			return
		}
		for j := i; j < n; j++ {
			p[i], p[j] = p[j], p[i]
			gen(i + 1)
			p[i], p[j] = p[j], p[i]
		}
	}
	gen(0)
	return res
}

func better(pref []int, item1, item2 int) bool {
	for _, v := range pref {
		if v == item1 {
			return true
		}
		if v == item2 {
			return false
		}
	}
	return false
}

func isOptimal(assign []int, prefs [][]int) bool {
	n := len(assign)
	for mask := 1; mask < (1 << n); mask++ {
		// build subset S
		indices := []int{}
		for i := 0; i < n; i++ {
			if mask&(1<<i) != 0 {
				indices = append(indices, i)
			}
		}
		// try all permutations of items among S
		perms := allPerms(len(indices))
		for _, perm := range perms {
			betterFlag := false
			ok := true
			for idx, agent := range indices {
				item := assign[indices[perm[idx]]]
				if better(prefs[agent], item, assign[agent]) {
					betterFlag = true
				} else if item != assign[agent] {
					if better(prefs[agent], assign[agent], item) {
						ok = false
						break
					}
				}
			}
			if ok && betterFlag {
				return false
			}
		}
	}
	return true
}

func countProfiles(n int, assign []int) int64 {
	perms := allPerms(n)
	total := int64(0)
	for _, p1 := range perms {
		for _, p2 := range perms {
			for _, p3 := range perms[:1] { // use at most 1 to reduce
				prefs := make([][]int, n)
				prefs[0] = p1
				if n > 1 {
					prefs[1] = p2
				}
				if n > 2 {
					prefs[2] = p3
				}
				if isOptimal(assign, prefs) {
					total++
				}
			}
		}
	}
	return total % mod
}

func solveH(assign []int) string {
	n := len(assign)
	val := countProfiles(n, assign)
	return fmt.Sprint(val)
}

func genCases() []string {
	rand.Seed(8)
	cases := make([]string, 100)
	for i := 0; i < 100; i++ {
		n := rand.Intn(3) + 1
		perm := rand.Perm(n)
		for j := range perm {
			perm[j]++
		}
		sb := strings.Builder{}
		sb.WriteString(fmt.Sprint(n))
		sb.WriteByte('\n')
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(perm[j]))
		}
		sb.WriteByte('\n')
		cases[i] = sb.String()
	}
	return cases
}

func runCase(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := genCases()
	for i, tc := range cases {
		lines := strings.Split(strings.TrimSpace(tc), "\n")
		var n int
		fmt.Sscan(lines[0], &n)
		items := make([]int, n)
		parts := strings.Fields(lines[1])
		for j := 0; j < n; j++ {
			fmt.Sscan(parts[j], &items[j])
		}
		want := solveH(items)
		got, err := runCase(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Runtime error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "Wrong answer on case %d\nInput:\n%sExpected: %s Got: %s\n", i+1, tc, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
