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

const testcasesRaw = `100
7 4
2 2 6 3 5 5 10
2 39
1 10
6 11
7 7 9 6 9 8
5 18
1 1 6 8 6
4 28
9 3 9 3
2 15
1 3
3 12
3 9 9
3 33
9 3 8
7 27
9 6 10 6 6 8 3
7 26
8 9 4 8 5 8 9
5 23
8 8 6 10 9
6 30
8 4 6 3 10 5
7 31
5 5 9 9 9 9 10
5 27
5 4 8 9 6
6 40
2 6 1 4 2 1
5 4
5 10 4 2 9
2 18
4 4
1 28
1
1 24
6
2 16
1 2
1 5
1
1 2
6
3 9
3 3 9
6 1
7 10 1 4 3 1
1 23
10
6 8
5 6 8 1 5 8
5 39
1 5 7 10 3
4 15
2 6 2 1
4 9
9 10 7 8
5 21
3 6 5 5 10
4 2
9 3 1 5
1 9
3
2 7
8 4
5 3
4 4 8 2 5
1 38
4
5 40
6 5 7 5 9
7 1
3 1 7 7 3 2 9
6 6
4 2 2 1 3 4
1 14
1
5 30
8 5 9 7 4
6 14
7 7 9 1 10 10
1 27
9
5 12
2 8 6 1 9
1 40
6
3 24
5 1 7
1 7
5
2 2
8 1
4 32
8 4 10 10
1 1
5
1 24
5
6 5
4 8 4 2 10 6
4 30
3 6 7 2
3 8
2 2 10
7 22
7 4 2 1 10 8 1
6 32
5 6 8 3 6 5
4 34
8 7 8 5
4 15
3 8 10 5
5 28
2 10 10 2 2
3 12
9 3 7
1 6
1
2 19
7 4
6 22
8 3 9 5 2 3
5 28
2 6 9 4 9
3 11
3 8 4
4 23
10 3 8 8
6 2
10 7 3 7 9 1
4 18
7 5 7 8
3 36
6 2 4
5 40
4 7 7 1 6
4 34
8 3 2 1
4 14
10 10 7 4
1 25
9
7 13
5 10 10 4 8 10 3
1 40
7
4 17
9 10 3 8
6 14
2 6 1 8 9 2
7 38
8 6 8 5 9 8 1
1 40
6
2 26
5 3
1 11
8
4 30
5 3 1 5
5 30
1 6 1 9 7
5 29
4 5 8 3 8
6 35
5 2 5 6 5 6
6 20
7 9 2 9 4 7
5 34
3 9 2 5 1
2 30
9 4
5 18
1 2 2 7 6
2 21
6 2
3 30
6 3 8
4 19
8 3 8 4
3 21
3 2 4
4 13
6 3 6 3
7 9
4 5 9 7 7 6 5
6 39
9 10 6 7 5 9
5 5
6 5 7 8 3
3 23
8 8 2
2 21
7 3
1 7
6
2 23
2 7`

func expected(times []int, M int) []int {
	res := make([]int, len(times))
	prev := make([]int, 0, len(times))
	sum := 0
	for i, t := range times {
		need := sum + t - M
		if need <= 0 {
			res[i] = 0
		} else {
			tmp := append([]int(nil), prev...)
			sort.Sort(sort.Reverse(sort.IntSlice(tmp)))
			removed := 0
			acc := 0
			for _, v := range tmp {
				acc += v
				removed++
				if acc >= need {
					break
				}
			}
			res[i] = removed
		}
		prev = append(prev, t)
		sum += t
	}
	return res
}

func runCase(bin string, n, M int, times []int) error {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, M)
	for i, v := range times {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	input := sb.String()
	cmd := exec.Command(bin)
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	expectedRes := expected(times, M)
	outs := strings.Fields(strings.TrimSpace(out.String()))
	if len(outs) != len(times) {
		return fmt.Errorf("expected %d numbers, got %d", len(times), len(outs))
	}
	for i, tok := range outs {
		val, err := strconv.Atoi(tok)
		if err != nil || val != expectedRes[i] {
			return fmt.Errorf("position %d expected %d got %s", i, expectedRes[i], tok)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierC1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	sc := bufio.NewScanner(strings.NewReader(testcasesRaw))
	sc.Split(bufio.ScanWords)
	if !sc.Scan() {
		fmt.Println("invalid test data")
		os.Exit(1)
	}
	t, err := strconv.Atoi(sc.Text())
	if err != nil {
		fmt.Println("invalid test data:", err)
		os.Exit(1)
	}
	for i := 0; i < t; i++ {
		if !sc.Scan() {
			fmt.Println("invalid test data")
			os.Exit(1)
		}
		n, err := strconv.Atoi(sc.Text())
		if err != nil {
			fmt.Println("invalid n in test data:", err)
			os.Exit(1)
		}
		if !sc.Scan() {
			fmt.Println("invalid test data")
			os.Exit(1)
		}
		M, err := strconv.Atoi(sc.Text())
		if err != nil {
			fmt.Println("invalid M in test data:", err)
			os.Exit(1)
		}
		times := make([]int, n)
		for j := 0; j < n; j++ {
			if !sc.Scan() {
				fmt.Println("invalid test data")
				os.Exit(1)
			}
			val, err := strconv.Atoi(sc.Text())
			if err != nil {
				fmt.Println("invalid time in test data:", err)
				os.Exit(1)
			}
			times[j] = val
		}
		if err := runCase(bin, n, M, times); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
