package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func runBinary(bin string, input []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func buildRef() (string, error) {
	_, cur, _, _ := runtime.Caller(0)
	dir := filepath.Dir(cur)
	src := filepath.Join(dir, "1151B.go")
	refBin := filepath.Join(os.TempDir(), "1151B_ref.bin")
	cmd := exec.Command("go", "build", "-o", refBin, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v\n%s", err, string(out))
	}
	return refBin, nil
}

func genTest() []byte {
	n := rand.Intn(4) + 1
	m := rand.Intn(4) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			val := rand.Intn(1024)
			sb.WriteString(fmt.Sprintf("%d", val))
			if j+1 < m {
				sb.WriteByte(' ')
			}
		}
		sb.WriteByte('\n')
	}
	return []byte(sb.String())
}

func parseInput(in []byte) (int, int, [][]int) {
	r := bytes.NewReader(in)
	var n, m int
	fmt.Fscan(r, &n, &m)
	a := make([][]int, n)
	for i := 0; i < n; i++ {
		a[i] = make([]int, m)
		for j := 0; j < m; j++ {
			fmt.Fscan(r, &a[i][j])
		}
	}
	return n, m, a
}

func checkOutput(out string, hasSolution bool, n, m int, a [][]int) error {
	tokens := strings.Fields(out)
	if len(tokens) == 0 {
		return fmt.Errorf("empty output")
	}
	if !hasSolution {
		if tokens[0] != "NIE" {
			return fmt.Errorf("expected NIE, got %q", tokens[0])
		}
		return nil
	}
	if tokens[0] != "TAK" {
		return fmt.Errorf("expected TAK, got %q", tokens[0])
	}
	if len(tokens) != n+1 {
		return fmt.Errorf("expected %d numbers, got %d", n, len(tokens)-1)
	}
	xor := 0
	for i := 0; i < n; i++ {
		v, err := strconv.Atoi(tokens[i+1])
		if err != nil || v < 1 || v > m {
			return fmt.Errorf("invalid choice at position %d", i+1)
		}
		xor ^= a[i][v-1]
	}
	if xor == 0 {
		return fmt.Errorf("xor is zero")
	}
	return nil
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}
	cand := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	for i := 1; i <= 100; i++ {
		in := genTest()
		n, m, a := parseInput(in)
		expected, err := runBinary(ref, in)
		if err != nil {
			fmt.Printf("reference failed on test %d: %v\n", i, err)
			os.Exit(1)
		}
		tokens := strings.Fields(expected)
		hasSolution := len(tokens) > 0 && tokens[0] == "TAK"
		got, err := runBinary(cand, in)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n", i, err)
			os.Exit(1)
		}
		if err := checkOutput(got, hasSolution, n, m, a); err != nil {
			fmt.Printf("wrong answer on test %d\ninput:\n%serror: %v\noutput:%s\n", i, string(in), err, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
