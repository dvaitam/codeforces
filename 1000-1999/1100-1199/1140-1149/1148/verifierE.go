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

func run(bin string, in []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func buildRef() (string, error) {
	ref := "./refE.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1148E.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
	}
	return ref, nil
}

func genCase() (string, []int, []int) {
	n := rand.Intn(5) + 1
	a := make([]int, n)
	b := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = rand.Intn(20)
		b[i] = rand.Intn(20)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(a[i]))
	}
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(b[i]))
	}
	sb.WriteByte('\n')
	return sb.String(), a, b
}

func simulate(a []int, ops [][3]int) ([]int, error) {
	for _, op := range ops {
		i, j, d := op[0], op[1], op[2]
		if i < 1 || i > len(a) || j < 1 || j > len(a) {
			return nil, fmt.Errorf("index out of range")
		}
		if a[i-1] > a[j-1] {
			i, j = j, i
		}
		if 2*d < 0 || 2*d > a[j-1]-a[i-1] {
			return nil, fmt.Errorf("invalid d")
		}
		a[i-1] += d
		a[j-1] -= d
	}
	return a, nil
}

func equalMultiset(x, y []int) bool {
	if len(x) != len(y) {
		return false
	}
	x2 := append([]int(nil), x...)
	y2 := append([]int(nil), y...)
	sort.Ints(x2)
	sort.Ints(y2)
	for i := range x2 {
		if x2[i] != y2[i] {
			return false
		}
	}
	return true
}

func parseOutput(out string) (bool, [][3]int, error) {
	out = strings.TrimSpace(out)
	if strings.HasPrefix(out, "NO") {
		return false, nil, nil
	}
	lines := strings.Split(out, "\n")
	if len(lines) < 2 {
		return false, nil, fmt.Errorf("bad output")
	}
	if strings.TrimSpace(lines[0]) != "YES" {
		return false, nil, fmt.Errorf("expected YES or NO")
	}
	m, err := strconv.Atoi(strings.TrimSpace(lines[1]))
	if err != nil {
		return false, nil, fmt.Errorf("bad m")
	}
	if m < 0 || m > 5*len(lines) {
		return false, nil, fmt.Errorf("bad m")
	}
	if len(lines)-2 != m {
		return false, nil, fmt.Errorf("expected %d operations", m)
	}
	ops := make([][3]int, m)
	for i := 0; i < m; i++ {
		parts := strings.Fields(lines[i+2])
		if len(parts) != 3 {
			return false, nil, fmt.Errorf("op %d format", i+1)
		}
		a1, _ := strconv.Atoi(parts[0])
		b1, _ := strconv.Atoi(parts[1])
		d, _ := strconv.Atoi(parts[2])
		ops[i] = [3]int{a1, b1, d}
	}
	return true, ops, nil
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

	rand.Seed(time.Now().UnixNano())
	for t := 0; t < 100; t++ {
		input, a, b := genCase()
		refOut, err := run(ref, []byte(input))
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d: %v\n", t+1, err)
			os.Exit(1)
		}
		possible := !strings.HasPrefix(strings.TrimSpace(refOut), "NO")
		candOut, err := run(bin, []byte(input))
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%s", t+1, err, input)
			os.Exit(1)
		}
		yes, ops, err := parseOutput(candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "bad output on test %d: %v\ninput:\n%soutput:\n%s", t+1, err, input, candOut)
			os.Exit(1)
		}
		if yes != possible {
			fmt.Fprintf(os.Stderr, "incorrect YES/NO on test %d\ninput:\n%s", t+1, input)
			os.Exit(1)
		}
		if yes {
			arr := append([]int(nil), a...)
			arr2, err := simulate(arr, ops)
			if err != nil {
				fmt.Fprintf(os.Stderr, "invalid operations on test %d: %v\ninput:\n%soutput:\n%s", t+1, err, input, candOut)
				os.Exit(1)
			}
			if !equalMultiset(arr2, b) {
				fmt.Fprintf(os.Stderr, "wrong final array on test %d\ninput:\n%soutput:\n%s", t+1, input, candOut)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed.")
}
