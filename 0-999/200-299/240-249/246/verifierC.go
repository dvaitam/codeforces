package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func runCase(bin string, n int, k int64, arr []int) error {
	var input strings.Builder
	input.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i, v := range arr {
		if i > 0 {
			input.WriteByte(' ')
		}
		input.WriteString(fmt.Sprint(v))
	}
	input.WriteByte('\n')

	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}

	scanner := bufio.NewScanner(strings.NewReader(out.String()))
	scanner.Split(bufio.ScanWords)

	arrMap := make(map[int]int)
	for _, v := range arr {
		arrMap[v]++
	}

	seenSums := make(map[int]bool)

	for i := int64(0); i < k; i++ {
		if !scanner.Scan() {
			return fmt.Errorf("expected %d lines, got %d", k, i)
		}
		c, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return fmt.Errorf("invalid number of elements at line %d: %v", i+1, err)
		}
		if c <= 0 || c > n {
			return fmt.Errorf("invalid subset size %d at line %d", c, i+1)
		}

		subsetSum := 0
		used := make(map[int]int)
		for j := 0; j < c; j++ {
			if !scanner.Scan() {
				return fmt.Errorf("expected %d elements at line %d, got %d", c, i+1, j)
			}
			val, err := strconv.Atoi(scanner.Text())
			if err != nil {
				return fmt.Errorf("invalid element at line %d: %v", i+1, err)
			}
			used[val]++
			if used[val] > arrMap[val] {
				return fmt.Errorf("element %d used too many times or not in array at line %d", val, i+1)
			}
			subsetSum += val
		}

		if seenSums[subsetSum] {
			return fmt.Errorf("duplicate sum %d at line %d", subsetSum, i+1)
		}
		seenSums[subsetSum] = true
	}

	if scanner.Scan() {
		return fmt.Errorf("extra output after %d lines", k)
	}

	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(4) + 1
		maxK := n * (n + 1) / 2
		k := int64(rng.Intn(maxK) + 1)
		arr := make([]int, n)
		usedValues := make(map[int]bool)
		for j := 0; j < n; j++ {
			for {
				val := rng.Intn(100) + 1
				if !usedValues[val] {
					usedValues[val] = true
					arr[j] = val
					break
				}
			}
		}
		if err := runCase(bin, n, k, arr); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed for n=%d, k=%d, arr=%v: %v\n", i+1, n, k, arr, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
