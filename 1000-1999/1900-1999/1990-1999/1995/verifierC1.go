package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strings"
)

const eps = 1e-12

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
	file, err := os.Open("testcasesC1.txt")
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
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
