package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleE")
	cmd := exec.Command("go", "build", "-o", oracle, "932E.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func buildInput(fields []string) (string, error) {
	if len(fields) != 2 {
		return "", fmt.Errorf("expected 2 fields")
	}
	return fmt.Sprintf("%s %s\n", fields[0], fields[1]), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)
	const testcasesRaw = `242 5
106 7
491 3
93 2
21 7
563 5
820 1
228 9
550 6
284 3
847 2
269 4
966 1
849 5
820 5
199 3
318 5
643 6
89 10
346 7
519 4
183 4
485 5
92 9
861 5
8 5
587 5
869 9
200 7
434 10
296 7
463 3
239 5
266 1
84 1
474 5
532 9
664 8
718 6
149 4
69 7
936 4
651 8
283 3
365 7
765 10
329 9
204 6
104 1
726 4
285 10
631 4
126 6
958 3
298 8
27 1
366 2
918 5
753 6
19 6
296 6
990 3
794 7
883 10
697 2
301 10
197 8
300 3
257 7
614 3
340 10
10 6
46 8
174 6
803 6
298 10
100 8
213 7
939 4
117 1
64 1
755 3
610 3
622 1
560 8
597 4
330 1
126 9
300 7
668 4
490 4
248 8
421 8
38 4
432 8
255 7
851 4
511 4
33 1
261 5`

	scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		fields := strings.Fields(line)
		input, err := buildInput(fields)
		if err != nil {
			fmt.Fprintf(os.Stderr, "bad testcase %d: %v\n", idx, err)
			os.Exit(1)
		}
		expect, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error on test %d: %v\n", idx, err)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("test %d failed\nexpected: %s\n got: %s\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
