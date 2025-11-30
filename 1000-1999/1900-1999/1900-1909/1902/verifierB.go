package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `19 80 1 15
30 39 1 7
8 40 10 18
28 348 16 3
32 391 8 18
26 31 9 1
18 17 2 1
26 406 17 19
29 91 4 9
19 22 7 20
3 10 3 9
35 272 11 4
16 54 6 3
19 69 10 17
37 545 17 7
7 104 14 18
48 453 9 10
24 81 19 5
8 52 4 13
38 574 15 5
43 485 10 12
48 489 14 7
32 505 17 11
42 154 2 15
10 160 16 2
14 31 1 12
26 18 1 17
44 4 3 13
24 80 2 4
1 24 9 10
15 74 5 19
13 59 4 14
46 173 11 13
22 76 14 14
29 162 5 17
9 57 7 6
23 252 13 14
25 225 8 7
14 199 19 2
3 28 8 3
2 26 3 1
25 26 5 2
45 478 7 15
36 150 5 16
37 130 17 5
34 445 16 4
43 47 16 19
19 357 19 8
37 460 10 9
6 124 6 5
40 209 3 8
23 147 8 4
2 6 3 4
23 45 1 2
8 46 5 20
14 67 11 19
36 455 17 12
35 389 11 8
40 445 14 16
25 312 4 18
8 150 16 17
4 53 12 6
8 61 17 7
39 549 1 7
38 283 10 6
22 198 19 14
22 208 18 11
25 340 11 4
18 292 8 6
13 108 9 10
25 443 8 18
40 550 17 12
43 522 13 7
23 122 19 19
47 506 14 10
35 370 11 7
36 105 2 6
6 53 18 9
7 64 18 8
28 83 15 16
30 97 3 4
12 122 7 5
37 18 1 8
39 439 3 14
35 232 1 2
32 172 19 19
33 258 9 11
7 84 16 18
6 55 2 9
13 67 10 17
23 150 18 12
19 280 17 16
39 209 11 7
19 115 10 12
27 441 12 20
44 26 12 2
38 513 5 9
21 209 16 7
21 290 3 18
5 7 5 5
32 340 1 18
4 60 15 3
1 2 2 2`

func expected(n, P, l, t int64) string {
	tasks := (n + 6) / 7
	low, high := int64(0), n
	for low < high {
		mid := (low + high) / 2
		taskDone := 2 * mid
		if taskDone > tasks {
			taskDone = tasks
		}
		points := mid*l + taskDone*t
		if points >= P {
			high = mid
		} else {
			low = mid + 1
		}
	}
	return fmt.Sprintf("%d", n-low)
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func loadCases() ([]string, []string) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	var inputs []string
	var expects []string
	for idx, line := range lines {
		fields := strings.Fields(strings.TrimSpace(line))
		if len(fields) != 4 {
			fmt.Printf("invalid testcase line %d\n", idx+1)
			os.Exit(1)
		}
		vals := make([]int64, 4)
		for i := 0; i < 4; i++ {
			v, err := strconv.ParseInt(fields[i], 10, 64)
			if err != nil {
				fmt.Printf("invalid number on line %d\n", idx+1)
				os.Exit(1)
			}
			vals[i] = v
		}
		inputs = append(inputs, fmt.Sprintf("1\n%d %d %d %d\n", vals[0], vals[1], vals[2], vals[3]))
		expects = append(expects, expected(vals[0], vals[1], vals[2], vals[3]))
	}
	return inputs, expects
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	inputs, expects := loadCases()
	for idx, input := range inputs {
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expects[idx] {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, expects[idx], got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(inputs))
}
