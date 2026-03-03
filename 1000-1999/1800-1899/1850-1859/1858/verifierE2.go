package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func buildOracle(srcFile, binName string) (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	src := filepath.Join(dir, srcFile)
	bin := filepath.Join(os.TempDir(), binName)
	cmd := exec.Command("go", "build", "-o", bin, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return bin, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func deterministicCases() []string {
	return []string{
		"5\n+ 1\n+ 2\n?\n- 1\n?\n",
		"6\n+ 3\n+ 4\n?\n!\n+ 5\n?\n",
	}
}

func randomCase(rng *rand.Rand) string {
	q := rng.Intn(30) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", q)
	arrLen := 0
	type change struct {
		op  byte // '+' or '-'
		val int  // x for +, k for -
	}
	var hist []change
	for i := 0; i < q; i++ {
		canPop := arrLen > 0
		canRollback := len(hist) > 0

		if !canPop && !canRollback {
			// Must push or query
			if rng.Intn(2) == 0 {
				x := rng.Intn(1000000) + 1
				fmt.Fprintf(&sb, "+ %d\n", x)
				hist = append(hist, change{'+', x})
				arrLen++
			} else {
				sb.WriteString("?\n")
			}
			continue
		}

		switch rng.Intn(4) {
		case 0: // push
			x := rng.Intn(1000000) + 1
			fmt.Fprintf(&sb, "+ %d\n", x)
			hist = append(hist, change{'+', x})
			arrLen++
		case 1: // pop
			if canPop {
				k := rng.Intn(arrLen) + 1
				fmt.Fprintf(&sb, "- %d\n", k)
				hist = append(hist, change{'-', k})
				arrLen -= k
			} else {
				sb.WriteString("?\n")
			}
		case 2: // rollback
			if canRollback {
				sb.WriteString("!\n")
				last := hist[len(hist)-1]
				hist = hist[:len(hist)-1]
				if last.op == '+' {
					arrLen--
				} else {
					arrLen += last.val
				}
			} else {
				sb.WriteString("?\n")
			}
		case 3: // query
			sb.WriteString("?\n")
		}
	}
	return sb.String()
}

func verify(oracle, userBin string, cases []string) {
	for i, in := range cases {
		want, err := run(oracle, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on case %d: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		got, err := run(userBin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if want != got {
			fmt.Printf("test %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, in, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE2.go /path/to/binary")
		os.Exit(1)
	}
	userBin := os.Args[1]
	oracle, err := buildOracle("1858E2.go", "oracle1858E2.bin")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := deterministicCases()
	for len(cases) < 100 {
		cases = append(cases, randomCase(rng))
	}
	verify(oracle, userBin, cases)
}
