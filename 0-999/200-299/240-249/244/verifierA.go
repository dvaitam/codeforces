package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCase struct {
	n, k int
	a    []int
}

func generateRandomCase(rng *rand.Rand) testCase {
	n := rng.Intn(30) + 1
	k := rng.Intn(30) + 1
	used := make(map[int]bool)
	a := make([]int, k)
	limit := n * k
	for i := 0; i < k; i++ {
		for {
			v := rng.Intn(limit) + 1
			if !used[v] {
				used[v] = true
				a[i] = v
				break
			}
		}
	}
	return testCase{n: n, k: k, a: a}
}

func edgeCases() []testCase {
	return []testCase{
		{n: 1, k: 1, a: []int{1}},
		{n: 1, k: 30, a: func() []int {
			x := make([]int, 30)
			for i := range x {
				x[i] = i + 1
			}
			return x
		}()},
		{n: 30, k: 1, a: []int{1}},
		{n: 30, k: 30, a: func() []int {
			x := make([]int, 30)
			for i := range x {
				x[i] = i + 1
			}
			return x
		}()},
		{n: 10, k: 5, a: []int{1, 2, 3, 4, 5}},
		{n: 5, k: 10, a: func() []int {
			x := make([]int, 10)
			for i := range x {
				x[i] = i + 1
			}
			return x
		}()},
		{n: 29, k: 3, a: []int{10, 20, 30}},
		{n: 3, k: 29, a: func() []int {
			x := make([]int, 29)
			for i := range x {
				x[i] = i + 1
			}
			return x
		}()},
		{n: 10, k: 20, a: func() []int {
			x := make([]int, 20)
			for i := range x {
				x[i] = i + 1
			}
			return x
		}()},
		{n: 30, k: 29, a: func() []int {
			x := make([]int, 29)
			for i := range x {
				x[i] = i + 1
			}
			return x
		}()},
	}
}

func runCase(exe string, tc testCase) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.k))
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}

	scanner := bufio.NewScanner(strings.NewReader(out.String()))
	scanner.Split(bufio.ScanWords)
	nums := make([]int, 0)
	for scanner.Scan() {
		var v int
		_, err := fmt.Sscan(scanner.Text(), &v)
		if err != nil {
			return fmt.Errorf("failed to parse int: %v", err)
		}
		nums = append(nums, v)
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner error: %v", err)
	}

	if len(nums) != tc.n*tc.k {
		return fmt.Errorf("expected %d numbers, got %d", tc.n*tc.k, len(nums))
	}
	used := make(map[int]bool)
	for _, v := range nums {
		if v < 1 || v > tc.n*tc.k {
			return fmt.Errorf("number %d out of range", v)
		}
		if used[v] {
			return fmt.Errorf("duplicate number %d", v)
		}
		used[v] = true
	}
	idx := 0
	for child := 0; child < tc.k; child++ {
		found := false
		for j := 0; j < tc.n; j++ {
			if nums[idx+j] == tc.a[child] {
				found = true
			}
		}
		if !found {
			return fmt.Errorf("child %d missing required segment %d", child+1, tc.a[child])
		}
		idx += tc.n
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	cases := edgeCases()
	for i, tc := range cases {
		if err := runCase(exe, tc); err != nil {
			fmt.Fprintf(os.Stderr, "edge case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	for i := 0; i < 100; i++ {
		tc := generateRandomCase(rng)
		if err := runCase(exe, tc); err != nil {
			fmt.Fprintf(os.Stderr, "random case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
