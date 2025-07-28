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

type testCase struct {
	h []int64
	a []int64
}

func expected(tc testCase) string {
	n := len(tc.h)
	maxA := int64(0)
	for _, v := range tc.a {
		if v > maxA {
			maxA = v
		}
	}
	res := make([]int64, n)
	for x := int64(1); x <= maxA+1; x++ {
		times := make([]int64, n)
		max1, max2 := int64(-1), int64(-1)
		maxIdx := -1
		for i := 0; i < n; i++ {
			cur := tc.h[i] * ((tc.a[i] + x - 1) / x)
			times[i] = cur
			if cur > max1 {
				max2 = max1
				max1 = cur
				maxIdx = i
			} else if cur > max2 {
				max2 = cur
			}
		}
		cntMax := 0
		for i := 0; i < n; i++ {
			if times[i] == max1 {
				cntMax++
			}
		}
		if cntMax == 1 {
			diff := max1 - max2
			if diff > res[maxIdx] {
				res[maxIdx] = diff
			}
		}
	}
	var sb strings.Builder
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", res[i]))
	}
	return sb.String()
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out, errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCases() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]testCase, 0, 102)
	cases = append(cases, testCase{h: []int64{1}, a: []int64{1}})
	cases = append(cases, testCase{h: []int64{2, 3}, a: []int64{1, 2}})
	for len(cases) < 102 {
		n := rng.Intn(4) + 1
		h := make([]int64, n)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			h[i] = rng.Int63n(5) + 1
			a[i] = rng.Int63n(5) + 1
		}
		cases = append(cases, testCase{h: h, a: a})
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genCases()
	for i, tc := range tests {
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", len(tc.h)))
		for j := 0; j < len(tc.h); j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", tc.h[j]))
		}
		sb.WriteByte('\n')
		for j := 0; j < len(tc.a); j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", tc.a[j]))
		}
		sb.WriteByte('\n')
		input := sb.String()
		want := expected(tc)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, want, got, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
