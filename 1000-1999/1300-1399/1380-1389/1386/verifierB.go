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

const numTestsB = 100

func prepareBinary(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "binB*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		cmd := exec.Command("go", "build", "-o", tmp.Name(), path)
		if out, err := cmd.CombinedOutput(); err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("go build failed: %v: %s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, nil, nil
}

func prepareOracle() (string, func(), error) {
	tmp, err := os.CreateTemp("", "oracleB*")
	if err != nil {
		return "", nil, err
	}
	tmp.Close()
	cmd := exec.Command("go", "build", "-o", tmp.Name(), "1386B.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(tmp.Name())
		return "", nil, fmt.Errorf("go build oracle failed: %v: %s", err, out)
	}
	return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	return strings.TrimSpace(buf.String()), err
}

func genCase(rng *rand.Rand) string {
	sf := rng.Intn(5)
	pf := rng.Intn(5)
	gf := rng.Intn(5)
	if sf+pf+gf == 0 {
		sf = 1
	}
	n := rng.Intn(4) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n%d\n", sf, pf, gf, n)
	added := 0
	active := make([]int, 0)
	for i := 0; i < n; i++ {
		if len(active) == 0 || rng.Intn(2) == 0 {
			s := rng.Intn(5)
			p := rng.Intn(5)
			g := rng.Intn(5)
			if s+p+g == 0 {
				s = 1
			}
			added++
			active = append(active, added)
			fmt.Fprintf(&sb, "A %d %d %d\n", s, p, g)
		} else {
			idx := rng.Intn(len(active))
			r := active[idx]
			active = append(active[:idx], active[idx+1:]...)
			fmt.Fprintf(&sb, "R %d\n", r)
		}
	}
	return sb.String()
}

func runCase(bin, oracle, input string) error {
	exp, err := run(oracle, input)
	if err != nil {
		return fmt.Errorf("oracle error: %v", err)
	}
	got, err := run(bin, input)
	if err != nil {
		return fmt.Errorf("runtime error: %v", err)
	}
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	argIdx := 1
	if len(os.Args) >= 3 && os.Args[1] == "--" {
		argIdx = 2
	}
	if len(os.Args) != argIdx+1 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}
	bin, cleanup, err := prepareBinary(os.Args[argIdx])
	if err != nil {
		fmt.Println("compile error:", err)
		return
	}
	if cleanup != nil {
		defer cleanup()
	}
	oracle, cleanOracle, err := prepareOracle()
	if err != nil {
		fmt.Println("oracle compile error:", err)
		return
	}
	defer cleanOracle()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < numTestsB; i++ {
		in := genCase(rng)
		if err := runCase(bin, oracle, in); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			return
		}
	}
	fmt.Println("All tests passed")
}
