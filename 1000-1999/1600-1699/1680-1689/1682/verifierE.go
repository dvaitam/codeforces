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
	"time"
)

func prepareBinary(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp := filepath.Join(os.TempDir(), "verifE_bin")
		cmd := exec.Command("go", "build", "-o", tmp, path)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", nil, fmt.Errorf("go build failed: %v: %s", err, out)
		}
		return tmp, func() { os.Remove(tmp) }, nil
	}
	return path, nil, nil
}

func prepareOracle() (string, func(), error) {
	tmp := filepath.Join(os.TempDir(), "oracleE_bin")
	cmd := exec.Command("go", "build", "-o", tmp, "1682E.go")
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

func minimalSwaps(p []int) [][2]int {
	n := len(p)
	vis := make([]bool, n)
	res := make([][2]int, 0)
	for i := 0; i < n; i++ {
		if vis[i] {
			continue
		}
		j := i
		cycle := []int{}
		for !vis[j] {
			vis[j] = true
			cycle = append(cycle, j)
			j = p[j] - 1
		}
		if len(cycle) > 1 {
			pivot := cycle[0]
			for k := 1; k < len(cycle); k++ {
				res = append(res, [2]int{pivot + 1, cycle[k] + 1})
			}
		}
	}
	return res
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(4) + 2
	var perm []int
	var swaps [][2]int
	for {
		perm = rng.Perm(n)
		for i := range perm {
			perm[i]++
		}
		swaps = minimalSwaps(append([]int(nil), perm...))
		if len(swaps) > 0 {
			break
		}
	}
	m := len(swaps)
	rng.Shuffle(m, func(i, j int) { swaps[i], swaps[j] = swaps[j], swaps[i] })
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for i, v := range perm {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	for _, sw := range swaps {
		fmt.Fprintf(&sb, "%d %d\n", sw[0], sw[1])
	}
	return sb.String()
}

func runCase(bin, oracle, input string) error {
	exp, err := run(oracle, input)
	if err != nil {
		return fmt.Errorf("oracle error: %v", err)
	}
	got, err := run(bin, input)
	if err != nil {
		return fmt.Errorf("runtime error: %v", err)
	}
	if strings.TrimSpace(got) != strings.TrimSpace(exp) {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}
	bin, clean, err := prepareBinary(os.Args[1])
	if err != nil {
		fmt.Println("compile error:", err)
		return
	}
	if clean != nil {
		defer clean()
	}
	oracle, cleanO, err := prepareOracle()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer cleanO()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := genCase(rng)
		if err := runCase(bin, oracle, input); err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", i+1, err, input)
			return
		}
	}
	fmt.Println("All tests passed")
}
