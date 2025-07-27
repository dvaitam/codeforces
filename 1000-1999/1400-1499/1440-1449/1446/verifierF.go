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

type testCase struct {
	n   int
	k   int64
	pts [][2]float64
}

func runProg(bin, input string) (string, error) {
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
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expected(t testCase) (string, error) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", t.n, t.k))
	for _, p := range t.pts {
		sb.WriteString(fmt.Sprintf("%f %f\n", p[0], p[1]))
	}
	return runProg("1446F.go", sb.String())
}

func generateTests() []testCase {
	tests := []testCase{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 105; i++ {
		n := r.Intn(5) + 2
		totalPairs := int64(n * (n - 1) / 2)
		k := int64(r.Intn(int(totalPairs)) + 1)
		pts := make([][2]float64, n)
		for j := 0; j < n; j++ {
			pts[j][0] = r.Float64()*20 - 10
			pts[j][1] = r.Float64()*20 - 10
		}
		tests = append(tests, testCase{n: n, k: k, pts: pts})
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", t.n, t.k))
		for _, p := range t.pts {
			sb.WriteString(fmt.Sprintf("%f %f\n", p[0], p[1]))
		}
		want, err := expected(t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runProg(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		wantVal := strings.TrimSpace(want)
		gotVal := strings.TrimSpace(got)
		if v1, err1 := strconv.ParseFloat(gotVal, 64); err1 == nil {
			if v2, err2 := strconv.ParseFloat(wantVal, 64); err2 == nil {
				diff := v1 - v2
				if diff < 0 {
					diff = -diff
				}
				if diff > 1e-6 {
					fmt.Fprintf(os.Stderr, "test %d failed: expected %s got %s\n", i+1, wantVal, gotVal)
					os.Exit(1)
				}
				continue
			}
		}
		if gotVal != wantVal {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %q got %q\n", i+1, wantVal, gotVal)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
