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

func runExe(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("%v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func buildRef() (string, error) {
	ref := "./refF.bin"
	cmd := exec.Command("go", "build", "-o", ref, "946F.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
	}
	return ref, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	const testcasesRaw = `100
5 4
10100
2 1
11
2 6
00
1 3
1
3 2
100
2 9
10
2 0
00
2 2
01
3 3
000
4 4
0110
2 4
01
3 9
010
3 5
111
2 7
10
1 4
0
6 5
101110
4 0
0000
4 5
1110
5 5
10100
3 8
101
3 8
011
3 2
001
4 2
0010
2 9
11
4 10
1000
4 5
0001
1 2
1
3 2
011
2 7
01
4 5
1111
4 2
0101
3 2
011
5 8
10101
3 8
111
6 4
100111
6 4
111111
6 6
101110
5 2
01110
6 10
101110
2 2
10
2 10
00
4 3
0000
3 2
100
4 7
1100
2 5
01
4 4
1010
5 9
10111
3 2
111
4 2
0110
2 9
01
6 0
101011
4 6
0111
6 6
101111
4 0
1111
2 7
00
6 9
101001
1 7
1
2 7
00
4 4
0111
5 7
00011
5 5
10111
6 9
000110
6 8
001111
4 6
0110
3 0
111
3 2
011
3 8
010
1 3
1
1 1
0
3 6
011
1 1
0
2 9
00
3 7
000
4 2
0000
6 6
010010
1 6
0
3 10
111
6 5
001010
3 1
011
5 2
01110
2 0
10
5 10
00010
4 7
1100
5 6
00100
1 10
1
2 2
11
2 4
11
4 6
0111
3 8
000
4 1
0000
2 3
01
4 10
1011
3 9
111
4 1
0000
3 6
101
1 4
0
3 10
001
1 5
1
3 0
011
2 1
11`

	scan := bufio.NewScanner(strings.NewReader(testcasesRaw))
	scan.Split(bufio.ScanLines)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(strings.TrimSpace(scan.Text()))
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		parts := strings.Fields(scan.Text())
		if len(parts) != 2 {
			fmt.Println("bad line")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(parts[0])
		x, _ := strconv.Atoi(parts[1])
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		s := strings.TrimSpace(scan.Text())
		input := fmt.Sprintf("%d %d\n%s\n", n, x, s)
		want, err := runExe(ref, input)
		if err != nil {
			fmt.Printf("reference runtime error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runExe(bin, input)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Printf("case %d failed\ninput:\n%sexpected: %s\ngot: %s\n", i+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
