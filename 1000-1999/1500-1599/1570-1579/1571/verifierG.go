package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type attack struct {
	damage int64
	b      int
}

type warrior struct {
	atks []attack
}

type testG struct {
	n  int
	m  int
	ws []warrior
}

func buildOracle() (string, error) {
	exe := "oracleG.bin"
	cmd := exec.Command("go", "build", "-o", exe, "1571G.go")
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

func genTests() []testG {
	rng := rand.New(rand.NewSource(1571))
	tests := make([]testG, 0, 100)
	for len(tests) < 100 {
		n := rng.Intn(3) + 1
		m := rng.Intn(5) + 1
		ws := make([]warrior, n)
		for i := 0; i < n; i++ {
			k := rng.Intn(3) + 1
			used := make(map[int]bool)
			atks := make([]attack, k)
			for j := 0; j < k; j++ {
				var b int
				for {
					b = rng.Intn(m + 1)
					if !used[b] {
						used[b] = true
						break
					}
				}
				dmg := rng.Int63n(20) + 1
				atks[j] = attack{dmg, b}
			}
			// sort by b
			for j := 0; j < k; j++ {
				for t := j + 1; t < k; t++ {
					if atks[t].b < atks[j].b {
						atks[t], atks[j] = atks[j], atks[t]
					}
				}
			}
			ws[i] = warrior{atks}
		}
		tests = append(tests, testG{n, m, ws})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
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
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
		for _, w := range tc.ws {
			sb.WriteString(fmt.Sprintf("%d\n", len(w.atks)))
			for idx, a := range w.atks {
				if idx > 0 {
					sb.WriteByte(' ')
				}
				sb.WriteString(fmt.Sprintf("%d", a.damage))
			}
			sb.WriteByte('\n')
			for idx, a := range w.atks {
				if idx > 0 {
					sb.WriteByte(' ')
				}
				sb.WriteString(fmt.Sprintf("%d", a.b))
			}
			sb.WriteByte('\n')
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
