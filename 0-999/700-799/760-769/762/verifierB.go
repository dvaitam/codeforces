package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type Mouse struct {
	cost int64
	typ  string
}

func expected(a, b, c int, items []Mouse) (int, int64) {
	usb := []int64{}
	ps := []int64{}
	for _, m := range items {
		if m.typ == "USB" {
			usb = append(usb, m.cost)
		} else {
			ps = append(ps, m.cost)
		}
	}
	sort.Slice(usb, func(i, j int) bool { return usb[i] < usb[j] })
	sort.Slice(ps, func(i, j int) bool { return ps[i] < ps[j] })
	count := 0
	total := int64(0)
	u := 0
	for u < len(usb) && a > 0 {
		total += usb[u]
		count++
		u++
		a--
	}
	p := 0
	for p < len(ps) && b > 0 {
		total += ps[p]
		count++
		p++
		b--
	}
	rest := append(append([]int64{}, usb[u:]...), ps[p:]...)
	sort.Slice(rest, func(i, j int) bool { return rest[i] < rest[j] })
	r := 0
	for r < len(rest) && c > 0 {
		total += rest[r]
		count++
		r++
		c--
	}
	return count, total
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(42))

	for t := 0; t < 100; t++ {
		a := rng.Intn(6)
		b := rng.Intn(6)
		c := rng.Intn(6)
		m := rng.Intn(15)
		items := make([]Mouse, m)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d\n", a, b, c))
		sb.WriteString(fmt.Sprintf("%d\n", m))
		for i := 0; i < m; i++ {
			cost := rng.Int63n(1000) + 1
			typ := "USB"
			if rng.Intn(2) == 0 {
				typ = "PS/2"
			}
			items[i] = Mouse{cost: cost, typ: typ}
			sb.WriteString(fmt.Sprintf("%d %s\n", cost, typ))
		}
		input := sb.String()
		cnt, tot := expected(a, b, c, items)
		exp := fmt.Sprintf("%d %d", cnt, tot)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\ninput:\n%s\noutput:\n%s\n", t+1, err, input, got)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != exp {
			fmt.Printf("wrong answer on test %d\ninput:\n%s\nexpected: %s\ngot: %s\n", t+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
