package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

const BIG = 1 << 30

// computeOptimalValue computes the optimal partition value for the given array.
// Uses the same algorithm as the reference: sort, find the bit level where groups <= 4,
// then compute the best split within each group.
func computeOptimalValue(a []int) int {
	n := len(a)
	type pair struct{ val, idx int }
	b := make([]pair, n)
	for i := range a {
		b[i] = pair{a[i], i}
	}
	sort.Slice(b, func(i, j int) bool { return b[i].val < b[j].val })

	mask := 0
	var ss [][]int
	for bit := 29; bit >= 0; bit-- {
		pre := -1
		cnt := 0
		ok := true
		for j := 0; j < n; j++ {
			curMask := b[j].val & mask
			if curMask != pre {
				pre = curMask
				cnt = 0
			}
			cnt++
			if cnt > 4 {
				ok = false
				break
			}
		}
		if !ok {
			mask |= 1 << bit
			continue
		}
		pre = -1
		var cur []int
		for j := 0; j < n; j++ {
			curMask := b[j].val & mask
			if curMask != pre {
				if len(cur) > 0 {
					ss = append(ss, cur)
				}
				pre = curMask
				cur = nil
			}
			cur = append(cur, b[j].idx)
		}
		if len(cur) > 0 {
			ss = append(ss, cur)
		}
		break
	}

	result := BIG
	for _, p := range ss {
		l := len(p)
		if l <= 2 {
			// cost is BIG for sets with < 2 elements in the "other" group
			// single or pair groups don't constrain
			continue
		}
		// For groups of 3 or 4, find the best pair to put together (as '0'),
		// rest go to '1'. The partition value for this group is the minimum
		// of (min xor among '0' group, min xor among '1' group).
		// Actually, the value is min(cost(S0), cost(S1)) overall.
		// We need to compute what the optimal split achieves for each group.
		groupVal := 0
		if l == 3 {
			x01 := a[p[0]] ^ a[p[1]]
			x02 := a[p[0]] ^ a[p[2]]
			x12 := a[p[1]] ^ a[p[2]]
			// Best pair gives highest XOR (the pair goes to same set, single to other)
			// Actually the value is the max of the three XOR values
			groupVal = x01
			if x02 > groupVal {
				groupVal = x02
			}
			if x12 > groupVal {
				groupVal = x12
			}
		} else { // l == 4
			// Try all 3 ways to split 4 into 2+2
			// For each split, the value is min(xor of pair1, xor of pair2)
			vals := make([]int, 4)
			for i := range p {
				vals[i] = a[p[i]]
			}
			best := 0
			// 3 partitions: {0,1},{2,3}  {0,2},{1,3}  {0,3},{1,2}
			pairs := [][2][2]int{
				{{0, 1}, {2, 3}},
				{{0, 2}, {1, 3}},
				{{0, 3}, {1, 2}},
			}
			for _, pp := range pairs {
				v1 := vals[pp[0][0]] ^ vals[pp[0][1]]
				v2 := vals[pp[1][0]] ^ vals[pp[1][1]]
				m := v1
				if v2 < m {
					m = v2
				}
				if m > best {
					best = m
				}
			}
			groupVal = best
		}
		if groupVal < result {
			result = groupVal
		}
	}
	return result
}

// computePartitionValue computes the value of a given partition (min cost of both sets).
func computePartitionValue(a []int, partition string) int {
	var s0, s1 []int
	for i, c := range partition {
		if c == '0' {
			s0 = append(s0, a[i])
		} else {
			s1 = append(s1, a[i])
		}
	}
	cost0 := setCost(s0)
	cost1 := setCost(s1)
	if cost0 < cost1 {
		return cost0
	}
	return cost1
}

func setCost(s []int) int {
	if len(s) < 2 {
		return BIG
	}
	sort.Ints(s)
	minXor := BIG
	for i := 1; i < len(s); i++ {
		v := s[i-1] ^ s[i]
		if v < minXor {
			minXor = v
		}
	}
	return minXor
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return out.String(), fmt.Errorf("%v\n%s", err, errBuf.String())
	}
	return out.String(), nil
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(5) + 2
	vals := make([]int, n)
	used := make(map[int]bool)
	for i := 0; i < n; i++ {
		v := rng.Intn(1024)
		for used[v] {
			v = rng.Intn(1024)
		}
		used[v] = true
		vals[i] = v
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", vals[i]))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func parseInput(input string) []int {
	r := strings.NewReader(input)
	var n int
	fmt.Fscan(r, &n)
	a := make([]int, n)
	for i := range a {
		fmt.Fscan(r, &a[i])
	}
	return a
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := genCase(rng)
		a := parseInput(input)
		n := len(a)
		optVal := computeOptimalValue(a)

		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		gotStr := strings.TrimSpace(got)
		if len(gotStr) != n {
			fmt.Fprintf(os.Stderr, "test %d failed: expected length %d, got length %d\noutput: %s\ninput:\n%s", i+1, n, len(gotStr), gotStr, input)
			os.Exit(1)
		}
		for _, c := range gotStr {
			if c != '0' && c != '1' {
				fmt.Fprintf(os.Stderr, "test %d failed: invalid character in output: %s\ninput:\n%s", i+1, gotStr, input)
				os.Exit(1)
			}
		}
		gotVal := computePartitionValue(a, gotStr)
		if gotVal != optVal {
			fmt.Fprintf(os.Stderr, "test %d failed: optimal value %d but got partition value %d\noutput: %s\ninput:\n%s", i+1, optVal, gotVal, gotStr, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
