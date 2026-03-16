package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// Embedded testcases (same format as original file).
const embeddedTestcases = `100
10
3 0
3 5
3 -5
3 -2
3 -5
1 -4 0
2 1
3 3
1 4 -2
1 -2 1

5
1 1 -3
1 -3 4
3 2
1 -3 -5
1 -2 -2

3
1 -1 0
1 3 5
3 -2

3
3 -2
3 -1
1 0 1

3
1 -1 -4
2 1
3 4

1
3 5

6
1 -1 0
2 1
3 0
1 2 2
3 -3
1 -1 -5

6
3 -5
3 1
3 1
3 -5
3 -5
3 -3

10
1 -4 -2
2 1
3 0
3 -1
3 -4
3 0
3 -5
3 -4
1 0 3
3 0

3
3 -1
3 3
1 -1 5

6
3 -3
1 5 -3
3 -1
2 1
3 -5
1 4 3

7
1 -2 4
2 1
3 5
3 -3
1 5 -5
2 1
1 -3 4

3
3 1
1 -3 1
2 1

1
3 -1

3
3 4
1 3 2
2 1

8
3 -1
3 1
1 -4 1
3 -3
3 2
2 1
1 2 -1
3 3

9
3 -4
3 4
3 -5
3 0
3 5
3 2
3 -1
3 5
1 4 -5

8
3 -1
3 5
3 2
3 3
3 5
3 0
3 5
3 1

6
1 2 0
2 1
3 -3
1 0 2
2 1
3 5

7
1 4 4
3 5
2 1
3 3
3 -1
3 -5
1 -3 4

8
3 5
1 -2 5
1 5 -5
2 1
1 -5 -3
1 0 -3
2 2
3 -5

7
3 0
3 5
3 -4
3 -2
1 0 -5
2 1
3 1

2
3 3
3 -5

9
3 -1
1 -1 3
3 0
3 -1
2 1
3 1
3 5
3 0
3 -3

3
3 1
3 2
1 -3 4

2
3 5
1 1 -4

6
3 4
3 -3
3 5
3 1
3 1
1 2 -1

8
3 1
3 -3
3 4
3 -1
3 -1
3 -5
3 -1
3 -1

3
3 -5
1 5 4
2 1

5
1 -3 1
1 2 3
3 -1
1 2 -5
1 2 -1

3
3 -1
3 2
3 2

9
3 4
3 -4
1 -3 -1
2 1
3 -1
3 3
1 2 0
2 1
1 -5 -1

9
3 5
1 5 3
1 -5 1
2 2
2 1
3 5
3 -2
3 2
1 -1 -5

7
3 -1
3 2
3 -3
1 2 3
2 1
3 -3
3 4

9
1 -4 -2
2 1
1 5 -5
2 1
1 1 2
1 -4 -4
1 4 5
1 -3 -1
3 2

3
1 4 -3
2 1
1 4 -2

1
3 -4

7
3 -4
3 -5
3 4
1 1 4
1 -5 1
1 0 5
3 2

8
3 5
3 -5
1 -1 -3
3 5
3 3
2 1
3 -4
3 -2

5
3 3
1 -2 0
2 1
1 -3 -5
3 0

2
3 5
1 -4 -4

9
3 2
1 1 2
2 1
1 3 -5
3 1
1 -3 1
3 -2
3 -4
1 5 2

4
1 4 1
2 1
3 0
3 0

7
3 -3
3 0
3 -1
3 1
3 3
1 4 4
1 -3 1

2
1 -5 -5
1 -2 -2

1
3 2

6
1 1 2
2 1
3 2
3 -4
1 -3 5
1 -4 0

7
3 -3
3 -1
1 -1 -3
2 1
1 2 -5
2 1
3 -5

1
3 2

8
1 -4 5
2 1
3 5
1 -3 -5
1 -5 4
1 5 -4
3 1
3 0

5
1 2 0
2 1
3 2
1 -2 -3
3 -2

1
3 4

7
3 -1
1 2 -5
2 1
1 -5 -5
3 3
3 0
1 -4 -2

8
1 2 -5
3 -2
3 -5
2 1
1 -5 -1
2 1
3 5
1 -5 -2

3
3 3
3 1
1 3 -2

2
3 -4
1 3 -3

10
1 -2 4
1 2 3
2 2
2 1
3 5
1 -5 -1
3 2
1 3 -4
1 -1 2
1 -1 1

6
1 -3 3
1 2 -2
3 -4
2 1
2 2
3 -2

1
3 -2

5
3 -4
1 1 0
2 1
3 5
1 4 -5

8
1 5 -5
3 -5
3 1
2 1
3 3
1 0 -4
2 1
1 -1 2

6
3 -5
1 3 3
1 -5 -4
1 -3 0
2 2
3 -5

4
3 1
1 0 -2
2 1
1 -3 4

10
1 -4 -1
3 4
1 5 -3
2 2
2 1
3 0
3 -2
3 4
3 -3
3 -5

3
3 2
3 4
1 -4 -3

7
3 -5
1 -5 -5
3 -2
2 1
3 2
1 -2 -1
1 -1 3

9
3 -2
3 0
1 5 -5
3 2
2 1
1 -3 -2
2 1
1 -4 -3
2 1

5
1 0 2
3 2
3 4
1 -2 1
3 -5

2
3 4
3 4

3
3 -5
3 5
1 4 -1

9
3 3
3 3
3 1
3 4
3 3
3 -5
3 -3
3 0
3 -5

9
1 -1 -1
1 -5 -4
2 2
2 1
3 -2
1 -5 -4
3 -1
2 1
3 3

2
1 -2 0
3 1

7
1 5 1
2 1
3 3
1 5 4
2 1
3 -4
3 3

6
1 -4 -2
2 1
3 4
1 -3 -4
3 1
1 -2 5

2
3 4
3 5

9
3 1
3 -2
3 2
3 4
1 -1 1
2 1
3 0
1 1 -5
1 -4 -1

10
3 4
3 -3
3 -2
3 4
1 5 0
3 -4
1 2 -3
3 -4
2 1
2 2

4
3 4
3 -4
1 -2 5
1 -4 -4

6
1 0 4
1 4 2
3 0
1 3 4
3 -4
2 1

2
1 -2 0
2 1

8
3 -4
1 1 3
2 1
1 -4 4
2 1
3 -5
3 1
3 -4

8
3 -1
1 -4 2
2 1
3 1
1 -1 2
3 -4
3 -4
2 1

6
1 3 1
1 -2 1
1 -1 5
3 -4
1 -1 -4
3 2

7
1 3 1
1 -4 2
2 1
1 -5 -3
2 2
1 -1 3
2 2

9
1 -5 -2
1 2 -5
2 1
2 2
3 -5
1 -3 1
3 1
3 5
3 4

10
1 -2 -5
2 1
3 1
1 -3 0
1 -1 -2
2 1
2 2
3 1
1 -1 -3
2 1

5
1 5 5
1 -3 2
1 -4 5
2 1
2 2

4
3 -5
3 2
3 4
3 -3

8
1 -4 0
3 -5
1 -3 1
1 -1 -3
3 0
2 1
1 -1 3
3 3

2
3 -4
1 5 1

6
3 0
3 3
3 2
1 4 -3
2 1
3 -5

7
1 -2 2
3 -5
1 -5 -2
2 1
3 4
1 -1 -5
3 -2

1
3 1

6
1 3 -1
2 1
1 5 -4
2 1
3 4
3 1

7
3 -3
3 5
1 5 2
1 -2 3
3 1
1 3 -1
3 0

7
1 2 0
1 2 -5
3 4
1 4 2
1 -5 -4
1 0 -3
2 3

5
3 2
1 -4 -1
2 1
3 -1
3 -3

4
3 -1
3 2
1 4 -3
3 -5

6
3 -4
1 2 -5
3 4
2 1
1 4 1
1 1 3`

