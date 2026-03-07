package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func runExe(path string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func buildRef() (string, error) {
	ref := "./refB.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1635B.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
	}
	return ref, nil
}

func buildCase(arr []int) []byte {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", len(arr)))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return []byte(sb.String())
}

func genRandomCase(rng *rand.Rand) []byte {
	n := rng.Intn(19) + 2 // 2..20
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(1_000_000_000) + 1
	}
	return buildCase(arr)
}

func genTests() [][]byte {
	rng := rand.New(rand.NewSource(43))
	tests := make([][]byte, 0, 100)
	tests = append(tests, buildCase([]int{1, 2}))
	tests = append(tests, buildCase([]int{5, 3, 4, 2, 1}))
	for len(tests) < 100 {
		tests = append(tests, genRandomCase(rng))
	}
	return tests
}

func parseRefOps(s string) (int, error) {
	r := bufio.NewReader(strings.NewReader(strings.TrimSpace(s)))
	var m int
	if _, err := fmt.Fscan(r, &m); err != nil {
		return 0, fmt.Errorf("parse ops: %v", err)
	}
	return m, nil
}

func parseCandidate(s string) (int, []int, error) {
	r := bufio.NewReader(strings.NewReader(strings.TrimSpace(s)))
	var m int
	if _, err := fmt.Fscan(r, &m); err != nil {
		return 0, nil, fmt.Errorf("parse ops: %v", err)
	}
	var arr []int
	for {
		var v int
		if _, err := fmt.Fscan(r, &v); err != nil {
			break
		}
		arr = append(arr, v)
	}
	return m, arr, nil
}

func parseInput(s string) []int {
	r := bufio.NewReader(strings.NewReader(s))
	var t, n int
	fmt.Fscan(r, &t, &n)
	arr := make([]int, n)
	for i := range arr {
		fmt.Fscan(r, &arr[i])
	}
	return arr
}

func check(gotOps, expOps int, gotArr, inArr []int) error {
	n := len(inArr)
	if len(gotArr) != n {
		return fmt.Errorf("expected array of length %d, got %d", n, len(gotArr))
	}
	if gotOps != expOps {
		return fmt.Errorf("wrong op count: expected %d, got %d", expOps, gotOps)
	}
	diffs := 0
	for i := range inArr {
		if inArr[i] != gotArr[i] {
			diffs++
		}
	}
	if diffs != gotOps {
		return fmt.Errorf("array differs in %d positions from input, but claimed %d ops", diffs, gotOps)
	}
	for i := 1; i < n-1; i++ {
		if gotArr[i] > gotArr[i-1] && gotArr[i] > gotArr[i+1] {
			return fmt.Errorf("local maximum at index %d (value %d)", i, gotArr[i])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		return
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
		expOut, err := runExe(ref, tc)
		if err != nil {
			fmt.Printf("reference runtime error on test %d: %v\n%s", i+1, err, expOut)
			os.Exit(1)
		}
		expOps, err := parseRefOps(expOut)
		if err != nil {
			fmt.Printf("failed to parse reference output on test %d: %v\n", i+1, err)
			os.Exit(1)
		}

		gotOut, err := runExe(bin, tc)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n%s", i+1, err, gotOut)
			os.Exit(1)
		}
		gotOps, gotArr, err := parseCandidate(gotOut)
		if err != nil {
			fmt.Printf("failed to parse candidate output on test %d: %v\n", i+1, err)
			os.Exit(1)
		}

		inArr := parseInput(string(tc))
		if err := check(gotOps, expOps, gotArr, inArr); err != nil {
			fmt.Printf("test %d failed: %v\ninput:\n%sexpected ops: %d\ngot:\n%s\n", i+1, err, string(tc), expOps, gotOut)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
