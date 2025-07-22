package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type testCase struct {
	input string
}

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleE")
	cmd := exec.Command("go", "build", "-o", oracle, "269E.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func generateRandomCase(rng *rand.Rand) testCase {
	n := rng.Intn(4) + 1
	m := rng.Intn(4) + 1
	type side struct {
		name  string
		count int
	}
	sides := []side{{"L", n}, {"R", n}, {"T", m}, {"B", m}}
	pins := make(map[string][]int)
	for _, s := range sides {
		arr := rng.Perm(s.count)
		for i := range arr {
			arr[i]++
		}
		pins[s.name] = arr
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	stringsUsed := 0
	total := n + m
	for stringsUsed < total {
		// choose two different sides with available pins
		side1 := sides[rng.Intn(len(sides))]
		for len(pins[side1.name]) == 0 {
			side1 = sides[rng.Intn(len(sides))]
		}
		side2 := sides[rng.Intn(len(sides))]
		for side2.name == side1.name || len(pins[side2.name]) == 0 {
			side2 = sides[rng.Intn(len(sides))]
		}
		p1 := pins[side1.name][len(pins[side1.name])-1]
		pins[side1.name] = pins[side1.name][:len(pins[side1.name])-1]
		p2 := pins[side2.name][len(pins[side2.name])-1]
		pins[side2.name] = pins[side2.name][:len(pins[side2.name])-1]
		sb.WriteString(fmt.Sprintf("%s %d %s %d\n", side1.name, p1, side2.name, p2))
		stringsUsed++
	}
	return testCase{input: sb.String()}
}

func runCase(bin, oracle string, tc testCase) error {
	cmdO := exec.Command(oracle)
	cmdO.Stdin = strings.NewReader(tc.input)
	var outO bytes.Buffer
	cmdO.Stdout = &outO
	if err := cmdO.Run(); err != nil {
		return fmt.Errorf("oracle run error: %v", err)
	}
	expected := strings.TrimSpace(outO.String())

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected:\n%s\n\ngot:\n%s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
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
	cases := make([]testCase, 0, 101)
	for i := 0; i < 101; i++ {
		cases = append(cases, generateRandomCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, oracle, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
