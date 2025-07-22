package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const mod = 1000000007

func expected(n int, d int, x int64) []int {
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = i + 1
	}
	for i := 0; i < n; i++ {
		x = (x*7 + 13) % mod
		j := int(x % int64(i+1))
		a[i], a[j] = a[j], a[i]
	}
	b := make([]int, n)
	for i := 0; i < n; i++ {
		x = (x*7 + 13) % mod
		if x&1 == 1 {
			b[i] = 1
		}
	}
	c := make([]int, n)
	for i := 0; i < n; i++ {
		mx := 0
		for j := 0; j <= i; j++ {
			val := a[j] * b[i-j]
			if val > mx {
				mx = val
			}
		}
		c[i] = mx
	}
	return c
}

func generateCase(rng *rand.Rand) (string, []int) {
	n := rng.Intn(6) + 1
	d := rng.Intn(n) + 1
	x := rng.Int63n(mod)
	if x == 27777500 {
		x++
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, d, x))
	return sb.String(), expected(n, d, x)
}

func runCase(bin string, input string, expect []int) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	scanner := bufio.NewScanner(strings.NewReader(out.String()))
	scanner.Split(bufio.ScanWords)
	got := make([]int, 0, len(expect))
	for scanner.Scan() {
		v, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return fmt.Errorf("bad output: %v", err)
		}
		got = append(got, v)
	}
	if len(got) != len(expect) {
		return fmt.Errorf("expected %d numbers got %d", len(expect), len(got))
	}
	for i, v := range expect {
		if got[i] != v {
			return fmt.Errorf("pos %d expected %d got %d", i, v, got[i])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
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
