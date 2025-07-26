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
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	err := cmd.Run()
	if err != nil {
		return out.String() + errb.String(), err
	}
	return out.String(), nil
}

type Test struct {
	n     int
	m     []int64
	input string
}

func genTest(rng *rand.Rand) Test {
	n := rng.Intn(8) + 1
	mvals := make([]int64, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		maxv := int64(i)
		if maxv == 0 {
			maxv = 1
		}
		mvals[i] = rng.Int63n(maxv)
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", mvals[i]))
	}
	sb.WriteByte('\n')
	return Test{n: n, m: mvals, input: sb.String()}
}

func solve(t Test) string {
	n := t.n
	mvals := t.m
	d := make([]int64, n)
	if n > 0 {
		d[n-1] = mvals[n-1]
		if d[n-1] < 0 {
			d[n-1] = 0
		}
		for i := n - 2; i >= 0; i-- {
			need := d[i+1] - 1
			if need < 0 {
				need = 0
			}
			if mvals[i] > need {
				d[i] = mvals[i]
			} else {
				d[i] = need
			}
		}
	}
	var sum int64
	for i := 0; i < n; i++ {
		sum += d[i]
	}
	return fmt.Sprintf("%d", sum)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		t := genTest(rng)
		expected := solve(t)
		out, err := run(bin, t.input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n%s", i+1, err, out)
			os.Exit(1)
		}
		out = strings.TrimSpace(out)
		if out != expected {
			fmt.Printf("test %d failed\ninput:\n%s\nexpected:%s got:%s\n", i+1, t.input, expected, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
