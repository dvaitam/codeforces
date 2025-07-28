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

type Monster struct {
	h int64
	p int64
}

func expected(n int, k int64, h []int64, p []int64) string {
	monsters := make([]Monster, n)
	for i := 0; i < n; i++ {
		monsters[i] = Monster{h: h[i], p: p[i]}
	}
	sort.Slice(monsters, func(i, j int) bool { return monsters[i].h < monsters[j].h })
	suf := make([]int64, n)
	suf[n-1] = monsters[n-1].p
	for i := n - 2; i >= 0; i-- {
		if monsters[i].p < suf[i+1] {
			suf[i] = monsters[i].p
		} else {
			suf[i] = suf[i+1]
		}
	}
	damage := int64(0)
	cur := k
	idx := 0
	for cur > 0 {
		damage += cur
		for idx < n && monsters[idx].h <= damage {
			idx++
		}
		if idx == n {
			break
		}
		cur -= suf[idx]
	}
	if idx == n {
		return "YES"
	}
	return "NO"
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 1
	k := int64(rng.Intn(20) + 1)
	h := make([]int64, n)
	p := make([]int64, n)
	for i := 0; i < n; i++ {
		h[i] = int64(rng.Intn(50) + 1)
	}
	for i := 0; i < n; i++ {
		p[i] = int64(rng.Intn(20) + 1)
	}
	var input strings.Builder
	input.WriteString("1\n")
	input.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i := 0; i < n; i++ {
		if i > 0 {
			input.WriteByte(' ')
		}
		input.WriteString(fmt.Sprintf("%d", h[i]))
	}
	input.WriteString("\n")
	for i := 0; i < n; i++ {
		if i > 0 {
			input.WriteByte(' ')
		}
		input.WriteString(fmt.Sprintf("%d", p[i]))
	}
	input.WriteString("\n")
	return input.String(), expected(n, k, h, p)
}

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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		got, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
