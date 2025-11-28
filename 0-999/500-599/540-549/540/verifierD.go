package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type key struct{ r, s, p int }

var memo map[key]float64

func prob(r, s, p int) float64 {
	if r == 0 {
		return 0
	}
	if p == 0 {
		return 1
	}
	if s == 0 {
		return 0
	}
	k := key{r, s, p}
	if v, ok := memo[k]; ok {
		return v
	}
	rs := float64(r * s)
	sp := float64(s * p)
	pr := float64(p * r)
	total := rs + sp + pr
	d := (pr / total) * prob(r-1, s, p)
	d += (rs / total) * prob(r, s-1, p)
	d += (sp / total) * prob(r, s, p-1)
	memo[k] = d
	return d
}

func solveCase(r, s, p int) string {
	memo = make(map[key]float64)
	rProb := prob(r, s, p)
	memo = make(map[key]float64)
	sProb := prob(s, p, r)
	memo = make(map[key]float64)
	pProb := prob(p, r, s)
	return fmt.Sprintf("%.9f %.9f %.9f", rProb, sProb, pProb)
}

func isClose(a, b float64) bool {
	const eps = 1e-9 // Problem requires 10^-9
	absA, absB := math.Abs(a), math.Abs(b)
	diff := math.Abs(a - b)

	if diff <= eps {
		return true
	}
	// Relative error
	if diff <= eps*math.Max(absA, absB) {
		return true
	}
	return false
}

func genCase(rng *rand.Rand) (string, string) {
	r := rng.Intn(100) + 1 // r,s,p up to 100
	s := rng.Intn(100) + 1
	p := rng.Intn(100) + 1
	input := fmt.Sprintf("%d %d %d\n", r, s, p)
	expected := solveCase(r, s, p)
	return input, expected
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin) // Removed .go check
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	args := os.Args[1:]
	if len(args) > 0 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := args[0]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	const cases = 100
	for i := 0; i < cases; i++ {
		inp, expectStr := genCase(rng)
		gotStr, err := runBinary(bin, inp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:%s", i+1, err, inp)
			os.Exit(1)
		}

		// Parse expected values
		var eR, eS, eP float64
		_, err = fmt.Sscanf(expectStr, "%f %f %f", &eR, &eS, &eP)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: failed to parse expected output %q: %v\n", i+1, expectStr, err)
			os.Exit(1)
		}

		// Parse got values
		var gR, gS, gP float64
		_, err = fmt.Sscanf(gotStr, "%f %f %f", &gR, &gS, &gP)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: failed to parse candidate output %q: %v\n", i+1, gotStr, err)
			os.Exit(1)
		}

		// Compare with tolerance
		if !isClose(eR, gR) || !isClose(eS, gS) || !isClose(eP, gP) {
			fmt.Fprintf(os.Stderr, "test %d failed\nexpected: %.9f %.9f %.9f\ngot: %.12f %.12f %.12f\ninput:%s",
				i+1, eR, eS, eP, gR, gS, gP, inp)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", cases)
}
