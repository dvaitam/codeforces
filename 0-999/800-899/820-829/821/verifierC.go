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
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	cand := os.Args[1]
	ref := "./refC.bin"
	if err := exec.Command("go", "build", "-o", ref, "821C.go").Run(); err != nil {
		fmt.Println("failed to build reference solution:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	const testcasesRaw = `100
1
add 1
remove
1
add 2
remove
2
add 3
add 2
remove
remove
5
add 1
remove
add 3
remove
add 7
remove
add 9
remove
add 9
remove
4
add 5
remove
add 1
add 6
remove
add 3
remove
remove
5
add 3
add 1
add 3
add 9
add 9
remove
remove
remove
remove
remove
6
add 9
add 8
remove
remove
add 6
remove
add 6
add 8
add 7
remove
remove
remove
6
add 12
add 9
add 5
remove
remove
remove
add 6
remove
add 8
remove
add 8
remove
3
add 5
remove
add 5
remove
add 4
remove
6
add 4
remove
add 12
remove
add 10
add 8
add 12
remove
remove
remove
add 10
remove
5
add 7
add 4
add 6
remove
remove
add 6
remove
remove
add 4
remove
6
add 2
add 11
add 10
add 2
remove
add 5
add 4
remove
remove
remove
remove
remove
1
add 2
remove
6
add 1
add 6
add 11
add 2
remove
add 12
remove
add 3
remove
remove
remove
remove
2
add 2
remove
add 1
remove
4
add 1
remove
add 4
add 1
add 2
remove
remove
remove
3
add 3
add 3
add 5
remove
remove
remove
6
add 1
remove
add 7
remove
add 12
add 4
add 11
add 2
remove
remove
remove
remove
1
add 2
remove
2
add 4
add 3
remove
remove
2
add 3
add 4
remove
remove
6
add 1
remove
add 3
remove
add 5
add 3
add 8
remove
remove
remove
add 1
remove
2
add 2
remove
add 1
remove
3
add 1
remove
add 5
remove
add 5
remove
6
add 6
add 7
add 1
add 7
add 2
remove
add 2
remove
remove
remove
remove
remove
1
add 1
remove
2
add 2
add 1
remove
remove
5
add 8
add 9
remove
add 4
remove
add 9
add 10
remove
remove
remove
1
add 2
remove
5
add 10
add 2
remove
add 1
remove
remove
add 10
add 6
remove
remove
3
add 1
remove
add 4
add 3
remove
remove
2
add 1
remove
add 1
remove
4
add 8
add 2
add 1
add 2
remove
remove
remove
remove
2
add 4
add 3
remove
remove
4
add 8
add 6
add 2
add 2
remove
remove
remove
remove
1
add 2
remove
6
add 7
remove
add 12
add 10
remove
remove
add 12
remove
add 5
add 8
remove
remove
2
add 3
add 4
remove
remove
6
add 12
remove
add 8
remove
add 5
add 3
add 5
remove
remove
remove
add 2
remove
5
add 10
add 6
add 9
add 7
remove
remove
remove
remove
add 1
remove
2
add 3
add 3
remove
remove
4
add 3
remove
add 2
add 7
add 4
remove
remove
remove
6
add 9
add 3
add 12
add 6
remove
remove
add 8
remove
remove
add 12
remove
remove
2
add 4
remove
add 4
remove
3
add 4
add 6
remove
remove
add 4
remove
3
add 5
add 6
remove
add 6
remove
remove
2
add 2
add 4
remove
remove
6
add 1
add 9
remove
remove
add 11
add 2
add 4
remove
remove
remove
add 2
remove
4
add 4
add 4
add 3
add 7
remove
remove
remove
remove
4
add 5
remove
add 3
add 4
remove
remove
add 6
remove
1
add 2
remove
5
add 2
remove
add 8
remove
add 6
add 9
add 2
remove
remove
remove
5
add 6
add 7
add 3
add 8
add 8
remove
remove
remove
remove
remove
6
add 5
add 5
remove
remove
add 6
add 7
remove
add 11
add 11
remove
remove
remove
2
add 4
remove
add 3
remove
1
add 2
remove
3
add 3
add 3
remove
add 1
remove
remove
5
add 4
add 9
remove
add 9
remove
add 4
remove
remove
add 9
remove
3
add 1
remove
add 1
remove
add 4
remove
3
add 2
add 1
add 3
remove
remove
remove
2
add 4
add 3
remove
remove
4
add 3
remove
add 8
remove
add 4
remove
add 6
remove
2
add 1
remove
add 4
remove
2
add 3
add 2
remove
remove
2
add 2
add 4
remove
remove
4
add 6
add 6
remove
remove
add 5
remove
add 2
remove
3
add 3
add 2
add 3
remove
remove
remove
4
add 8
add 3
add 7
add 1
remove
remove
remove
remove
1
add 2
remove
2
add 3
add 4
remove
remove
1
add 2
remove
2
add 4
remove
add 4
remove
6
add 3
add 4
remove
add 3
remove
add 3
add 5
add 4
remove
remove
remove
remove
2
add 2
remove
add 1
remove
3
add 3
remove
add 5
add 2
remove
remove
2
add 1
add 3
remove
remove
2
add 3
remove
add 1
remove
3
add 6
remove
add 5
remove
add 2
remove
1
add 2
remove
4
add 8
remove
add 6
remove
add 2
remove
add 8
remove
4
add 8
add 7
remove
add 2
add 3
remove
remove
remove
1
add 1
remove
1
add 2
remove
3
add 4
remove
add 1
add 6
remove
remove
3
add 2
add 2
add 1
remove
remove
remove
3
add 6
add 4
add 4
remove
remove
remove
2
add 1
add 1
remove
remove
1
add 1
remove
5
add 1
remove
add 9
add 9
remove
add 8
add 9
remove
remove
remove
1
add 1
remove
1
add 2
remove
2
add 1
remove
add 4
remove
3
add 5
remove
add 4
remove
add 6
remove
5
add 9
remove
add 3
remove
add 3
remove
add 4
add 4
remove
remove
3
add 3
add 4
add 3
remove
remove
remove
3
add 1
remove
add 4
add 5
remove
remove
4
add 5
add 3
remove
remove
add 1
remove
add 4
remove
2
add 2
add 1
remove
remove
6
add 3
remove
add 12
add 12
add 1
remove
remove
add 3
remove
add 4
remove
remove
1
add 1
remove`

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
		n, _ := strconv.Atoi(scan.Text())
		var input bytes.Buffer
		input.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < 2*n; i++ {
			if !scan.Scan() {
				fmt.Printf("bad test %d\n", tc+1)
				os.Exit(1)
			}
			cmdWord := scan.Text()
			input.WriteString(cmdWord)
			if cmdWord == "add" {
				if !scan.Scan() {
					fmt.Printf("bad test %d\n", tc+1)
					os.Exit(1)
				}
				val := scan.Text()
				input.WriteByte(' ')
				input.WriteString(val)
			}
			input.WriteByte('\n')
		}
		want, err := run(ref, input.Bytes())
		if err != nil {
			fmt.Println("reference runtime error:", err)
			os.Exit(1)
		}
		got, err := run(cand, input.Bytes())
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", tc+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(want) != strings.TrimSpace(got) {
			fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", tc+1, input.String(), want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", T)
}
