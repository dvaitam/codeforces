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

type ghost struct {
	x, vx, vy int64
}

func solveCase(a, b int64, g []ghost) string {
	totalByW := make(map[int64]int64)
	totalByWV := make(map[[2]int64]int64)
	for _, p := range g {
		w := p.vy - a*p.vx
		totalByW[w]++
		key := [2]int64{w, p.vx}
		totalByWV[key]++
	}
	var collisions int64
	for _, cnt := range totalByW {
		collisions += cnt * (cnt - 1) / 2
	}
	for _, cnt := range totalByWV {
		collisions -= cnt * (cnt - 1) / 2
	}
	return fmt.Sprintf("%d\n", collisions*2)
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(6) + 1
	a := int64(rng.Intn(5) - 2)
	b := int64(rng.Intn(5) - 2)
	ghosts := make([]ghost, n)
	for i := 0; i < n; i++ {
		ghosts[i] = ghost{
			x:  int64(rng.Intn(11) - 5),
			vx: int64(rng.Intn(7) - 3),
			vy: int64(rng.Intn(7) - 3),
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, a, b))
	for _, g := range ghosts {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", g.x, g.vx, g.vy))
	}
	return sb.String(), solveCase(a, b, ghosts)
}

func runCase(bin, in, exp string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != strings.TrimSpace(exp) {
		return fmt.Errorf("expected %s got %s", exp, out.String())
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
