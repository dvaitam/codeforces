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

func expected(arr []int64) []int {
	n := len(arr)
	sum := int64(0)
	for _, v := range arr {
		sum += v
	}
	max1, idx1 := arr[0], 0
	for i := 1; i < n; i++ {
		if arr[i] > max1 {
			max1 = arr[i]
			idx1 = i
		}
	}
	max2 := int64(-1)
	for i := 0; i < n; i++ {
		if i == idx1 {
			continue
		}
		if arr[i] > max2 {
			max2 = arr[i]
		}
	}
	res := []int{}
	for i := 0; i < n; i++ {
		if i == idx1 {
			if sum-arr[i]-max2 == max2 {
				res = append(res, i+1)
			}
		} else {
			if sum-arr[i]-max1 == max1 {
				res = append(res, i+1)
			}
		}
	}
	return res
}

func runCase(exe string, input string, exp []int, n int) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(out.String())
	if len(fields) == 0 {
		return fmt.Errorf("no output")
	}
	k, err := strconv.Atoi(fields[0])
	if err != nil {
		return fmt.Errorf("invalid count: %q", fields[0])
	}
	if k != len(exp) {
		return fmt.Errorf("expected count %d got %d", len(exp), k)
	}
	if len(fields)-1 != k {
		return fmt.Errorf("expected %d indices, got %d", k, len(fields)-1)
	}
	seen := make(map[int]bool)
	for i := 0; i < k; i++ {
		idx, err := strconv.Atoi(fields[i+1])
		if err != nil {
			return fmt.Errorf("invalid index %q", fields[i+1])
		}
		if idx < 1 || idx > n {
			return fmt.Errorf("index out of range %d", idx)
		}
		if seen[idx] {
			return fmt.Errorf("duplicate index %d", idx)
		}
		seen[idx] = true
		valid := false
		for _, e := range exp {
			if e == idx {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("unexpected index %d", idx)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for tcase := 0; tcase < 100; tcase++ {
		n := rng.Intn(199) + 2
		arr := make([]int64, n)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", n)
		for i := 0; i < n; i++ {
			arr[i] = rng.Int63n(1e6) + 1
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(arr[i], 10))
		}
		sb.WriteByte('\n')
		input := sb.String()
		exp := expected(arr)
		if err := runCase(exe, input, exp, n); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", tcase+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
