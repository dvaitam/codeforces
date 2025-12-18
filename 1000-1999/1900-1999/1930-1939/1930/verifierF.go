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

func solveF(n, q int, e []int) []int {
	a := make([]int, 0, q)
	res := make([]int, q)
	ans := 0
	last := 0
	for i := 0; i < q; i++ {
		v := (e[i] + last) % n
		a = append(a, v)

		// Brute force update ans
		// Since q is small in the generator (<= 6), we can check all pairs.
		// For larger q, this oracle would be too slow, but it's fine here.
		for _, u := range a {
			// Check u & ~v
			val1 := u &^ v
			if val1 > ans {
				ans = val1
			}
			// Check v & ~u
			val2 := v &^ u
			if val2 > ans {
				ans = val2
			}
		}
		res[i] = ans
		last = ans
	}
	return res
}

func generateCaseF(rng *rand.Rand) (string, string) {
	n := rng.Intn(32) + 1
	q := rng.Intn(6) + 1
	e := make([]int, q)
	for i := range e {
		e[i] = rng.Intn(n)
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", n, q))
	for i, v := range e {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	input := sb.String()
	ans := solveF(n, q, e)
	var exp strings.Builder
	for i, v := range ans {
		if i > 0 {
			exp.WriteByte(' ')
		}
		exp.WriteString(strconv.Itoa(v))
	}
	exp.WriteByte('\n')
	return input, exp.String()
}

func runProg(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("%v", err)
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCaseF(rng)
		got, err := runProg(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\noutput:\n%s", i+1, err, got)
			os.Exit(1)
		}
		got = strings.TrimSpace(got)
		if got != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\ninput:\n%s", i+1, strings.TrimSpace(exp), got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}