func buildReference() (string, func(), error) {
	refSrc := os.Getenv("REFERENCE_SOURCE_PATH")
	if refSrc == "" {
		return "", nil, fmt.Errorf("REFERENCE_SOURCE_PATH not set")
	}

	content, err := os.ReadFile(refSrc)
	if err != nil {
		return "", nil, fmt.Errorf("cannot read reference source: %v", err)
	}

	tmpDir, err := os.MkdirTemp("", "678F-ref")
	if err != nil {
		return "", nil, err
	}
	cleanup := func() { os.RemoveAll(tmpDir) }

	binPath := filepath.Join(tmpDir, "ref_678F")

	if strings.Contains(string(content), "#include") {
		cppPath := filepath.Join(tmpDir, "ref.cpp")
		if err := os.WriteFile(cppPath, content, 0644); err != nil {
			cleanup()
			return "", nil, err
		}
		cmd := exec.Command("g++", "-O2", "-o", binPath, cppPath)
		if out, err := cmd.CombinedOutput(); err != nil {
			cleanup()
			return "", nil, fmt.Errorf("g++ build failed: %v\n%s", err, string(out))
		}
	} else {
		cmd := exec.Command("go", "build", "-o", binPath, refSrc)
		if out, err := cmd.CombinedOutput(); err != nil {
			cleanup()
			return "", nil, fmt.Errorf("go build failed: %v\n%s", err, string(out))
		}
	}

	return binPath, cleanup, nil
}

