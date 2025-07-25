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

func expected(b, d, s int64) int64 {
	mx := b
	if d > mx {
		mx = d
	}
	if s > mx {
		mx = s
	}
	res := int64(0)
	if mx-1 > b {
		res += mx - 1 - b
	}
	if mx-1 > d {
		res += mx - 1 - d
	}
	if mx-1 > s {
		res += mx - 1 - s
	}
	return res
}

func runCase(bin string, b, d, s int64) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	input := fmt.Sprintf("%d %d %d\n", b, d, s)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var ans int64
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &ans); err != nil {
		return fmt.Errorf("bad output: %s", out.String())
	}
	if ans != expected(b, d, s) {
		return fmt.Errorf("expected %d got %d", expected(b, d, s), ans)
	}
	return nil
}

func randCase(rng *rand.Rand) (int64, int64, int64) {
	b := rng.Int63n(20)
	d := rng.Int63n(20)
	s := rng.Int63n(20)
	if b+d+s == 0 {
		b = 1
	}
	return b, d, s
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []struct{ b, d, s int64 }{
		{1, 1, 1},
		{0, 0, 1},
		{10, 10, 10},
	}
	for len(cases) < 105 {
		b, d, s := randCase(rng)
		cases = append(cases, struct{ b, d, s int64 }{b, d, s})
	}
	for i, c := range cases {
		if err := runCase(bin, c.b, c.d, c.s); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput: %d %d %d\n", i+1, err, c.b, c.d, c.s)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
