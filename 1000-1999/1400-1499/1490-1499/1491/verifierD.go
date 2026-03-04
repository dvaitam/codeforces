package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type query struct {
	u int
	v int
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func canTravel(u, v int) bool {
	if u > v {
		return false
	}
	onesU, onesV := 0, 0
	for i := 0; i < 30; i++ {
		if (u>>i)&1 == 1 {
			onesU++
		}
		if (v>>i)&1 == 1 {
			onesV++
		}
		if onesV > onesU {
			return false
		}
	}
	return true
}

func parseCase(input string) ([]query, error) {
	fields := strings.Fields(input)
	if len(fields) == 0 {
		return nil, fmt.Errorf("empty input")
	}
	q, err := strconv.Atoi(fields[0])
	if err != nil || q < 1 {
		return nil, fmt.Errorf("invalid q")
	}
	if len(fields) != 1+2*q {
		return nil, fmt.Errorf("wrong token count for q=%d", q)
	}
	out := make([]query, 0, q)
	p := 1
	for i := 0; i < q; i++ {
		u, errU := strconv.Atoi(fields[p])
		v, errV := strconv.Atoi(fields[p+1])
		if errU != nil || errV != nil {
			return nil, fmt.Errorf("invalid query #%d", i+1)
		}
		out = append(out, query{u: u, v: v})
		p += 2
	}
	return out, nil
}

func genCase(rng *rand.Rand) string {
	q := rng.Intn(40) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", q))
	for i := 0; i < q; i++ {
		u := rng.Intn((1 << 30) - 1)
		v := rng.Intn((1 << 30) - 1)
		// Keep queries within statement bounds [1, 2^30).
		u++
		v++
		sb.WriteString(fmt.Sprintf("%d %d\n", u, v))
	}
	return sb.String()
}

func edgeCases() []string {
	return []string{
		"5\n1 1\n1 2\n2 1\n3 5\n10 10\n",
		"6\n5 9\n5 8\n8 9\n7 15\n7 8\n536870912 1073741823\n",
	}
}

func validateOutput(input, got string) error {
	qs, err := parseCase(input)
	if err != nil {
		return err
	}
	ans := strings.Fields(got)
	if len(ans) != len(qs) {
		return fmt.Errorf("wrong number of answers: got %d want %d", len(ans), len(qs))
	}
	for i, tk := range ans {
		upper := strings.ToUpper(tk)
		if upper != "YES" && upper != "NO" {
			return fmt.Errorf("query %d: invalid token %q", i+1, tk)
		}
		want := canTravel(qs[i].u, qs[i].v)
		if (upper == "YES") != want {
			return fmt.Errorf("query %d failed: u=%d v=%d expected %v got %s", i+1, qs[i].u, qs[i].v, want, upper)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	rng := rand.New(rand.NewSource(1))
	cases := edgeCases()
	for i := 0; i < 100; i++ {
		cases = append(cases, genCase(rng))
	}

	for tcase, input := range cases {
		got, err := run(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on case %d: %v\ninput:\n%s", tcase+1, err, input)
			os.Exit(1)
		}
		if err := validateOutput(input, got); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s\ngot:\n%s\n", tcase+1, err, input, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
