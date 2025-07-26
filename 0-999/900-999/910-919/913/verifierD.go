package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type Problem struct {
	a   int
	t   int
	idx int
}

const maxT = 10000

func check(problems []Problem, limit int, k int) (bool, []int) {
	if k == 0 {
		return true, []int{}
	}
	buckets := make([][]int, maxT+1)
	for _, p := range problems {
		if p.a >= k {
			buckets[p.t] = append(buckets[p.t], p.idx)
		}
	}
	count := 0
	for i := 1; i <= maxT; i++ {
		count += len(buckets[i])
	}
	if count < k {
		return false, nil
	}
	res := make([]int, 0, k)
	timeSum := 0
	for t := 1; t <= maxT && len(res) < k; t++ {
		for _, idx := range buckets[t] {
			res = append(res, idx)
			timeSum += t
			if len(res) == k {
				break
			}
		}
	}
	if len(res) < k || timeSum > limit {
		return false, nil
	}
	return true, res
}

func solve(n int, limit int, a, t []int) (int, []int) {
	problems := make([]Problem, n)
	for i := 0; i < n; i++ {
		problems[i] = Problem{a: a[i], t: t[i], idx: i + 1}
	}
	lo, hi := 0, n
	best := 0
	var bestSet []int
	for lo <= hi {
		mid := (lo + hi) / 2
		ok, set := check(problems, limit, mid)
		if ok {
			best = mid
			bestSet = set
			lo = mid + 1
		} else {
			hi = mid - 1
		}
	}
	return best, bestSet
}

func randomInput() (int, int, []int, []int) {
	n := rand.Intn(10) + 1
	limit := rand.Intn(100) + 1
	a := make([]int, n)
	t := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = rand.Intn(n) + 1
		t[i] = rand.Intn(20) + 1
	}
	return n, limit, a, t
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rand.Seed(1)
	const cases = 100
	for i := 0; i < cases; i++ {
		n, limit, a, t := randomInput()
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, limit))
		for j := 0; j < n; j++ {
			sb.WriteString(fmt.Sprintf("%d %d\n", a[j], t[j]))
		}
		input := sb.String()
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n", i+1, err)
			fmt.Printf("program output:\n%s\n", string(out))
			return
		}
		got := strings.TrimSpace(string(out))
		best, ids := solve(n, limit, a, t)
		var expected strings.Builder
		expected.WriteString(fmt.Sprintf("%d\n", best))
		expected.WriteString(fmt.Sprintf("%d\n", best))
		for j, id := range ids {
			if j > 0 {
				expected.WriteByte(' ')
			}
			expected.WriteString(fmt.Sprintf("%d", id))
		}
		if best > 0 {
			expected.WriteByte('\n')
		}
		want := strings.TrimSpace(expected.String())
		if got != want {
			fmt.Printf("case %d failed:\ninput:\n%sexpected:\n%s\nGot:\n%s\n", i+1, input, want, got)
			return
		}
	}
	fmt.Printf("OK %d cases\n", cases)
}
