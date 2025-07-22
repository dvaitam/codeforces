package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func solve(n, a, d int, ts, vs []int) []float64 {
	tMax := math.Sqrt(float64(2*d) / float64(a))
	res := make([]float64, n)
	cur := 0.0
	for i := 0; i < n; i++ {
		v := float64(vs[i])
		t1 := v / float64(a)
		d1 := t1 * v / 2
		dt := 0.0
		if d1 <= float64(d) {
			dt = (float64(d)-d1)/v + t1
		} else {
			dt = tMax
		}
		arrive := float64(ts[i]) + dt
		if arrive < cur {
			arrive = cur
		}
		res[i] = arrive
		cur = arrive
	}
	return res
}

func generateCase(rng *rand.Rand) (string, []float64) {
	n := rng.Intn(5) + 1
	a := rng.Intn(9) + 1
	d := rng.Intn(30) + 1
	ts := make([]int, n)
	vs := make([]int, n)
	cur := 0
	for i := 0; i < n; i++ {
		cur += rng.Intn(5) + 1
		ts[i] = cur
		vs[i] = rng.Intn(20) + 1
	}
	exp := solve(n, a, d, ts, vs)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, a, d))
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", ts[i], vs[i]))
	}
	return sb.String(), exp
}

func run(bin, input string) ([]float64, error) {
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
		return nil, fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	fields := strings.Fields(out.String())
	res := make([]float64, len(fields))
	for i, f := range fields {
		v, err := strconv.ParseFloat(f, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid float on line %d", i+1)
		}
		res[i] = v
	}
	return res, nil
}

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if len(out) != len(exp) {
			fmt.Fprintf(os.Stderr, "case %d wrong number of lines\ninput:\n%s", i+1, in)
			os.Exit(1)
		}
		for j := range exp {
			if math.Abs(out[j]-exp[j]) > 1e-3 {
				fmt.Fprintf(os.Stderr, "case %d mismatch on line %d: expected %.4f got %.4f\ninput:\n%s", i+1, j+1, exp[j], out[j], in)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
