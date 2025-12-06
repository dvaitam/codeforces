package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
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

func buildBinary(source string, output string) error {
	cmd := exec.Command("go", "build", "-o", output, source)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("build %s failed: %v\n%s", source, err, out)
	}
	return nil
}

func runProg(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("timeout")
	}
	if err != nil {
		return out.String() + errBuf.String(), err
	}
	return strings.TrimSpace(out.String()), nil
}

func genTests() []testG {
	numTests := 20
	if val := os.Getenv("NUM_TESTS"); val != "" {
		if n, err := strconv.Atoi(val); err == nil && n > 0 {
			numTests = n
		}
	}

	rng := rand.New(rand.NewSource(1571))
	tests := make([]testG, 0, numTests)

	// Manual sample test case
	tests = append(tests, testG{
		n: 2, m: 3,
		ws: []warrior{
			{atks: []attack{{10, 3}}},
			{atks: []attack{{10, 3}}},
		},
	})

	for len(tests) < numTests {
		n := rng.Intn(3) + 1
		m := rng.Intn(5) + 1
		ws := make([]warrior, n)
		for i := 0; i < n; i++ {
			k := rng.Intn(3) + 1
			if k > m+1 {
				k = m + 1
			}
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
		fmt.Println("usage: go run verifierG.go /path/to/binary_or_source")
		return
	}
	candidateSource := os.Args[1]
	candidateBin := "./candidate.bin"
	oracleBin := "./oracleG.bin"
	// Clean up binaries on exit (best effort)
	defer func() {
		os.Remove(candidateBin)
		os.Remove(oracleBin)
	}()

	// Build oracle
	fmt.Println("Building oracle...")
	if err := buildBinary("1571G.go", oracleBin); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Build candidate if needed
	runBin := candidateSource
	if strings.HasSuffix(candidateSource, ".go") {
		fmt.Println("Building candidate...")
		if err := buildBinary(candidateSource, candidateBin); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		runBin = "./" + candidateBin
	} else {
		if !strings.Contains(runBin, "/") {
			runBin = "./" + runBin
		}
	}

	if !strings.Contains(oracleBin, "/") {
		oracleBin = "./" + oracleBin
	}

	tests := genTests()
	for i, tc := range tests {
		fmt.Printf("Running test %d/%d... ", i+1, len(tests)) 
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
		want, err := runProg(oracleBin, input)
		if err != nil {
			fmt.Printf("FAIL\noracle runtime error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runProg(runBin, input)
		if err != nil {
			fmt.Printf("FAIL\ncandidate runtime error on case %d: %v\n%s", i+1, err, got)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(want) {
			fmt.Printf("FAIL\ncase %d failed\ninput:\n%sexpected: %s\ngot: %s\n", i+1, input, want, got)
			os.Exit(1)
		}
		fmt.Println("PASS")
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
