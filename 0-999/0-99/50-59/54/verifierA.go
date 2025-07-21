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

func expected(n, k int, holidays []int) int {
	last := 0
	ans := 0
	idx := 0
	for idx < len(holidays) {
		h := holidays[idx]
		for last+k < h {
			last += k
			ans++
		}
		ans++
		last = h
		idx++
	}
	for last+k <= n {
		last += k
		ans++
	}
	return ans
}

func runCase(exe string, input string, expectedOut int) error {
	cmd := exec.Command(exe)
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
	if got != expectedOut {
		return fmt.Errorf("expected %d got %d", expectedOut, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(365) + 1
		k := rng.Intn(n) + 1
		c := rng.Intn(n + 1)
		daysMap := make(map[int]struct{})
		for len(daysMap) < c {
			d := rng.Intn(n) + 1
			daysMap[d] = struct{}{}
		}
		days := make([]int, 0, c)
		for d := range daysMap {
			days = append(days, d)
		}
		sort.Ints(days)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
		sb.WriteString(fmt.Sprintf("%d", len(days)))
		for _, d := range days {
			sb.WriteString(fmt.Sprintf(" %d", d))
		}
		sb.WriteString("\n")
		inp := sb.String()
		expectedOut := expected(n, k, days)
		if err := runCase(exe, inp, expectedOut); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, inp)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
