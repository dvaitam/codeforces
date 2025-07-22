package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func abs64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func exgcd(a, b int64) (int64, int64, int64) {
	if b == 0 {
		return 1, 0, a
	}
	x1, y1, d := exgcd(b, a%b)
	return y1, x1 - (a/b)*y1, d
}

func solveCase(n, s int64, arr []int64) string {
	var f [3]int64
	for _, v := range arr {
		f[v-3]++
	}
	const inf int64 = 1e18
	ans := inf
	var ansX, ansY, ansZ int64
	var step1, step2 int64
	var px, py, pdiv int64
	for i := int64(0); i*f[1] <= s; i++ {
		x0, y0, d := exgcd(f[0], f[2])
		px, py, pdiv = x0, y0, d
		remain := s - i*f[1]
		if pdiv != 0 && remain%pdiv == 0 {
			px *= remain / pdiv
			py *= remain / pdiv
			lcm := f[0] * f[2] / pdiv
			t1 := lcm / f[0]
			t2 := lcm / f[2]
			shift := -px / t1
			px += t1 * shift
			py -= t2 * shift
			if px < 0 {
				px += t1
				py -= t2
			}
			step1 = t1
			step2 = t2
			pain := func(r2 int64) {
				start := r2 - 5
				if start < 0 {
					start = 0
				}
				end := r2 + 5
				for r := start; r <= end; r++ {
					c0 := px + r*step1
					c1 := i
					c2 := py - r*step2
					if c0 <= c1 && c1 <= c2 {
						cost := abs64(c0*f[0]-c1*f[1]) + abs64(c1*f[1]-c2*f[2])
						if cost < ans {
							ans = cost
							ansX, ansY, ansZ = c0, c1, c2
						}
					}
				}
			}
			pain(0)
			pain((py - i) / step2)
			pain((i - px) / step1)
			if step1*f[0] != 0 {
				pain((i*f[1] - px*f[0]) / (step1 * f[0]))
			}
			if step2*f[2] != 0 {
				pain((i*f[1] - py*f[2]) / (step2 * f[2]))
			}
		}
	}
	if ans == inf {
		return "-1"
	}
	return fmt.Sprintf("%d %d %d", ansX, ansY, ansZ)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	data, err := os.ReadFile("testcasesE.txt")
	if err != nil {
		fmt.Println("could not read testcasesE.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("bad file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	expected := make([]string, t)
	for i := 0; i < t; i++ {
		scan.Scan()
		n64, _ := strconv.ParseInt(scan.Text(), 10, 64)
		scan.Scan()
		s64, _ := strconv.ParseInt(scan.Text(), 10, 64)
		arr := make([]int64, n64)
		for j := int64(0); j < n64; j++ {
			scan.Scan()
			val, _ := strconv.ParseInt(scan.Text(), 10, 64)
			arr[j] = val
		}
		expected[i] = solveCase(n64, s64, arr)
	}
	cmd := exec.Command(os.Args[1])
	cmd.Stdin = bytes.NewReader(data)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("execution failed:", err)
		os.Exit(1)
	}
	outScan := bufio.NewScanner(bytes.NewReader(out))
	outScan.Split(bufio.ScanWords)
	for i := 0; i < t; i++ {
		if !outScan.Scan() {
			fmt.Printf("missing output for test %d\n", i+1)
			os.Exit(1)
		}
		got := outScan.Text()
		if got != expected[i] {
			fmt.Printf("test %d failed: expected %s got %s\n", i+1, expected[i], got)
			os.Exit(1)
		}
	}
	if outScan.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed!")
}
