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

func brute(n, l, s int, x []int) (uint64, []int, bool) {
	idxs := make([]int, 0, n-1)
	for i := 1; i <= n; i++ {
		if i != s {
			idxs = append(idxs, i)
		}
	}
	bestTime := uint64(^uint64(0))
	var bestOrder []int
	perm := make([]int, len(idxs))
	copy(perm, idxs)
	for {
		left := l
		right := n - 1 - l
		cur := s
		var tm uint64
		valid := true
		for _, nxt := range perm {
			if nxt < cur {
				if left == 0 {
					valid = false
					break
				}
				left--
				tm += uint64(x[cur] - x[nxt])
			} else {
				if right == 0 {
					valid = false
					break
				}
				right--
				tm += uint64(x[nxt] - x[cur])
			}
			cur = nxt
		}
		if valid && left == 0 && right == 0 {
			if tm < bestTime {
				bestTime = tm
				bestOrder = append([]int(nil), perm...)
			}
		}
		if !nextPermutation(perm) {
			break
		}
	}
	if bestOrder == nil {
		return 0, nil, false
	}
	return bestTime, bestOrder, true
}

func nextPermutation(a []int) bool {
	i := len(a) - 2
	for i >= 0 && a[i] >= a[i+1] {
		i--
	}
	if i < 0 {
		return false
	}
	j := len(a) - 1
	for a[j] <= a[i] {
		j--
	}
	a[i], a[j] = a[j], a[i]
	for k, l := i+1, len(a)-1; k < l; k, l = k+1, l-1 {
		a[k], a[l] = a[l], a[k]
	}
	return true
}

func generateCase(rng *rand.Rand) (string, uint64, []int, bool) {
	n := rng.Intn(6) + 2
	l := rng.Intn(n - 1)
	s := rng.Intn(n) + 1
	x := make([]int, n+1)
	cur := 0
	for i := 1; i <= n; i++ {
		cur += rng.Intn(5) + 1
		x[i] = cur
	}
	time, order, ok := brute(n, l, s, x)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, l, s)
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(x[i]))
	}
	sb.WriteByte('\n')
	return sb.String(), time, order, ok
}

func runCase(bin string, input string, expTime uint64, expOrder []int, has bool) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outLines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if !has {
		if strings.TrimSpace(outLines[0]) != "-1" {
			return fmt.Errorf("expected -1")
		}
		return nil
	}
	var gotTime uint64
	fmt.Sscan(strings.TrimSpace(outLines[0]), &gotTime)
	if gotTime != expTime {
		return fmt.Errorf("expected time %d got %d", expTime, gotTime)
	}
	if len(outLines) < 2 {
		return fmt.Errorf("missing sequence line")
	}
	gotStrs := strings.Fields(strings.TrimSpace(outLines[1]))
	if len(gotStrs) != len(expOrder) {
		return fmt.Errorf("sequence length mismatch")
	}
	for i, s := range gotStrs {
		var v int
		fmt.Sscan(s, &v)
		if v != expOrder[i] {
			return fmt.Errorf("seq pos %d expected %d got %d", i+1, expOrder[i], v)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, ans, order, ok := generateCase(rng)
		if err := runCase(bin, in, ans, order, ok); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
