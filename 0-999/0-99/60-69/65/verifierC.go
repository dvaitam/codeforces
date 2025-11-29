package main

import (
	"bufio"
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

type P struct{ x, y, z float64 }

func (p P) add(q P) P       { return P{p.x + q.x, p.y + q.y, p.z + q.z} }
func (p P) sub(q P) P       { return P{p.x - q.x, p.y - q.y, p.z - q.z} }
func (p P) mul(s float64) P { return P{p.x * s, p.y * s, p.z * s} }
func (p P) norm() float64   { return math.Sqrt(p.x*p.x + p.y*p.y + p.z*p.z) }

func run(bin, input string) (string, error) {
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

type Result struct {
	possible bool
	t        float64
	p        P
}

func expected(n int, v []P, vh, vs float64, h P) Result {
	eps := 1e-11
	can := func(a P, t float64) bool {
		return a.sub(h).norm()/vh < t+eps
	}
	var onde P
	onde.x = -1e10
	var endIdx int
	var tAccum float64
	for i := 0; i < n; i++ {
		segDist := v[i].sub(v[i+1]).norm()
		dt := segDist / vs
		if can(v[i+1], tAccum+dt) {
			lo, hi := 0.0, 1.0
			for it := 0; it < 100; it++ {
				mid := (lo + hi) / 2.0
				m := v[i].add(v[i+1].sub(v[i]).mul(mid))
				tPerson := tAccum + v[i].sub(m).norm()/vs
				if can(m, tPerson) {
					hi = mid
				} else {
					lo = mid
				}
			}
			onde = v[i].add(v[i+1].sub(v[i]).mul(lo))
			endIdx = i
			break
		}
		tAccum += dt
	}
	if onde.x <= -1e8 {
		return Result{possible: false}
	}
	tPerson := tAccum + v[endIdx].sub(onde).norm()/vs
	tHeli := onde.sub(h).norm() / vh
	tFinal := math.Max(tPerson, tHeli)
	return Result{possible: true, t: tFinal, p: onde}
}

func genCase(rng *rand.Rand) (string, Result) {
	n := rng.Intn(3) + 1
	v := make([]P, n+1)
	for i := 0; i <= n; i++ {
		v[i] = P{float64(rng.Intn(21) - 10), float64(rng.Intn(21) - 10), float64(rng.Intn(21) - 10)}
		if i > 0 {
			for v[i] == v[i-1] {
				v[i] = P{float64(rng.Intn(21) - 10), float64(rng.Intn(21) - 10), float64(rng.Intn(21) - 10)}
			}
		}
	}
	vh := float64(rng.Intn(10) + 5)
	vs := float64(rng.Intn(int(vh)) + 1)
	h := P{float64(rng.Intn(21) - 10), float64(rng.Intn(21) - 10), float64(rng.Intn(21) - 10)}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i <= n; i++ {
		sb.WriteString(fmt.Sprintf("%.0f %.0f %.0f\n", v[i].x, v[i].y, v[i].z))
	}
	sb.WriteString(fmt.Sprintf("%.0f %.0f\n", vh, vs))
	sb.WriteString(fmt.Sprintf("%.0f %.0f %.0f\n", h.x, h.y, h.z))
	input := sb.String()
	exp := expected(n, v, vh, vs, h)
	return input, exp
}

func parseResult(s string) (Result, error) {
	sc := bufio.NewScanner(strings.NewReader(s))
	sc.Split(bufio.ScanWords)
	if !sc.Scan() {
		return Result{}, fmt.Errorf("empty output")
	}
	tok := sc.Text()
	if tok == "NO" {
		return Result{possible: false}, nil
	}
	if tok != "YES" {
		return Result{}, fmt.Errorf("expected YES or NO, got %s", tok)
	}
	if !sc.Scan() {
		return Result{}, fmt.Errorf("missing time")
	}
	t, err := strconv.ParseFloat(sc.Text(), 64)
	if err != nil {
		return Result{}, fmt.Errorf("invalid time: %v", err)
	}
	var p P
	if !sc.Scan() {
		return Result{}, fmt.Errorf("missing x")
	}
	if p.x, err = strconv.ParseFloat(sc.Text(), 64); err != nil {
		return Result{}, err
	}
	if !sc.Scan() {
		return Result{}, fmt.Errorf("missing y")
	}
	if p.y, err = strconv.ParseFloat(sc.Text(), 64); err != nil {
		return Result{}, err
	}
	if !sc.Scan() {
		return Result{}, fmt.Errorf("missing z")
	}
	if p.z, err = strconv.ParseFloat(sc.Text(), 64); err != nil {
		return Result{}, err
	}
	return Result{possible: true, t: t, p: p}, nil
}

func near(a, b float64) bool {
	diff := math.Abs(a - b)
	if diff < 1e-6 {
		return true
	}
	if math.Abs(b) > 1e-9 && diff/math.Abs(b) < 1e-6 {
		return true
	}
	return false
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		outStr, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		got, err := parseResult(outStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed to parse output: %v\ninput:\n%s\noutput:\n%s", i+1, err, in, outStr)
			os.Exit(1)
		}
		if got.possible != exp.possible {
			fmt.Fprintf(os.Stderr, "case %d failed: expected possible=%v got %v\ninput:\n%s", i+1, exp.possible, got.possible, in)
			os.Exit(1)
		}
		if exp.possible {
			if !near(got.t, exp.t) {
				fmt.Fprintf(os.Stderr, "case %d failed: expected time %.9f got %.9f\ninput:\n%s", i+1, exp.t, got.t, in)
				os.Exit(1)
			}
			// check point distance
			dist := got.p.sub(exp.p).norm()
			if dist > 1e-5 { // slightly relaxed check for coordinates consistency
				// Note: The problem mainly checks time and point validity. 
				// Since multiple points might map to same time (unlikely here due to monotonicity),
				// strictly checking time is usually sufficient.
				// However, let's warn or fail if coordinates are wildly off.
				// Given floating point drift, 1e-5 is safe.
				fmt.Fprintf(os.Stderr, "case %d failed: expected point %v got %v (dist %.9f)\ninput:\n%s", i+1, exp.p, got.p, dist, in)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
