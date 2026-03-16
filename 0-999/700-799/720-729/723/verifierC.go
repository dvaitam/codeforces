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

const testcasesCRaw = `100
13 13
14 2 9 17 16 13 10 16 12 7 17 5 10
5 1
5 3 5 6 5
5 3
2 2 6 8 2
12 7
6 10 11 4 9 8 8 9 5 1 9 1
3 3
7 1 8
11 4
6 2 4 4 4 3 9 8 2 2 6
17 16
4 10 18 10 4 18 11 18 7 20 18 19 10 15 3 20 13
11 10
4 5 3 4 14 3 1 10 11 5 8
3 1
6 2 2
2 1
6 5
13 12
17 9 17 8 7 14 9 15 16 12 3 11 4
16 11
7 8 1 9 4 8 12 6 11 14 2 4 5 8 2 3
1 1
6
7 5
10 2 7 2 6 2 1
20 1
2 2 6 1 4 2 6 1 6 1 5 4 5 1 3 1 2 1 6 3
12 7
3 1 9 8 1 10 2 12 7 4 5 6
16 6
11 4 1 11 3 3 6 9 5 2 10 8 11 3 1 8
14 10
14 9 15 5 11 6 7 14 11 5 3 9 12 1
15 12
3 11 2 9 5 8 16 12 10 12 5 10 13 14 3
1 1
6
11 3
4 4 8 7 7 1 7 7 1 3 8
3 2
6 2 4
17 16
18 20 1 2 16 11 10 15 2 14 7 18 21 3 5 1 13
14 6
1 4 1 1 11 9 10 2 4 2 10 11 4 5
9 3
2 8 7 2 1 5 8 2 5
5 5
6 2 3 5 1
2 1
2 6
9 9
6 6 10 14 1 14 12 12 10
16 15
14 12 18 6 7 13 19 10 1 5 5 9 11 11 12 3
11 10
1 1 5 3 3 10 5 6 7 9 3
10 2
4 6 2 1 3 2 7 5 6 1
10 7
6 5 7 2 2 9 8 8 6 6
4 4
2 8 7 1
10 6
11 3 3 11 10 7 11 2 2 2
7 6
4 1 7 1 2 7 9
17 10
8 15 8 13 10 12 11 4 7 2 6 4 5 10 13 3 7
7 3
2 2 1 8 4 2 8
13 5
4 1 4 10 3 2 4 8 7 6 9 3 2
20 16
5 19 13 21 14 17 16 11 16 16 21 7 18 20 8 1 11 11 11 2
17 5
5 10 3 7 10 5 8 2 2 9 1 2 4 3 1 5 1
15 6
3 3 11 8 6 9 7 9 9 1 10 2 11 9 10
3 3
7 4 5
18 14
16 13 19 8 1 1 6 10 17 19 9 11 3 16 9 10 14 13
13 1
2 6 2 2 3 6 3 1 1 4 4 2 4
20 3
3 6 7 1 8 7 8 1 2 8 3 1 1 3 6 2 6 4 7 8
4 1
5 6 4 5
11 11
4 10 5 13 10 4 7 2 13 15 12
7 4
6 2 1 1 8 5 1
17 7
4 2 11 9 12 9 7 9 5 2 3 7 10 7 2 2 7
3 1
4 2 6
1 1
4
14 1
4 3 6 3 1 3 1 1 3 6 1 3 3 2
1 1
3
3 3
3 4 1
7 6
11 2 1 5 6 1 10
8 3
3 8 2 8 6 5 3 1
7 3
6 8 5 5 6 3 2
4 3
3 7 3 3
8 6
9 4 4 3 5 6 7 11
2 1
5 1
13 2
6 1 2 4 3 5 4 6 2 5 4 3 6
12 2
2 4 6 3 6 5 1 4 4 1 4 6
11 8
4 6 5 8 2 3 13 2 5 2 9
20 5
8 7 3 7 7 3 4 8 6 9 3 6 8 2 8 4 5 1 8 10
15 1
2 3 1 6 3 5 5 2 4 6 4 1 6 4 2
18 13
9 1 4 9 2 1 9 13 17 13 15 4 9 12 10 7 3 2
3 2
3 5 3
4 2
7 2 1 4
10 5
9 3 10 9 4 9 2 7 9 7
9 5
8 6 10 3 3 2 2 7 7
19 15
5 18 10 12 16 14 7 16 16 17 11 16 2 15 10 5 16 2 20
7 1
3 4 4 1 5 1 6
3 3
7 1 6
2 1
5 1
9 5
4 3 10 5 4 2 7 8 6
13 3
6 7 7 3 8 3 6 3 4 3 8 6 7
14 13
16 13 8 7 15 7 2 13 2 8 3 6 12 2
6 2
5 3 5 1 6 7
17 10
13 15 15 6 7 8 1 11 12 9 11 11 9 15 12 15 7
19 15
16 9 16 7 11 9 2 2 2 6 12 1 10 1 5 3 14 8 20
13 9
4 8 4 6 10 2 10 2 13 6 6 9 8
11 5
1 9 1 4 6 2 4 9 6 4 4
9 5
5 9 7 5 8 6 4 1 5
18 3
1 8 8 8 1 7 8 8 8 2 2 2 4 2 3 7 4 8
20 3
7 7 1 3 4 8 4 3 5 6 6 7 2 5 4 5 8 8 5 5
8 1
1 5 6 1 2 6 4 2
7 3
1 7 1 2 7 5 2
19 12
8 10 8 8 3 17 10 11 8 12 16 10 6 5 1 17 11 12 1
5 4
3 3 9 2 3
7 7
8 10 12 4 4 12 3
8 7
6 10 10 3 11 8 2 10
1 1
4
15 5
1 4 9 3 8 8 9 6 2 5 3 10 7 4 6
10 7
1 4 1 6 12 12 4 6 8 11
8 5
6 3 5 1 6 10 9 1
5 3
1 8 1 1 4
2 1
2 6
11 2
7 1 3 6 4 2 2 4 4 2 3
10 3
6 7 7 1 7 5 8 1 2 7
13 3
1 3 3 2 6 4 3 4 1 3 3 2 7
`

