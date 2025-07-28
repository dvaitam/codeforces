package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type enemy struct {
	x int
	y int
	p int
}

type testH struct {
	n  int
	a  int
	b  int
	es []enemy
}

func buildOracle() (string, error) {
	exe := "oracleH.bin"
	cmd := exec.Command("go", "build", "-o", exe, "1571H.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return "./" + exe, nil
}

func runProg(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return out.String() + errBuf.String(), err
	}
	return strings.TrimSpace(out.String()), nil
}

func genTests() []testH {
	rng := rand.New(rand.NewSource(1571))
	tests := make([]testH, 0, 100)
	for len(tests) < 100 {
		a := rng.Intn(8) + 2
		b := rng.Intn(8) + 2
		n := rng.Intn(3) + 1
		used := make(map[[2]int]bool)
		es := make([]enemy, n)
		for i := 0; i < n; i++ {
			var x, y int
			for {
				x = rng.Intn(a-1) + 1
				y = rng.Intn(b-1) + 1
				if !used[[2]int{x, y}] {
					used[[2]int{x, y}] = true
					break
				}
			}
			p := rng.Intn(999999) + 1
			es[i] = enemy{x, y, p}
		}
		tests = append(tests, testH{n, a, b, es})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierH.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	tests := genTests()
	for i, tc := range tests {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.n, tc.a, tc.b))
		for _, e := range tc.es {
			sb.WriteString(fmt.Sprintf("%d %d %d\n", e.x, e.y, e.p))
		}
		input := sb.String()
		want, err := runProg(oracle, input)
		if err != nil {
			fmt.Printf("oracle runtime error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Printf("candidate runtime error on case %d: %v\n%s", i+1, err, got)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(want) {
			fmt.Printf("case %d failed\ninput:\n%sexpected: %s\ngot: %s\n", i+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
