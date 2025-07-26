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

func solveD(a []int, m, k int) int {
	sort.Ints(a)
	queue := make([]int, 0, len(a))
	left := 0
	removed := 0
	for _, t := range a {
		limit := t - m + 1
		for left < len(queue) && queue[left] < limit {
			left++
		}
		queue = append(queue, t)
		if len(queue)-left >= k {
			queue = queue[:len(queue)-1]
			removed++
		}
	}
	return removed
}

func genCaseD(rng *rand.Rand) (string, int) {
	n := rng.Intn(20) + 1
	m := rng.Intn(10) + 1
	k := rng.Intn(n) + 1
	times := rng.Perm(100)[:n]
	for i := range times {
		times[i]++
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, m, k)
	for i, v := range times {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteString("\n")
	return sb.String(), solveD(times, m, k)
}

func runCaseD(bin, in string, exp int) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, buf.String())
	}
	gotStr := strings.TrimSpace(buf.String())
	got, err := strconv.Atoi(gotStr)
	if err != nil {
		return fmt.Errorf("invalid output: %v", err)
	}
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
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
		in, exp := genCaseD(rng)
		if err := runCaseD(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
