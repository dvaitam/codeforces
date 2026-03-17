package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func equivAnswer(got, expect string) bool {
	if got == expect {
		return true
	}
	// Accept negated version (swap bipartite coloring)
	gTokens := strings.Fields(got)
	eTokens := strings.Fields(expect)
	if len(gTokens) != len(eTokens) {
		return false
	}
	for i := range gTokens {
		gv, err1 := strconv.Atoi(gTokens[i])
		ev, err2 := strconv.Atoi(eTokens[i])
		if err1 != nil || err2 != nil {
			return false
		}
		if gv != ev && gv != -ev {
			return false
		}
	}
	// Check all are same sign or all negated
	if len(gTokens) == 0 {
		return true
	}
	g0, _ := strconv.Atoi(gTokens[0])
	e0, _ := strconv.Atoi(eTokens[0])
	neg := (g0 == -e0)
	for i := range gTokens {
		gv, _ := strconv.Atoi(gTokens[i])
		ev, _ := strconv.Atoi(eTokens[i])
		if neg {
			if gv != -ev {
				return false
			}
		} else {
			if gv != ev {
				return false
			}
		}
	}
	return true
}

func buildOracle() (string, error) {
	src := os.Getenv("REFERENCE_SOURCE_PATH")
	if src == "" {
		return "", fmt.Errorf("REFERENCE_SOURCE_PATH not set")
	}
	// Detect C++ by checking content
	content, err := os.ReadFile(src)
	if err != nil {
		return "", fmt.Errorf("read reference: %v", err)
	}
	bin := "/tmp/oracle1656E.bin"
	if strings.Contains(string(content), "#include") {
		cppSrc := "/tmp/oracle1656E.cpp"
		if err := os.WriteFile(cppSrc, content, 0644); err != nil {
			return "", err
		}
		cmd := exec.Command("g++", "-O2", "-o", bin, cppSrc)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
		}
	} else {
		cmd := exec.Command("go", "build", "-o", bin, src)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
		}
	}
	return bin, nil
}

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

func genCase(r *rand.Rand) string {
	n := r.Intn(5) + 2
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	edges := make([][2]int, n-1)
	for i := 1; i < n; i++ {
		u := r.Intn(i) + 1
		v := i + 1
		edges[i-1] = [2]int{u, v}
	}
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := genCase(r)
		expect, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\n", i+1, err)
			fmt.Fprintln(os.Stderr, input)
			os.Exit(1)
		}
		got, err := run(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if !equivAnswer(strings.TrimSpace(got), strings.TrimSpace(expect)) {
			fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", i+1, input, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
