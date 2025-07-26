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

func solve(t Test) string {
	zero := [10]int{}
	one := [10]int{}
	for i := 0; i < 10; i++ {
		zero[i] = 0
		one[i] = 1
	}
	for _, op := range t.ops {
		for j := 0; j < 10; j++ {
			b := (op.x >> j) & 1
			switch op.op {
			case "|":
				zero[j] |= b
				one[j] |= b
			case "&":
				zero[j] &= b
				one[j] &= b
			case "^":
				zero[j] ^= b
				one[j] ^= b
			}
		}
	}
	andMask, orMask, xorMask := 0, 0, 0
	for j := 0; j < 10; j++ {
		a0 := zero[j]
		a1 := one[j]
		if a0 == 0 && a1 == 0 {
		} else if a0 == 1 && a1 == 1 {
			andMask |= 1 << j
			orMask |= 1 << j
		} else if a0 == 0 && a1 == 1 {
			andMask |= 1 << j
		} else {
			andMask |= 1 << j
			xorMask |= 1 << j
		}
	}
	return fmt.Sprintf("3\n& %d\n| %d\n^ %d", andMask, orMask, xorMask)
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
		expected := solve(t)
		out, err := run(bin, t.input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\noutput:\n%s", i+1, err, out)
			os.Exit(1)
		}
		out = strings.TrimSpace(out)
		if out != expected {
			fmt.Printf("test %d failed\ninput:\n%s\nexpected:\n%s\nactual:\n%s\n", i+1, t.input, expected, out)
			os.Exit(1)
		}
	}
	fmt.Println("ok 100 tests")
}
