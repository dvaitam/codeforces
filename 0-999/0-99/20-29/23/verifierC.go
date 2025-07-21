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

type item struct {
	a   int
	s   int
	pos int
}

func solveC(q int, arr []item) string {
	items := make([]item, len(arr))
	copy(items, arr)
	sort.Slice(items, func(i, j int) bool { return items[i].a < items[j].a })
	var sum0, sum1 int64
	for i := 0; i < len(items); i++ {
		if i%2 == 0 {
			sum0 += int64(items[i].s)
		} else {
			sum1 += int64(items[i].s)
		}
	}
	var sb strings.Builder
	sb.WriteString("YES\n")
	if sum0 >= sum1 {
		for i := 0; i < len(items); i += 2 {
			fmt.Fprintf(&sb, "%d ", items[i].pos)
		}
	} else {
		for i := 1; i < len(items); i += 2 {
			fmt.Fprintf(&sb, "%d ", items[i].pos)
		}
		fmt.Fprintf(&sb, "%d ", items[len(items)-1].pos)
	}
	sb.WriteByte('\n')
	return sb.String()
}

func generateCase(rng *rand.Rand) (string, string) {
	q := rng.Intn(5) + 1
	cnt := 2*q - 1
	arr := make([]item, cnt)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", q)
	for i := 0; i < cnt; i++ {
		arr[i].a = rng.Intn(20)
		arr[i].s = rng.Intn(20)
		arr[i].pos = i + 1
		fmt.Fprintf(&sb, "%d %d\n", arr[i].a, arr[i].s)
	}
	return sb.String(), solveC(q, arr)
}

func runCase(bin, input, expected string) error {
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
	outStr := strings.TrimSpace(out.String())
	expStr := strings.TrimSpace(expected)
	if outStr != expStr {
		return fmt.Errorf("expected %q got %q", expStr, outStr)
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
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
