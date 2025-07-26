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

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solve(a []int64) string {
	n := len(a)
	sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
	dupCount := 0
	dupIdx := -1
	for i := 1; i < n; i++ {
		if a[i] == a[i-1] {
			dupCount++
			dupIdx = i
		}
	}
	if dupCount > 1 {
		return "cslnb"
	}
	if dupCount == 1 {
		v := a[dupIdx]
		if v == 0 {
			return "cslnb"
		}
		target := v - 1
		idx := sort.Search(len(a), func(i int) bool { return a[i] >= target })
		if idx < n && a[idx] == target {
			return "cslnb"
		}
	}
	var moves int64
	for i := 0; i < n; i++ {
		moves += a[i] - int64(i)
	}
	if moves%2 == 1 {
		return "sjfnb"
	}
	return "cslnb"
}

func genTest(rng *rand.Rand) (string, string) {
	n := rng.Intn(4) + 1
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		a[i] = int64(rng.Intn(10))
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", a[i]))
	}
	sb.WriteByte('\n')
	return sb.String(), solve(a)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, expected := genTest(rng)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		if out != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", i+1, expected, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
