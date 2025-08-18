package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

func run(prog, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(prog, ".go") {
		cmd = exec.Command("go", "run", prog)
	} else {
		cmd = exec.Command(prog)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func isValidInput(arr []int64) bool {
	if len(arr) == 0 {
		return false
	}
	g := arr[0]
	for _, v := range arr[1:] {
		g = gcd(g, v)
	}
	return g == arr[0]
}

func validateOutput(arr []int64, out string) error {
	valid := isValidInput(arr)
	out = strings.TrimSpace(out)
	if !valid {
		if out != "-1" {
			return fmt.Errorf("expected -1 for invalid case, got %q", out)
		}
		return nil
	}
	// valid case: parse header length and sequence
	lines := strings.Split(out, "\n")
	if len(lines) == 0 {
		return fmt.Errorf("empty output")
	}
	header := strings.Fields(lines[0])
	if len(header) == 0 {
		return fmt.Errorf("empty header line")
	}
	L, err := strconv.Atoi(header[0])
	if err != nil {
		return fmt.Errorf("cannot parse length: %v", err)
	}
	n := len(arr)
	if !(L == 2*n-1 || L == 2*(n-1)) {
		return fmt.Errorf("invalid length %d for n=%d", L, n)
	}
	seqFields := []string{}
	for i := 1; i < len(lines); i++ {
		seqFields = append(seqFields, strings.Fields(lines[i])...)
	}
	if len(seqFields) != L {
		return fmt.Errorf("expected %d numbers, got %d", L, len(seqFields))
	}
	base := arr[0]
	cntBase := 0
	other := make(map[int64]int)
	for _, s := range seqFields {
		v, e := strconv.ParseInt(s, 10, 64)
		if e != nil {
			return fmt.Errorf("non-integer token %q", s)
		}
		if v == base {
			cntBase++
		} else {
			other[v]++
		}
	}
	// base count check
	if L == 2*n-1 {
		if cntBase != n {
			return fmt.Errorf("base %d count %d != %d", base, cntBase, n)
		}
	} else { // 2(n-1)
		if cntBase != n-1 {
			return fmt.Errorf("base %d count %d != %d", base, cntBase, n-1)
		}
	}
	// other numbers must match arr[1..] exactly once each
	for i := 1; i < n; i++ {
		v := arr[i]
		if other[v] != 1 {
			return fmt.Errorf("value %d appears %d times (expected 1)", v, other[v])
		}
		delete(other, v)
	}
	if len(other) != 0 {
		return fmt.Errorf("unexpected values present: %v", other)
	}
	return nil
}

func genValidCase(rng *rand.Rand) []int64 {
	m := rng.Intn(5) + 1
	base := int64(rng.Intn(1000) + 1)
	set := make(map[int64]struct{})
	arr := make([]int64, 0, m)
	arr = append(arr, base)
	set[base] = struct{}{}
	for len(arr) < m {
		v := base * int64(rng.Intn(20)+1)
		if _, ok := set[v]; ok {
			continue
		}
		set[v] = struct{}{}
		arr = append(arr, v)
	}
	sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
	return arr
}

func genInvalidCase(rng *rand.Rand) []int64 {
	arr := genValidCase(rng)
	if len(arr) == 1 {
		arr[0]++
		return arr
	}
	arr[0]++
	sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
	return arr
}

func buildInput(arr []int64) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(arr))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		var arr []int64
		if t%2 == 0 {
			arr = genValidCase(rng)
		} else {
			arr = genInvalidCase(rng)
		}
		input := buildInput(arr)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", t+1, err, input)
			os.Exit(1)
		}
		if err := validateOutput(arr, got); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s\noutput:%s", t+1, err, input, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
