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

func buildOracle() (string, error) {
	_, self, _, _ := runtime.Caller(0)
	dir := filepath.Dir(self)
	oracle := filepath.Join(dir, "oracleG")
	cmd := exec.Command("go", "build", "-o", oracle, "1566G.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(2) + 4
	m := rng.Intn(3) + 4
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	edges := make(map[[2]int]bool)
	for len(edges) < m {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v {
			continue
		}
		key := [2]int{u, v}
		if u > v {
			key = [2]int{v, u}
		}
		if edges[key] {
			continue
		}
		edges[key] = true
		w := rng.Intn(10) + 1
		sb.WriteString(fmt.Sprintf("%d %d %d\n", u, v, w))
	}
	q := rng.Intn(3)
	sb.WriteString(fmt.Sprintf("%d\n", q))
	existing := make(map[[2]int]bool)
	for e := range edges {
		existing[e] = true
	}
	for i := 0; i < q; i++ {
		if rng.Intn(2) == 0 && len(existing) > 0 {
			// delete
			var e [2]int
			for k := range existing {
				e = k
				break
			}
			sb.WriteString(fmt.Sprintf("0 %d %d\n", e[0], e[1]))
			delete(existing, e)
		} else {
			// add
			var u, v int
			for {
				u = rng.Intn(n) + 1
				v = rng.Intn(n) + 1
				if u != v {
					key := [2]int{u, v}
					if u > v {
						key = [2]int{v, u}
					}
					if !existing[key] {
						existing[key] = true
						break
					}
				}
			}
			w := rng.Intn(10) + 1
			sb.WriteString(fmt.Sprintf("1 %d %d %d\n", u, v, w))
		}
	}
	return sb.String()
}

func run(prog, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(prog, ".go") {
		cmd = exec.Command("go", "run", prog)
	} else {
		cmd = exec.Command(prog)
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

func runCase(bin, oracle, input string) error {
	got, err := run(bin, input)
	if err != nil {
		return err
	}
	expect, err := run(oracle, input)
	if err != nil {
		return fmt.Errorf("oracle error: %v", err)
	}
	if strings.TrimSpace(got) != strings.TrimSpace(expect) {
		return fmt.Errorf("expected %s got %s", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
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
	for i := 0; i < 100; i++ {
		input := generateCase(rng)
		if err := runCase(bin, oracle, input); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
