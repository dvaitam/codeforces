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

func expected(nums []int) string {
	freq := make(map[int]int)
	for _, v := range nums {
		freq[v]++
	}
	leg := 0
	for v, c := range freq {
		if c >= 4 {
			leg = v
			break
		}
	}
	if leg == 0 {
		return "Alien"
	}
	rem := make([]int, 0, 2)
	removed := 0
	for _, v := range nums {
		if v == leg && removed < 4 {
			removed++
		} else {
			rem = append(rem, v)
		}
	}
	if len(rem) != 2 {
		return "Alien"
	}
	if rem[0] == rem[1] {
		return "Elephant"
	}
	return "Bear"
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesA.txt")
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
		parts := strings.Fields(line)
		if len(parts) != 6 {
			fmt.Fprintf(os.Stderr, "test %d malformed\n", idx)
			os.Exit(1)
		}
		nums := make([]int, 6)
		for i := 0; i < 6; i++ {
			v, _ := strconv.Atoi(parts[i])
			nums[i] = v
		}
		expect := expected(nums)
		input := line + "\n"
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
		if got != expect {
			fmt.Printf("test %d failed: expected %s got %s\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
