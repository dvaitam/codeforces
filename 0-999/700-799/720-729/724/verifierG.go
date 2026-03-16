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
4
3

2 1
1 10
8 1

5 3
0 0 1 3 3
8 9 0 8 3

5 3
3 7 9 4 0
2 6 5 4 2

2 2
1 1
6 1

3 2
9 4 0
7 8 1

4 0
8 4 10 9
5 9 3 1

1 5
3
4

1 1
1
6

3 3
10 5 2
5 5 3

3 5
10 10 1
9 10 2

5 5
3 2 7 6 4
10 8 3 10 5

1 1
0
5

4 2
1 3 9 5
3 10 7 6

4 1
4 2 3 8
8 4 9 6

5 3
5 3 2 8 7
1 0 1 2 10

2 5
6 9
1 6

4 4
7 8 4 8
0 10 1 10

5 2
10 5 1 4 6
2 7 0 4 8

2 4
1 10
4 10

5 4
3 2 5 2 8
8 0 9 5 7

1 0
5
4

2 0
3 9
1 1

4 0
8 2 2 10
7 8 2 4

5 4
6 3 8 3 4
6 10 10 5 7

5 3
1 3 3 1 5
0 9 8 3 9

2 0
1 10
0 3

1 0
5
1

5 1
4 10 7 3 8
2 9 9 7 3

4 3
3 1 1 10
6 5 6 6

4 5
0 10 10 10
1 0 6 5

1 1
3
3

5 3
2 6 2 4 7
3 1 7 8 1

1 5
8
0

1 1
2
6

4 3
3 6 0 2
6 0 6 4

4 2
6 8 10 7
2 3 4 3

1 4
8
0

3 0
0 9 7
8 8 2

1 4
1
2

1 4
1
10

2 3
1 9
3 9

5 0
9 1 6 10 9
9 8 5 4 3

3 1
4 6 2
10 10 4

4 2
1 0 7 9
9 1 1 8

2 4
4 2
5 1

2 2
4 2
7 8

3 4
10 8 0
10 8 4

1 1
4
1

1 5
8
2

3 2
9 3 5
3 10 10

3 4
7 4 0
1 10 6

3 0
0 5 2
10 4 2

4 4
6 8 0 1
1 2 8 0

3 4
8 2 6
2 0 4

3 0
5 3 10
3 10 1

3 4
6 9 2
3 2 2

4 0
2 5 6 10
3 4 2 1

4 0
7 3 3 7
5 4 3 3

1 5
3
6

3 2
1 4 5
10 8 6

5 2
0 1 4 2 9
4 0 1 9 6

3 5
5 6 9
8 1 6

5 1
4 0 6 0 8
8 10 10 3 5

4 0
10 5 9 5
10 1 4 8

3 5
6 5 6
4 8 2

2 3
10 6
10 2

5 4
4 6 8 0 4
4 3 6 9 9

3 3
7 7 10
3 8 7

2 5
1 4
8 10

5 2
1 3 10 4 3
3 2 0 0 3

4 4
1 7 6 10
9 3 6 7

4 1
2 10 0 1
6 3 2 8

4 0
8 3 1 7
2 7 10 8

5 4
5 7 9 8 6
8 7 2 7 7

3 1
10 4 8
7 10 3

3 3
1 4 3
4 5 5

5 0
2 2 3 6 2
3 1 6 6 5

5 3
6 0 3 6 6
9 0 9 6 7

1 2
4
6

4 4
8 9 3 7
3 4 6 7

1 3
5
10

4 5
2 7 2 9
8 0 6 9

5 5
0 1 10 6 2
7 2 0 4 6

3 1
7 5 5
6 4 6

3 0
7 0 8
0 5 3

1 5
0
0

2 1
0 9
2 3

2 3
10 1
9 3

4 5
4 5 2 9
9 1 2 4

1 4
0
4

5 5
6 6 3 1 9
10 3 1 4 10

5 0
9 0 5 8 6
10 5 1 8 10

3 0
6 7 1
6 5 10

4 5
2 6 2 8
10 4 9 8

4 3
6 9 4 5
3 1 4 7

2 3
9 9
10 6

3 0
7 5 2
7 3 5

3 2
4 9 4
8 0 8

2 0
3 6
7 8`

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
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]
	if err := verify(candidate, "724G.go"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
