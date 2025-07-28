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

func runProg(bin string, input string) (string, error) {
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
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveCase(n int, a []int, queries [][2]int) []int {
	maxIdx := 1
	for i := 1; i <= n; i++ {
		if a[i] > a[maxIdx] {
			maxIdx = i
		}
	}
	wins := make([][]int, n+1)
	champion := 1
	for j := 2; j <= n; j++ {
		round := j - 1
		if a[champion] > a[j] {
			wins[champion] = append(wins[champion], round)
		} else {
			wins[j] = append(wins[j], round)
			champion = j
		}
	}
	res := make([]int, len(queries))
	for idx, q := range queries {
		i := q[0]
		k := q[1]
		limit := k
		if limit > n-1 {
			limit = n - 1
		}
		rounds := wins[i]
		cnt := sort.Search(len(rounds), func(x int) bool { return rounds[x] > limit })
		if i == maxIdx && k > n-1 {
			cnt += k - (n - 1)
		}
		res[idx] = cnt
	}
	return res
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesC.txt")
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
		parts := strings.Fields(line)
		if len(parts) < 2 {
			fmt.Fprintf(os.Stderr, "bad line %d\n", idx)
			os.Exit(1)
		}
		pos := 0
		n, _ := strconv.Atoi(parts[pos])
		pos++
		qNum, _ := strconv.Atoi(parts[pos])
		pos++
		if len(parts) != 2+n+2*qNum {
			fmt.Fprintf(os.Stderr, "bad line %d length", idx)
			os.Exit(1)
		}
		a := make([]int, n+1)
		for i := 1; i <= n; i++ {
			val, _ := strconv.Atoi(parts[pos])
			pos++
			a[i] = val
		}
		queries := make([][2]int, qNum)
		for i := 0; i < qNum; i++ {
			first, _ := strconv.Atoi(parts[pos])
			pos++
			second, _ := strconv.Atoi(parts[pos])
			pos++
			queries[i] = [2]int{first, second}
		}
		inputBuilder := &strings.Builder{}
		fmt.Fprintf(inputBuilder, "1\n%d %d\n", n, qNum)
		for i := 1; i <= n; i++ {
			if i > 1 {
				inputBuilder.WriteByte(' ')
			}
			fmt.Fprintf(inputBuilder, "%d", a[i])
		}
		inputBuilder.WriteByte('\n')
		for _, qu := range queries {
			fmt.Fprintf(inputBuilder, "%d %d\n", qu[0], qu[1])
		}
		out, err := runProg(bin, inputBuilder.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		expectedRes := solveCase(n, a, queries)
		outTokens := strings.Fields(out)
		if len(outTokens) != len(expectedRes) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d numbers got %d\n", idx, len(expectedRes), len(outTokens))
			os.Exit(1)
		}
		for i, tok := range outTokens {
			val, err := strconv.Atoi(tok)
			if err != nil || val != expectedRes[i] {
				fmt.Fprintf(os.Stderr, "case %d failed: expected %v got %v\n", idx, expectedRes, outTokens)
				os.Exit(1)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
