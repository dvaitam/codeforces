package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

func expected(coins []int64, query int64) int64 {
	// sort coins ascending
	sort.Slice(coins, func(i, j int) bool { return coins[i] < coins[j] })
	// compress
	var vals []int64
	var cnt []int64
	for _, c := range coins {
		if len(vals) == 0 || vals[len(vals)-1] != c {
			vals = append(vals, c)
			cnt = append(cnt, 1)
		} else {
			cnt[len(cnt)-1]++
		}
	}
	rem := query
	used := int64(0)
	for i := len(vals) - 1; i >= 0 && rem > 0; i-- {
		take := rem / vals[i]
		if take > cnt[i] {
			take = cnt[i]
		}
		rem -= take * vals[i]
		used += take
	}
	if rem != 0 {
		return -1
	}
	return used
}

func runCase(bin string, coins []int64, queries []int64) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", len(coins), len(queries)))
	for i, v := range coins {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	sb.WriteByte('\n')
	for i, q := range queries {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(q, 10))
	}
	sb.WriteByte('\n')
	input := sb.String()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	tokens := strings.Fields(strings.TrimSpace(out.String()))
	if len(tokens) != len(queries) {
		return fmt.Errorf("expected %d numbers got %d", len(queries), len(tokens))
	}
	for i, tok := range tokens {
		got, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid integer %q", tok)
		}
		want := expected(coins, queries[i])
		if got != want {
			return fmt.Errorf("query %d expected %d got %d", i+1, want, got)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesD.txt")
	if err != nil {
		fmt.Println("could not open testcasesD.txt:", err)
		os.Exit(1)
	}
	defer f.Close()
	scan := bufio.NewScanner(f)
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("invalid test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		if !scan.Scan() {
			fmt.Println("invalid test file")
			os.Exit(1)
		}
		m, _ := strconv.Atoi(scan.Text())
		coins := make([]int64, n)
		for j := 0; j < n; j++ {
			if !scan.Scan() {
				fmt.Println("invalid test file")
				os.Exit(1)
			}
			v, _ := strconv.ParseInt(scan.Text(), 10, 64)
			coins[j] = v
		}
		queries := make([]int64, m)
		for j := 0; j < m; j++ {
			if !scan.Scan() {
				fmt.Println("invalid test file")
				os.Exit(1)
			}
			q, _ := strconv.ParseInt(scan.Text(), 10, 64)
			queries[j] = q
		}
		if err := runCase(bin, coins, queries); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
