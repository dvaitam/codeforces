package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded reference solution (1312C.go).
const solutionSource = `package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for tc := 0; tc < t; tc++ {
       var n int
       var k uint64
       fmt.Fscan(reader, &n, &k)
       a := make([]uint64, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &a[i])
       }
       used := make([]bool, 64)
       ok := true
       for i := 0; i < n && ok; i++ {
           x := a[i]
           pos := 0
           for x > 0 {
               rem := x % k
               if rem > 1 {
                   ok = false
                   break
               }
               if rem == 1 {
                   if pos >= len(used) || used[pos] {
                       ok = false
                       break
                   }
                   used[pos] = true
               }
               x /= k
               pos++
           }
       }
       if ok {
           fmt.Fprintln(writer, "YES")
       } else {
           fmt.Fprintln(writer, "NO")
       }
   }
}
`

const testcasesRaw = `1 2 5
3 3 47 42 19
3 3 38 2 37
2 5 40 25
5 4 34 28 32 17 2
1 4 29
3 5 27 33 10
5 3 15 14 1 11 20
2 3 32 32
3 3 28 50 26
5 4 50 37 22 23 28
2 5 45 47
4 3 31 17 31 32
5 4 42 29 29 22 36
5 5 31 42 14 20 44
2 4 49 30
3 4 45 32 35
5 5 19 46 13 31 32
3 2 50 21 46
1 3 47
1 2 36
1 4 37
2 2 48 33
2 4 15 13
1 5 45
1 2 23
3 3 15 43 1
1 2 4
1 2 46
1 4 16
2 3 47 11
5 2 24 37 2 50 15
2 2 0 22
5 2 18 21 31 1 19
4 2 16 48 25 39
2 5 14 5
3 2 1 28 50
2 5 31 32
3 3 21 16 16
5 5 41 1 44 35 8
1 4 2
2 3 10 6
4 3 32 45 2 15
2 5 4 16
1 3 39
5 4 16 43 27 17 33
1 3 2
4 5 10 7 32 46
1 3 6
1 2 11
2 2 13 1
5 5 29 19 34 41 24
2 3 46 27
4 2 37 37 3 26
5 3 6 42 30 23 1
5 2 39 23 18 44 23
3 2 43 26 6
1 4 12
1 5 3
4 5 29 13 37 39
1 2 18
1 4 19
1 3 48
4 3 7 36 23 25
4 3 48 22 25 7
3 2 7 5 39
3 5 13 44 6
1 5 49
1 5 18
3 5 9 23 17
4 5 46 46 26 31
3 5 14 10 31
5 4 35 27 44 43 44
1 2 4
3 3 34 9 26
1 2 43
1 3 18
4 3 45 42 43 21
4 3 33 18 7 9
5 5 6 21 33 15 45
5 4 10 10 29 45 15
4 4 50 48 36 46
2 5 28 46
1 5 47
2 5 32 3
4 4 25 16 45 46
4 5 23 35 21 45
1 3 34
5 3 25 42 24 40 0
3 5 33 45 29
2 2 1 25
2 5 13 6
4 3 17 47 37 37
2 5 39 8
1 5 30
3 3 29 45 13
1 4 0
4 2 48 37 31 43
3 5 17 32 29
1 2 39`

var _ = solutionSource

type testCase struct {
	n int
	k uint64
	a []uint64
}

func runProg(prog, input string) (string, error) {
	cmd := exec.Command(prog)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("runtime error: %v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

func lineToInput(line string) string {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return ""
	}
	n := fields[0]
	k := fields[1]
	rest := strings.Join(fields[2:], " ")
	return fmt.Sprintf("1\n%s %s\n%s\n", n, k, rest)
}

func parseTests() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	tests := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			return nil, fmt.Errorf("invalid test line %d", idx+1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("bad n on line %d: %v", idx+1, err)
		}
		kVal, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("bad k on line %d: %v", idx+1, err)
		}
		if len(fields) != 2+n {
			return nil, fmt.Errorf("line %d: expected %d numbers, got %d", idx+1, n, len(fields)-2)
		}
		a := make([]uint64, n)
		for i := 0; i < n; i++ {
			val, err := strconv.ParseUint(fields[2+i], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("bad value on line %d: %v", idx+1, err)
			}
			a[i] = val
		}
		tests = append(tests, testCase{n: n, k: kVal, a: a})
	}
	return tests, nil
}

func expected(tc testCase) string {
	used := make([]bool, 64)
	for _, x0 := range tc.a {
		x := x0
		pos := 0
		for x > 0 {
			rem := x % tc.k
			if rem > 1 {
				return "NO"
			}
			if rem == 1 {
				if pos >= len(used) || used[pos] {
					return "NO"
				}
				used[pos] = true
			}
			x /= tc.k
			pos++
		}
	}
	return "YES"
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests, err := parseTests()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for idx, tc := range tests {
		var input strings.Builder
		input.WriteString("1\n")
		input.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.k))
		for i, v := range tc.a {
			input.WriteString(strconv.FormatUint(v, 10))
			if i+1 < len(tc.a) {
				input.WriteByte(' ')
			}
		}
		input.WriteByte('\n')
		want := expected(tc)
		got, err := runProg(bin, input.String())
		if err != nil {
			fmt.Printf("candidate error on test %d: %v\n%s\n", idx+1, err, got)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%s\nGot:\n%s\n", idx+1, input.String(), want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
