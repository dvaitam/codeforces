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

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	f, err := os.Open("testcasesB.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		fields := strings.Fields(line)
		if len(fields) < 1 {
			fmt.Printf("test %d invalid\n", idx)
			os.Exit(1)
		}
		var d, sumTime int
		fmt.Sscanf(fields[0], "%d,%d", &d, &sumTime)
		if len(fields)-1 < d {
			fmt.Printf("test %d invalid pair count\n", idx)
			os.Exit(1)
		}
		mins := make([]int, d)
		maxs := make([]int, d)
		for i := 0; i < d; i++ {
			fmt.Sscanf(fields[i+1], "%d,%d", &mins[i], &maxs[i])
		}
		var buf bytes.Buffer
		fmt.Fprintf(&buf, "%d %d\n", d, sumTime)
		for i := 0; i < d; i++ {
			fmt.Fprintf(&buf, "%d %d\n", mins[i], maxs[i])
		}
		cmd := exec.Command(binary)
		cmd.Stdin = bytes.NewReader(buf.Bytes())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", idx, err)
			os.Exit(1)
		}
		outFields := strings.Fields(strings.TrimSpace(string(out)))
		if len(outFields) == 0 {
			fmt.Printf("Test %d: empty output\n", idx)
			os.Exit(1)
		}
		minSum := 0
		maxSum := 0
		for i := 0; i < d; i++ {
			minSum += mins[i]
			maxSum += maxs[i]
		}
		feasible := !(sumTime < minSum || sumTime > maxSum)
		if !feasible {
			if strings.ToUpper(outFields[0]) != "NO" {
				fmt.Printf("Test %d failed: expected NO got %s\n", idx, strings.Join(outFields, " "))
				os.Exit(1)
			}
			continue
		}
		if strings.ToUpper(outFields[0]) != "YES" {
			fmt.Printf("Test %d failed: expected YES got %s\n", idx, strings.Join(outFields, " "))
			os.Exit(1)
		}
		nums := []int{}
		for _, t := range outFields[1:] {
			v, err := strconv.Atoi(t)
			if err != nil {
				continue
			}
			nums = append(nums, v)
		}
		if len(nums) != d {
			fmt.Printf("Test %d failed: expected %d numbers got %d\n", idx, d, len(nums))
			os.Exit(1)
		}
		total := 0
		for i := 0; i < d; i++ {
			v := nums[i]
			if v < mins[i] || v > maxs[i] {
				fmt.Printf("Test %d failed: value %d out of range\n", idx, v)
				os.Exit(1)
			}
			total += v
		}
		if total != sumTime {
			fmt.Printf("Test %d failed: sum mismatch expected %d got %d\n", idx, sumTime, total)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
