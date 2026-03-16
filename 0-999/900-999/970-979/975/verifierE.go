package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
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

func buildReference() (string, func(), error) {
	refSrc := os.Getenv("REFERENCE_SOURCE_PATH")
	if refSrc == "" {
		return "", nil, fmt.Errorf("REFERENCE_SOURCE_PATH not set")
	}

	content, err := os.ReadFile(refSrc)
	if err != nil {
		return "", nil, fmt.Errorf("cannot read reference source: %v", err)
	}

	tmpDir, err := os.MkdirTemp("", "975E-ref")
	if err != nil {
		return "", nil, err
	}
	cleanup := func() { os.RemoveAll(tmpDir) }

	binPath := filepath.Join(tmpDir, "ref_975E")

	if strings.Contains(string(content), "#include") {
		cppPath := filepath.Join(tmpDir, "ref.cpp")
		if err := os.WriteFile(cppPath, content, 0644); err != nil {
			cleanup()
			return "", nil, err
		}
		cmd := exec.Command("g++", "-O2", "-o", binPath, cppPath)
		if out, err := cmd.CombinedOutput(); err != nil {
			cleanup()
			return "", nil, fmt.Errorf("g++ build failed: %v\n%s", err, string(out))
		}
	} else {
		cmd := exec.Command("go", "build", "-o", binPath, refSrc)
		if out, err := cmd.CombinedOutput(); err != nil {
			cleanup()
			return "", nil, fmt.Errorf("go build failed: %v\n%s", err, string(out))
		}
	}

	return binPath, cleanup, nil
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

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in := genCase(rng)

		exp, err := runBin(refBin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on case %d: %v\n", i+1, err)
			os.Exit(1)
		}

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
