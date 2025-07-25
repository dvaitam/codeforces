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

const H = 40

type query struct {
	typ int
	v   int
}

type testCase struct {
	q     int
	qs    []query
	input string
}

func solve(qs []query) []float64 {
	n := len(qs) + 2
	dp := make([]float64, n*H)
	mul := make([]float64, n*H)
	p := make([]int, n)
	var path [H]int
	dp[0*H+0] = 1.0
	p[0] = -1
	for i := 0; i < n; i++ {
		for j := 0; j < H; j++ {
			mul[i*H+j] = 1.0
		}
	}
	cur := 1
	var outputs []float64
	for _, qu := range qs {
		if qu.typ == 1 {
			v := qu.v
			p[cur] = v
			dp[cur*H+0] = 1.0
			u := cur
			pcur := 0
			for h := 0; h < H && u != -1; h++ {
				path[pcur] = u
				pcur++
				u = p[u]
			}
			for k := 2; k < pcur; k++ {
				node := path[k]
				prev := path[k-1]
				mul[node*H+k] /= (1 - dp[prev*H+(k-1)]/2)
			}
			for k := 1; k < pcur; k++ {
				node := path[k]
				prev := path[k-1]
				mul[node*H+k] *= (1 - dp[prev*H+(k-1)]/2)
				dp[node*H+k] = 1 - mul[node*H+k]
			}
			cur++
		} else {
			v := qu.v
			sum := 0.0
			base := v * H
			for j := 1; j < H; j++ {
				sum += dp[base+j]
			}
			outputs = append(outputs, sum)
		}
	}
	return outputs
}

func genCase(rng *rand.Rand) testCase {
	q := rng.Intn(8) + 2 // 2..9
	var qs []query
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", q)
	cur := 1
	for i := 0; i < q; i++ {
		if i == q-1 { // ensure at least one type 2
			v := rng.Intn(cur) + 1
			qs = append(qs, query{2, v - 1})
			fmt.Fprintf(&sb, "2 %d\n", v)
			continue
		}
		typ := rng.Intn(2) + 1
		if typ == 1 {
			v := rng.Intn(cur) + 1
			qs = append(qs, query{1, v - 1})
			fmt.Fprintf(&sb, "1 %d\n", v)
			cur++
		} else {
			v := rng.Intn(cur) + 1
			qs = append(qs, query{2, v - 1})
			fmt.Fprintf(&sb, "2 %d\n", v)
		}
	}
	return testCase{q, qs, sb.String()}
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func parseFloats(lines []string) ([]float64, error) {
	res := make([]float64, len(lines))
	for i, l := range lines {
		v, err := strconv.ParseFloat(strings.TrimSpace(l), 64)
		if err != nil {
			return nil, err
		}
		res[i] = v
	}
	return res, nil
}

func validate(tc testCase, out string) error {
	outputs := strings.Split(strings.TrimSpace(out), "\n")
	expect := solve(tc.qs)
	if len(outputs) != len(expect) {
		return fmt.Errorf("expected %d lines, got %d", len(expect), len(outputs))
	}
	vals, err := parseFloats(outputs)
	if err != nil {
		return err
	}
	for i := range vals {
		if math.Abs(vals[i]-expect[i]) > 1e-6*math.Max(1, math.Abs(expect[i])) {
			return fmt.Errorf("line %d expected %.7f got %.7f", i+1, expect[i], vals[i])
		}
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
	const cases = 100
	for i := 0; i < cases; i++ {
		c := genCase(rng)
		out, err := run(bin, c.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:\n%s", i+1, err, c.input)
			os.Exit(1)
		}
		if err := validate(c, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%soutput:\n%s", i+1, err, c.input, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", cases)
}
