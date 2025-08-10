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

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return out.String() + stderr.String(), err
	}
	return out.String(), nil
}

type Op struct {
	op string
	x  int
}

type Test struct {
	ops   []Op
	input string
}

func genTest(rng *rand.Rand) Test {
	n := rng.Intn(20) + 1
	ops := make([]Op, n)
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(n))
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		typ := []string{"&", "|", "^"}[rng.Intn(3)]
		val := rng.Intn(1024)
		ops[i] = Op{op: typ, x: val}
		sb.WriteString(fmt.Sprintf("%s %d\n", typ, val))
	}
	return Test{ops: ops, input: sb.String()}
}

func apply(ops []Op, x int) int {
	for _, op := range ops {
		switch op.op {
		case "&":
			x &= op.x
		case "|":
			x |= op.x
		case "^":
			x ^= op.x
		}
	}
	return x
}

func parseOutput(out string) ([]Op, error) {
	tokens := strings.Fields(out)
	if len(tokens) == 0 {
		return nil, fmt.Errorf("empty output")
	}
	m, err := strconv.Atoi(tokens[0])
	if err != nil {
		return nil, fmt.Errorf("invalid number of ops: %v", err)
	}
	if m < 0 || m > 3 {
		return nil, fmt.Errorf("invalid m=%d", m)
	}
	if len(tokens) != 1+2*m {
		return nil, fmt.Errorf("expected %d tokens, got %d", 1+2*m, len(tokens))
	}
	ops := make([]Op, m)
	idx := 1
	for i := 0; i < m; i++ {
		typ := tokens[idx]
		idx++
		if typ != "&" && typ != "|" && typ != "^" {
			return nil, fmt.Errorf("invalid op %q", typ)
		}
		val, err := strconv.Atoi(tokens[idx])
		idx++
		if err != nil {
			return nil, fmt.Errorf("invalid value: %v", err)
		}
		ops[i] = Op{op: typ, x: val}
	}
	return ops, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		t := genTest(rng)
		out, err := run(bin, t.input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\noutput:\n%s", i+1, err, out)
			os.Exit(1)
		}
		ops, err := parseOutput(out)
		if err != nil {
			fmt.Printf("test %d parse error: %v\noutput:\n%s", i+1, err, out)
			os.Exit(1)
		}
		for x := 0; x < 1024; x++ {
			expected := apply(t.ops, x)
			actual := apply(ops, x)
			if expected != actual {
				fmt.Printf("test %d failed\ninput:\n%sx = %d\nexpected = %d\nactual = %d\noutput:\n%s", i+1, t.input, x, expected, actual, out)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
