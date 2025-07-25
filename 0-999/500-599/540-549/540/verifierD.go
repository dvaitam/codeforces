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

func genCase(rng *rand.Rand) (string, string) {
	r := rng.Intn(10) + 1
	s := rng.Intn(10) + 1
	p := rng.Intn(10) + 1
	input := fmt.Sprintf("%d %d %d\n", r, s, p)
	expected := solveCase(r, s, p)
	return input, expected
}

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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	const cases = 100
	for i := 0; i < cases; i++ {
		inp, expect := genCase(rng)
		got, err := run(bin, inp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed\nexpected: %s\ngot: %s\ninput:%s", i+1, expect, got, inp)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", cases)
}
