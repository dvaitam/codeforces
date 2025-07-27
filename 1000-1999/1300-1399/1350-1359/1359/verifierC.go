package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		if out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput(); err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func runProg(exe, input string) (string, error) {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type testCase struct{ h, c, t int64 }

func generateCase(rng *rand.Rand) testCase {
	h := int64(rng.Intn(1_000_000-1) + 2)
	c := int64(rng.Intn(int(h-1)) + 1)
	t := c + int64(rng.Intn(int(h-c)+1))
	return testCase{h, c, t}
}

func expected(tc testCase) string {
	h, c, t := tc.h, tc.c, tc.t
	if t >= h {
		return "1"
	}
	if 2*t <= h+c {
		return "2"
	}
	k := (h - t) / (2*t - h - c)
	n1 := 2*k + 1
	diff1 := math.Abs(float64((k+1)*h+k*c)/float64(n1) - float64(t))
	k++
	n2 := 2*k + 1
	diff2 := math.Abs(float64((k+1)*h+k*c)/float64(n2) - float64(t))
	if diff1 <= diff2 {
		return fmt.Sprintf("%d", n1)
	}
	return fmt.Sprintf("%d", n2)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	exe, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCase(rng)
		input := fmt.Sprintf("1\n%d %d %d\n", tc.h, tc.c, tc.t)
		want := expected(tc)
		got, err := runProg(exe, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d mismatch\nexpected:%s\n got:%s\ninput:%s", i+1, want, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
