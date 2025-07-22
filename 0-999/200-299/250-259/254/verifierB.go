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
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return out.String(), nil
}

var monthDays = []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

func expected(events [][4]int) int {
	const offset = 500
	const maxDays = 1000
	load := make([]int, maxDays)
	prefix := make([]int, 13)
	for i := 1; i <= 12; i++ {
		prefix[i] = prefix[i-1] + monthDays[i-1]
	}
	for _, e := range events {
		m, d, p, t := e[0], e[1], e[2], e[3]
		doy := prefix[m-1] + (d - 1)
		start := doy - t + offset
		end := doy - 1 + offset
		if start < 0 {
			start = 0
		}
		if end >= maxDays {
			end = maxDays - 1
		}
		for i := start; i <= end; i++ {
			load[i] += p
		}
	}
	ans := 0
	for _, v := range load {
		if v > ans {
			ans = v
		}
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, [][4]int) {
	n := rng.Intn(10) + 1
	events := make([][4]int, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		m := rng.Intn(12) + 1
		d := rng.Intn(monthDays[m-1]) + 1
		p := rng.Intn(10) + 1
		t := rng.Intn(10)
		events[i] = [4]int{m, d, p, t}
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", m, d, p, t))
	}
	return sb.String(), events
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, events := generateCase(rng)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		var got int
		if _, err := fmt.Fscan(strings.NewReader(out), &got); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed to parse output: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		exp := expected(events)
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\ninput:\n%soutput:\n%s", i+1, exp, got, input, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
