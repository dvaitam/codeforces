package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func buildRef() (string, error) {
	src := os.Getenv("REFERENCE_SOURCE_PATH")
	if src == "" {
		_, file, _, _ := runtime.Caller(0)
		dir := filepath.Dir(file)
		src = filepath.Join(dir, "623D.go")
	}
	tmp, err := os.CreateTemp("", "refD-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	cmd := exec.Command("go", "build", "-o", tmp.Name(), src)
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out)
	}
	return tmp.Name(), nil
}

func randProbs(rng *rand.Rand, n int) []int {
	probs := make([]int, n)
	remain := 100
	for i := 0; i < n-1; i++ {
		maxv := remain - (n - 1 - i)
		v := rng.Intn(maxv) + 1
		probs[i] = v
		remain -= v
	}
	probs[n-1] = remain
	return probs
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(5) + 1
	ps := randProbs(rng, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range ps {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runCase(bin, ref, input string) error {
	cmdRef := exec.Command(ref)
	cmdRef.Stdin = strings.NewReader(input)
	var refOut bytes.Buffer
	cmdRef.Stdout = &refOut
	cmdRef.Stderr = &refOut
	if err := cmdRef.Run(); err != nil {
		return fmt.Errorf("reference runtime error: %v\n%s", err, refOut.String())
	}
	expected := strings.TrimSpace(refOut.String())

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got == expected {
		return nil
	}
	// Tolerance-based float comparison
	gotF, errG := strconv.ParseFloat(got, 64)
	expF, errE := strconv.ParseFloat(expected, 64)
	if errG != nil || errE != nil {
		return fmt.Errorf("expected %q got %q", expected, got)
	}
	diff := math.Abs(gotF - expF)
	rel := diff / math.Max(1.0, math.Abs(expF))
	if diff > 1e-6 && rel > 1e-6 {
		return fmt.Errorf("expected %q got %q (diff=%e)", expected, got, diff)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in := generateCase(rng)
		if err := runCase(bin, ref, in); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
