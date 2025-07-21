package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type result struct {
	win      bool
	own, opp int64
}

func better(a, b result) bool {
	if a.win != b.win {
		return a.win && !b.win
	}
	if a.own != b.own {
		return a.own > b.own
	}
	return a.opp < b.opp
}

func solveCase(words []string) (string, string) {
	subMap := make(map[string]int)
	num := make([]int64, 0)
	for _, w := range words {
		seen := make(map[string]bool)
		L := len(w)
		for l := 0; l < L; l++ {
			for r := l + 1; r <= L; r++ {
				s := w[l:r]
				if !seen[s] {
					seen[s] = true
					if id, ok := subMap[s]; ok {
						num[id]++
					} else {
						id = len(num)
						subMap[s] = id
						num = append(num, 1)
					}
				}
			}
		}
	}
	m := len(num)
	subs := make([]string, m)
	for s, id := range subMap {
		subs[id] = s
	}
	weight := make([]int64, m)
	maxLen := 0
	buckets := make(map[int][]int)
	for id, s := range subs {
		sum := int64(0)
		for i := 0; i < len(s); i++ {
			sum += int64(s[i] - 'a' + 1)
		}
		weight[id] = sum * num[id]
		L := len(s)
		if L > maxLen {
			maxLen = L
		}
		buckets[L] = append(buckets[L], id)
	}
	nexts := make([][]int, m)
	for id, s := range subs {
		for c := byte('a'); c <= 'z'; c++ {
			t1 := string(c) + s
			if j, ok := subMap[t1]; ok {
				nexts[id] = append(nexts[id], j)
			}
			t2 := s + string(c)
			if j, ok := subMap[t2]; ok {
				nexts[id] = append(nexts[id], j)
			}
		}
	}
	dp := make([]result, m)
	for L := maxLen; L >= 1; L-- {
		for _, id := range buckets[L] {
			bestSet := false
			var best result
			for _, j := range nexts[id] {
				child := dp[j]
				w := weight[j]
				res := result{win: !child.win, own: w + child.opp, opp: child.own}
				if !bestSet || better(res, best) {
					bestSet = true
					best = res
				}
			}
			if bestSet {
				dp[id] = best
			} else {
				dp[id] = result{win: false, own: 0, opp: 0}
			}
		}
	}
	var initSet bool
	var initRes result
	for _, id := range buckets[1] {
		child := dp[id]
		w := weight[id]
		res := result{win: !child.win, own: w + child.opp, opp: child.own}
		if !initSet || better(res, initRes) {
			initSet = true
			initRes = res
		}
	}
	if initRes.win {
		return "First", fmt.Sprintf("%d %d", initRes.own, initRes.opp)
	}
	return "Second", fmt.Sprintf("%d %d", initRes.own, initRes.opp)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(4) + 1
	words := make([]string, n)
	for i := 0; i < n; i++ {
		l := rng.Intn(3) + 1
		b := make([]byte, l)
		for j := range b {
			b[j] = byte('a' + rng.Intn(3))
		}
		words[i] = string(b)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		sb.WriteString(words[i])
		if i+1 < n {
			sb.WriteByte('\n')
		}
	}
	if n > 0 {
		sb.WriteByte('\n')
	}
	w, res := solveCase(words)
	return sb.String(), w + "\n" + res
}

func runCase(bin, input, expected string) error {
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
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != strings.TrimSpace(expected) {
		return fmt.Errorf("expected\n%s\ngot\n%s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
