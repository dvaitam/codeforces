package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type query struct{ k, pos int }

func solve(a []int, qs []query) []int {
	res := make([]int, len(qs))
	for qi, q := range qs {
		n := len(a)
		pairs := make([]struct{ v, idx int }, n)
		for i, v := range a {
			pairs[i] = struct{ v, idx int }{v, i}
		}
		sort.Slice(pairs, func(i, j int) bool {
			if pairs[i].v == pairs[j].v {
				return pairs[i].idx < pairs[j].idx
			}
			return pairs[i].v > pairs[j].v
		})
		sel := pairs[:q.k]
		sort.Slice(sel, func(i, j int) bool { return sel[i].idx < sel[j].idx })
		res[qi] = a[sel[q.pos-1].idx]
	}
	return res
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func runCase(bin string, a []int, qs []query) error {
	input := fmt.Sprintf("%d\n", len(a))
	for i, v := range a {
		if i > 0 {
			input += " "
		}
		input += fmt.Sprintf("%d", v)
	}
	input += "\n" + fmt.Sprintf("%d\n", len(qs))
	for _, q := range qs {
		input += fmt.Sprintf("%d %d\n", q.k, q.pos)
	}
	expect := solve(a, qs)
	out, err := run(bin, input)
	if err != nil {
		return err
	}
	fields := strings.Fields(out)
	if len(fields) != len(qs) {
		return fmt.Errorf("expected %d numbers got %d", len(qs), len(fields))
	}
	for i, f := range fields {
		var v int
		if _, err := fmt.Sscan(f, &v); err != nil {
			return fmt.Errorf("bad int %q", f)
		}
		if v != expect[i] {
			return fmt.Errorf("ans %d expected %d got %d", i+1, expect[i], v)
		}
	}
	return nil
}

func genCase(rng *rand.Rand) ([]int, []query) {
	n := rng.Intn(8) + 1
	a := make([]int, n)
	for i := range a {
		a[i] = rng.Intn(50)
	}
	m := rng.Intn(5) + 1
	qs := make([]query, m)
	for i := range qs {
		k := rng.Intn(n) + 1
		pos := rng.Intn(k) + 1
		qs[i] = query{k, pos}
	}
	return a, qs
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	a := []int{10, 20, 30, 20}
	qs := []query{{2, 1}, {2, 2}}
	if err := runCase(bin, a, qs); err != nil {
		fmt.Fprintf(os.Stderr, "case 1 failed: %v\n", err)
		os.Exit(1)
	}
	for i := 1; i < 100; i++ {
		a, qs := genCase(rng)
		if err := runCase(bin, a, qs); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
