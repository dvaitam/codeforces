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

type Env struct {
	w, h, idx int
}

func solve(cardW, cardH int, envs []Env) (int, []int) {
	filtered := make([]Env, 0)
	for _, e := range envs {
		if e.w > cardW && e.h > cardH {
			filtered = append(filtered, e)
		}
	}
	sort.Slice(filtered, func(i, j int) bool {
		if filtered[i].w == filtered[j].w {
			return filtered[i].h < filtered[j].h
		}
		return filtered[i].w < filtered[j].w
	})
	n := len(filtered)
	if n == 0 {
		return 0, nil
	}
	dp := make([]int, n)
	prev := make([]int, n)
	for i := range dp {
		dp[i] = 1
		prev[i] = -1
	}
	bestLen := 0
	bestIdx := -1
	for i := 0; i < n; i++ {
		for j := 0; j < i; j++ {
			if filtered[j].w < filtered[i].w && filtered[j].h < filtered[i].h {
				if dp[j]+1 > dp[i] {
					dp[i] = dp[j] + 1
					prev[i] = j
				}
			}
		}
		if dp[i] > bestLen {
			bestLen = dp[i]
			bestIdx = i
		}
	}
	if bestLen == 0 {
		return 0, nil
	}
	seq := make([]int, 0, bestLen)
	for bestIdx != -1 {
		seq = append(seq, filtered[bestIdx].idx)
		bestIdx = prev[bestIdx]
	}
	for i, j := 0, len(seq)-1; i < j; i, j = i+1, j-1 {
		seq[i], seq[j] = seq[j], seq[i]
	}
	return bestLen, seq
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	file, err := os.Open("testcasesD.txt")
	if err != nil {
		panic(err)
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
		var n, w, h int
		fmt.Sscanf(fields[0], "%d,%d,%d", &n, &w, &h)
		if len(fields)-1 < n {
			fmt.Printf("test %d invalid\n", idx)
			os.Exit(1)
		}
		envs := make([]Env, n)
		for i := 0; i < n; i++ {
			fmt.Sscanf(fields[i+1], "%d,%d", &envs[i].w, &envs[i].h)
			envs[i].idx = i + 1
		}
		expLen, _ := solve(w, h, envs)
		var buf bytes.Buffer
		fmt.Fprintf(&buf, "%d %d %d\n", n, w, h)
		for i := 0; i < n; i++ {
			fmt.Fprintf(&buf, "%d %d\n", envs[i].w, envs[i].h)
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
			fmt.Printf("Test %d empty output\n", idx)
			os.Exit(1)
		}
		gotLen, err := strconv.Atoi(outFields[0])
		if err != nil {
			fmt.Printf("Test %d invalid first number\n", idx)
			os.Exit(1)
		}
		if gotLen != expLen {
			fmt.Printf("Test %d failed: expected length %d got %d\n", idx, expLen, gotLen)
			os.Exit(1)
		}
		if gotLen == 0 {
			continue
		}
		nums := []int{}
		for _, t := range outFields[1:] {
			v, err := strconv.Atoi(t)
			if err != nil {
				continue
			}
			nums = append(nums, v)
		}
		if len(nums) != gotLen {
			fmt.Printf("Test %d failed: expected %d indices got %d\n", idx, gotLen, len(nums))
			os.Exit(1)
		}
		prevW := w
		prevH := h
		seen := make(map[int]bool)
		for _, id := range nums {
			if id < 1 || id > n {
				fmt.Printf("Test %d invalid index\n", idx)
				os.Exit(1)
			}
			if seen[id] {
				fmt.Printf("Test %d duplicate index\n", idx)
				os.Exit(1)
			}
			seen[id] = true
			e := envs[id-1]
			if e.w <= prevW || e.h <= prevH {
				fmt.Printf("Test %d invalid chain order\n", idx)
				os.Exit(1)
			}
			prevW = e.w
			prevH = e.h
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
