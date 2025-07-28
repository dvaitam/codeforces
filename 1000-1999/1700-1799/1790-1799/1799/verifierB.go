package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const numTestsB = 100

func prepareBinary(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp := filepath.Join(os.TempDir(), "binB")
		cmd := exec.Command("go", "build", "-o", tmp, path)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", nil, fmt.Errorf("go build failed: %v: %s", err, out)
		}
		return tmp, func() { os.Remove(tmp) }, nil
	}
	return path, nil, nil
}

func prepareOracle() (string, func(), error) {
	tmp := filepath.Join(os.TempDir(), "oracleB")
	cmd := exec.Command("go", "build", "-o", tmp, "1799B.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", nil, fmt.Errorf("go build oracle failed: %v: %s", err, out)
	}
	return tmp, func() { os.Remove(tmp) }, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	return strings.TrimSpace(buf.String()), err
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(5) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d\n", n)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", rng.Intn(9)+1)
	}
	sb.WriteByte('\n')
	return sb.String()
}

func simulate(arr []int, ops [][2]int) ([]int, bool) {
	for _, op := range ops {
		i, j := op[0]-1, op[1]-1
		if i < 0 || i >= len(arr) || j < 0 || j >= len(arr) || i == j || arr[j] == 0 {
			return nil, false
		}
		arr[i] = (arr[i] + arr[j] - 1) / arr[j]
	}
	return arr, true
}

func parseOutput(out string, n int) (int, [][2]int, error) {
	scanner := bufio.NewScanner(strings.NewReader(out))
	if !scanner.Scan() {
		return 0, nil, fmt.Errorf("no output")
	}
	var q int
	if _, err := fmt.Sscan(scanner.Text(), &q); err != nil {
		return 0, nil, fmt.Errorf("bad q")
	}
	if q == -1 {
		if scanner.Scan() {
			return 0, nil, fmt.Errorf("extra output after -1")
		}
		return -1, nil, nil
	}
	if q > 30*n {
		return 0, nil, fmt.Errorf("too many operations")
	}
	ops := make([][2]int, q)
	for i := 0; i < q; i++ {
		if !scanner.Scan() {
			return 0, nil, fmt.Errorf("missing op")
		}
		var x, y int
		if _, err := fmt.Sscan(scanner.Text(), &x, &y); err != nil {
			return 0, nil, fmt.Errorf("bad op line")
		}
		ops[i] = [2]int{x, y}
	}
	if scanner.Scan() {
		return 0, nil, fmt.Errorf("extra output")
	}
	return q, ops, nil
}

func runCase(bin, oracle, input string) error {
	expectStr, err := run(oracle, input)
	if err != nil {
		return fmt.Errorf("oracle error: %v", err)
	}
	// parse input
	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Scan() // t=1
	var n int
	fmt.Sscan(scanner.Text())
	scanner.Scan()
	arr := []int{}
	for _, s := range strings.Fields(scanner.Text()) {
		var v int
		fmt.Sscan(s, &v)
		arr = append(arr, v)
	}
	gotStr, err := run(bin, input)
	if err != nil {
		return fmt.Errorf("runtime error: %v", err)
	}
	q, ops, err := parseOutput(gotStr, n)
	if err != nil {
		return err
	}
	if q == -1 {
		if expectStr != "-1" {
			return fmt.Errorf("expected solution but got -1")
		}
		return nil
	}
	if expectStr == "-1" {
		return fmt.Errorf("expected -1 but provided operations")
	}
	final, ok := simulate(append([]int(nil), arr...), ops)
	if !ok {
		return fmt.Errorf("invalid operations")
	}
	eq := true
	for i := 1; i < len(final); i++ {
		if final[i] != final[0] {
			eq = false
			break
		}
	}
	if !eq {
		return fmt.Errorf("array not equal after operations")
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}
	binPath := os.Args[1]
	bin, cleanup, err := prepareBinary(binPath)
	if err != nil {
		fmt.Println("compile error:", err)
		return
	}
	if cleanup != nil {
		defer cleanup()
	}
	oracle, cleanOracle, err := prepareOracle()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer cleanOracle()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < numTestsB; i++ {
		input := generateCase(rng)
		if err := runCase(bin, oracle, input); err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", i+1, err, input)
			return
		}
	}
	fmt.Println("All tests passed")
}
