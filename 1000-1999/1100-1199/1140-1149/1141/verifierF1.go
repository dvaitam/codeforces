package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Run candidate solution
func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("execution error: %v", err)
	}
	return out.String(), nil
}

// Correct solver for F1 (N <= 50)
func solveCorrect(n int, a []int) int {
	type Interval struct{ l, r int }
	bySum := make(map[int][]Interval)
	for i := 0; i < n; i++ {
		sum := 0
		for j := i; j < n; j++ {
			sum += a[j]
			bySum[sum] = append(bySum[sum], Interval{i + 1, j + 1})
		}
	}

	maxK := 0
	for _, intervals := range bySum {
		// Greedy interval scheduling
		sort.Slice(intervals, func(i, j int) bool {
			return intervals[i].r < intervals[j].r
		})
		count := 0
		lastR := -1
		for _, iv := range intervals {
			if iv.l > lastR {
				count++
				lastR = iv.r
			}
		}
		if count > maxK {
			maxK = count
		}
	}
	return maxK
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(50) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		val := rng.Intn(201) - 100
		sb.WriteString(fmt.Sprintf("%d", val))
		if i+1 < n {
			sb.WriteByte(' ')
		}
	}
	sb.WriteByte('\n')
	return sb.String()
}

func verify(n int, a []int, output string, expectedK int) error {
	sc := bufio.NewScanner(strings.NewReader(output))
	sc.Split(bufio.ScanWords)

	if !sc.Scan() {
		return fmt.Errorf("no output produced")
	}
	kStr := sc.Text()
	k, err := strconv.Atoi(kStr)
	if err != nil {
		return fmt.Errorf("invalid k: %v", err)
	}

	if k != expectedK {
		return fmt.Errorf("expected k=%d, got k=%d", expectedK, k)
	}

	if k == 0 {
		return nil
	}

	type Seg struct{ l, r int }
	var segs []Seg

	for i := 0; i < k; i++ {
		if !sc.Scan() {
			return fmt.Errorf("expected l for segment %d", i+1)
		}
		lStr := sc.Text()
		if !sc.Scan() {
			return fmt.Errorf("expected r for segment %d", i+1)
		}
		rStr := sc.Text()
		l, err := strconv.Atoi(lStr)
		if err != nil {
			return fmt.Errorf("invalid l: %v", err)
		}
		r, err := strconv.Atoi(rStr)
		if err != nil {
			return fmt.Errorf("invalid r: %v", err)
		}

		if l < 1 || r > n || l > r {
			return fmt.Errorf("invalid segment boundaries: %d %d (n=%d)", l, r, n)
		}
		segs = append(segs, Seg{l, r})
	}

	// Check overlap
	sort.Slice(segs, func(i, j int) bool {
		return segs[i].l < segs[j].l
	})

	for i := 0; i < k-1; i++ {
		if segs[i].r >= segs[i+1].l {
			return fmt.Errorf("overlapping segments: [%d, %d] and [%d, %d]", segs[i].l, segs[i].r, segs[i+1].l, segs[i+1].r)
		}
	}

	// Check sums equality
	// Calculate sum of first segment
	targetSum := 0
	for i := segs[0].l - 1; i < segs[0].r; i++ {
		targetSum += a[i]
	}

	for i := 1; i < k; i++ {
		s := 0
		for j := segs[i].l - 1; j < segs[i].r; j++ {
			s += a[j]
		}
		if s != targetSum {
			return fmt.Errorf("segment sums mismatch: seg #1 sum=%d, seg %d sum=%d", targetSum, i+1, s)
		}
	}

	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF1.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 100; i++ {
		inputStr := genCase(rng)

		// Parse input to pass to solver
		var n int
		var a []int
		{
			sc := bufio.NewScanner(strings.NewReader(inputStr))
			sc.Split(bufio.ScanWords)
			sc.Scan()
			n, _ = strconv.Atoi(sc.Text())
			a = make([]int, n)
			for j := 0; j < n; j++ {
				sc.Scan()
				a[j], _ = strconv.Atoi(sc.Text())
			}
		}

		expectedK := solveCorrect(n, a)
		gotStr, err := run(candidate, inputStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Test %d failed execution: %v\nInput:\n%s", i+1, err, inputStr)
			os.Exit(1)
		}

		if err := verify(n, a, gotStr, expectedK); err != nil {
			fmt.Fprintf(os.Stderr, "Test %d failed.\nInput:\n%s\nError: %v\nOutput:\n%s\n", i+1, inputStr, err, gotStr)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
