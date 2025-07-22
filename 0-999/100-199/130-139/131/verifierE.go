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

type queen struct{ r, c int }

func solveCase(n, m int, qs []queen) string {
	idx := make([]int, m)
	for i := 0; i < m; i++ {
		idx[i] = i
	}
	att := make([]int, m)
	// rows
	sort.Slice(idx, func(i, j int) bool {
		if qs[idx[i]].r != qs[idx[j]].r {
			return qs[idx[i]].r < qs[idx[j]].r
		}
		return qs[idx[i]].c < qs[idx[j]].c
	})
	for i := 0; i+1 < m; i++ {
		a, b := idx[i], idx[i+1]
		if qs[a].r == qs[b].r {
			att[a]++
			att[b]++
		}
	}
	// cols
	for i := 0; i < m; i++ {
		idx[i] = i
	}
	sort.Slice(idx, func(i, j int) bool {
		if qs[idx[i]].c != qs[idx[j]].c {
			return qs[idx[i]].c < qs[idx[j]].c
		}
		return qs[idx[i]].r < qs[idx[j]].r
	})
	for i := 0; i+1 < m; i++ {
		a, b := idx[i], idx[i+1]
		if qs[a].c == qs[b].c {
			att[a]++
			att[b]++
		}
	}
	// main diag r-c
	for i := 0; i < m; i++ {
		idx[i] = i
	}
	sort.Slice(idx, func(i, j int) bool {
		ai, aj := idx[i], idx[j]
		di := qs[ai].r - qs[ai].c
		dj := qs[aj].r - qs[aj].c
		if di != dj {
			return di < dj
		}
		return qs[ai].r < qs[aj].r
	})
	for i := 0; i+1 < m; i++ {
		a, b := idx[i], idx[i+1]
		if qs[a].r-qs[a].c == qs[b].r-qs[b].c {
			att[a]++
			att[b]++
		}
	}
	// anti diag r+c
	for i := 0; i < m; i++ {
		idx[i] = i
	}
	sort.Slice(idx, func(i, j int) bool {
		ai, aj := idx[i], idx[j]
		si := qs[ai].r + qs[ai].c
		sj := qs[aj].r + qs[aj].c
		if si != sj {
			return si < sj
		}
		return qs[ai].r < qs[aj].r
	})
	for i := 0; i+1 < m; i++ {
		a, b := idx[i], idx[i+1]
		if qs[a].r+qs[a].c == qs[b].r+qs[b].c {
			att[a]++
			att[b]++
		}
	}
	res := make([]int, 9)
	for i := 0; i < m; i++ {
		if att[i] >= 0 && att[i] <= 8 {
			res[att[i]]++
		}
	}
	var sb strings.Builder
	for i := 0; i < 9; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(res[i]))
	}
	return sb.String()
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 1
	maxM := n * n
	if maxM > 50 {
		maxM = 50
	}
	m := rng.Intn(maxM) + 1
	used := make(map[int]bool)
	qs := make([]queen, 0, m)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for len(qs) < m {
		r := rng.Intn(n) + 1
		c := rng.Intn(n) + 1
		key := r*n + c
		if used[key] {
			continue
		}
		used[key] = true
		qs = append(qs, queen{r, c})
		sb.WriteString(fmt.Sprintf("%d %d\n", r, c))
	}
	input := sb.String()
	expected := solveCase(n, m, qs)
	return input, expected
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
	outStr := strings.TrimSpace(out.String())
	if outStr != expected {
		return fmt.Errorf("expected %s got %s", expected, outStr)
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
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
