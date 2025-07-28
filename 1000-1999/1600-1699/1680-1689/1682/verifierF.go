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
		tmp := filepath.Join(os.TempDir(), "verifF_bin")
		cmd := exec.Command("go", "build", "-o", tmp, path)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", nil, fmt.Errorf("go build failed: %v: %s", err, out)
		}
		return tmp, func() { os.Remove(tmp) }, nil
	}
	return path, nil, nil
}

func prepareOracle() (string, func(), error) {
	tmp := filepath.Join(os.TempDir(), "oracleF_bin")
	cmd := exec.Command("go", "build", "-o", tmp, "1682F.go")
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

func genCase(rng *rand.Rand) string {
	n := rng.Intn(4) + 2
	a := make([]int64, n)
	a[0] = int64(rng.Intn(10))
	for i := 1; i < n; i++ {
		a[i] = a[i-1] + int64(rng.Intn(5)+1)
	}
	prefix := make([]int64, n+1)
	b := make([]int64, n)
	for i := 1; i < n; i++ {
		delta := int64(rng.Intn(5) + 1)
		if rng.Intn(2) == 0 {
			delta = -delta
		}
		b[i-1] = delta
		prefix[i] = prefix[i-1] + delta
	}
	idx := rng.Intn(n)
	prefix[n] = prefix[idx]
	b[n-1] = prefix[n] - prefix[n-1]
	if b[n-1] == 0 {
		b[n-1] = 1
		prefix[n] = prefix[n-1] + b[n-1]
	}
	pairs := [][2]int{}
	for i := 0; i < n; i++ {
		for j := i + 1; j <= n; j++ {
			if prefix[i] == prefix[j] {
				pairs = append(pairs, [2]int{i + 1, j})
			}
		}
	}
	rng.Shuffle(len(pairs), func(i, j int) { pairs[i], pairs[j] = pairs[j], pairs[i] })
	q := rng.Intn(len(pairs)) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, q)
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	sb.WriteByte('\n')
	for i, v := range b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	sb.WriteByte('\n')
	for i := 0; i < q; i++ {
		fmt.Fprintf(&sb, "%d %d\n", pairs[i][0], pairs[i][1])
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
		fmt.Println("usage: go run verifierF.go /path/to/binary")
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
