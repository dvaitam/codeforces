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

func buildOracle() (string, error) {
	exe := "oracleH"
	cmd := exec.Command("go", "build", "-o", exe, "1305H.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return "./" + exe, nil
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(4) + 1
	m := rng.Intn(4) + 1
	l := make([]int, n)
	r := make([]int, n)
	for i := 0; i < n; i++ {
		li := rng.Intn(m + 1)
		ri := li + rng.Intn(m-li+1)
		l[i] = li
		r[i] = ri
	}
	q := rng.Intn(m + 1)
	ps := rand.Perm(m)[:q]
	sort.Ints(ps)
	svals := make([]int, q)
	for i := 0; i < q; i++ {
		maxv := n
		if i > 0 && svals[i-1] < maxv {
			maxv = svals[i-1]
		}
		svals[i] = rng.Intn(maxv + 1)
	}
	totalL := 0
	totalR := 0
	for i := 0; i < n; i++ {
		totalL += l[i]
		totalR += r[i]
	}
	t := rng.Intn(totalR-totalL+1) + totalL

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", l[i], r[i]))
	}
	sb.WriteString(fmt.Sprintf("%d\n", q))
	for i := 0; i < q; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", ps[i]+1, svals[i]))
	}
	sb.WriteString(fmt.Sprintf("%d\n", t))
	return sb.String()
}

func runCase(bin, oracle, input string) error {
	cmdO := exec.Command(oracle)
	cmdO.Stdin = strings.NewReader(input)
	var outO bytes.Buffer
	cmdO.Stdout = &outO
	if err := cmdO.Run(); err != nil {
		return fmt.Errorf("oracle runtime error: %v", err)
	}
	expected := strings.TrimSpace(outO.String())

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\nstderr: %s", err, errb.String())
	}
	result := strings.TrimSpace(out.String())
	if result != expected {
		return fmt.Errorf("expected %q got %q", expected, result)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := generateCase(rng)
		if err := runCase(bin, oracle, input); err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
