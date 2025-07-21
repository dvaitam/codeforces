package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Query struct {
	typ int
	a   int
	b   int
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
	return out.String(), nil
}

func simulate(n int, a []int, qs []Query) []string {
	res := []string{}
	for _, q := range qs {
		if q.typ == 0 {
			a[q.a-1] = q.b
		} else {
			pos := q.a
			steps := 0
			last := pos
			for pos <= n {
				last = pos
				pos += a[pos-1]
				steps++
			}
			res = append(res, fmt.Sprintf("%d %d", last, steps))
		}
	}
	return res
}

func generateCaseE(rng *rand.Rand) (string, []string) {
	n := rng.Intn(10) + 1
	m := rng.Intn(15) + 1
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(n) + 1
	}
	qs := make([]Query, m)
	for i := range qs {
		if rng.Intn(2) == 0 {
			qs[i].typ = 0
			qs[i].a = rng.Intn(n) + 1
			qs[i].b = rng.Intn(n) + 1
		} else {
			qs[i].typ = 1
			qs[i].a = rng.Intn(n) + 1
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	for _, q := range qs {
		if q.typ == 0 {
			sb.WriteString(fmt.Sprintf("0 %d %d\n", q.a, q.b))
		} else {
			sb.WriteString(fmt.Sprintf("1 %d\n", q.a))
		}
	}
	return sb.String(), simulate(n, arr, qs)
}

func runCase(bin, input string, expected []string) error {
	out, err := run(bin, input)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(strings.NewReader(out))
	for i := 0; i < len(expected); i++ {
		if !scanner.Scan() {
			return fmt.Errorf("not enough output, expected %d lines", len(expected))
		}
		line := strings.TrimSpace(scanner.Text())
		if line != expected[i] {
			return fmt.Errorf("line %d: expected '%s' got '%s'", i+1, expected[i], line)
		}
	}
	if scanner.Scan() {
		return fmt.Errorf("extra output: %s", scanner.Text())
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCaseE(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
