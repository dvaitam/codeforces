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

const testcasesARaw = `100
27 24
24 26
1 16
30 32
15 25
29 19
30 30
11 37
28 13
16 8
9 8
24 6
19 16
29 34
22 38
28 9
9 6
23 4
28 21
15 35
3 22
13 20
19 40
29 13
30 35
15 28
27 33
8 3
25 35
29 0
2 25
22 40
0 39
15 21
7 20
22 4
6 36
7 15
25 9
25 34
14 5
2 20
28 32
29 31
3 19
17 18
22 7
17 21
26 34
6 38
17 37
9 28
2 38
25 24
10 36
7 18
5 12
26 11
1 39
21 16
15 4
2 8
28 9
29 2
26 5
28 34
21 25
26 33
8 33
25 15
27 13
28 37
26 26
18 17
14 31
21 22
2 20
19 7
15 37
20 21
27 12
7 1
23 17
3 14
11 10
10 27
26 3
3 9
27 14
1 36
20 34
19 4
0 7
20 12
19 36
3 25
2 23
26 7
1 38
0 12
`

func solveA(x, y int) string {
	turn := true
	for {
		if turn {
			if x >= 2 && y >= 2 {
				x -= 2
				y -= 2
			} else if x >= 1 && y >= 12 {
				x -= 1
				y -= 12
			} else if y >= 22 {
				y -= 22
			} else {
				return "Hanako"
			}
		} else {
			if y >= 22 {
				y -= 22
			} else if x >= 1 && y >= 12 {
				x -= 1
				y -= 12
			} else if x >= 2 && y >= 2 {
				x -= 2
				y -= 2
			} else {
				return "Ciel"
			}
		}
		turn = !turn
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	data := []byte(testcasesARaw)
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	cases := make([][2]int, t)
	expected := make([]string, t)
	for i := 0; i < t; i++ {
		scan.Scan()
		x, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		y, _ := strconv.Atoi(scan.Text())
		cases[i] = [2]int{x, y}
		expected[i] = solveA(x, y)
	}
	for i, c := range cases {
		in := fmt.Sprintf("%d %d\n", c[0], c[1])
		cmd := exec.Command(os.Args[1])
		cmd.Stdin = strings.NewReader(in)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("execution failed on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got := strings.TrimSpace(string(out))
		if got != expected[i] {
			fmt.Printf("test %d failed: expected %s got %s\n", i+1, expected[i], got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
