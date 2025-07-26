package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const numTestsE = 100

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

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	return strings.TrimSpace(buf.String()), err
}

func solveE(s string) string {
	n := len(s)
	tmp := n >> 2
	ans := make([]byte, tmp)
	for i := 0; i < tmp; i++ {
		k1 := i << 1
		k2 := k1 | 1
		k3 := n - ((i + 1) << 1)
		k4 := k3 + 1
		if s[k1] == s[k3] || s[k1] == s[k4] {
			ans[i] = s[k1]
		} else {
			ans[i] = s[k2]
		}
	}
	var buf bytes.Buffer
	buf.Write(ans)
	if n%4 != 0 {
		buf.WriteByte(s[tmp<<1])
	}
	for i := tmp - 1; i >= 0; i-- {
		buf.WriteByte(ans[i])
	}
	return buf.String()
}

func genCaseE(rng *rand.Rand) string {
	n := rng.Intn(50) + 2
	b := make([]byte, n)
	last := byte(0)
	for i := 0; i < n; i++ {
		for {
			c := byte('a' + rng.Intn(3))
			if c != last {
				b[i] = c
				last = c
				break
			}
		}
	}
	return string(b)
}

func runCaseE(bin string, s string) error {
	input := s + "\n"
	expected := solveE(s)
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
		fmt.Println("usage: go run verifierE.go /path/to/binary")
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
	rng := rand.New(rand.NewSource(5))
	for t := 0; t < numTestsE; t++ {
		s := genCaseE(rng)
		if err := runCaseE(bin, s); err != nil {
			fmt.Printf("case %d failed: %v\n", t+1, err)
			return
		}
	}
	fmt.Println("All tests passed")
}
