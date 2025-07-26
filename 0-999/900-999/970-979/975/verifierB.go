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

func solveCase(a []int64) string {
	var best int64
	for i := 0; i < 14; i++ {
		if a[i] == 0 {
			continue
		}
		b := make([]int64, 14)
		copy(b, a)
		x := b[i]
		b[i] = 0
		add := x / 14
		for j := 0; j < 14; j++ {
			b[j] += add
		}
		r := int(x % 14)
		for j := 1; j <= r; j++ {
			idx := (i + j) % 14
			b[idx]++
		}
		var cur int64
		for j := 0; j < 14; j++ {
			if b[j]%2 == 0 {
				cur += b[j]
			}
		}
		if cur > best {
			best = cur
		}
	}
	return fmt.Sprintf("%d\n", best)
}

func genCase(rng *rand.Rand) (string, string) {
	a := make([]int64, 14)
	for i := 0; i < 14; i++ {
		a[i] = rng.Int63n(30)
	}
	var sb strings.Builder
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	return sb.String(), solveCase(a)
}

func runCase(bin, in, exp string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != strings.TrimSpace(exp) {
		return fmt.Errorf("expected %s got %s", exp, out.String())
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
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
