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

func run(bin, input string) (string, error) {
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
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCaseC(rng *rand.Rand) string {
	t := 1
	n := rng.Intn(3) + 1
	m := rng.Intn(3) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n%d %d\n", t, n, m))
	for i := 0; i < n; i++ {
		val := rng.Intn(10)
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(val))
	}
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(rng.Intn(10)))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refPath := os.Getenv("REFERENCE_SOURCE_PATH")
	if refPath == "" {
		fmt.Fprintln(os.Stderr, "REFERENCE_SOURCE_PATH not set")
		os.Exit(1)
	}

	refBin, err := buildBinary(refPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := genCaseC(rng)
		refOut, err := run(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		candOut, err := run(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(candOut) != strings.TrimSpace(refOut) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\ninput:\n%s", i+1, refOut, candOut, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

func buildBinary(src string) (string, error) {
	tmp, err := os.CreateTemp("", "1936C-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	cmd := exec.Command("go", "build", "-o", tmp.Name(), src)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}
