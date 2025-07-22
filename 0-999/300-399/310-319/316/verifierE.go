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

const mod = 1000000000

func solveCase(n int, a []int, ops [][]int) string {
	f := make([]int, n+1)
	if n >= 0 {
		f[0] = 1
	}
	if n >= 1 {
		f[1] = 1
	}
	for i := 2; i <= n; i++ {
		f[i] = f[i-1] + f[i-2]
		if f[i] >= mod {
			f[i] -= mod
		}
	}
	var out strings.Builder
	for _, op := range ops {
		t := op[0]
		switch t {
		case 1:
			x, v := op[1], op[2]
			if x >= 1 && x <= n {
				a[x] = v
			}
		case 2:
			l, r := op[1], op[2]
			sum := 0
			for x := l; x <= r; x++ {
				idx := x - l
				sum += f[idx] * a[x]
				if sum >= mod {
					sum %= mod
				}
			}
			fmt.Fprintf(&out, "%d\n", sum%mod)
		}
	}
	return out.String()
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	m := rng.Intn(5) + 1
	a := make([]int, n+1)
	for i := 1; i <= n; i++ {
		a[i] = rng.Intn(5)
	}
	ops := make([][]int, m)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", a[i])
	}
	sb.WriteByte('\n')
	for i := 0; i < m; i++ {
		if rng.Intn(2) == 0 {
			x := rng.Intn(n) + 1
			v := rng.Intn(5)
			ops[i] = []int{1, x, v}
			fmt.Fprintf(&sb, "1 %d %d\n", x, v)
		} else {
			l := rng.Intn(n) + 1
			r := rng.Intn(n-l+1) + l
			ops[i] = []int{2, l, r}
			fmt.Fprintf(&sb, "2 %d %d\n", l, r)
		}
	}
	expected := solveCase(n, append([]int(nil), a...), ops)
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
