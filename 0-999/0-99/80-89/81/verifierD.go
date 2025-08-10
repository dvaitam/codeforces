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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// possible determines whether a valid arrangement exists.
func possible(n, m int, a []int) bool {
	low, high := 1, n
	bestK := -1
	for low <= high {
		mid := low + (high-low)/2
		if mid == 0 {
			low = 1
			continue
		}
		sum := 0
		for _, v := range a {
			sum += min(v, mid)
		}
		if sum >= n {
			bestK = mid
			high = mid - 1
		} else {
			low = mid + 1
		}
	}
	if bestK == -1 || 2*bestK > n {
		return false
	}
	return true
}

// validateOutput checks if the contestant output satisfies the problem constraints.
func validateOutput(n, m int, a []int, out string) error {
	out = strings.TrimSpace(out)
	if out == "-1" {
		if possible(n, m, a) {
			return fmt.Errorf("solution exists but got -1")
		}
		return nil
	}

	parts := strings.Fields(out)
	if len(parts) != n {
		return fmt.Errorf("expected %d numbers, got %d", n, len(parts))
	}
	used := make([]int, m+1)
	prev := -1
	first := -1
	for i, p := range parts {
		v, err := strconv.Atoi(p)
		if err != nil {
			return fmt.Errorf("output contains non-integer: %v", p)
		}
		if v < 1 || v > m {
			return fmt.Errorf("album %d out of range", v)
		}
		used[v]++
		if used[v] > a[v-1] {
			return fmt.Errorf("album %d used more than available", v)
		}
		if i == 0 {
			first = v
		} else if v == prev {
			return fmt.Errorf("adjacent photos from album %d", v)
		}
		prev = v
	}
	if n > 1 && prev == first {
		return fmt.Errorf("first and last photos from same album")
	}
	total := 0
	for i := 1; i <= m; i++ {
		total += used[i]
	}
	if total != n {
		return fmt.Errorf("expected %d photos, got %d", n, total)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesD.txt")
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
		if len(parts) < 2 {
			fmt.Printf("test %d: invalid line\n", idx)
			os.Exit(1)
		}
		nVal, _ := strconv.Atoi(parts[0])
		mVal, _ := strconv.Atoi(parts[1])
		if len(parts) != mVal+2 {
			fmt.Printf("test %d: wrong number of album sizes\n", idx)
			os.Exit(1)
		}
		arr := make([]int, mVal)
		for i := 0; i < mVal; i++ {
			v, _ := strconv.Atoi(parts[2+i])
			arr[i] = v
		}
		input := fmt.Sprintf("%d %d\n%s\n", nVal, mVal, strings.Join(parts[2:], " "))
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
		got := out.String()
		if err := validateOutput(nVal, mVal, arr, got); err != nil {
			fmt.Printf("test %d failed: %v\n", idx, err)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
