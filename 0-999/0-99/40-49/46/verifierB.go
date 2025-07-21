package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

var sizes = []string{"S", "M", "L", "XL", "XXL"}

func expected(ns, nm, nl, nxl, nxxl int, req []string) []string {
	counts := []int{ns, nm, nl, nxl, nxxl}
	res := make([]string, len(req))
	for i, d := range req {
		var idx int
		switch d {
		case "S":
			idx = 0
		case "M":
			idx = 1
		case "L":
			idx = 2
		case "XL":
			idx = 3
		case "XXL":
			idx = 4
		}
		assigned := -1
		for dist := 0; assigned == -1 && dist < len(counts); dist++ {
			if idx+dist < len(counts) && counts[idx+dist] > 0 {
				assigned = idx + dist
				break
			}
			if dist > 0 && idx-dist >= 0 && counts[idx-dist] > 0 {
				assigned = idx - dist
				break
			}
		}
		counts[assigned]--
		res[i] = sizes[assigned]
	}
	return res
}

func generateCase(rng *rand.Rand) (string, []string) {
	k := rng.Intn(7) + 1
	counts := make([]int, 5)
	for i := 0; i < 5; i++ {
		counts[i] = rng.Intn(6)
	}
	sum := 0
	for _, c := range counts {
		sum += c
	}
	if sum < k {
		counts[0] += k - sum
	}
	req := make([]string, k)
	for i := 0; i < k; i++ {
		req[i] = sizes[rng.Intn(5)]
	}
	sb := strings.Builder{}
	fmt.Fprintf(&sb, "%d %d %d %d %d\n%d\n", counts[0], counts[1], counts[2], counts[3], counts[4], k)
	for i := 0; i < k; i++ {
		sb.WriteString(req[i])
		sb.WriteByte('\n')
	}
	expect := expected(counts[0], counts[1], counts[2], counts[3], counts[4], req)
	return sb.String(), expect
}

func runCase(exe, input string, expect []string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	lines := strings.Fields(out.String())
	if len(lines) != len(expect) {
		return fmt.Errorf("expected %d lines got %d", len(expect), len(lines))
	}
	for i, l := range lines {
		if l != expect[i] {
			return fmt.Errorf("line %d expected %s got %s", i+1, expect[i], l)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
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
