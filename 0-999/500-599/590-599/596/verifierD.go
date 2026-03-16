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
	oracle := filepath.Join(dir, "oracleD")
	cmd := exec.Command("go", "build", "-o", oracle, "596D.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstderr: %s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

const testcasesDRaw = `2 5 0.544
0 4

4 5 0.066
-5 -1 2 3

2 2 0.996
2 3

5 4 0.397
-3 -2 1 3 5

1 1 0.159
4

1 3 0.78
-1

4 5 0.719
-3 1 2 4

3 1 0.036
-2 -1 2

4 3 0.421
0 1 3 4

5 4 0.584
-5 -1 0 4 5

2 3 0.965
-4 4

2 5 0.267
-4 2

4 1 0.344
-5 -4 -3 1

3 4 0.769
-5 -4 4

5 1 0.378
-2 -1 0 3 4

1 3 0.007
-4

5 5 0.031
-3 -2 -1 1 4

1 3 0.314
-3

4 4 0.46
1 3 4 5

5 1 0.62
-2 -1 1 3 5

3 4 0.978
-1 0 3

1 4 0.987
0

1 4 0.616
5

2 1 0.634
0 2

3 3 0.609
-5 -1 2

5 1 0.955
-5 -1 0 2 5

3 5 0.601
-3 0 4

3 3 0.787
-5 -4 4

2 3 0.5
-1 5

2 3 0.187
1 5

1 1 0.601
0

2 4 0.81
-4 -3

3 2 0.885
-2 -1 2

1 1 0.53
-2

3 5 0.184
-1 0 5

1 5 0.345
-3

4 3 0.518
-1 0 2 5

4 3 0.42
-5 -3 -2 1

1 4 0.943
4

5 4 0.559
-5 -2 2 3 5

3 5 0.341
-4 -2 4

3 1 0.81
-5 -2 3

4 5 0.049
-4 -3 2 3

3 2 0.663
-5 1 3

5 1 0.341
-5 -1 0 2 3

2 2 0.122
-4 -3

2 3 0.922
-5 -3

4 5 0.863
-5 -2 -1 4

5 5 0.423
-5 -4 -3 0 2

1 1 0.483
-5

1 5 0.502
0

2 3 0.072
1 5

4 5 0.304
-2 -1 0 1

1 2 0.556
1

1 5 0.179
0

4 5 0.65
-5 1 3 5

5 4 0.053
-5 0 1 2 5

2 2 0.536
-4 4

4 2 0.426
-5 -1 0 3

1 4 0.691
5

5 4 0.668
-5 -4 0 3 4

4 2 0.236
-5 -4 1 3

5 1 0.659
-5 -4 -3 0 1

1 1 0.674
-1

5 3 0.799
-5 -4 -2 3 4

5 1 0.935
-5 -3 0 3 4

1 2 0.998
5

2 4 0.616
-1 1

3 5 0.397
0 1 3

1 4 0.5
1

2 4 0.691
4 5

5 4 0.156
-4 -3 1 2 3

4 5 0.719
-3 -2 -1 4

5 3 0.935
-1 1 3 4 5

5 5 0.984
-5 -2 -1 1 2

2 2 0.57
-2 0

4 2 0.418
-2 2 4 5

5 1 0.481
-5 -4 -2 1 2

2 1 0.959
-2 -1

2 3 0.137
4 5

1 3 0.17
-5

3 2 0.423
-5 -4 -1

3 4 0.58
-5 0 5

1 3 0.331
1

4 1 0.21
-3 1 2 4

5 3 0.119
-4 -1 1 2 5

5 3 0.097
-4 -1 0 2 5

3 5 0.536
-4 2 5

5 3 0.06
-3 -1 0 4 5

4 1 0.108
-3 0 3 5

5 4 0.555
-4 -3 -1 2 5

1 2 0.754
3

5 4 0.359
-5 -3 -1 1 2

5 3 0.247
-4 0 1 2 3

3 4 0.847
-4 -3 -1

5 1 0.681
-4 -3 -2 1 4

4 2 0.592
-3 -2 1 3

5 2 0.568
-5 -2 -1 0 2

4 4 0.316
-1 2 4 5

5 3 0.896
-5 -2 2 4 5

`

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	scanner := bufio.NewScanner(strings.NewReader(testcasesDRaw))
	idx := 0
	var lines []string
	process := func(ls []string) {
		if len(ls) == 0 {
			return
		}
		idx++
		input := strings.Join(ls, "\n") + "\n"
		exp, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error on case %d: %v\n", idx, err)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("case %d failed:\nexpected: %s\n got: %s\n", idx, exp, got)
			os.Exit(1)
		}
	}

	for scanner.Scan() {
		line := strings.TrimRight(scanner.Text(), "\r")
		if strings.TrimSpace(line) == "" {
			process(lines)
			lines = nil
			continue
		}
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	process(lines)
	fmt.Printf("All %d tests passed\n", idx)
}
