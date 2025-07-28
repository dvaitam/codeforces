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

type caseC struct {
	n, m    int
	grid    []string
	q       int
	queries [][2]int
}

func genCase(rng *rand.Rand) caseC {
	n := rng.Intn(5) + 1
	m := rng.Intn(5) + 1
	grid := make([]string, n)
	for i := 0; i < n; i++ {
		row := make([]byte, m)
		for j := 0; j < m; j++ {
			if rng.Intn(2) == 0 {
				row[j] = '.'
			} else {
				row[j] = 'X'
			}
		}
		grid[i] = string(row)
	}
	q := rng.Intn(5) + 1
	qs := make([][2]int, q)
	for i := 0; i < q; i++ {
		x1 := rng.Intn(m) + 1
		x2 := rng.Intn(m) + 1
		if x1 > x2 {
			x1, x2 = x2, x1
		}
		qs[i] = [2]int{x1, x2}
	}
	return caseC{n, m, grid, q, qs}
}

func expected(tc caseC) []string {
	bad := make([]int, tc.m+1)
	for i := 1; i < tc.n; i++ {
		for j := 1; j < tc.m; j++ {
			if tc.grid[i-1][j] == 'X' && tc.grid[i][j-1] == 'X' {
				bad[j+1] = 1
			}
		}
	}
	pref := make([]int, tc.m+1)
	for j := 1; j <= tc.m; j++ {
		pref[j] = pref[j-1] + bad[j]
	}
	ans := make([]string, tc.q)
	for i, qu := range tc.queries {
		x1, x2 := qu[0], qu[1]
		if pref[x2]-pref[x1] == 0 {
			ans[i] = "YES"
		} else {
			ans[i] = "NO"
		}
	}
	return ans
}

func runCase(bin string, tc caseC) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
	for _, row := range tc.grid {
		sb.WriteString(row)
		sb.WriteByte('\n')
	}
	sb.WriteString(fmt.Sprintf("%d\n", tc.q))
	for _, qu := range tc.queries {
		sb.WriteString(fmt.Sprintf("%d %d\n", qu[0], qu[1]))
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
	fields := strings.Fields(out.String())
	if len(fields) != tc.q {
		return fmt.Errorf("expected %d answers got %d", tc.q, len(fields))
	}
	exp := expected(tc)
	for i := 0; i < tc.q; i++ {
		got := strings.ToUpper(fields[i])
		if got != exp[i] {
			return fmt.Errorf("query %d expected %s got %s", i+1, exp[i], got)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
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
