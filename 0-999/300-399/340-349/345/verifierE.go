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

func run(bin string, input string) (string, error) {
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

func intersects(a, v, u int, x, y float64) bool {
	af := float64(a)
	vf := float64(v)
	uf := float64(u)
	alpha := 1.0 - (uf*uf)/(vf*vf)
	eps := 1e-12
	xif := x
	yif := y
	f := func(t float64) float64 {
		return alpha*t*t - 2*xif*t + xif*xif + yif*yif
	}
	if v == u {
		if x > 0 {
			return f(af) <= eps
		}
		return f(0) <= eps
	}
	if alpha > 0 {
		x0 := xif / alpha
		if x0 < 0 {
			x0 = 0
		} else if x0 > af {
			x0 = af
		}
		return f(x0) <= eps
	}
	return f(0) <= eps || f(af) <= eps
}

func countCats(a, v, u int, cats [][2]int) int {
	cnt := 0
	for _, c := range cats {
		x := float64(c[0])
		y := float64(c[1])
		if intersects(a, v, u, x, y) {
			cnt++
		}
	}
	return cnt
}

func runCase(bin string, a, v, u int, cats [][2]int) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d %d\n", a, v, u, len(cats)))
	for _, c := range cats {
		sb.WriteString(fmt.Sprintf("%d %d\n", c[0], c[1]))
	}
	out, err := run(bin, sb.String())
	if err != nil {
		return err
	}
	got, err := strconv.Atoi(strings.TrimSpace(out))
	if err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	exp := countCats(a, v, u, cats)
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		a := rng.Intn(50) + 1
		v := rng.Intn(100) + 1
		u := rng.Intn(100) + 1
		n := rng.Intn(5) + 1
		cats := make([][2]int, n)
		for j := 0; j < n; j++ {
			cats[j][0] = rng.Intn(201) - 100
			cats[j][1] = rng.Intn(201) - 100
		}
		if err := runCase(bin, a, v, u, cats); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
