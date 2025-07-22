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

type Part struct{ a, b int }
type Actor struct{ c, d, k int }

func canPerform(act Actor, p Part) bool {
	return act.c <= p.a && p.b <= act.d
}

func search(idx int, parts []Part, actors []Actor, assign []int, rem []int) bool {
	if idx == len(parts) {
		return true
	}
	for i := range actors {
		if rem[i] > 0 && canPerform(actors[i], parts[idx]) {
			rem[i]--
			assign[idx] = i + 1
			if search(idx+1, parts, actors, assign, rem) {
				return true
			}
			rem[i]++
		}
	}
	return false
}

func expectedC(parts []Part, actors []Actor) (bool, []int) {
	rem := make([]int, len(actors))
	for i, a := range actors {
		rem[i] = a.k
	}
	assign := make([]int, len(parts))
	if search(0, parts, actors, assign, rem) {
		return true, assign
	}
	return false, nil
}

func genCaseC(rng *rand.Rand) (string, string) {
	n := rng.Intn(3) + 1
	m := rng.Intn(3) + 1
	parts := make([]Part, n)
	for i := 0; i < n; i++ {
		x := rng.Intn(10) + 1
		y := x + rng.Intn(5)
		parts[i] = Part{x, y}
	}
	actors := make([]Actor, m)
	for i := 0; i < m; i++ {
		c := rng.Intn(10) + 1
		d := c + rng.Intn(7)
		k := rng.Intn(2) + 1
		actors[i] = Actor{c, d, k}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		fmt.Fprintf(&sb, "%d %d\n", parts[i].a, parts[i].b)
	}
	fmt.Fprintf(&sb, "%d\n", m)
	for i := 0; i < m; i++ {
		fmt.Fprintf(&sb, "%d %d %d\n", actors[i].c, actors[i].d, actors[i].k)
	}
	ok, assign := expectedC(parts, actors)
	if !ok {
		return sb.String(), "NO\n"
	}
	var exp strings.Builder
	exp.WriteString("YES\n")
	for i, v := range assign {
		if i > 0 {
			exp.WriteByte(' ')
		}
		fmt.Fprintf(&exp, "%d", v)
	}
	exp.WriteByte('\n')
	return sb.String(), exp.String()
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != strings.TrimSpace(expected) {
		return fmt.Errorf("expected \n%s\ngot \n%s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCaseC(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
