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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solve(n, m, k int64, p []int64) int64 {
	removed := int64(0)
	ops := int64(0)
	i := 0
	for i < len(p) {
		curr := (p[i] - removed - 1) / k
		cnt := 0
		for i < len(p) && (p[i]-removed-1)/k == curr {
			cnt++
			i++
		}
		removed += int64(cnt)
		ops++
	}
	return ops
}

func genTest(rng *rand.Rand) (string, int64) {
	n := int64(rng.Intn(100) + 1)
	m := int64(rng.Intn(int(n)) + 1)
	k := int64(rng.Intn(int(n)) + 1)
	// generate m distinct indices 1..n
	set := make(map[int64]struct{})
	for len(set) < int(m) {
		v := int64(rng.Intn(int(n)) + 1)
		set[v] = struct{}{}
	}
	p := make([]int64, 0, m)
	for v := range set {
		p = append(p, v)
	}
	// sort ascending
	for i := 0; i < len(p); i++ {
		for j := i + 1; j < len(p); j++ {
			if p[j] < p[i] {
				p[i], p[j] = p[j], p[i]
			}
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, k))
	for i, v := range p {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String(), solve(n, m, k, p)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
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
		if out != fmt.Sprint(expected) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\ninput:%s", i+1, expected, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
