package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `1 2 0
9 3 1 1 0 1 0 0 1 0 0
19 3 1 1 0 1 0 1 0 0 1 0 0 1 1 1 1 1 0 0 1
16 3 1 0 1 0 1 1 0 1 1 0 0 1 0 0 0 0
8 1 1 0 0 1 1 1 1 0
14 1 1 0 1 0 1 0 1 0 0 0 1 0 0 1
13 3 0 0 0 0 0 0 1 1 1 0 0 1 0
19 4 1 1 1 1 1 0 0 0 1 1 1 0 0 1 0 1 1 1 0
3 5 0 0 0
5 1 1 0 1 1 0
5 4 1 1 0 0 0
14 2 1 1 1 1 0 0 0 0 1 1 1 0 1 1
10 4 1 1 0 0 0 0 1 1 0 0
16 4 0 1 0 0 1 1 0 1 1 1 0 0 1 0 0 0
20 5 0 1 0 1 0 1 1 0 0 1 1 0 1 0 1 1 0 0 0 1
15 3 0 1 1 0 0 0 1 1 0 0 0 0 1 1 0
5 4 1 0 0 1 0
4 4 0 0 1 1
14 1 1 1 1 0 1 0 0 1 0 1 1 0 0 0
12 1 0 0 0 0 0 0 1 1 0 0 0 0
15 1 1 1 1 0 0 0 1 1 1 1 1 1 0 0 0
18 5 1 0 1 0 0 0 1 0 0 0 1 1 1 0 0 0 1 0
3 2 1 1 1
5 5 1 1 1 0 0
15 3 0 1 1 0 1 1 1 0 1 1 1 0 0 0 1
4 5 0 1 1 0
14 4 0 0 1 1 0 1 1 0 1 0 1 0 1 1
1 2 1
4 3 0 1 1 0
16 2 1 1 0 0 1 0 0 1 1 1 1 0 1 1 1 0
20 1 0 0 1 1 1 0 0 0 0 1 1 1 0 0 0 1 1 1 1 1
12 5 0 0 0 0 1 1 1 0 1 1 1 1
7 4 1 1 1 0 1 1 0
16 1 0 0 1 1 1 0 0 0 1 0 1 0 0 0 1 1
8 2 1 0 0 1 1 1 1 0
11 4 1 0 1 0 1 0 0 0 1 1 1
14 4 1 0 0 1 0 0 1 0 0 0 0 1 0 0
8 5 1 0 1 1 1 1 0 1
19 4 1 1 1 0 1 1 0 0 0 0 1 0 1 0 0 0 1 0 1
18 2 1 0 1 0 0 1 1 1 1 1 0 0 0 1 0 0 1 0
7 3 1 1 1 1 1 1 0
2 3 0 0
15 4 1 0 1 1 1 1 0 0 0 0 0 0 1 0 1
20 1 1 1 0 0 0 1 0 0 1 1 1 1 0 1 0 1 1 1 1 1
8 1 0 0 0 1 0 0 1 0
18 5 1 0 0 1 1 0 1 0 1 0 0 0 1 1 0 1 1 1
19 2 0 0 1 1 0 0 1 0 0 0 0 0 1 0 0 0 0 1 1
20 5 0 1 0 0 1 1 1 1 0 0 0 1 1 1 0 1 0 1 0 1
10 4 0 0 0 1 0 1 1 0 1 1
6 3 0 1 0 0 1 0
16 1 0 0 0 0 0 1 0 0 1 1 0 0 1 1 0 1
10 2 1 1 1 0 1 1 1 0 0 1
13 2 0 0 0 0 1 0 0 0 0 0 0 0 1
20 1 1 0 0 1 1 1 0 0 1 0 1 1 0 1 0 1 0 1 1 1
16 4 0 1 1 1 0 1 0 1 1 0 1 1 1 0 0 1
2 3 1 0
17 5 1 1 1 0 0 1 0 1 0 0 1 0 0 0 1 1 1
10 3 1 1 1 0 0 0 0 1 0 1
1 4 1
20 2 0 1 1 0 0 0 0 0 1 1 1 1 1 0 1 1 1 1 1 0
14 2 0 0 1 1 1 0 1 0 0 1 0 0 1 1
1 2 0
9 2 1 0 1 0 1 0 0 0 1
3 4 1 1 1
14 3 1 0 1 0 1 0 1 0 1 1 1 1 0 1
12 5 0 1 1 0 1 1 0 0 0 1 1 0
3 2 0 1 0
11 1 1 1 0 1 1 0 0 1 1 0 0
5 2 1 0 0 1 0
12 5 0 0 0 1 0 1 1 1 1 1 0 1
18 2 0 0 0 1 0 0 1 0 1 0 1 0 1 0 0 1 0 0
4 4 1 1 1 1
17 3 0 0 1 0 0 0 0 1 0 0 1 0 1 1 0 0 1
3 4 0 1 1
20 5 1 1 1 1 0 0 1 1 0 1 0 1 0 0 0 1 0 0 1 0
11 4 1 1 0 0 1 1 0 0 0 1 1
4 4 0 0 1 0
1 4 1
6 5 0 0 0 0 0 1
1 4 1
2 2 0 1
11 2 0 1 0 0 1 1 1 0 1 0 0
19 1 1 0 0 0 0 1 0 0 0 0 0 0 1 0 1 1 1 0 0
12 5 0 1 0 0 1 1 0 1 1 0 1 1
8 4 0 1 1 0 1 1 1 0
20 5 1 0 1 0 1 1 0 0 0 1 0 1 1 1 0 1 1 0 0 1
16 4 0 1 0 1 1 0 1 1 0 0 1 0 1 1 1 0
10 5 1 0 1 1 1 0 1 0 0 0
20 3 0 0 1 1 1 0 0 1 0 0 0 0 0 0 0 0 1 1 1 1
4 4 1 0 1 1
18 2 1 1 0 1 1 1 1 1 1 1 0 0 0 0 0 0 0 1
20 2 1 0 0 0 0 1 1 1 0 1 1 1 1 0 0 0 0 0 0 1
20 2 1 1 1 0 0 1 0 1 0 1 0 1 1 0 1 0 1 0 1 1
18 4 1 1 0 0 1 0 1 0 1 0 0 1 1 1 0 1 1 0
4 3 0 1 0 0
2 2 0 0
2 4 1 0
16 5 1 1 0 0 0 0 1 1 1 0 0 1 1 1 0 0
10 5 1 0 1 1 0 0 1 0 1 1
10 4 0 0 1 1 1 0 1 1 1 1`

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func solve(n, r int, a []int) int {
	cnt := 0
	pos := 1
	for pos <= n {
		temp := min(pos+r-1, n)
		lower := pos - r + 1
		if lower < 1 {
			lower = 1
		}
		found := -1
		for i := temp; i >= lower; i-- {
			if a[i-1] == 1 {
				found = i
				break
			}
		}
		if found == -1 {
			return -1
		}
		cnt++
		pos = found + r
	}
	return cnt
}

func runCase(bin string, n, r int, arr []int) error {
	var input strings.Builder
	input.WriteString(fmt.Sprintf("%d %d\n", n, r))
	for i, v := range arr {
		if i > 0 {
			input.WriteByte(' ')
		}
		input.WriteString(strconv.Itoa(v))
	}
	input.WriteByte('\n')
	cmd := exec.Command(bin)
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	}
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	expected := fmt.Sprintf("%d", solve(n, r, arr))
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 2 {
			fmt.Printf("invalid test case line %d\n", idx+1)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(parts[0])
		r, _ := strconv.Atoi(parts[1])
		if len(parts) != n+2 {
			fmt.Printf("invalid test case line %d\n", idx+1)
			os.Exit(1)
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			arr[i], _ = strconv.Atoi(parts[2+i])
		}
		if err := runCase(bin, n, r, arr); err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		idx++
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
