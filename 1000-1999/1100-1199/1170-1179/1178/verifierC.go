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

const mod = 998244353
const numTestsC = 100

func prepareBinary(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp := filepath.Join(os.TempDir(), "verifC_bin")
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

func modPow(a, e int) int {
	res := 1
	a %= mod
	for e > 0 {
		if e&1 == 1 {
			res = int(int64(res) * int64(a) % mod)
		}
		a = int(int64(a) * int64(a) % mod)
		e >>= 1
	}
	return res
}

func solveC(w, h int) string {
	return strconv.Itoa(modPow(2, w+h))
}

func genCaseC(rng *rand.Rand) (int, int) {
	w := rng.Intn(1000) + 1
	h := rng.Intn(1000) + 1
	return w, h
}

func runCaseC(bin string, w, h int) error {
	input := fmt.Sprintf("%d %d\n", w, h)
	expected := solveC(w, h)
	out, err := run(bin, input)
	if err != nil {
		return fmt.Errorf("runtime error: %v", err)
	}
	if out != expected {
		return fmt.Errorf("expected %s got %s", expected, out)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
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
	rng := rand.New(rand.NewSource(3))
	for t := 0; t < numTestsC; t++ {
		w, h := genCaseC(rng)
		if err := runCaseC(bin, w, h); err != nil {
			fmt.Printf("case %d failed: %v\n", t+1, err)
			return
		}
	}
	fmt.Println("All tests passed")
}
