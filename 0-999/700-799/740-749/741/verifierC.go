package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleC")
	cmd := exec.Command("go", "build", "-o", oracle, "741C.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(10) + 1
	seats := rng.Perm(2 * n)
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(n))
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		a := seats[2*i] + 1
		b := seats[2*i+1] + 1
		sb.WriteString(fmt.Sprintf("%d %d\n", a, b))
	}
	return sb.String()
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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

func check(input, output string) error {
	in := bufio.NewReader(strings.NewReader(input))
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return fmt.Errorf("invalid input format")
	}
	m := 2 * n
	a := make([]int, n)
	b := make([]int, n)
	for i := 0; i < n; i++ {
		if _, err := fmt.Fscan(in, &a[i], &b[i]); err != nil {
			return fmt.Errorf("invalid input format")
		}
		a[i]--
		b[i]--
	}

	out := bufio.NewReader(strings.NewReader(output))
	color := make([]int, m)
	for i := range color {
		color[i] = -1
	}
	for i := 0; i < n; i++ {
		var ca, cb int
		if _, err := fmt.Fscan(out, &ca, &cb); err != nil {
			return fmt.Errorf("invalid output format")
		}
		if ca < 1 || ca > 2 || cb < 1 || cb > 2 {
			return fmt.Errorf("invalid food type")
		}
		ca--
		cb--
		if ca == cb {
			return fmt.Errorf("pair %d has same food", i+1)
		}
		color[a[i]] = ca
		color[b[i]] = cb
	}
	for i := 0; i < m; i++ {
		if color[i] == -1 {
			return fmt.Errorf("seat %d has no food assigned", i+1)
		}
	}
	for i := 0; i < m; i++ {
		c1 := color[i]
		c2 := color[(i+1)%m]
		c3 := color[(i+2)%m]
		if c1 == c2 && c2 == c3 {
			return fmt.Errorf("three consecutive chairs starting at %d have same food", i+1)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in := generateCase(rng)
		exp, err := run(oracle, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error on case %d: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		got, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		exp = strings.TrimSpace(exp)
		got = strings.TrimSpace(got)
		if exp == "-1" {
			if got != "-1" {
				fmt.Printf("case %d failed\nexpected:\n%s\n\ngot:\n%s\n", i+1, exp, got)
				os.Exit(1)
			}
			continue
		}
		if got == "-1" {
			fmt.Printf("case %d failed\nexpected a valid solution but got -1\ninput:\n%s", i+1, in)
			os.Exit(1)
		}
		if err := check(in, got); err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%soutput:\n%s\n", i+1, err, in, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
