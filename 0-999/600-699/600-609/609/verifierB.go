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

func expectedPairs(n, m int, arr []int) int64 {
	cnt := make([]int, m+1)
	for _, g := range arr {
		cnt[g]++
	}
	rem := n
	var ans int64
	for i := 1; i <= m; i++ {
		rem -= cnt[i]
		ans += int64(cnt[i] * rem)
	}
	return ans
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	const testcasesBRaw = `7 3 1 3 2 1 2 1 2
12 3 2 2 2 1 2 1 3 2 3 3 2 3
7 3 1 1 3 2 2 2 2
17 4 4 1 4 1 4 2 1 2 4 2 2 3 2 4 1 4 2
16 3 3 1 1 2 1 3 2 3 2 2 1 1 1 3 1 1
8 3 3 1 2 1 3 3 1 2
15 5 5 1 3 1 1 2 3 2 1 4 4 4 4 5 4
7 5 5 1 4 1 3 5 2
10 5 4 1 4 2 5 3 2 2 2 3
13 5 4 2 3 2 3 1 4 2 2 4 5 4 1
15 3 3 1 2 3 1 3 2 3 2 1 2 1 3 1 3
8 3 1 3 3 2 2 1 2 3
10 3 1 2 3 2 2 1 1 3 1 2
11 3 3 3 2 3 3 1 2 3 3 3 1
4 3 3 2 3 1
19 2 1 2 1 2 1 2 2 2 2 2 2 2 2 2 2 1 2 2 1
14 5 3 3 4 3 2 5 1 3 3 3 1 1 1 2
6 4 3 2 4 1 1 4
13 2 2 1 2 1 1 2 1 1 1 2 1 2 1
20 4 4 4 2 2 2 2 4 3 4 2 3 4 3 3 4 1 3 1 1 3
20 3 2 3 2 1 2 3 1 2 1 2 3 1 3 2 1 1 1 2 1 3
5 5 1 5 3 2 4
4 2 2 2 1 2
5 3 3 3 1 1 2
7 4 4 4 3 3 1 2 2
10 5 1 3 1 1 5 1 4 5 2 2
16 2 2 1 1 1 1 2 2 2 1 1 1 1 2 2 1 1
5 3 2 1 2 3 2
11 2 2 1 1 1 1 1 1 1 2 1 2
11 2 1 1 2 2 2 1 2 1 1 1 1
9 5 3 1 5 4 1 3 1 2 5
6 2 2 2 2 2 1 2
10 3 1 2 3 3 2 3 2 1 2 1
16 3 2 3 3 3 3 3 3 2 1 2 3 1 1 2 1 2
20 2 2 1 2 2 2 1 1 2 2 2 1 1 2 1 2 2 2 2 1 1
6 3 1 2 3 2 2 1
19 2 2 2 2 1 2 2 1 1 1 1 2 1 2 2 2 1 1 1 1
6 4 2 4 3 3 4 1
18 2 2 1 1 1 2 2 1 2 1 2 2 2 2 1 2 1 1 2
17 3 1 1 2 1 1 1 2 3 2 1 1 3 2 1 2 3 1
10 4 3 4 1 2 2 3 2 4 1 1
7 5 3 4 4 2 1 2 5
16 3 1 3 3 2 3 2 3 1 1 1 1 3 1 3 2 3
20 3 2 1 3 1 2 2 1 1 1 3 2 3 1 3 3 2 2 2 2 2
14 3 3 1 1 2 3 2 1 2 3 1 2 2 3 3
6 2 1 2 1 1 2 2
18 5 4 3 3 4 1 5 4 4 1 5 2 3 4 3 2 4 1 5
16 3 3 3 1 3 2 2 1 2 1 2 2 3 1 3 3 3
11 3 3 2 2 1 3 3 1 3 3 3 3
17 4 1 3 4 1 1 3 1 3 4 3 2 4 2 1 1 1 4
13 4 2 2 3 4 2 1 3 1 1 2 3 1 1
18 2 2 2 1 2 1 1 2 1 2 1 1 1 2 1 2 2 2 1
13 5 1 4 5 5 3 2 2 1 5 5 2 3 3
12 4 2 2 3 4 4 2 1 3 3 3 1 2
18 4 1 3 4 4 4 1 3 2 3 4 2 3 3 2 2 2 4 1
11 5 3 2 4 4 3 3 5 3 1 3 4
11 5 5 5 4 2 4 3 1 1 4 3 2
10 5 4 1 3 4 3 3 1 5 1 2
3 2 2 1 2
9 2 2 1 1 1 1 2 1 2 2
9 3 3 2 1 2 2 1 1 3 1
17 5 1 5 5 2 4 5 5 3 4 4 5 3 2 2 3 5 4
10 4 2 3 4 4 4 1 1 2 2 4
20 4 1 3 1 3 1 3 3 3 3 4 2 3 2 4 2 3 3 4 4 4
10 5 4 1 4 5 2 1 2 5 3 4
10 4 1 4 4 1 4 1 2 3 2 4
5 2 1 1 2 1 1
10 2 2 2 2 2 1 2 1 2 2 1
10 5 2 4 4 4 3 5 1 2 2 3
7 4 4 3 1 2 2 4 3
8 5 5 4 1 3 2 1 3 5
7 5 4 5 5 1 2 3 4
8 5 1 5 2 4 1 3 4 2
12 3 1 2 1 2 2 1 3 2 3 2 2 3
17 3 3 2 1 3 1 3 2 1 3 2 3 1 2 3 3 1 3
8 2 2 2 1 1 2 1 1 1
19 3 2 2 2 1 1 1 3 3 3 2 2 2 2 1 2 3 3 2 3
12 5 1 5 3 3 1 4 4 1 1 1 2 2
11 4 3 2 1 4 3 2 1 2 1 2 2
17 5 5 3 3 1 5 5 1 1 4 4 4 1 5 1 3 5 2
20 5 5 4 4 5 1 3 3 1 3 3 1 1 1 5 2 2 3 3 5 4
13 2 1 2 2 2 1 1 2 2 2 2 2 2 1
18 3 3 3 2 2 2 2 1 2 3 1 3 1 3 1 3 3 2 1
20 2 1 2 2 1 1 2 1 1 1 2 2 1 2 2 1 1 1 1 2 2
17 4 4 3 4 2 1 1 4 4 2 2 1 1 3 2 1 1 1
13 2 2 2 2 2 1 1 2 1 2 2 2 1 1
19 4 3 3 2 3 3 2 3 4 1 3 4 3 4 3 3 4 2 2 4
5 5 1 3 2 5 4
20 3 1 2 2 1 1 1 3 2 3 1 1 2 2 1 2 2 3 1 1 2
8 5 5 2 3 5 2 4 4 1
20 5 5 3 2 4 2 3 3 1 3 5 4 4 5 1 2 3 2 4 1 4
10 4 2 2 2 2 1 3 4 4 4 3
7 2 2 1 1 2 2 1 2
12 2 2 1 1 2 2 1 2 2 1 1 2 1
8 3 3 2 2 1 1 1 2 3
16 5 4 2 1 2 5 1 4 5 1 5 1 4 3 3 5 3
6 3 3 1 1 3 2 3
17 5 3 4 2 5 5 5 2 4 4 5 4 1 2 3 1 1 4
10 5 3 3 5 3 4 3 2 4 1 4
9 5 2 1 2 3 5 5 2 4 2`

	scanner := bufio.NewScanner(strings.NewReader(testcasesBRaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		parts := strings.Fields(line)
		if len(parts) < 2 {
			fmt.Printf("test %d invalid\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(parts[0])
		m, _ := strconv.Atoi(parts[1])
		if len(parts)-2 < n {
			fmt.Printf("test %d invalid length\n", idx)
			os.Exit(1)
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			v, _ := strconv.Atoi(parts[2+i])
			arr[i] = v
		}
		exp := expectedPairs(n, m, arr)
		var buf bytes.Buffer
		fmt.Fprintf(&buf, "%d %d", n, m)
		for i := 0; i < n; i++ {
			fmt.Fprintf(&buf, " %d", arr[i])
		}
		buf.WriteByte('\n')
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewReader(buf.Bytes())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", idx, err)
			os.Exit(1)
		}
		got := strings.TrimSpace(string(out))
		if got != fmt.Sprint(exp) {
			fmt.Printf("Test %d failed: expected %d got %s\n", idx, exp, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
