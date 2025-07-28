package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type testCaseB struct {
	a []int
	b []int
}

func expectedB(a, b []int) string {
	bs := append([]int(nil), b...)
	sort.Sort(sort.Reverse(sort.IntSlice(bs)))
	ans := make([]int, 0, len(a)+len(b))
	p := 0
	for _, bi := range bs {
		for p < len(a) && a[p] >= bi {
			ans = append(ans, a[p])
			p++
		}
		ans = append(ans, bi)
	}
	for p < len(a) {
		ans = append(ans, a[p])
		p++
	}
	var sb strings.Builder
	for i, v := range ans {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	return sb.String()
}

func genTestsB() []testCaseB {
	rand.Seed(2)
	tests := make([]testCaseB, 0, 100)
	for len(tests) < 100 {
		n := rand.Intn(5) + 1
		m := rand.Intn(5) + 1
		a := make([]int, n)
		for i := range a {
			a[i] = rand.Intn(50)
		}
		sort.Sort(sort.Reverse(sort.IntSlice(a)))
		b := make([]int, m)
		for i := range b {
			b[i] = rand.Intn(50)
		}
		tests = append(tests, testCaseB{a: a, b: b})
	}
	return tests
}

func runCase(bin string, tc testCaseB) error {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", len(tc.a), len(tc.b)))
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	for i, v := range tc.b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')

	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	expect := expectedB(tc.a, tc.b)
	if got != expect {
		return fmt.Errorf("expected %s got %s (a=%v b=%v)", expect, got, tc.a, tc.b)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsB()
	for i, tc := range tests {
		if err := runCase(bin, tc); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
