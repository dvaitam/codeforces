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

type operation struct {
	typ int
	l   int
	r   int
	x   int
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	m := rng.Intn(20) + 1
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(20) + 1
	}
	initial := make([]int, n)
	copy(initial, arr)
	var ops []operation
	var outputs []int
	for i := 0; i < m; i++ {
		typ := rng.Intn(3) + 1
		switch typ {
		case 1:
			l := rng.Intn(n) + 1
			r := rng.Intn(n-l+1) + l
			sum := 0
			for j := l - 1; j < r; j++ {
				sum += arr[j]
			}
			outputs = append(outputs, sum)
			ops = append(ops, operation{typ: 1, l: l, r: r})
		case 2:
			l := rng.Intn(n) + 1
			r := rng.Intn(n-l+1) + l
			x := rng.Intn(20) + 1
			for j := l - 1; j < r; j++ {
				arr[j] %= x
			}
			ops = append(ops, operation{typ: 2, l: l, r: r, x: x})
		case 3:
			k := rng.Intn(n) + 1
			x := rng.Intn(20) + 1
			arr[k-1] = x
			ops = append(ops, operation{typ: 3, l: k, x: x})
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i, v := range initial {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	for _, op := range ops {
		switch op.typ {
		case 1:
			sb.WriteString(fmt.Sprintf("1 %d %d\n", op.l, op.r))
		case 2:
			sb.WriteString(fmt.Sprintf("2 %d %d %d\n", op.l, op.r, op.x))
		case 3:
			sb.WriteString(fmt.Sprintf("3 %d %d\n", op.l, op.x))
		}
	}
	var out strings.Builder
	for i, v := range outputs {
		if i > 0 {
			out.WriteByte('\n')
		}
		out.WriteString(fmt.Sprintf("%d", v))
	}
	out.WriteByte('\n')
	return sb.String(), out.String()
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
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if got != exp {
		return fmt.Errorf("expected '%s' got '%s'", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
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
