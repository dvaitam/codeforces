package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func run(bin string, input []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	cand := os.Args[1]
	ref := "./refB.bin"
	if err := exec.Command("go", "build", "-o", ref, "821B.go").Run(); err != nil {
		fmt.Println("failed to build reference solution:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	const testcasesRaw = `100
5 19
3 9
4 16
15 16
13 7
4 16
1 13
14 20
1 15
9 8
19 4
11 1
1 1
18 1
13 7
14 1
17 8
15 16
18 8
12 8
8 15
10 1
14 18
4 6
10 4
11 17
14 17
7 10
10 19
16 17
13 19
2 16
8 13
14 6
12 18
12 3
15 17
4 6
17 13
12 16
1 16
2 10
20 19
19 13
6 6
17 8
1 7
18 18
8 13
17 12
19 12
15 9
18 20
1 13
17 5
17 18
7 14
2 16
12 19
18 7
17 14
16 12
14 12
1 18
18 20
20 11
15 20
1 8
6 18
19 6
3 18
9 2
3 3
1 15
1 9
8 9
4 20
6 12
10 3
6 6
9 17
6 9
10 15
11 16
16 4
1 10
13 11
14 7
9 4
9 17
7 20
14 1
8 1
13 5
2 6
15 17
14 18
8 17
15 8
17 1`

	scan := bufio.NewScanner(strings.NewReader(testcasesRaw))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("empty test file")
		os.Exit(1)
	}
	T, _ := strconv.Atoi(scan.Text())
	for tc := 0; tc < T; tc++ {
		if !scan.Scan() {
			fmt.Printf("bad test %d\n", tc+1)
			os.Exit(1)
		}
		m, _ := strconv.Atoi(scan.Text())
		if !scan.Scan() {
			fmt.Printf("bad test %d\n", tc+1)
			os.Exit(1)
		}
		b, _ := strconv.Atoi(scan.Text())
		input := fmt.Sprintf("%d %d\n", m, b)
		want, err := run(ref, []byte(input))
		if err != nil {
			fmt.Println("reference runtime error:", err)
			os.Exit(1)
		}
		got, err := run(cand, []byte(input))
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", tc+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(want) != strings.TrimSpace(got) {
			fmt.Printf("test %d failed\ninput: %sexpected: %s\ngot: %s\n", tc+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", T)
}
