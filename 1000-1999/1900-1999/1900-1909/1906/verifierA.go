package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const testcases = `BBA
BCB
BBB

BCA
CAB
AAC

BCC
CAB
ACA

CBB
CAB
BBC

CAC
BBC
BAC

AAC
BCC
CAC

BBA
CBC
AAC

AAA
CBA
ABC

BAB
CBC
ACB

CAC
CCB
BAC

BBC
ABA
AAA

CCB
BAA
CAA

AAC
CCB
CCB

CAA
CCB
CBB

BCC
CBA
BCA

BCC
BAA
ACB

ACA
BAB
BAA

ACA
ACC
CCC

AAA
CAC
CAB

ABA
ACA
AAC

ABA
CAC
ACB

CAB
AAA
CBB

BAA
CBA
CAC

BAB
BCB
CAC

CAA
CAA
BCB

ACB
CAA
BCB

CCB
CBB
CBA

CCA
BCA
BCA

CBA
ABB
CBC

BCC
CAC
BBC

BCA
ACA
CBA

AAC
BBC
CCB

ABC
CBC
CAA

BAB
CAB
CBC

CAA
BBB
BAB

ACC
ACA
ABC

BBA
AAC
ACC

CAA
ACC
ABB

CAA
BBC
AAB

BAB
ACC
CCB

AAB
AAA
ACB

CBB
CAC
CCC

BCC
BCB
BCA

ABC
BAA
ABB

BBC
ABC
AAB

AAC
BBB
CAB

ABC
AAB
ACC

ABB
BBB
AAC

BBB
BAB
ACB

BAB
BCC
AAC

CBC
AAA
ACA

ABA
ABC
CBB

BCC
CAB
ABA

BCA
BAB
AAC

ACB
CAA
BBB

ACA
ACA
AAB

BBC
AAC
BAC

BCC
BCB
CBB

BCC
ACC
AAB

CCB
BAC
ABC

ABC
BCC
BAA

CAA
AAA
BAB

BAA
CBB
CBC

CAC
ACC
CAC

BAB
CCB
BBC

CAA
CAC
ABC

CBB
ABB
BBB

BAA
CAA
BCB

AAB
BAB
CCA

CCA
BBA
CBB

BAA
BAA
ACC

ACB
ACC
CBA

BBA
ACC
BCC

BCA
CCC
BAB

BCC
ACA
ABB

BAB
BCA
AAB

BAC
CCC
AAA

CCC
CBC
BAA

BCB
AAB
AAB

ACA
BBC
BAB

BCB
ABA
ABC

ABB
AAA
BAC

AAA
ACC
CAC

ABB
CAC
AAA

BAB
BCB
AAA

BBB
BBC
CBA

CAA
CCB
ABA

AAB
CAA
ABB

BCA
ACA
BAC

AAB
BCB
CAC

BBC
BAA
BCB

CCA
BBA
BCB

BAB
BBA
AAB

ACC
CAC
BBA

BBA
ABB
CAB

BCC
ABA
BAC

BCB
AAB
ACB

CCA
BCB
ACB

ACB
BCA
ABC

AAB
BCC
CBB`

func runBinary(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solve(input string) string {
	in := bufio.NewScanner(strings.NewReader(input))
	grid := make([]string, 0, 3)
	for i := 0; i < 3; i++ {
		in.Scan()
		grid = append(grid, strings.TrimSpace(in.Text()))
	}
	dirs := [][2]int{{-1, -1}, {-1, 0}, {-1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 0}, {1, 1}}
	best := "~~~"
	for x1 := 0; x1 < 3; x1++ {
		for y1 := 0; y1 < 3; y1++ {
			for _, d1 := range dirs {
				x2, y2 := x1+d1[0], y1+d1[1]
				if x2 < 0 || x2 >= 3 || y2 < 0 || y2 >= 3 {
					continue
				}
				for _, d2 := range dirs {
					x3, y3 := x2+d2[0], y2+d2[1]
					if x3 < 0 || x3 >= 3 || y3 < 0 || y3 >= 3 {
						continue
					}
					if (x3 == x1 && y3 == y1) || (x3 == x2 && y3 == y2) {
						continue
					}
					word := string([]byte{grid[x1][y1], grid[x2][y2], grid[x3][y3]})
					if word < best {
						best = word
					}
				}
			}
		}
	}
	return best
}

func main() {
	if len(os.Args) < 2 || len(os.Args) > 3 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]

	blocks := strings.Split(testcases, "\n\n")
	count := 0
	for idx, block := range blocks {
		block = strings.TrimSpace(block)
		if block == "" {
			continue
		}
		lines := strings.Split(block, "\n")
		if len(lines) != 3 {
			fmt.Fprintf(os.Stderr, "invalid testcase block at %d\n", idx+1)
			os.Exit(1)
		}
		for _, l := range lines {
			if len(strings.TrimSpace(l)) != 3 {
				fmt.Fprintf(os.Stderr, "invalid line length in case %d\n", idx+1)
				os.Exit(1)
			}
		}
		input := strings.Join(lines, "\n") + "\n"
		exp := solve(input)
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", idx+1, exp, got, input)
			os.Exit(1)
		}
		count++
	}
	fmt.Printf("All %d tests passed\n", count)
}
