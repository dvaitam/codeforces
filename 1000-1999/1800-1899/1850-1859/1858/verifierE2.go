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
	history := 0
	for i := 0; i < q; i++ {
		opType := rng.Intn(4)
		if opType == 0 || arrLen == 0 { // push
			x := rng.Intn(1000000) + 1
			fmt.Fprintf(&sb, "+ %d\n", x)
			arrLen++
			history++
			continue
		}
		if opType == 1 && arrLen > 0 { // pop
			k := rng.Intn(arrLen) + 1
			fmt.Fprintf(&sb, "- %d\n", k)
			arrLen -= k
			history++
			continue
		}
		if opType == 2 && history > 0 { // rollback
			sb.WriteString("!\n")
			history--
			continue
		}
		sb.WriteString("?\n")
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
