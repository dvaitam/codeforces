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

func solveCase(n, k int, arr []int) int64 {
	groups := make(map[int][]int)
	for _, v := range arr {
		r := v % k
		groups[r] = append(groups[r], v)
	}
	oddGroups := 0
	for _, g := range groups {
		if len(g)%2 == 1 {
			oddGroups++
		}
	}
	if (n%2 == 0 && oddGroups > 0) || (n%2 == 1 && oddGroups > 1) {
		return -1
	}
	ops := int64(0)
	for _, g := range groups {
		sort.Ints(g)
		for i := 0; i+1 < len(g); i += 2 {
			ops += int64((g[i+1] - g[i]) / k)
		}
	}
	return ops
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 1
	k := rng.Intn(9) + 1
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(20) + 1
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d %d\n", n, k)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", arr[i])
	}
	sb.WriteByte('\n')
	expected := fmt.Sprintf("%d\n", solveCase(n, k, arr))
	return sb.String(), expected
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	if strings.TrimSpace(out.String()) != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %s got %s", expected, out.String())
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
