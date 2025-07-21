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

type caseH struct {
	n    int
	grid [3][]int
}

func genCase(rng *rand.Rand) caseH {
	n := rng.Intn(4) + 2
	g := [3][]int{}
	for i := 0; i < 3; i++ {
		g[i] = make([]int, n)
		for j := 0; j < n; j++ {
			idx := rng.Intn(n) + 1
			if rng.Intn(2) == 0 {
				g[i][j] = idx
			} else {
				g[i][j] = -idx
			}
		}
	}
	return caseH{n, g}
}

func canWin(tc caseH) bool {
	n := tc.n
	for mask := 0; mask < (1 << n); mask++ {
		assign := make([]int, n)
		for i := 0; i < n; i++ {
			if mask&(1<<i) != 0 {
				assign[i] = 1
			} else {
				assign[i] = -1
			}
		}
		ok := true
		for j := 0; j < n && ok; j++ {
			col := []int{
				val(tc.grid[0][j], assign),
				val(tc.grid[1][j], assign),
				val(tc.grid[2][j], assign),
			}
			sort.Ints(col)
			if col[1] != 1 {
				ok = false
			}
		}
		if ok {
			return true
		}
	}
	return false
}

func val(x int, a []int) int {
	if x > 0 {
		return a[x-1]
	}
	return -a[-x-1]
}

func runCase(bin string, tc caseH) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("1\n%d\n", tc.n))
	for i := 0; i < 3; i++ {
		for j := 0; j < tc.n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", tc.grid[i][j]))
		}
		sb.WriteByte('\n')
	}
	input := sb.String()
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.ToLower(strings.TrimSpace(out.String()))
	exp := "no"
	if canWin(tc) {
		exp = "yes"
	}
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		tc := genCase(rng)
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
