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

type op struct {
	t int
	a int
	x int
}

func solveOps(ops []op) []float64 {
	seq := []int{0}
	res := make([]float64, len(ops))
	for i, o := range ops {
		switch o.t {
		case 1:
			for j := 0; j < o.a && j < len(seq); j++ {
				seq[j] += o.x
			}
		case 2:
			seq = append(seq, o.x)
		case 3:
			if len(seq) > 1 {
				seq = seq[:len(seq)-1]
			}
		}
		sum := 0
		for _, v := range seq {
			sum += v
		}
		res[i] = float64(sum) / float64(len(seq))
	}
	return res
}

func generateCase(rng *rand.Rand) (string, []string) {
	n := rng.Intn(30) + 1
	ops := make([]op, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	seqLen := 1
	for i := 0; i < n; i++ {
		t := rng.Intn(3) + 1
		if seqLen == 1 && t == 3 {
			t = rng.Intn(2) + 1 // avoid deletion when size 1
		}
		ops[i].t = t
		switch t {
		case 1:
			ops[i].a = rng.Intn(seqLen) + 1
			ops[i].x = rng.Intn(201) - 100
			sb.WriteString(fmt.Sprintf("1 %d %d\n", ops[i].a, ops[i].x))
		case 2:
			ops[i].x = rng.Intn(201) - 100
			sb.WriteString(fmt.Sprintf("2 %d\n", ops[i].x))
			seqLen++
		case 3:
			sb.WriteString("3\n")
			seqLen--
		}
	}
	expect := solveOps(ops)
	var out strings.Builder
	for _, v := range expect {
		out.WriteString(fmt.Sprintf("%.8f\n", v))
	}
	return sb.String(), strings.Split(strings.TrimSuffix(out.String(), "\n"), "\n")
}

func runCase(bin string, input string, exp []string) error {
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
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != len(exp) {
		return fmt.Errorf("expected %d lines got %d", len(exp), len(lines))
	}
	for i, line := range lines {
		var val float64
		if _, err := fmt.Sscan(line, &val); err != nil {
			return fmt.Errorf("bad float on line %d: %v", i+1, err)
		}
		var expVal float64
		fmt.Sscan(exp[i], &expVal)
		if diff := val - expVal; diff < -1e-6 || diff > 1e-6 {
			return fmt.Errorf("line %d expected %.8f got %.8f", i+1, expVal, val)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, expLines := generateCase(rng)
		if err := runCase(bin, in, expLines); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
