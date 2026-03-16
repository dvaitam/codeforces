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

func runRef(input string) (string, error) {
	cmd := exec.Command("go", "run", "919C.go")
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func runBin(bin, input string) (string, error) {
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
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 || len(os.Args) > 3 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	if bin == "--" && len(os.Args) == 3 {
		bin = os.Args[2]
	}
	const testcasesRaw = `100
5 2 4
*.
**
..
**
.*
4 1 4
.
*
.
*
5 4 1
.*.*
.**.
*.**
..*.
.***
4 1 3
*
*
.
.
1 1 1
*
5 3 1
...
***
...
***
***
5 2 2
..
**
*.
**
**
2 4 4
**.*
*.*.
5 2 1
**
*.
..
*.
*.
1 5 1
**..*
2 1 2
*
*
4 2 3
**
.*
.*
..
2 4 3
****
..*.
5 1 4
.
.
.
.
*
1 2 1
*.
5 3 3
**.
***
**.
**.
...
3 1 2
.
.
.
4 2 4
.*
.*
..
**
5 4 3
*...
*..*
..**
****
*..*
5 1 1
*
*
*
.
*
4 4 4
....
..*.
*.**
...*
2 2 2
**
*.
5 3 5
.**
***
...
..*
..*
1 5 5
*..**
1 5 3
.*...
5 3 3
.**
*..
.**
..*
...
2 2 2
..
..
4 3 2
*..
***
*..
.**
5 3 1
*.*
.**
*..
.*.
**.
3 3 3
.*.
*..
*.*
1 1 1
.
1 2 2
..
3 4 2
.**.
**.*
**.*
3 5 5
*..**
.....
*....
2 5 2
.*.*.
.***.
1 4 1
**.*
2 5 3
.****
*.***
1 5 5
*.**.
4 3 4
..*
.**
.**
*..
5 3 2
*..
*..
.**
***
***
2 1 1
.
*
2 3 1
**.
.**
2 5 2
...**
***.*
3 3 3
**.
*..
.**
4 1 4
.
.
*
.
2 3 3
*..
.*.
4 1 4
.
*
.
.
2 4 1
****
***.
3 1 2
.
*
.
4 4 4
*.**
.**.
**..
.**.
1 2 1
*.
3 1 2
*
.
*
4 2 2
**
..
..
*.
5 1 5
*
.
*
.
.
5 2 3
.*
**
**
.*
..
5 2 5
.*
..
*.
*.
*.
4 5 1
.*...
*****
*..*.
...*.
1 4 1
**..
3 1 2
.
*
*
5 5 3
***..
**.*.
*...*
..*.*
***..
5 4 3
...*
*.*.
.*..
**..
....
5 3 1
**.
..*
*..
*..
***
5 2 4
..
.*
..
..
*.
1 1 1
.
2 4 1
***.
.*.*
1 2 2
*.
4 4 2
**.*
.**.
***.
*.*.
5 3 4
...
*.*
**.
**.
.**
4 2 3
.*
*.
**
..
3 2 3
**
*.
**
1 4 3
*.*.
2 1 2
.
.
4 5 4
*..*.
.....
..***
*.**.
3 5 3
.**.*
*****
*....
5 5 5
...*.
*....
***.*
***..
....*
5 2 4
**
..
*.
*.
*.
5 4 3
.*.*
.***
**..
*.*.
*..*
4 5 4
.**..
*.*..
.....
**.**
3 1 3
.
.
.
4 3 3
..*
**.
.**
.**
2 2 2
.*
**
4 5 2
.***.
*****
**.**
.*.*.
2 3 2
..*
***
5 5 2
***.*
..***
.*..*
.*.*.
..*..
1 1 1
.
2 1 1
.
*
2 1 2
.
.
2 2 2
**
.*
2 4 1
*.**
**..
3 1 3
.
*
.
2 5 2
*.*.*
..*.*
4 2 4
*.
**
*.
*.
4 2 3
**
.*
.*
*.
4 2 3
**
..
**
**
5 2 5
..
*.
..
*.
*.
4 4 2
**..
*.*.
*.*.
*.*.
2 4 3
....
..*.
1 5 4
***..
3 5 3
*****
.....
**.*.
1 4 2
.**.`

	scan := bufio.NewScanner(strings.NewReader(testcasesRaw))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Fprintln(os.Stderr, "invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		if !scan.Scan() {
			fmt.Fprintln(os.Stderr, "bad file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		if !scan.Scan() {
			fmt.Fprintln(os.Stderr, "bad file")
			os.Exit(1)
		}
		m, _ := strconv.Atoi(scan.Text())
		if !scan.Scan() {
			fmt.Fprintln(os.Stderr, "bad file")
			os.Exit(1)
		}
		k, _ := strconv.Atoi(scan.Text())
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, k))
		for i := 0; i < n; i++ {
			if !scan.Scan() {
				fmt.Fprintln(os.Stderr, "bad file")
				os.Exit(1)
			}
			row := scan.Text()
			sb.WriteString(row + "\n")
		}
		input := sb.String()
		exp, err := runRef(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on case %d: %v\n", caseIdx+1, err)
			os.Exit(1)
		}
		got, err := runBin(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\n", caseIdx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected:%s\ngot:%s\n", caseIdx+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
