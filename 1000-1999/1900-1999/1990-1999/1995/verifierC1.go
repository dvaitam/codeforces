package main

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strings"
)

const eps = 1e-12

// embeddedTestcasesC1 matches the previous contents of testcasesC1.txt.
const embeddedTestcasesC1 = `2 4 8
3 5 4 1
4 6 9 6 2
1 5 
5 7 4 6 7 3
2 5 4
6 8 6 5 7 10 3
1 7 
3 9 8 4
6 6 6 7 5 6 7
6 5 2 8 5 2 8
2 6 4
6 3 6 8 4 2 7
4 8 9 8 10
5 4 7 9 5 8
2 6 9
6 1 2 8 6 7 4
4 1 10 1 7
1 5 
2 6 3
1 3 
6 6 1 4 1 1 7
5 1 3 2 10 10
2 2 8
2 1 9
5 7 2 9 3 4
2 7 7
4 1 7 4 7
1 10 
3 1 10 6
6 6 6 8 3 10 9
1 5 
1 2 
3 1 3 10
6 3 7 4 10 6 4
4 9 9 2 9
1 8 
1 9 
4 8 3 8 9
3 3 7 5
4 2 10 9 6
2 8 4
3 8 7 1
4 1 9 7 8
2 7 4
3 8 8 3
2 8 5
3 8 10 3
5 2 4 5 9 2
1 3 
3 1 6 10
2 8 7
6 4 7 8 4 3 7
1 6 
6 9 4 6 3 9 9
6 10 10 7 3 10 8
6 10 4 10 8 1 4
4 10 6 9 4
1 1 
1 3 
2 8 10
2 2 9
4 5 6 8 4
2 6 9
3 8 7 8
6 6 8 2 3 4 2
2 7 8
6 9 4 6 1 1 7
2 2 7
5 10 7 9 4 1
6 3 9 8 5 7 9
5 2 2 10 8 1
2 7 10
3 5 2 1
2 10 1
5 9 6 3 8 5
1 10 
4 5 4 7 7
6 9 6 4 10 5 2
1 3 
3 10 10 8
1 9 
4 7 6 7 1
1 7 
5 3 3 1 2 7
2 1 5
3 7 2 6
3 6 10 4
6 9 9 1 6 9 9
1 3 
4 1 5 6 3
3 2 10 9
4 2 5 10 2
5 2 10 3 8 1
5 4 6 2 9 4
2 4 7
2 10 8
3 3 5 6
5 4 4 8 6 3
3 9 1 2
`

func delta(x, y int64) int64 {
	if x == y {
		return 0
	}
	if x == 1 {
		return -60
	}
	val := math.Log(float64(x)) / math.Log(float64(y))
	val = math.Log2(val)
	nearest := math.Round(val)
	if math.Abs(val-nearest) <= eps {
		return int64(nearest)
	}
	return int64(math.Ceil(val))
}

func expectedC1(x int64, arr []int64) int64 {
	var ans, temp int64
	for _, y := range arr {
		if y == 1 && x != 1 {
			ans = -1
		}
		if ans != -1 {
			d := delta(x, y)
			if nt := temp + d; nt > 0 {
				temp = nt
			} else {
				temp = 0
			}
			ans += temp
		}
		x = y
	}
	return ans
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	lines := strings.Split(embeddedTestcasesC1, "\n")
	idx := 0
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		idx++
		fields := strings.Fields(line)
		if len(fields) < 2 {
			fmt.Printf("test %d: invalid line\n", idx)
			os.Exit(1)
		}
		var n int
		var x int64
		fmt.Sscan(fields[0], &n)
		fmt.Sscan(fields[1], &x)
		if len(fields) != 2+n-1 {
			fmt.Printf("test %d: invalid number of values\n", idx)
			os.Exit(1)
		}
		arr := make([]int64, n-1)
		for i := 0; i < n-1; i++ {
			fmt.Sscan(fields[2+i], &arr[i])
		}
		expect := expectedC1(x, arr)
		input := fmt.Sprintf("1\n%d %d\n", n, x)
		for i := 0; i < n-1; i++ {
			input += fmt.Sprintf("%d ", arr[i])
		}
		input += "\n"
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		cmd.Stdout = &out
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		res := strings.TrimSpace(out.String())
		var got int64
		if _, err := fmt.Sscan(res, &got); err != nil {
			fmt.Printf("test %d: failed to parse output %q\n", idx, res)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("test %d failed: expected %d got %d\n", idx, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
