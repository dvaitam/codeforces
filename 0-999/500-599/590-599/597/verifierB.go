package main

import (
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

type interval struct {
	l int
	r int
}

func generateCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(50) + 1 // 1..50
	arr := make([]interval, n)
	for i := range arr {
		l := rng.Intn(1000) + 1
		r := rng.Intn(1000) + 1
		if l > r {
			l, r = r, l
		}
		arr[i] = interval{l, r}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for _, it := range arr {
		sb.WriteString(fmt.Sprintf("%d %d\n", it.l, it.r))
	}
	return sb.String(), expected(arr)
}

func expected(arr []interval) int {
	sort.Slice(arr, func(i, j int) bool {
		if arr[i].r == arr[j].r {
			return arr[i].l < arr[j].l
		}
		return arr[i].r < arr[j].r
	})
	count := 0
	last := 0
	for _, it := range arr {
		if it.l > last {
			count++
			last = it.r
		}
	}
	return count
}

func runCase(exe, input string, expected int) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	got, err := strconv.Atoi(outStr)
	if err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if got != expected {
		return fmt.Errorf("expected %d got %d", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(exe, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
