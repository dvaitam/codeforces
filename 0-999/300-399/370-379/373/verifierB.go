package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func run(bin string, input string) (string, error) {
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
	return strings.TrimSpace(out.String()), nil
}

func digits(n uint64) uint64 {
	var d uint64
	for n > 0 {
		d++
		n /= 10
	}
	return d
}

func expected(w, m, k uint64) uint64 {
	rem := w
	curr := m
	var ans uint64
	for rem > 0 {
		d := digits(curr)
		costPer := d * k
		if rem < costPer {
			break
		}
		// determine end of this digit block
		pow10 := uint64(1)
		for i := uint64(0); i < d; i++ {
			pow10 *= 10
		}
		end := pow10 - 1
		count := end - curr + 1
		maxTake := rem / costPer
		take := count
		if maxTake < count {
			take = maxTake
		}
		ans += take
		rem -= take * costPer
		if take < count {
			break
		}
		curr = end + 1
	}
	return ans
}

func runCase(bin string, w, m, k uint64) error {
	input := fmt.Sprintf("%d %d %d\n", w, m, k)
	out, err := run(bin, input)
	if err != nil {
		return err
	}
	var got uint64
	if _, err := fmt.Sscan(out, &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	exp := expected(w, m, k)
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(99)
	total := 200
	for i := 0; i < 150; i++ {
		w := uint64(rand.Intn(100) + 1)
		m := uint64(rand.Intn(50) + 1)
		k := uint64(rand.Intn(5) + 1)
		if err := runCase(bin, w, m, k); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v (w=%d m=%d k=%d)\n", i+1, err, w, m, k)
			os.Exit(1)
		}
	}
	// additional big number cases
	big := []struct{ w, m, k uint64 }{
		{1e16, 1, 1},
		{1e16, 1e15, 1},
		{1e16, 1e16 - 1e3, 123456789},
		{9999999999999999, 123456789012345, 1000000000},
		{10000000000000000, 9999999999999999, 1},
	}
	for i, c := range big {
		if err := runCase(bin, c.w, c.m, c.k); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v (w=%d m=%d k=%d)\n", 150+i+1, err, c.w, c.m, c.k)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", total)
}
