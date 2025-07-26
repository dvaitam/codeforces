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

func solveCase(n int, x, k int64, arr []int64) string {
	sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
	ans := int64(0)
	for _, val := range arr {
		t := val / x
		b := t - k
		L := b*x + 1
		R := (b + 1) * x
		if R > val {
			R = val
		}
		if L < 1 {
			L = 1
		}
		if R >= L {
			l := sort.Search(len(arr), func(i int) bool { return arr[i] >= L })
			r := sort.Search(len(arr), func(i int) bool { return arr[i] > R })
			ans += int64(r - l)
		}
	}
	return fmt.Sprintf("%d", ans)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	data, err := os.ReadFile("testcasesB.txt")
	if err != nil {
		fmt.Println("could not read testcasesB.txt:", err)
		os.Exit(1)
	}
	scanner := bufio.NewScanner(bytes.NewReader(data))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		fields := strings.Fields(line)
		if len(fields) < 4 {
			fmt.Println("bad test case")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(fields[0])
		x64, _ := strconv.ParseInt(fields[1], 10, 64)
		k64, _ := strconv.ParseInt(fields[2], 10, 64)
		if len(fields)-3 != n {
			fmt.Println("bad test case length")
			os.Exit(1)
		}
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			v, _ := strconv.ParseInt(fields[3+i], 10, 64)
			arr[i] = v
		}
		expected := solveCase(n, x64, k64, arr)
		input := fmt.Sprintf("%d %d %d\n%s\n", n, x64, k64, strings.Join(fields[3:], " "))
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != expected {
			fmt.Printf("test %d failed: expected %s got %s\n", idx, expected, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
