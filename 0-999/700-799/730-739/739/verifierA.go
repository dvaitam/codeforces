package main

import (
    "bytes"
    "fmt"
    "math/rand"
    "os"
    "os/exec"
    "strconv"
    "strings"
)

type Test struct {
	input string
}

func runExe(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func genTests() []Test {
	rand.Seed(0)
	tests := make([]Test, 0, 101)
	for i := 0; i < 100; i++ {
		n := rand.Intn(20) + 1
		m := rand.Intn(20) + 1
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, m)
		for j := 0; j < m; j++ {
			l := rand.Intn(n) + 1
			r := l + rand.Intn(n-l+1)
			fmt.Fprintf(&sb, "%d %d\n", l, r)
		}
		tests = append(tests, Test{sb.String()})
	}
	tests = append(tests, Test{"1 1\n1 1\n"})
	return tests
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("Usage: go run verifierA.go /path/to/binary")
        return
    }
    bin := os.Args[1]

    tests := genTests()
    for i, tc := range tests {
        got, err := runExe(bin, tc.input)
        if err != nil {
            fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
            os.Exit(1)
        }
        // Validate candidate output instead of string-matching reference
        if err := validate(tc.input, got); err != nil {
            fmt.Printf("Test %d failed\nInput:\n%sError: %v\nOutput:\n%s\n", i+1, tc.input, err, got)
            os.Exit(1)
        }
    }
    fmt.Println("all tests passed")
}

func validate(input, output string) error {
    // parse input
    lines := strings.Split(strings.TrimSpace(input), "\n")
    header := strings.Fields(lines[0])
    if len(header) < 2 { return fmt.Errorf("invalid input header") }
    n, _ := strconv.Atoi(header[0])
    m, _ := strconv.Atoi(header[1])
    segs := make([][2]int, m)
    minLen := n
    for i := 0; i < m; i++ {
        parts := strings.Fields(lines[1+i])
        if len(parts) < 2 { return fmt.Errorf("invalid segment line") }
        l, _ := strconv.Atoi(parts[0])
        r, _ := strconv.Atoi(parts[1])
        segs[i] = [2]int{l, r}
        if r-l+1 < minLen { minLen = r-l+1 }
    }
    // parse output ints
    fields := strings.Fields(output)
    if len(fields) < 1+n { return fmt.Errorf("not enough numbers in output") }
    k, err := strconv.Atoi(fields[0])
    if err != nil || k <= 0 || k > n { return fmt.Errorf("invalid k") }
    if k != minLen { return fmt.Errorf("k=%d does not equal min segment length %d", k, minLen) }
    a := make([]int, n)
    for i := 0; i < n; i++ {
        v, err := strconv.Atoi(fields[1+i])
        if err != nil { return fmt.Errorf("invalid label at pos %d", i+1) }
        if v < 0 || v >= k { return fmt.Errorf("label out of range at pos %d", i+1) }
        a[i] = v
    }
    // Check that in every segment [l,r], all residues 0..k-1 appear at least once
    for _, seg := range segs {
        l, r := seg[0]-1, seg[1]-1
        seen := make([]bool, k)
        cnt := 0
        for i := l; i <= r; i++ {
            if !seen[a[i]] { seen[a[i]] = true; cnt++ }
        }
        if cnt < k { return fmt.Errorf("segment [%d,%d] does not contain all labels", seg[0], seg[1]) }
    }
    return nil
}
