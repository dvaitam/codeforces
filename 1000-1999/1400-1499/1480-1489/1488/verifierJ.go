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

func buildOracle() (string, error) {
	exe := "oracleJ"
	cmd := exec.Command("go", "build", "-o", exe, "1488J.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle: %v\n%s", err, out)
	}
	return exe, nil
}

func runProg(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(4) + 1
	m := rng.Intn(8) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for i := 0; i < n; i++ {
		w := rng.Intn(5) + 1
		fmt.Fprintf(&sb, "%d", w)
		if i+1 < n {
			sb.WriteByte(' ')
		}
	}
	sb.WriteByte('\n')
	stock := make([]int, n)
	for i := 0; i < m; i++ {
		t := rng.Intn(3) + 1
		if t == 1 {
			idx := rng.Intn(n) + 1
			c := rng.Intn(5) + 1
			stock[idx-1] += c
			fmt.Fprintf(&sb, "1 %d %d\n", idx, c)
		} else if t == 2 {
			idx := rng.Intn(n) + 1
			if stock[idx-1] == 0 {
				stock[idx-1]++
			}
			c := rng.Intn(stock[idx-1]) + 1
			stock[idx-1] -= c
			fmt.Fprintf(&sb, "2 %d %d\n", idx, c)
		} else {
			l := rng.Intn(n) + 1
			r := rng.Intn(n-l+1) + l
			k := rng.Intn(10) + 1
			fmt.Fprintf(&sb, "3 %d %d %d\n", l, r, k)
		}
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierJ.go /path/to/binary")
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
		exp, err := runProg("./"+oracle, input)
		if err != nil {
			fmt.Printf("oracle error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(exp) {
			fmt.Printf("case %d mismatch\nexpected:\n%s\n got:\n%s\n", i+1, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
