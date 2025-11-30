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

const mod = 1000000007

// Embedded testcases from testcasesA.txt.
const embeddedTestcasesA = `13
14
2
9
17
16
13
10
16
12
19
7
17
5
10
5
4
20
9
18
20
5
10
4
3
11
16
18
4
12
14
11
20
7
18
16
15
17
9
2
18
1
3
13
1
20
16
11
8
11
3
7
19
8
8
5
18
15
3
3
11
17
16
4
10
18
10
4
18
11
18
7
20
18
19
10
15
3
20
13
11
19
8
10
6
7
6
2
20
9
16
3
3
5
5
2
3
18
13
17`

func expected(n int) int64 {
	var f0, f1 int64
	for i := 1; i <= n; i++ {
		add := int64(1)
		if i%2 == 0 {
			add = (add + f1) % mod
			f0 = (f0 + add) % mod
		} else {
			add = (add + f0) % mod
			f1 = (f1 + add) % mod
		}
	}
	return (f0 + f1) % mod
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	scanner := bufio.NewScanner(strings.NewReader(embeddedTestcasesA))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		n, _ := strconv.Atoi(line)
		exp := expected(n)
		input := fmt.Sprintf("%d\n", n)
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewBufferString(input)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", idx, err)
			os.Exit(1)
		}
		gotStr := strings.TrimSpace(string(out))
		got, err2 := strconv.ParseInt(gotStr, 10, 64)
		if err2 != nil || got != exp {
			fmt.Printf("Test %d failed: expected %d got %s\n", idx, exp, gotStr)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
