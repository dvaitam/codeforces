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

type Day struct {
	d int64
	t int64
}

type Test struct {
	p    int64
	m    int64
	days []Day
}

func (t Test) Input() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", len(t.days), t.p, t.m))
	for _, d := range t.days {
		sb.WriteString(fmt.Sprintf("%d %d\n", d.d, d.t))
	}
	return sb.String()
}

func buildRef() (string, error) {
	ref := "./refF.bin"
	cmd := exec.Command("go", "build", "-o", ref, "926F.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build ref failed: %v: %s", err, out)
	}
	return ref, nil
}

func runExe(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func genTests() []Test {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]Test, 0, 101)
	for i := 0; i < 100; i++ {
		n := rng.Intn(5) + 1
		m := int64(rng.Intn(30) + n)
		p := int64(rng.Intn(20) + 1)
		days := make([]Day, n)
		used := make(map[int64]struct{})
		last := int64(0)
		for j := 0; j < n; j++ {
			last += int64(rng.Intn(3) + 1)
			if last > m {
				last = m
			}
			if _, ok := used[last]; ok {
				last++
			}
			used[last] = struct{}{}
			days[j] = Day{d: last, t: int64(rng.Intn(20) + 1)}
		}
		tests = append(tests, Test{p: p, m: m, days: days})
	}
	tests = append(tests, Test{p: 1, m: 1, days: []Day{{1, 1}}})
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := genTests()
	for i, tc := range tests {
		input := tc.Input()
		exp, err := runExe(ref, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runExe(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if exp != got {
			fmt.Printf("Test %d failed\nInput:%sExpected:%s\nGot:%s\n", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
