package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func expected(a []int) int {
	sort.Ints(a)
	prev := -1
	left, right := 0, len(a)-1
	eaten := 0
	for left <= right {
		idx := -1
		for i := left; i <= right; i++ {
			if a[i] > prev {
				idx = i
				break
			}
		}
		if idx == -1 {
			break
		}
		prev = a[idx]
		eaten++
		if idx == left {
			left++
		} else if idx == right {
			right--
		} else {
			copy(a[idx:right], a[idx+1:right+1])
			right--
		}
		if left <= right {
			right--
		}
	}
	return eaten
}

func generateCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(6) + 1
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(10) + 1
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	// copy arr for expectation as expected modifies slice
	cp := append([]int(nil), arr...)
	return sb.String(), expected(cp)
}

func runCase(bin string, input string, expected int) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if got != expected {
		return fmt.Errorf("expected %d got %d", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
