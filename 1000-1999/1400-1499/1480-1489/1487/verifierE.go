package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func buildOracle() (string, error) {
	dir, _ := os.Getwd()
	oracle := filepath.Join(dir, "oracleE")
	cmd := exec.Command("go", "build", "-o", oracle, "1487E.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func run(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func genPairs(rng *rand.Rand, na, nb int) [][2]int {
	type pair = [2]int
	all := make([]pair, 0, na*nb)
	for x := 1; x <= na; x++ {
		for y := 1; y <= nb; y++ {
			all = append(all, pair{x, y})
		}
	}
	rng.Shuffle(len(all), func(i, j int) { all[i], all[j] = all[j], all[i] })
	m := rng.Intn(len(all) + 1)
	return all[:m]
}

func generateCase(rng *rand.Rand) string {
	n1 := rng.Intn(5) + 1
	n2 := rng.Intn(5) + 1
	n3 := rng.Intn(5) + 1
	n4 := rng.Intn(5) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d %d\n", n1, n2, n3, n4))
	for i := 0; i < n1; i++ {
		sb.WriteString(fmt.Sprintf("%d ", rng.Intn(10)+1))
	}
	sb.WriteByte('\n')
	for i := 0; i < n2; i++ {
		sb.WriteString(fmt.Sprintf("%d ", rng.Intn(10)+1))
	}
	sb.WriteByte('\n')
	for i := 0; i < n3; i++ {
		sb.WriteString(fmt.Sprintf("%d ", rng.Intn(10)+1))
	}
	sb.WriteByte('\n')
	for i := 0; i < n4; i++ {
		sb.WriteString(fmt.Sprintf("%d ", rng.Intn(10)+1))
	}
	sb.WriteByte('\n')
	pairs12 := genPairs(rng, n1, n2)
	sb.WriteString(fmt.Sprintf("%d\n", len(pairs12)))
	for _, p := range pairs12 {
		sb.WriteString(fmt.Sprintf("%d %d\n", p[0], p[1]))
	}
	pairs23 := genPairs(rng, n2, n3)
	sb.WriteString(fmt.Sprintf("%d\n", len(pairs23)))
	for _, p := range pairs23 {
		sb.WriteString(fmt.Sprintf("%d %d\n", p[0], p[1]))
	}
	pairs34 := genPairs(rng, n3, n4)
	sb.WriteString(fmt.Sprintf("%d\n", len(pairs34)))
	for _, p := range pairs34 {
		sb.WriteString(fmt.Sprintf("%d %d\n", p[0], p[1]))
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	rng := rand.New(rand.NewSource(1))
    for t := 0; t < 100; t++ {
        tc := generateCase(rng)
        input := tc
        exp, err := run(oracle, input)
        if err != nil {
            fmt.Printf("oracle runtime error: %v\n", err)
            os.Exit(1)
        }
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n%s", t+1, err, got)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("case %d failed:\ninput:\n%sexpected %s got %s\n", t+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
