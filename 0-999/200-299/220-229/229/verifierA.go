package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type testCaseA struct {
	n, m int
	rows []string
	ans  int
}

func solveA(tc testCaseA) int {
	const INF = int(1e9)
	total := make([]int, tc.m)
	for _, row := range tc.rows {
		has1 := false
		for i := 0; i < tc.m; i++ {
			if row[i] == '1' {
				has1 = true
				break
			}
		}
		if !has1 {
			return -1
		}
		dist := make([]int, tc.m)
		for i := range dist {
			dist[i] = INF
		}
		last := -INF
		for j := 0; j < 2*tc.m; j++ {
			idx := j % tc.m
			if row[idx] == '1' {
				last = j
			}
			if last >= 0 {
				d := j - last
				if d < dist[idx] {
					dist[idx] = d
				}
			}
		}
		last = 2*tc.m + INF
		for j := 2*tc.m - 1; j >= 0; j-- {
			idx := j % tc.m
			if row[idx] == '1' {
				last = j
			}
			d := last - j
			if d < dist[idx] {
				dist[idx] = d
			}
		}
		for i := 0; i < tc.m; i++ {
			total[i] += dist[i]
		}
	}
	ans := INF
	for i := 0; i < tc.m; i++ {
		if total[i] < ans {
			ans = total[i]
		}
	}
	if ans >= INF {
		return -1
	}
	return ans
}

func genCaseA(rng *rand.Rand) testCaseA {
	n := rng.Intn(4) + 1
	m := rng.Intn(6) + 1
	rows := make([]string, n)
	for i := 0; i < n; i++ {
		b := make([]byte, m)
		has1 := false
		for j := 0; j < m; j++ {
			if rng.Intn(2) == 0 {
				b[j] = '0'
			} else {
				b[j] = '1'
				has1 = true
			}
		}
		if !has1 && rng.Intn(5) != 0 {
			b[rng.Intn(m)] = '1'
		}
		rows[i] = string(b)
	}
	tc := testCaseA{n: n, m: m, rows: rows}
	tc.ans = solveA(tc)
	return tc
}

func runCaseA(bin string, tc testCaseA) error {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.m)
	for _, r := range tc.rows {
		fmt.Fprintln(&sb, r)
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(strings.TrimSpace(out.String()))
	if len(fields) != 1 {
		return fmt.Errorf("expected 1 number got %d", len(fields))
	}
	val, err := strconv.Atoi(fields[0])
	if err != nil {
		return fmt.Errorf("invalid integer: %v", err)
	}
	if val != tc.ans {
		return fmt.Errorf("expected %d got %d", tc.ans, val)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := genCaseA(rng)
		if err := runCaseA(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