func runProg(exe, input string) (string, error) {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return out.String() + errBuf.String(), fmt.Errorf("%v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierF /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	scan := bufio.NewScanner(strings.NewReader(embeddedTestcases))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t := 0
	fmt.Sscan(scan.Text(), &t)

	for caseIdx := 1; caseIdx <= t; caseIdx++ {
		if !scan.Scan() {
			fmt.Printf("missing q for case %d\n", caseIdx)
			os.Exit(1)
		}
		q, _ := strconv.Atoi(scan.Text())
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", q))
		for i := 0; i < q; i++ {
			if !scan.Scan() {
				fmt.Printf("bad test file at case %d\n", caseIdx)
				os.Exit(1)
			}
			tok := scan.Text()
			tt, _ := strconv.Atoi(tok)
			sb.WriteString(tok)
			switch tt {
			case 1:
				var a, b int
				if !scan.Scan() {
					fmt.Printf("bad test file at case %d\n", caseIdx)
					os.Exit(1)
				}
				fmt.Sscan(scan.Text(), &a)
				if !scan.Scan() {
					fmt.Printf("bad test file at case %d\n", caseIdx)
					os.Exit(1)
				}
				fmt.Sscan(scan.Text(), &b)
				sb.WriteString(fmt.Sprintf(" %d %d", a, b))
			case 2:
				var idx int
				if !scan.Scan() {
					fmt.Printf("bad test file at case %d\n", caseIdx)
					os.Exit(1)
				}
				fmt.Sscan(scan.Text(), &idx)
				sb.WriteString(fmt.Sprintf(" %d", idx))
			default:
				var v int
				if !scan.Scan() {
					fmt.Printf("bad test file at case %d\n", caseIdx)
					os.Exit(1)
				}
				fmt.Sscan(scan.Text(), &v)
				sb.WriteString(fmt.Sprintf(" %d", v))
			}
			sb.WriteByte('\n')
		}
		input := sb.String()

		expect, err := runProg(refBin, input)
		if err != nil {
			fmt.Printf("case %d: reference runtime error: %v\n%s", caseIdx, err, expect)
			os.Exit(1)
		}

		got, err := runProg(bin, input)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n%s", caseIdx, err, got)
			os.Exit(1)
		}

		// Compare using fields (whitespace-insensitive)
		gotFields := strings.Fields(got)
		expFields := strings.Fields(expect)
		match := len(gotFields) == len(expFields)
		if match {
			for k := range gotFields {
				if gotFields[k] != expFields[k] {
					match = false
					break
				}
			}
		}
		if !match {
			fmt.Printf("case %d failed:\nexpected:\n%s\ngot:\n%s\n", caseIdx, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
