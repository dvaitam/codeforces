package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func compileRef() (string, error) {
	out := "refE.bin"
	cmd := exec.Command("go", "build", "-o", out, "1109E.go")
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return "./" + out, nil
}

func runBin(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
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

func genCase(rng *rand.Rand) string {
	n := rng.Intn(4) + 1
	mod := rng.Intn(50) + 2
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, mod))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(rng.Intn(20) + 1))
	}
	sb.WriteByte('\n')
	q := rng.Intn(5) + 1
	sb.WriteString(fmt.Sprintf("%d\n", q))
	for i := 0; i < q; i++ {
		op := rng.Intn(3) + 1
		if op == 1 {
			l := rng.Intn(n)
			r := rng.Intn(n-l) + l
			x := rng.Intn(20) + 1
			sb.WriteString(fmt.Sprintf("1 %d %d %d\n", l+1, r+1, x))
		} else if op == 2 {
			u := rng.Intn(n)
			x := rng.Intn(20) + 1
			sb.WriteString(fmt.Sprintf("2 %d %d\n", u+1, x))
		} else {
			l := rng.Intn(n)
			r := rng.Intn(n-l) + l
			sb.WriteString(fmt.Sprintf("3 %d %d\n", l+1, r+1))
		}
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref, err := compileRef()
	if err != nil {
		fmt.Println("failed to compile reference:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in := genCase(rng)
		exp, err := runBin(ref, in)
		if err != nil {
			fmt.Printf("reference runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		out, err := runBin(bin, in)
		if err != nil {
			fmt.Printf("test %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Printf("test %d failed\nInput:\n%sExpected: %s\nGot: %s\n", i+1, in, exp, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
