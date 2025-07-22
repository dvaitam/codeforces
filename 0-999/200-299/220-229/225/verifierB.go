package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func runBinary(bin string, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("timeout")
		}
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func kbonacci(k int, limit int64) []int64 {
	f := []int64{0, 1}
	for {
		sum := int64(0)
		for j := 1; j <= k && j < len(f); j++ {
			sum += f[len(f)-j]
		}
		f = append(f, sum)
		if sum > limit {
			break
		}
	}
	return f
}

func contains(seq []int64, v int64) bool {
	for _, x := range seq {
		if x == v {
			return true
		}
	}
	return false
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	r := rand.New(rand.NewSource(2))
	for t := 1; t <= 100; t++ {
		s := int64(r.Intn(1000000-1) + 1)
		k := r.Intn(8) + 2
		input := fmt.Sprintf("%d %d\n", s, k)
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("Test %d runtime error: %v\nInput:%s", t, err, input)
			return
		}
		fields := strings.Fields(out)
		if len(fields) < 1 {
			fmt.Printf("Test %d invalid output\nInput:%sGot:%s", t, input, out)
			return
		}
		m, errm := strconv.Atoi(fields[0])
		if errm != nil || m < 2 || m != len(fields)-1 {
			fmt.Printf("Test %d invalid m\nInput:%sGot:%s", t, input, out)
			return
		}
		nums := make([]int64, m)
		for i := 0; i < m; i++ {
			val, err := strconv.ParseInt(fields[i+1], 10, 64)
			if err != nil {
				fmt.Printf("Test %d invalid number\nInput:%sGot:%s", t, input, out)
				return
			}
			nums[i] = val
		}
		seq := kbonacci(k, s)
		set := make(map[int64]bool)
		var sum int64
		for _, v := range nums {
			if set[v] {
				fmt.Printf("Test %d numbers not distinct\n", t)
				return
			}
			set[v] = true
			if !contains(seq, v) {
				fmt.Printf("Test %d number %d not k-bonacci\n", t, v)
				return
			}
			sum += v
		}
		if sum != s {
			fmt.Printf("Test %d wrong sum expected %d got %d\n", t, s, sum)
			return
		}
	}
	fmt.Println("All tests passed")
}
