package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func runProgram(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		tmpDir, err := os.MkdirTemp("", "cand")
		if err != nil {
			return "", err
		}
		defer os.RemoveAll(tmpDir)
		data, err := os.ReadFile(bin)
		if err != nil {
			return "", err
		}
		tmpSrc := filepath.Join(tmpDir, filepath.Base(bin))
		if err := os.WriteFile(tmpSrc, data, 0644); err != nil {
			return "", err
		}
		cmd = exec.Command("go", "run", tmpSrc)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}


func generateCase(rng *rand.Rand) string {
	var sb strings.Builder
	n := rng.Intn(99999) + 2    // 2 to 100000
	m := rng.Intn(99999) + 2    // 2 to 100000
	k := rng.Intn(100) + 1 // 1 to 100 queries

	fmt.Fprintf(&sb, "%d %d %d\n", n, m, k)

	for i := 0; i < k; i++ {
		x := rng.Intn(n-1) + 1 // 1 to n-1
		y := rng.Intn(m-1) + 1 // 1 to m-1
		fmt.Fprintf(&sb, "%d %d\n", x, y)
	}
	return sb.String()
}

func verify(candidate, refSrc string) error {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	srcPath := filepath.Join(dir, refSrc)
	tmpDir, err := os.MkdirTemp("", "refbuild")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)
	tmpSrc := filepath.Join(tmpDir, filepath.Base(srcPath))
	data, err := os.ReadFile(srcPath)
	if err != nil {
		return err
	}
	if err := os.WriteFile(tmpSrc, data, 0644); err != nil {
		return err
	}
	refBin := filepath.Join(tmpDir, "refbin")
	cmd := exec.Command("go", "build", "-o", refBin, tmpSrc)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to build reference: %v\n%s", err, out)
	}
	candPath := candidate
	var cleanup func()
	if strings.HasSuffix(candidate, ".go") {
		tmpCdir, err := os.MkdirTemp("", "candbuild")
		if err != nil {
			return err
		}
		data, err := os.ReadFile(candidate)
		if err != nil {
			os.RemoveAll(tmpCdir)
			return err
		}
		tmpSrc := filepath.Join(tmpCdir, filepath.Base(candidate))
		if err := os.WriteFile(tmpSrc, data, 0644); err != nil {
			os.RemoveAll(tmpCdir)
			return err
		}
		candBin := filepath.Join(tmpCdir, "candbin")
		if out, err := exec.Command("go", "build", "-o", candBin, tmpSrc).CombinedOutput(); err != nil {
			os.RemoveAll(tmpCdir)
			return fmt.Errorf("failed to build candidate: %v\n%s", err, out)
		}
		candPath = candBin
		cleanup = func() { os.RemoveAll(tmpCdir) }
	}
	if cleanup != nil {
		defer cleanup()
	}

	var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 100; i++ {
		in := generateCase(rng)
		candOut, err := runProgram(candPath, in)
		if err != nil {
			return fmt.Errorf("case %d: %v", i+1, err)
		}
		refOut, err := runProgram(refBin, in)
		if err != nil {
			return fmt.Errorf("reference failed on case %d: %v", i+1, err)
		}

		canonicalize := func(s string) string {
			fields := strings.Fields(s)
			return strings.Join(fields, " ")
		}
		if canonicalize(candOut) != canonicalize(refOut) {
			return fmt.Errorf("case %d failed: expected %q got %q\nInput:\n%s", i+1, refOut, candOut, in)
		}
	}
	fmt.Printf("All %d tests passed\n", 100)
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]
	if err := verify(candidate, "724C.go"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
