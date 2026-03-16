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

const testcasesRaw = `1 0

3 1
3 1 17

1 0

1 0

1 0

2 0

5 1
5 4 7

4 4
3 1 5
4 3 8
1 4 3
2 1 1

2 1
1 2 8

4 2
2 3 11
1 2 17

2 0

4 3
3 2 10
1 2 1
3 4 8

1 0

5 2
2 4 12
3 2 7

5 4
3 5 13
5 4 11
1 2 20
2 4 19

1 0

4 4
4 3 17
1 3 13
2 4 0
3 2 16

1 0

5 4
3 2 17
5 1 19
3 4 0
1 3 9

2 0

2 0

1 0

2 1
1 2 16

5 3
2 5 6
3 4 20
5 4 3

2 0

1 0

1 0

5 1
1 2 2

1 0

1 0

3 3
1 3 4
3 2 7
2 1 7

2 0

5 3
2 4 5
3 4 7
1 4 17

1 0

5 0

1 0

2 1
1 2 1

2 1
1 2 8

4 2
3 2 1
1 3 1

1 0

5 4
2 1 16
1 5 2
2 4 3
5 2 18

5 0

5 0

4 4
2 3 7
3 4 4
3 1 0
4 1 2

5 1
5 3 4

3 0

2 1
2 1 14

5 2
1 5 9
1 2 8

1 0

5 1
5 2 10

2 1
2 1 20

3 1
3 2 17

4 4
1 2 17
1 3 18
2 4 4
3 2 7

1 0

5 3
5 2 7
4 1 5
3 4 7

3 1
3 1 12

1 0

2 0

4 2
3 2 7
1 2 12

3 2
1 2 11
2 3 17

3 0

1 0

2 1
1 2 18

2 1
1 2 0

5 4
2 3 13
1 3 19
3 5 9
4 3 12

3 1
1 2 12

2 1
2 1 9

3 1
2 3 19

3 3
3 1 16
2 3 5
1 2 13

5 1
4 2 4

1 0

4 1
2 4 1

5 1
1 4 4

4 4
3 4 19
2 4 14
3 2 20
4 1 9

2 1
2 1 6

1 0

4 2
1 2 13
4 1 18

4 3
1 3 9
2 4 7
3 4 15

1 0

3 3
3 1 14
2 3 18
2 1 14

2 0

3 3
2 1 14
1 3 17
2 3 8

3 1
2 1 18

1 0

5 3
4 2 2
5 2 3
3 5 3

5 0

3 3
3 2 2
2 1 13
3 1 13

2 1
2 1 2

3 3
1 2 18
2 3 0
3 1 2

2 1
2 1 15

4 3
3 2 12
3 4 17
2 1 6

3 0

5 1
4 3 18

5 4
3 1 6
3 2 11
1 5 4
4 1 0

5 2
4 3 5
1 3 15

1 0

4 3
3 1 7
1 4 19
2 4 14

4 2
4 3 18
3 1 5

2 0`

func readTests() ([]string, error) {
	raw := strings.TrimSpace(testcasesRaw)
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

	tests, err := readTests()
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]
	if err := verify(candidate, "724E.go"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
