package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

func solveCase(line string) string {
	fields := strings.Fields(line)
	idx := 0
	n, _ := strconv.Atoi(fields[idx])
	idx++
	edges := make([][2]int, n-1)
	for i := 0; i < n-1; i++ {
		a, _ := strconv.Atoi(fields[idx])
		idx++
		b, _ := strconv.Atoi(fields[idx])
		idx++
		edges[i] = [2]int{a - 1, b - 1}
	}
	best := -1
	pairs := make(map[[2]int]struct{})
	states := make([]int, n)
	var dfs func(int)
	dfs = func(i int) {
		if i == n {
			a, b := 0, 0
			for _, v := range states {
				if v == 1 {
					a++
				} else if v == 2 {
					b++
				}
			}
			if a == 0 || b == 0 {
				return
			}
			for _, e := range edges {
				x, y := states[e[0]], states[e[1]]
				if (x == 1 && y == 2) || (x == 2 && y == 1) {
					return
				}
			}
			sum := a + b
			if sum > best {
				best = sum
				pairs = map[[2]int]struct{}{{a, b}: {}}
			} else if sum == best {
				pairs[[2]int{a, b}] = struct{}{}
			}
			return
		}
		for t := 0; t < 3; t++ {
			states[i] = t
			dfs(i + 1)
		}
	}
	dfs(0)
	res := make([][2]int, 0, len(pairs))
	for p := range pairs {
		res = append(res, p)
	}
	sort.Slice(res, func(i, j int) bool { return res[i][0] < res[j][0] })
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(res)))
	for _, p := range res {
		sb.WriteString(fmt.Sprintf("%d %d\n", p[0], p[1]))
	}
	return strings.TrimSpace(sb.String())
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("runtime error: %v", err)
	}
	return strings.TrimSpace(string(out)), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesE.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		expected := solveCase(line)
		got, err := run(bin, line+"\n")
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx, expected, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
