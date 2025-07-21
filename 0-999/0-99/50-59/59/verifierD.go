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

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveCaseD(n int, results []int, teams [][3]int, k int) string {
	total := 3 * n
	rank := make([]int, total+1)
	for i, id := range results {
		rank[id] = i + 1
	}
	isCap := false
	capIdx := -1
	mateA, mateB := 0, 0
	for i := 0; i < n; i++ {
		team := teams[i]
		cap := team[0]
		for j := 1; j < 3; j++ {
			if rank[team[j]] < rank[cap] {
				cap = team[j]
			}
		}
		if cap == k {
			isCap = true
			capIdx = i
			a, b := -1, -1
			for j := 0; j < 3; j++ {
				if team[j] == k {
					continue
				}
				if a == -1 {
					a = team[j]
				} else {
					b = team[j]
				}
			}
			if a > b {
				a, b = b, a
			}
			mateA, mateB = a, b
			break
		}
	}
	var ans []int
	if !isCap {
		for i := 1; i <= total; i++ {
			if i == k {
				continue
			}
			ans = append(ans, i)
		}
		sort.Ints(ans)
	} else {
		takenBefore := make([]bool, total+1)
		for i := 0; i < capIdx; i++ {
			for j := 0; j < 3; j++ {
				takenBefore[teams[i][j]] = true
			}
		}
		var sBefore []int
		for i := 1; i <= total; i++ {
			if i == k {
				continue
			}
			if takenBefore[i] {
				sBefore = append(sBefore, i)
			}
		}
		sort.Ints(sBefore)
		var sR []int
		for i := 1; i <= total; i++ {
			if i == k || i == mateA || i == mateB {
				continue
			}
			if !takenBefore[i] {
				sR = append(sR, i)
			}
		}
		sort.Ints(sR)
		idxB := 0
		haveA, haveB := false, false
		for !haveA {
			if idxB < len(sBefore) && sBefore[idxB] < mateA {
				ans = append(ans, sBefore[idxB])
				idxB++
			} else {
				ans = append(ans, mateA)
				haveA = true
			}
		}
		for !haveB {
			if idxB < len(sBefore) && sBefore[idxB] < mateB {
				ans = append(ans, sBefore[idxB])
				idxB++
			} else {
				ans = append(ans, mateB)
				haveB = true
			}
		}
		iB, iR := idxB, 0
		for iB < len(sBefore) && iR < len(sR) {
			if sBefore[iB] < sR[iR] {
				ans = append(ans, sBefore[iB])
				iB++
			} else {
				ans = append(ans, sR[iR])
				iR++
			}
		}
		for iB < len(sBefore) {
			ans = append(ans, sBefore[iB])
			iB++
		}
		for iR < len(sR) {
			ans = append(ans, sR[iR])
			iR++
		}
	}
	var sb strings.Builder
	for i, v := range ans {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	return sb.String()
}

func generateCaseD(rng *rand.Rand) (string, string) {
	n := rng.Intn(4) + 1
	total := 3 * n
	perm := rng.Perm(total)
	results := make([]int, total)
	for i := 0; i < total; i++ {
		results[i] = perm[i] + 1
	}
	remaining := make([]int, total)
	for i := 0; i < total; i++ {
		remaining[i] = results[i]
	}
	rank := make(map[int]int, total)
	for i, id := range results {
		rank[id] = i
	}
	teams := make([][3]int, n)
	used := make(map[int]bool, total)
	for i := 0; i < n; i++ {
		best := -1
		bestRank := total + 1
		for _, id := range remaining {
			if used[id] {
				continue
			}
			if rank[id] < bestRank {
				bestRank = rank[id]
				best = id
			}
		}
		pool := []int{}
		for _, id := range remaining {
			if used[id] || id == best {
				continue
			}
			pool = append(pool, id)
		}
		m1 := pool[rng.Intn(len(pool))]
		pool2 := []int{}
		for _, id := range pool {
			if id != m1 {
				pool2 = append(pool2, id)
			}
		}
		m2 := pool2[rng.Intn(len(pool2))]
		teams[i] = [3]int{best, m1, m2}
		used[best] = true
		used[m1] = true
		used[m2] = true
	}
	k := rng.Intn(total) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range results {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		fmt.Fprintf(&sb, "%d %d %d\n", teams[i][0], teams[i][1], teams[i][2])
	}
	sb.WriteString(fmt.Sprintf("%d\n", k))
	expect := solveCaseD(n, results, teams, k)
	return sb.String(), expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCaseD(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected\n%s\ngot\n%s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
