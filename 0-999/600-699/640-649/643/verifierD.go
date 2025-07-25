package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type query struct {
	typ int
	a   int
	b   int
}

type testCase struct {
	n, q  int
	t     []int64
	fol   []int
	qs    []query
	input string
}

func computeIncome(t []int64, fol []int) []int64 {
	n := len(t)
	deg := make([]int, n)
	for _, f := range fol {
		deg[f]++
	}
	x := make([]int64, n)
	for i := 0; i < n; i++ {
		x[i] = t[i] / int64(deg[i]+2)
	}
	followerPart := make([]int64, n)
	for i := 0; i < n; i++ {
		followerPart[fol[i]] += x[i]
	}
	ownerPart := make([]int64, n)
	for i := 0; i < n; i++ {
		ownerPart[i] = t[i] - int64(deg[i]+1)*x[i]
	}
	inc := make([]int64, n)
	for i := 0; i < n; i++ {
		inc[i] = ownerPart[i] + followerPart[i] + x[fol[i]]
	}
	return inc
}

func genCase(rng *rand.Rand) testCase {
	n := rng.Intn(4) + 2 // 2..5
	q := rng.Intn(4) + 2 // 2..5
	t := make([]int64, n)
	for i := range t {
		t[i] = int64(rng.Intn(10) + 1)
	}
	fol := make([]int, n)
	for i := range fol {
		fol[i] = rng.Intn(n)
	}
	var qs []query
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, q)
	for i, v := range t {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	for i, v := range fol {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v+1)
	}
	sb.WriteByte('\n')
	for i := 0; i < q; i++ {
		typ := rng.Intn(3) + 1
		if typ == 1 {
			a := rng.Intn(n) + 1
			b := rng.Intn(n) + 1
			qs = append(qs, query{1, a, b})
			fmt.Fprintf(&sb, "1 %d %d\n", a, b)
			fol[a-1] = b - 1
		} else if typ == 2 {
			a := rng.Intn(n) + 1
			qs = append(qs, query{2, a, 0})
			fmt.Fprintf(&sb, "2 %d\n", a)
		} else {
			qs = append(qs, query{3, 0, 0})
			fmt.Fprintln(&sb, "3")
		}
	}
	return testCase{n, q, t, fol, qs, sb.String()}
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

func parseInts(line string) ([]int64, error) {
	fields := strings.Fields(strings.TrimSpace(line))
	res := make([]int64, len(fields))
	for i, f := range fields {
		v, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return nil, err
		}
		res[i] = v
	}
	return res, nil
}

func validate(tc testCase, out string) error {
	outputs := strings.Split(strings.TrimSpace(out), "\n")
	idx := 0
	fol := make([]int, len(tc.fol))
	copy(fol, tc.fol)
	for _, q := range tc.qs {
		if q.typ == 1 {
			fol[q.a-1] = q.b - 1
		} else {
			if idx >= len(outputs) {
				return fmt.Errorf("missing output line for query")
			}
			incomes := computeIncome(tc.t, fol)
			if q.typ == 2 {
				want := incomes[q.a-1]
				val, err := strconv.ParseInt(strings.TrimSpace(outputs[idx]), 10, 64)
				if err != nil {
					return err
				}
				if val != want {
					return fmt.Errorf("expected %d got %d", want, val)
				}
			} else {
				vals, err := parseInts(outputs[idx])
				if err != nil {
					return err
				}
				if len(vals) != 2 {
					return fmt.Errorf("expected two integers")
				}
				minv, maxv := incomes[0], incomes[0]
				for _, v := range incomes[1:] {
					if v < minv {
						minv = v
					}
					if v > maxv {
						maxv = v
					}
				}
				if vals[0] != minv || vals[1] != maxv {
					return fmt.Errorf("expected %d %d got %d %d", minv, maxv, vals[0], vals[1])
				}
			}
			idx++
		}
	}
	if idx != len(outputs) {
		return fmt.Errorf("extra output lines")
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
