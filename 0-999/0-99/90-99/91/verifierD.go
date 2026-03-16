package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const testcasesDRaw = `5
4 1 3 2 5

7
2 6 1 3 4 7 5

4
2 4 3 1

6
1 4 5 3 6 2

2
1 2

6
6 2 4 1 5 3

1
1

1
1

2
1 2

4
1 3 4 2

5
4 2 5 3 1

5
5 2 1 3 4

7
2 4 7 5 3 1 6

5
2 5 1 4 3

7
1 7 4 3 5 6 2

7
5 3 2 4 1 6 7

3
3 2 1

7
5 3 6 2 7 4 1

3
2 1 3

4
3 2 1 4

2
1 2

7
6 4 7 5 2 3 1

7
1 3 6 2 4 5 7

2
1 2

5
5 3 4 1 2

1
1

2
1 2

6
4 3 5 1 2 6

3
1 3 2

4
4 3 2 1

4
1 4 3 2

1
1

5
1 2 3 4 5

2
1 2

3
2 3 1

6
5 4 1 2 3 6

6
6 1 5 3 2 4

7
5 3 7 2 6 1 4

1
1

1
1

6
3 6 4 1 2 5

6
2 3 6 1 5 4

7
4 3 2 1 6 7 5

2
2 1

4
2 3 4 1

5
3 1 4 5 2

5
5 2 4 1 3

7
6 2 4 3 7 1 5

5
4 5 1 3 2

3
3 2 1

3
3 1 2

3
2 1 3

3
2 3 1

7
3 7 2 1 4 5 6

7
2 6 1 3 7 5 4

6
2 6 3 1 4 5

6
5 3 4 2 6 1

5
4 1 3 5 2

1
1

7
5 2 7 3 6 4 1

1
1

6
1 5 6 2 4 3

5
4 3 5 1 2

5
5 3 4 2 1

6
4 3 1 6 5 2

2
2 1

3
1 3 2

6
1 4 6 5 2 3

3
3 2 1

7
3 2 6 4 5 1 7

7
1 5 3 7 6 4 2

5
2 4 5 1 3

2
1 2

4
4 1 3 2

5
3 4 5 2 1

1
1

3
2 3 1

5
1 4 3 2 5

3
3 2 1

2
2 1

5
1 5 3 4 2

7
7 5 6 3 2 1 4

6
6 5 2 4 3 1

1
1

3
2 1 3

3
3 1 2

3
2 3 1

7
5 6 3 1 4 7 2

6
6 5 1 4 3 2

5
3 5 2 1 4

5
5 2 4 1 3

6
4 6 1 2 3 5

3
1 2 3

2
1 2

3
2 3 1

1
1

4
1 3 4 2

3
3 1 2

7
5 6 1 7 3 4 2

4
1 3 2 4

`

func buildOracle() (string, error) {
	exe := "oracleD"
	cmd := exec.Command("go", "build", "-o", exe, "91D.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle: %v\n%s", err, out)
	}
	return exe, nil
}

func runProg(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	scanner := bufio.NewScanner(strings.NewReader(testcasesDRaw))
	idx := 0
	var lines []string
	process := func() {
		if len(lines) == 0 {
			return
		}
		idx++
		input := strings.Join(lines, "\n") + "\n"
		exp, err := runProg("./"+oracle, input)
		if err != nil {
			fmt.Printf("oracle error on case %d: %v\n", idx, err)
			os.Exit(1)
		}
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("case %d mismatch\nexpected: %s\n got: %s\n", idx, exp, got)
			os.Exit(1)
		}
	}
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			process()
			lines = nil
			continue
		}
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	process()
	fmt.Printf("All %d tests passed\n", idx)
}
