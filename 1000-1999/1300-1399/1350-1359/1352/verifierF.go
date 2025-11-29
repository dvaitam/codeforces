package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded testcases from testcasesF.txt.
const testcasesRaw = `100
1 2 0
5 3 3
1 0 0
0 3 4
2 0 1
4 4 2
2 1 0
2 1 0
5 2 2
1 1 2
2 5 5
2 0 4
2 5 3
4 1 1
1 3 2
0 4 2
0 2 4
5 2 4
1 3 3
4 2 3
3 1 1
2 2 0
0 0 3
5 2 4
4 5 3
5 2 1
5 1 0
3 1 5
5 3 2
1 2 3
5 4 2
5 4 1
2 0 0
5 1 2
4 4 1
0 2 1
2 3 0
0 2 5
0 2 5
5 2 0
2 2 2
1 5 3
4 5 0
2 4 1
3 2 1
2 3 4
1 2 4
0 2 0
3 1 2
2 2 4
0 3 1
3 1 0
0 0 0
5 1 4
5 1 4
0 4 3
4 1 2
0 0 4
2 3 5
1 3 1
1 3 3
3 0 1
3 3 1
5 3 1
3 1 0
0 2 2
1 4 1
1 3 2
1 2 0
2 4 0
4 3 5
5 5 5
0 3 3
0 3 1
4 1 2
2 5 3
5 2 3
4 1 5
5 2 2
3 3 0
2 5 5
1 0 3
4 1 2
5 0 1
5 5 3
4 3 5
3 3 1
0 1 1
0 4 2
0 3 3
1 4 0
1 1 5
4 2 4
3 4 3
0 0 0
5 4 0
3 4 2
4 1 0
2 0 4
0 2 2`

type testCase struct {
	n0 int
	n1 int
	n2 int
}

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

// referenceSolution embeds the construction logic from 1352F.go to produce one valid string.
func referenceSolution(n0, n1p, n2p int) string {
	// map variables as in original source
	n1 := n2p
	n2 := n1p
	n3 := n0

	var sBuilder []byte
	var tBuilder []byte
	tBuilder = append(tBuilder, '1')
	cFlag := false

	if n2 == 0 {
		if n1 > 0 {
			sBuilder = append(sBuilder, '1')
			for i := 0; i < n1; i++ {
				sBuilder = append(sBuilder, '1')
			}
		} else {
			sBuilder = append(sBuilder, '0')
			for i := 0; i < n3; i++ {
				sBuilder = append(sBuilder, '0')
			}
		}
	} else if n2%2 == 1 {
		cnt := n2
		for cnt > 0 {
			if cFlag {
				tBuilder = append(tBuilder, '1')
				cFlag = false
			} else {
				tBuilder = append(tBuilder, '0')
				cFlag = true
			}
			cnt--
		}
		for i := 0; i < n1; i++ {
			sBuilder = append(sBuilder, '1')
		}
		sBuilder = append(sBuilder, tBuilder...)
		for i := 0; i < n3; i++ {
			sBuilder = append(sBuilder, '0')
		}
	} else {
		cnt := n2 - 1
		for cnt > 0 {
			if cFlag {
				tBuilder = append(tBuilder, '1')
				cFlag = false
			} else {
				tBuilder = append(tBuilder, '0')
				cFlag = true
			}
			cnt--
		}
		for i := 0; i < n3; i++ {
			tBuilder = append(tBuilder, '0')
		}
		tBuilder = append(tBuilder, '1')
		for i := 0; i < n1; i++ {
			sBuilder = append(sBuilder, '1')
		}
		sBuilder = append(sBuilder, tBuilder...)
	}
	return string(sBuilder)
}

func checkString(s string, n0, n1, n2 int) bool {
	if len(s) != n0+n1+n2+1 {
		return false
	}
	for i := range s {
		if s[i] != '0' && s[i] != '1' {
			return false
		}
	}
	c0 := 0
	c1 := 0
	c2 := 0
	for i := 0; i+1 < len(s); i++ {
		if s[i] == '0' && s[i+1] == '0' {
			c0++
		} else if s[i] == '1' && s[i+1] == '1' {
			c2++
		} else {
			c1++
		}
	}
	return c0 == n0 && c1 == n1 && c2 == n2
}

func runCase(bin string, n0, n1, n2 int) error {
	input := fmt.Sprintf("1\n%d %d %d\n", n0, n1, n2)
	out, err := runBinary(bin, input)
	if err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out)
	}
	s := strings.TrimSpace(out)
	exp := referenceSolution(n0, n1, n2)
	if s != exp {
		return fmt.Errorf("output mismatch: expected %q got %q", exp, s)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	fields := strings.Fields(testcasesRaw)
	if len(fields) == 0 {
		fmt.Println("no testcases provided")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(fields[0])
	if len(fields) != 1+t*3 {
		fmt.Println("embedded testcases malformed")
		os.Exit(1)
	}
	idx := 1
	for i := 0; i < t; i++ {
		n0, _ := strconv.Atoi(fields[idx])
		n1, _ := strconv.Atoi(fields[idx+1])
		n2, _ := strconv.Atoi(fields[idx+2])
		idx += 3
		if err := runCase(bin, n0, n1, n2); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
