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

func expectedD(l, r int64, k int) (int64, []int64) {
	result := l
	S := []int64{l}
	maxK := 4
	if k < maxK {
		maxK = k
	}
	for K := 2; K <= maxK; K++ {
		for i := 0; i < 2; i++ {
			start := l + int64(i)
			if start+int64(K)-1 > r {
				continue
			}
			temp := int64(0)
			V := make([]int64, 0, K)
			for j := start; j < start+int64(K); j++ {
				V = append(V, j)
				temp ^= j
			}
			if temp < result {
				result = temp
				S = V
			}
		}
	}
	if k >= 3 {
		msb := 0
		for b := 62; b >= 0; b-- {
			if (l & (1 << uint(b))) != 0 {
				msb = b
				break
			}
		}
		A := (int64(1) << uint(msb+1)) | (int64(1) << uint(msb))
		B := (int64(1) << uint(msb+1)) | (l ^ (int64(1) << uint(msb)))
		if A <= r {
			result = 0
			S = []int64{l, A, B}
		}
	}
	return result, S
}

func runCase(bin string, l, r int64, k int) error {
	input := fmt.Sprintf("%d %d %d\n", l, r, k)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	scanner := bufio.NewScanner(strings.NewReader(out.String()))
	scanner.Split(bufio.ScanWords)
	var nums []int64
	for scanner.Scan() {
		var x int64
		fmt.Sscan(scanner.Text(), &x)
		nums = append(nums, x)
	}
	if len(nums) < 2 {
		return fmt.Errorf("bad output")
	}
	result := nums[0]
	cnt := nums[1]
	if int(cnt) != len(nums)-2 {
		return fmt.Errorf("expected %d numbers but got %d", cnt, len(nums)-2)
	}
	vals := nums[2:]
	expectVal, expectSet := expectedD(l, r, k)
	if result != expectVal {
		return fmt.Errorf("expected value %d got %d", expectVal, result)
	}
	if len(vals) != len(expectSet) {
		return fmt.Errorf("expected set %v got %v", expectSet, vals)
	}
	for i := range vals {
		if vals[i] != expectSet[i] {
			return fmt.Errorf("expected set %v got %v", expectSet, vals)
		}
	}
	return nil
}

func generateCase(rng *rand.Rand) (int64, int64, int) {
	l := rng.Int63n(1000) + 1
	r := l + rng.Int63n(20)
	k := rng.Intn(10) + 1
	if int64(k) > r-l+1 {
		k = int(r - l + 1)
	}
	return l, r, k
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	edges := []struct {
		l, r int64
		k    int
	}{
		{1, 1, 1},
		{1, 10, 3},
		{100, 120, 4},
	}
	for i, e := range edges {
		if err := runCase(bin, e.l, e.r, e.k); err != nil {
			fmt.Fprintf(os.Stderr, "edge case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}

	for i := 0; i < 100; i++ {
		l, r, k := generateCase(rng)
		if err := runCase(bin, l, r, k); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
