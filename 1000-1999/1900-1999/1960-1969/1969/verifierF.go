package main

import (
	"fmt"
	"math/bits"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

const testcasesF = `10 6 6 3 6 6 6 5 1 4 2 6
2 2 1 2
8 4 4 1 2 1 2 4 3 2
8 4 1 2 4 2 2 1 1 2
4 4 2 3 3 2
10 4 2 2 4 3 1 3 4 2 2 3
2 2 2 1
10 6 1 3 3 3 4 6 3 2 4 4
4 2 2 1 2 2
2 2 2 2
8 8 1 8 1 3 4 2 4 8
6 6 3 5 3 4 1 5
6 6 1 4 1 2 3 5
10 6 2 3 3 6 5 1 3 6 3 3
4 2 1 2 2 1
2 2 2 1
4 4 3 3 4 4
4 2 1 2 2 1
4 4 2 4 1 2
8 6 2 1 4 3 2 4 5 2
10 8 8 6 8 5 5 8 7 3 2 7
10 4 4 3 2 1 4 3 3 1 3 1
6 6 5 6 6 3 4 3
6 6 6 2 5 1 4 5
6 6 6 3 4 3 5 6
6 6 3 6 3 6 4 3
4 4 3 3 2 2
4 4 4 3 1 4
4 4 4 3 3 1
4 4 4 2 2 2
2 2 1 1
2 2 1 2
4 4 2 1 4 4
6 6 6 5 1 5 2 2
6 2 2 2 2 2 1 2
2 2 2 1
6 6 5 3 5 3 3 2
8 8 6 8 3 3 7 8 4 3
10 2 2 1 2 1 2 1 2 2 2 2
4 4 3 4 4 4
4 4 3 3 4 3
8 2 2 2 2 2 1 2 1 1
10 8 4 5 1 3 7 1 8 5 4 8
2 2 2 2
4 4 3 4 4 1
2 2 2 2
10 6 5 3 6 5 1 4 3 3 6 6
10 4 1 1 3 4 1 2 1 4 4 4
8 8 7 4 5 8 2 5 1 7
10 6 6 4 3 2 2 4 6 5 4 3
10 4 4 1 1 2 3 1 1 1 3 4
2 2 2 1
2 2 1 1
4 4 4 2 2 2
8 4 1 2 1 1 4 1 3 1
10 10 2 7 10 3 1 7 2 6 10 8
8 6 6 3 1 2 6 3 2 5
10 6 5 5 5 2 3 1 5 2 3 3
10 4 2 3 4 4 2 2 1 3 1 1
2 2 2 1
8 8 8 6 2 1 7 1 3 7
4 2 1 2 1 1
10 8 6 4 5 6 6 7 7 3 6 5
8 6 5 1 5 5 2 6 2 4
2 2 1 1
4 4 2 1 4 4
6 2 2 2 2 2 2 2
8 2 1 1 1 1 2 2 2 1
10 6 1 3 2 3 6 2 1 4 1 3
6 6 1 6 1 6 4 4
4 2 2 2 2 1
4 2 1 1 1 1
10 8 6 5 3 8 6 5 2 8 3 4
4 4 1 4 3 1
8 2 2 2 1 1 1 2 1 1
4 4 1 4 1 2
2 2 2 1
2 2 2 2
6 2 1 1 1 2 1 1
2 2 1 1
10 4 1 2 1 4 3 4 4 4 2 1
6 6 2 5 1 2 3 4
2 2 2 2
4 4 1 4 2 1
8 8 4 5 1 8 4 1 4 5
6 2 1 2 2 2 1 2
4 4 1 4 2 1
2 2 2 1
8 8 3 6 2 5 3 3 5 8
6 6 1 1 5 6 5 1
2 2 1 1
6 6 3 5 1 2 5 4
4 4 2 3 3 1
4 4 2 1 3 2
4 4 4 4 1 3
6 6 2 4 5 4 2 4
2 2 2 2
10 2 1 1 2 1 1 1 1 1 2 2
8 8 4 4 5 4 5 6 4 6
6 4 1 4 4 3 2 2`

func bruteMaxCoins(a []int, n, k int) int {
	type state struct {
		mask uint64
		ptr  int
	}
	memo := make(map[state]int)
	var dfs func(mask uint64, ptr int) int
	dfs = func(mask uint64, ptr int) int {
		st := state{mask, ptr}
		if v, ok := memo[st]; ok {
			return v
		}
		if mask == 0 && ptr >= n {
			return 0
		}
		indices := make([]int, 0, bits.OnesCount64(mask))
		for i := 0; i < n; i++ {
			if mask&(1<<i) != 0 {
				indices = append(indices, i)
			}
		}
		best := 0
		for i := 0; i < len(indices); i++ {
			for j := i + 1; j < len(indices); j++ {
				m := mask &^ (1 << indices[i]) &^ (1 << indices[j])
				p := ptr
				if p < n {
					m |= 1 << p
					p++
				}
				if p < n {
					m |= 1 << p
					p++
				}
				gain := 0
				if a[indices[i]] == a[indices[j]] {
					gain = 1
				}
				val := gain + dfs(m, p)
				if val > best {
					best = val
				}
			}
		}
		memo[st] = best
		return best
	}
	startMask := (uint64(1) << k) - 1
	return dfs(startMask, k)
}

func solveGreedy(a []int, n, k int) int {
	next := make([]int, n)
	last := make(map[int]int)
	for i := n - 1; i >= 0; i-- {
		if v, ok := last[a[i]]; ok {
			next[i] = v
		} else {
			next[i] = n
		}
		last[a[i]] = i
	}

	hand := make([]int, k)
	for i := 0; i < k; i++ {
		hand[i] = i
	}
	ptr := k
	coins := 0

	for len(hand) > 0 {
		groups := make(map[int][]int)
		for _, idx := range hand {
			groups[a[idx]] = append(groups[a[idx]], idx)
		}
		var choose [2]int
		chosen := false
		for _, idxs := range groups {
			if len(idxs) >= 2 {
				sort.Slice(idxs, func(i, j int) bool { return next[idxs[i]] > next[idxs[j]] })
				choose[0], choose[1] = idxs[0], idxs[1]
				chosen = true
				break
			}
		}
		if !chosen {
			sort.Slice(hand, func(i, j int) bool { return next[hand[i]] > next[hand[j]] })
			choose[0], choose[1] = hand[0], hand[1]
		}
		newHand := make([]int, 0, len(hand)-2)
		for _, idx := range hand {
			if idx == choose[0] || idx == choose[1] {
				continue
			}
			newHand = append(newHand, idx)
		}
		hand = newHand
		if a[choose[0]] == a[choose[1]] {
			coins++
		}
		if ptr < n {
			hand = append(hand, ptr)
			ptr++
		}
		if ptr < n {
			hand = append(hand, ptr)
			ptr++
		}
	}
	return coins
}

func parseLine(line string) (int, int, []int, error) {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return 0, 0, nil, fmt.Errorf("not enough fields")
	}
	n, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, 0, nil, err
	}
	k, err := strconv.Atoi(fields[1])
	if err != nil {
		return 0, 0, nil, err
	}
	if len(fields) != 2+n {
		return 0, 0, nil, fmt.Errorf("expected %d numbers, got %d", 2+n, len(fields))
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		v, err := strconv.Atoi(fields[2+i])
		if err != nil {
			return 0, 0, nil, err
		}
		a[i] = v
	}
	return n, k, a, nil
}

func buildInput(n, k int, a []int) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, k)
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierF /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	lines := strings.Split(testcasesF, "\n")
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		n, k, a, err := parseLine(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d parse error: %v\n", i+1, err)
			os.Exit(1)
		}
		input := buildInput(n, k, a)
		want := solveGreedy(a, n, k)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n%s\n", i+1, err, got)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strconv.Itoa(want) {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%s\nexpected: %d\ngot: %s\n", i+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(lines))
}
