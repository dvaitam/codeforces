package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"
)

type testCaseA struct {
	arr []int
}

func generateCase(rng *rand.Rand) (string, testCaseA) {
	n := rng.Intn(20) + 1
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(100) + 1
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", arr[i])
	}
	sb.WriteByte('\n')
	return sb.String(), testCaseA{arr: arr}
}

func solveLocal(tc testCaseA) string {
	a := append([]int(nil), tc.arr...)
	b := append([]int(nil), tc.arr...)
	sort.Ints(b)
	cntA := make(map[int][2]int)
	cntB := make(map[int][2]int)
	for i, v := range a {
		p := i % 2
		val := cntA[v]
		val[p]++
		cntA[v] = val
	}
	for i, v := range b {
		p := i % 2
		val := cntB[v]
		val[p]++
		cntB[v] = val
	}
	ok := true
	if len(cntA) != len(cntB) {
		ok = false
	} else {
		for k, va := range cntA {
			if vb, okb := cntB[k]; !okb || va != vb {
				ok = false
				break
			}
		}
	}
	if ok {
		return "YES\n"
	}
	return "NO\n"
}

func runProg(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func runRef(input string) (string, error) {
	_, self, _, _ := runtime.Caller(0)
	dir := filepath.Dir(self)
	ref := filepath.Join(dir, "1545A.go")
	return runProg(ref, input)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go <binary>")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, tc := generateCase(rng)
		expect := solveLocal(tc)
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n%s", i+1, err, got)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "test %d failed\ninput:\n%sexpected:%sactual:%s", i+1, input, expect, got)
			os.Exit(1)
		}
		// also cross-check with reference implementation to ensure generator valid
		refOut, err := runRef(input)
		if err == nil && refOut != expect {
			fmt.Fprintf(os.Stderr, "reference mismatch on test %d\ninput:%sref:%sour:%s", i+1, input, refOut, expect)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
