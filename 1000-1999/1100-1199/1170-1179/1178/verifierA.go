package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const numTestsA = 100

func prepareBinary(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp := filepath.Join(os.TempDir(), "verifA_bin")
		cmd := exec.Command("go", "build", "-o", tmp, path)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", nil, fmt.Errorf("go build failed: %v: %s", err, out)
		}
		return tmp, func() { os.Remove(tmp) }, nil
	}
	return path, nil, nil
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

func solveA(a []int) string {
	n := len(a)
	total := 0
	for _, v := range a {
		total += v
	}
	leader := a[0]
	supporters := []int{}
	curr := leader
	for i := 1; i < n; i++ {
		if a[i]*2 <= leader {
			curr += a[i]
			supporters = append(supporters, i+1)
		}
	}
	if curr*2 <= total {
		return "0"
	}
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(supporters) + 1))
	sb.WriteByte('\n')
	sb.WriteString("1")
	for _, idx := range supporters {
		sb.WriteByte(' ')
		sb.WriteString(strconv.Itoa(idx))
	}
	return sb.String()
}

func genCaseA(rng *rand.Rand) []int {
	n := rng.Intn(99) + 2
	a := make([]int, n)
	for i := range a {
		a[i] = rng.Intn(100) + 1
	}
	return a
}

func runCaseA(bin string, a []int) error {
	var sb strings.Builder
	n := len(a)
	sb.WriteString(strconv.Itoa(n))
	sb.WriteByte('\n')
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	input := sb.String()
	expected := solveA(a)
	out, err := run(bin, input)
	if err != nil {
		return fmt.Errorf("runtime error: %v", err)
	}
	if out != expected {
		return fmt.Errorf("expected %q got %q", expected, out)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		return
	}
	bin, cleanup, err := prepareBinary(os.Args[1])
	if err != nil {
		fmt.Println("compile error:", err)
		return
	}
	if cleanup != nil {
		defer cleanup()
	}
	rng := rand.New(rand.NewSource(1))
	for t := 0; t < numTestsA; t++ {
		a := genCaseA(rng)
		if err := runCaseA(bin, a); err != nil {
			fmt.Printf("case %d failed: %v\n", t+1, err)
			return
		}
	}
	fmt.Println("All tests passed")
}
