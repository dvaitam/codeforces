package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

func runCandidate(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expected(n int, strengths [3]int, criminals []int) int {
	s := []int{strengths[0], strengths[1], strengths[2]}
	sort.Ints(s)
	a, b, c := s[0], s[1], s[2]
	sort.Ints(criminals)
	used := make([]bool, n)
	remain := n
	pa := n - 1
	for pa >= 0 && criminals[pa] > a {
		pa--
	}
	pb := n - 1
	for pb >= 0 && criminals[pb] > b {
		pb--
	}
	r := n - 1
	hours := 0
	remove := func(limit int, p *int) bool {
		for *p >= 0 && (used[*p] || criminals[*p] > limit) {
			(*p)--
		}
		if *p >= 0 {
			used[*p] = true
			remain--
			(*p)--
			return true
		}
		return false
	}
	for remain > 0 {
		for r >= 0 && used[r] {
			r--
		}
		if r < 0 {
			break
		}
		idx := r
		x := criminals[idx]
		used[idx] = true
		remain--
		r--
		hours++
		if x > a+b+c {
			return -1
		}
		if x > b+c {
			continue
		}
		if x > c {
			if x > a+c {
				remove(a, &pa)
			} else {
				remove(b, &pb)
			}
			continue
		}
		remove(b, &pb)
		remove(a, &pa)
	}
	return hours
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(5))
	for t := 0; t < 100; t++ {
		n := rng.Intn(20) + 1
		strengths := [3]int{rng.Intn(50) + 1, rng.Intn(50) + 1, rng.Intn(50) + 1}
		criminals := make([]int, n)
		for i := 0; i < n; i++ {
			criminals[i] = rng.Intn(50) + 1
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", n)
		fmt.Fprintf(&sb, "%d %d %d\n", strengths[0], strengths[1], strengths[2])
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", criminals[i])
		}
		sb.WriteByte('\n')
		input := sb.String()
		want := expected(n, strengths, append([]int(nil), criminals...))
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", t+1, err)
			os.Exit(1)
		}
		got, err := strconv.Atoi(strings.TrimSpace(out))
		if err != nil || got != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\n", t+1, want, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
