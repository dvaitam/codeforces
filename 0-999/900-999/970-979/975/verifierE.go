package main

import (
	"bytes"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

type PT struct{ x, y float64 }

type query struct {
	typ  int
	f, t int
	v    int
}

type testCase struct {
	pts []PT
	qs  []query
}

func buildCase(tc testCase) string {
	var sb strings.Builder
	n := len(tc.pts)
	q := len(tc.qs)
	sb.WriteString(fmt.Sprintf("%d %d\n", n, q))
	for _, p := range tc.pts {
		sb.WriteString(fmt.Sprintf("%d %d\n", int(p.x), int(p.y)))
	}
	for _, qu := range tc.qs {
		if qu.typ == 1 {
			sb.WriteString(fmt.Sprintf("1 %d %d\n", qu.f, qu.t))
		} else {
			sb.WriteString(fmt.Sprintf("2 %d\n", qu.v))
		}
	}
	return sb.String()
}

func randomPolygon(rng *rand.Rand) []PT {
	n := rng.Intn(3) + 3 // 3..5
	angles := make([]float64, n)
	for i := 0; i < n; i++ {
		angles[i] = rng.Float64() * 2 * math.Pi
	}
	sort.Float64s(angles)
	pts := make([]PT, n)
	for i, a := range angles {
		r := float64(rng.Intn(9) + 1)
		x := math.Round(r * math.Cos(a))
		y := math.Round(r * math.Sin(a))
		pts[i] = PT{x, y}
	}
	return pts
}

func genCase(rng *rand.Rand) string {
	tc := testCase{}
	tc.pts = randomPolygon(rng)
	q := rng.Intn(5) + 1
	tc.qs = make([]query, q)
	pinned := [2]int{1, 2}
	has2 := false
	for i := 0; i < q; i++ {
		if rng.Intn(2) == 0 {
			idx := rng.Intn(2)
			f := pinned[idx]
			t := rng.Intn(len(tc.pts)) + 1
			pinned[idx] = t
			tc.qs[i] = query{typ: 1, f: f, t: t}
		} else {
			v := rng.Intn(len(tc.pts)) + 1
			tc.qs[i] = query{typ: 2, v: v}
			has2 = true
		}
	}
	if !has2 {
		tc.qs[0] = query{typ: 2, v: 1}
	}
	return buildCase(tc)
}

// Embedded correct solver for 975E
func solve975E(input string) string {
	data := []byte(input)
	p := 0
	nextInt := func() int64 {
		for p < len(data) && data[p] <= ' ' {
			p++
		}
		sign := int64(1)
		if p < len(data) && data[p] == '-' {
			sign = -1
			p++
		}
		var v int64
		for p < len(data) && data[p] > ' ' {
			v = v*10 + int64(data[p]-'0')
			p++
		}
		return sign * v
	}

	n := int(nextInt())
	q := int(nextInt())

	xi := make([]int64, n+1)
	yi := make([]int64, n+1)
	x := make([]float64, n+1)
	y := make([]float64, n+1)

	for i := 1; i <= n; i++ {
		xv := nextInt()
		yv := nextInt()
		xi[i] = xv
		yi[i] = yv
		x[i] = float64(xv)
		y[i] = float64(yv)
	}

	var area2, numX, numY big.Int
	var b1, b2, term big.Int

	for i := 1; i <= n; i++ {
		j := i + 1
		if j > n {
			j = 1
		}
		cross := xi[i]*yi[j] - xi[j]*yi[i]
		area2.Add(&area2, b1.SetInt64(cross))

		term.Mul(b1.SetInt64(cross), b2.SetInt64(xi[i]+xi[j]))
		numX.Add(&numX, &term)

		term.Mul(b1.SetInt64(cross), b2.SetInt64(yi[i]+yi[j]))
		numY.Add(&numY, &term)
	}

	var denom big.Int
	denom.Mul(&area2, big.NewInt(3))

	prec := uint(256)
	denF := new(big.Float).SetPrec(prec).SetInt(&denom)

	cxF := new(big.Float).SetPrec(prec).SetInt(&numX)
	cyF := new(big.Float).SetPrec(prec).SetInt(&numY)

	cxF.Quo(cxF, denF)
	cyF.Quo(cyF, denF)

	cx, _ := cxF.Float64()
	cy, _ := cyF.Float64()

	stableC := make([]float64, n+1)
	stableS := make([]float64, n+1)

	for i := 1; i <= n; i++ {
		ux := cx - x[i]
		uy := cy - y[i]
		r := math.Hypot(ux, uy)
		stableC[i] = -uy / r
		stableS[i] = -ux / r
	}

	curC, curS := 1.0, 0.0
	tx, ty := 0.0, 0.0
	p1, p2 := 1, 2

	out := make([]byte, 0, q*40)

	for ; q > 0; q-- {
		typ := int(nextInt())
		if typ == 1 {
			f := int(nextInt())
			t := int(nextInt())

			g := p1
			if p1 == f {
				g = p2
			}

			wgx := curC*x[g] - curS*y[g] + tx
			wgy := curS*x[g] + curC*y[g] + ty

			curC = stableC[g]
			curS = stableS[g]

			tx = wgx - (curC*x[g] - curS*y[g])
			ty = wgy - (curS*x[g] + curC*y[g])

			p1, p2 = g, t
		} else {
			v := int(nextInt())
			xv := curC*x[v] - curS*y[v] + tx
			yv := curS*x[v] + curC*y[v] + ty

			out = strconv.AppendFloat(out, xv, 'f', 10, 64)
			out = append(out, ' ')
			out = strconv.AppendFloat(out, yv, 'f', 10, 64)
			out = append(out, '\n')
		}
	}

	return strings.TrimSpace(string(out))
}

func runBin(bin, input string) (string, error) {
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

func compareOutputs(got, exp string) error {
	gotFields := strings.Fields(got)
	expFields := strings.Fields(exp)
	if len(gotFields) != len(expFields) {
		return fmt.Errorf("expected %d numbers got %d", len(expFields), len(gotFields))
	}
	for i := 0; i < len(expFields); i++ {
		g, err := strconv.ParseFloat(gotFields[i], 64)
		if err != nil {
			return fmt.Errorf("bad output token %q: %v", gotFields[i], err)
		}
		e, _ := strconv.ParseFloat(expFields[i], 64)
		if math.Abs(g-e) > 1e-4*math.Max(1, math.Abs(e)) {
			return fmt.Errorf("value mismatch at position %d: expected %v got %v", i, e, g)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierE /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in := genCase(rng)

		exp := solve975E(in)

		got, err := runBin(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}

		if err := compareOutputs(got, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
