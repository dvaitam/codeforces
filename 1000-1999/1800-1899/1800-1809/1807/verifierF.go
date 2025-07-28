package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type Case struct {
	n, m   int
	i1, j1 int
	i2, j2 int
	d      string
}

func runProg(bin, input string) (string, error) {
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
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCases() []Case {
	rng := rand.New(rand.NewSource(1807))
	dirs := []string{"DR", "DL", "UR", "UL"}
	cases := make([]Case, 100)
	for i := range cases {
		n := rng.Intn(8) + 2
		m := rng.Intn(8) + 2
		i1 := rng.Intn(n) + 1
		j1 := rng.Intn(m) + 1
		i2 := rng.Intn(n) + 1
		j2 := rng.Intn(m) + 1
		d := dirs[rng.Intn(len(dirs))]
		cases[i] = Case{n, m, i1, j1, i2, j2, d}
	}
	return cases
}

type state struct {
	r, c, dx, dy int
}

func expected(c Case) int {
	dir := map[string][2]int{"DR": {1, 1}, "DL": {1, -1}, "UR": {-1, 1}, "UL": {-1, -1}}
	dx, dy := dir[c.d][0], dir[c.d][1]
	r, curr := c.i1, c.j1
	bounces := 0
	visited := make(map[state]bool)
	for {
		if r == c.i2 && curr == c.j2 {
			return bounces
		}
		st := state{r, curr, dx, dy}
		if visited[st] {
			return -1
		}
		visited[st] = true
		nr, nc := r+dx, curr+dy
		bounce := false
		if nr < 1 || nr > c.n {
			dx = -dx
			bounce = true
		}
		if nc < 1 || nc > c.m {
			dy = -dy
			bounce = true
		}
		if bounce {
			bounces++
		}
		r += dx
		curr += dy
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := genCases()
	for i, c := range cases {
		input := fmt.Sprintf("1\n%d %d %d %d %d %d %s\n", c.n, c.m, c.i1, c.j1, c.i2, c.j2, c.d)
		exp := fmt.Sprintf("%d", expected(c))
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != exp {
			fmt.Printf("case %d failed: expected %s got %s (case %+v)\n", i+1, exp, got, c)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
