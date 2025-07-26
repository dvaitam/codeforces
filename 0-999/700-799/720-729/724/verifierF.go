package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
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

func readTests(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	raw := strings.TrimSpace(string(data))
	if !strings.Contains(raw, "\n\n") {
		lines := strings.Split(raw, "\n")
		tests := make([]string, 0, len(lines))
		for _, ln := range lines {
			ln = strings.TrimSpace(ln)
			if ln == "" {
				continue
			}
			tests = append(tests, ln+"\n")
		}
		return tests, nil
	}
	parts := strings.Split(raw, "\n\n")
	tests := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.Trim(p, "\n")
		if p == "" {
			continue
		}
		tests = append(tests, p+"\n")
	}
	return tests, nil
}

func verify(candidate, refSrc, testFile string) error {
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

	tests, err := readTests(filepath.Join(dir, testFile))
	if err != nil {
		return err
	}
	for i, in := range tests {
		candOut, err := runProgram(candPath, in)
		if err != nil {
			return fmt.Errorf("case %d: %v", i+1, err)
		}
		refOut, err := runProgram(refBin, in)
		if err != nil {
			return fmt.Errorf("reference failed on case %d: %v", i+1, err)
		}
		if strings.TrimSpace(candOut) != strings.TrimSpace(refOut) {
			return fmt.Errorf("case %d failed: expected %q got %q", i+1, refOut, candOut)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]
	if err := verify(candidate, "724F.go", "testcasesF.txt"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
