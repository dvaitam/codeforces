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

func buildRef() (string, error) {
	ref := "./refE.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1443E.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return ref, nil
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return out.String() + stderr.String(), err
	}
	return out.String(), nil
}

type Test struct {
	input   string
	answers int
}

func genTest(rng *rand.Rand) Test {
	n := rng.Intn(4) + 2 // 2..5
	q := rng.Intn(8) + 1 // 1..8
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, q))
	answers := 0
	for i := 0; i < q; i++ {
		if rng.Intn(2) == 0 {
			l := rng.Intn(n) + 1
			r := rng.Intn(n-l+1) + l
			sb.WriteString(fmt.Sprintf("1 %d %d\n", l, r))
			answers++
		} else {
			x := rng.Intn(5) + 1
			sb.WriteString(fmt.Sprintf("2 %d\n", x))
		}
	}
	return Test{input: sb.String(), answers: answers}
}

func parseInts64(s string) ([]int64, error) {
	fields := strings.Fields(strings.TrimSpace(s))
	res := make([]int64, len(fields))
	for i, f := range fields {
		v, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return nil, err
		}
		res[i] = v
	}
	return res, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := genTest(rng)
		exp, err := run(ref, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference error on test %d: %v\n%s", i+1, err, exp)
			os.Exit(1)
		}
		got, err := run(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\noutput:\n%s", i+1, err, got)
			os.Exit(1)
		}
		expVals, err := parseInts64(exp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "bad reference output on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		gotVals, err := parseInts64(got)
		if err != nil {
			fmt.Fprintf(os.Stderr, "bad output on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if len(expVals) != tc.answers || len(gotVals) != tc.answers {
			fmt.Printf("test %d expected %d numbers got %d\ninput:\n%s", i+1, tc.answers, len(gotVals), tc.input)
			os.Exit(1)
		}
		for j := 0; j < tc.answers; j++ {
			if expVals[j] != gotVals[j] {
				fmt.Printf("test %d failed at answer %d\ninput:\n%s\nexpected %d got %d\n", i+1, j+1, tc.input, expVals[j], gotVals[j])
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