type Pair struct{ x, y int }

func expectedC(n, m int, arr []int) (int, int, []int) {
	a := make([]int, n)
	for i := range arr {
		a[i] = arr[i] - 1
	}
	cnt := make([]int, m)
	for _, v := range a {
		if v >= 0 && v < m {
			cnt[v]++
		}
	}
	per := n / m
	more := n % m
	id := make([]int, m)
	for i := range id {
		id[i] = i
	}
	sort.Slice(id, func(i, j int) bool { return cnt[id[i]] > cnt[id[j]] })
	tar := make([]int, m)
	for idx, v := range id {
		if idx < more {
			tar[v] = per + 1
		} else {
			tar[v] = per
		}
	}
	ans := 0
	for i := 0; i < m; i++ {
		for cnt[i] > tar[i] {
			idx := -1
			for j := 0; j < n; j++ {
				if a[j] == i {
					idx = j
					break
				}
			}
			for j := 0; j < m; j++ {
				if cnt[j] < tar[j] {
					a[idx] = j
					cnt[j]++
					cnt[i]--
					ans++
					break
				}
			}
		}
	}
	for i := 0; i < n; i++ {
		if a[i] < 0 || a[i] >= m {
			for j := 0; j < m; j++ {
				if cnt[j] < tar[j] {
					a[i] = j
					cnt[j]++
					ans++
					break
				}
			}
		}
	}
	out := make([]int, n)
	for i := range a {
		out[i] = a[i] + 1
	}
	return per, ans, out
}

func runCase(exe, input, exp string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp = strings.TrimSpace(exp)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	scan := bufio.NewScanner(strings.NewReader(testcasesCRaw))
	scan.Split(bufio.ScanLines)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(strings.TrimSpace(scan.Text()))
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		parts := strings.Fields(scan.Text())
		if len(parts) != 2 {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(parts[0])
		m, _ := strconv.Atoi(parts[1])
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		nums := strings.Fields(scan.Text())
		if len(nums) != n {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		arr := make([]int, n)
		for j, s := range nums {
			arr[j], _ = strconv.Atoi(s)
		}
		per, ans, outArr := expectedC(n, m, arr)
		input := fmt.Sprintf("%d %d\n%s\n", n, m, strings.Join(nums, " "))
		outStrs := make([]string, n)
		for j, v := range outArr {
			outStrs[j] = strconv.Itoa(v)
		}
		exp := fmt.Sprintf("%d %d\n%s\n", per, ans, strings.Join(outStrs, " "))
		if err := runCase(exe, input, exp); err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
