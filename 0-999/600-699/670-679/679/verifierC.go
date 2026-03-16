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

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleC")
	cmd := exec.Command("go", "build", "-o", oracle, "679C.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func runBinary(bin, input string) (string, string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	err := cmd.Run()
	return strings.TrimSpace(out.String()), errb.String(), err
}

func atoi(s string) int {
	v, _ := strconv.Atoi(s)
	return v
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	const testcasesCRaw = `5 5 ..... X.X.. ....X ...X. X...X
3 2 ..X ..X .XX
4 2 .... XX.X ..XX XXXX
2 1 X. X.
5 3 ..XXX X.XXX ...X. ..XXX .XX..
4 1 X.XX ..XX XX.X ....
4 4 X... .... .XX. .XX.
1 1 X
1 1 .
5 3 ..X.. ..XXX .XXXX ..XX. ..XXX
3 3 .X. .X. .XX
4 2 X.X. .X.. XXXX X..X
4 3 X.X. XX.X X..X ....
2 1 X. .X
5 5 XXX.X .X.X. X.X.. .X..X XX...
5 2 ..XXX ..X.X XXXXX ...XX X..X.
4 3 X... .... X.XX ..XX
5 4 ...X. XXXX. ...XX X.XXX XXX..
2 1 XX ..
4 4 .X.. XX.X XX.. X...
5 5 .X.X. XX..X X.X.X X...X XX.XX
1 1 .
4 3 .... XX.. XX.. X..X
2 1 XX X.
4 3 X.X. .X.X X... X...
1 1 .
1 1 X
1 1 .
2 2 .X XX
2 1 .X XX
3 2 X.. .X. X..
2 2 .. .X
3 2 ... X.. .XX
5 4 .XXX. .XX.X X.XXX .XXX. ..X..
4 4 .XX. .XX. XX.X .X.X
5 4 ..X.X .XX.X .XX.. X..XX XX.XX
3 3 ... .XX X..
2 1 XX X.
5 5 ..XXX XXX.. ..XXX .XXXX .XXXX
1 1 X
2 2 .. .X
4 4 ...X .X.. .XX. .X..
4 4 XX.X XX.X .X.. .XXX
4 4 X..X ..X. ...X ...X
5 1 XXXX. XXXXX .XX.. ..X.X ...X.
5 4 .X.X. .XXXX X...X ..X.. XXXXX
4 3 ..X. .XXX .XXX X...
2 1 .X .X
5 1 XX... X..XX XX.X. XXXXX .....
4 2 .X.X ..XX .X.X ...X
3 1 XXX ... XX.
2 2 .. ..
2 2 .. ..
4 3 .X.. XXXX ...X XX.X
2 2 .X XX
1 1 .
3 3 .XX .XX .X.
3 3 ..X .X. ...
1 1 X
1 1 X
4 2 .XX. XX.X XX.X XX..
4 4 .... .X.. .... .X.X
2 1 XX X.
1 1 .
3 2 .X. X.X XXX
4 1 XXX. X.XX .XXX ..X.
3 3 X.X XX. .X.
5 5 X..X. ..XXX XXXXX ....X .X.XX
5 2 .XX.. ...XX XXX.X XXXX. X...X
3 2 .X. .X. .XX
1 1 .
3 1 X.X .X. ..X
1 1 X
3 2 XXX .X. X.X
2 2 XX X.
5 4 X.XX. XX... XX... .X.X. XX.XX
2 1 XX ..
2 1 X. .X
2 2 .. .X
2 2 XX XX
2 2 .. ..
3 1 .X. X.X .X.
5 1 X...X XXXXX ..X.. ..X.. X.XX.
5 2 X.X.X XXXXX ..XX. X.X.. .X..X
5 1 XXXX. .XX.. .XX.X ..X.. XX...
2 1 .X .X
4 1 ..XX ..X. .XXX .X..
5 1 X.... X.... ..X.X XX..X .X..X
5 3 .XX.X X.X.X X.XXX .X.X. XX...
4 1 XXX. XX.. XXX. X.XX
2 2 X. .X
2 2 XX .X
5 4 .XXX. X...X ..XXX ..X.. .....
5 2 XXXX. XX.XX .XX.X XXXXX X....
5 5 ...X. X.... XXX.X XXX.. ....X
5 2 XXX.. X.X.X .XX.X .X.XX XXX..
5 3 .X.X. .XXX. XX..X .X... ....X
5 5 X.XXX ....X XX..X XX..X X.XX.
2 2 .X XX
4 2 .XXX .XXX XXXX .XX.`

	scanner := bufio.NewScanner(strings.NewReader(testcasesCRaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		fields := strings.Fields(line)
		if len(fields) < 2 {
			fmt.Printf("bad case %d\n", idx)
			os.Exit(1)
		}
		n := atoi(fields[0])
		k := atoi(fields[1])
		if len(fields) != 2+n {
			fmt.Printf("bad case %d\n", idx)
			os.Exit(1)
		}
		rows := fields[2:]
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d %d\n", n, k))
		for _, r := range rows {
			input.WriteString(r)
			input.WriteByte('\n')
		}
		inputStr := input.String()

		exp, errStr, err := runBinary(oracle, inputStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error on case %d: %v\n%s", idx, err, errStr)
			os.Exit(1)
		}
		got, errStr2, err := runBinary(bin, inputStr)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, errStr2)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("test %d failed\nexpected:\n%s\n\ngot:\n%s\n", idx, exp, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
