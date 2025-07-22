package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func expectedAnswerC(arr []int64) []int64 {
	for j := 31; j >= 0; j-- {
		var subset []int64
		for _, v := range arr {
			if (v>>uint(j))&1 == 1 {
				subset = append(subset, v)
			}
		}
		if len(subset) == 0 {
			continue
		}
		common := subset[0]
		for k := 1; k < len(subset); k++ {
			common &= subset[k]
		}
		tz := 0
		for common&1 == 0 {
			tz++
			common >>= 1
		}
		if tz == j {
			return subset
		}
	}
	return []int64{arr[0]}
}

func generateCaseC(rng *rand.Rand) []int64 {
	n := rng.Intn(8) + 1
	nums := make([]int64, n)
	cur := int64(0)
	for i := 0; i < n; i++ {
		cur += int64(rng.Intn(20) + 1)
		nums[i] = cur
	}
	return nums
}

func runCaseC(bin string, arr []int64) error {
	input := fmt.Sprintf("%d\n", len(arr))
	for i, v := range arr {
		if i > 0 {
			input += " "
		}
		input += fmt.Sprint(v)
	}
	input += "\n"
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(strings.TrimSpace(out.String()))
	if len(fields) < 1 {
		return fmt.Errorf("no output")
	}
	k, err := strconv.Atoi(fields[0])
	if err != nil {
		return fmt.Errorf("bad k: %v", err)
	}
	if k != len(fields)-1 {
		return fmt.Errorf("expected %d numbers got %d", k, len(fields)-1)
	}
	expected := expectedAnswerC(arr)
	if k != len(expected) {
		return fmt.Errorf("expected length %d got %d", len(expected), k)
	}
	for i := 0; i < k; i++ {
		val, err := strconv.ParseInt(fields[i+1], 10, 64)
		if err != nil {
			return fmt.Errorf("bad number: %v", err)
		}
		if val != expected[i] {
			return fmt.Errorf("expected %v got %v", expected, fields[1:])
		}
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
		arr := generateCaseC(rng)
		if err := runCaseC(bin, arr); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%v\n", i+1, err, arr)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
