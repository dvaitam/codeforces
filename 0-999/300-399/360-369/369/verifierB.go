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

type caseB struct {
	n, k, l, r, sAll, sK int
	input                string
}

func generateCase(rng *rand.Rand) caseB {
	n := rng.Intn(10) + 1
	k := rng.Intn(n) + 1
	l := rng.Intn(11) // 0..10
	r := l + rng.Intn(11)
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = l + rng.Intn(r-l+1)
	}
	sAll := 0
	for _, v := range arr {
		sAll += v
	}
	sorted := append([]int(nil), arr...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] > sorted[j] })
	sK := 0
	for i := 0; i < k; i++ {
		sK += sorted[i]
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d %d %d %d\n", n, k, l, r, sAll, sK))
	return caseB{n, k, l, r, sAll, sK, sb.String()}
}

func runCase(bin string, c caseB) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(c.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	reader := strings.NewReader(out.String())
	arr := make([]int, c.n)
	for i := 0; i < c.n; i++ {
		if _, err := fmt.Fscan(reader, &arr[i]); err != nil {
			return fmt.Errorf("failed to read number %d: %v\n%s", i+1, err, out.String())
		}
	}
	for _, v := range arr {
		if v < c.l || v > c.r {
			return fmt.Errorf("value %d outside [%d,%d]", v, c.l, c.r)
		}
	}
	sum := 0
	for _, v := range arr {
		sum += v
	}
	if sum != c.sAll {
		return fmt.Errorf("sum %d != %d", sum, c.sAll)
	}
	sort.Slice(arr, func(i, j int) bool { return arr[i] > arr[j] })
	sK := 0
	for i := 0; i < c.k; i++ {
		sK += arr[i]
	}
	if sK != c.sK {
		return fmt.Errorf("top %d sum %d != %d", c.k, sK, c.sK)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		c := generateCase(rng)
		if err := runCase(bin, c); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, c.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
