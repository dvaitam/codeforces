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

type item struct {
	s int64
	p int64
	q int64
}

func expectedD(a, b []int64) int64 {
	n := len(a)
	items := make([]item, n)
	var sum int64
	for i := 0; i < n; i++ {
		diff := a[i] - b[i]
		if diff < 0 {
			diff = -diff
		}
		sum += diff
		s := a[i] + b[i]
		items[i] = item{s: s, p: s - diff, q: s + diff}
	}
	sort.Slice(items, func(i, j int) bool { return items[i].s < items[j].s })
	pre := make([]int64, n)
	cur := int64(-1 << 60)
	for i := 0; i < n; i++ {
		if items[i].q > cur {
			cur = items[i].q
		}
		pre[i] = cur
	}
	suf := make([]int64, n)
	cur = int64(-1 << 60)
	for i := n - 1; i >= 0; i-- {
		if items[i].p > cur {
			cur = items[i].p
		}
		suf[i] = cur
	}
	var best int64
	for i := 0; i < n; i++ {
		if i > 0 {
			cand := items[i].p - pre[i-1]
			if cand > best {
				best = cand
			}
		}
		if i+1 < n {
			cand := suf[i+1] - items[i].q
			if cand > best {
				best = cand
			}
		}
	}
	if best < 0 {
		best = 0
	}
	return sum + best
}

func generateCase(rng *rand.Rand) (string, int64) {
	n := rng.Intn(10) + 2
	a := make([]int64, n)
	b := make([]int64, n)
	for i := 0; i < n; i++ {
		a[i] = rng.Int63n(1000) + 1
	}
	for i := 0; i < n; i++ {
		b[i] = rng.Int63n(1000) + 1
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	for i, v := range b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	exp := expectedD(a, b)
	return sb.String(), exp
}

func runCase(bin, input string, exp int64) error {
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
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(out.String())
	if len(fields) != 1 {
		return fmt.Errorf("expected single integer output")
	}
	got, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid integer: %v", err)
	}
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
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
