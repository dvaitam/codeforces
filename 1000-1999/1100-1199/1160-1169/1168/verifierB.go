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

// Embedded source for the reference solution (was 1168B.go).
const solutionSource = `package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   n := len(s)
   // total intervals
   total := int64(n) * int64(n+1) / 2
   // count small intervals of length up to 8
   var smallCount int64
   maxL := 8
   if n < maxL {
       maxL = n
   }
   for L := 1; L <= maxL; L++ {
       smallCount += int64(n - L + 1)
   }
   // intervals of length >=9 automatically contain a monochromatic 3-term AP (van der Waerden W(2,3)=9)
   longCount := total - smallCount
   // count small intervals that actually contain a progression
   sBytes := []byte(s)
   var smallGood int64
   for L := 1; L <= maxL; L++ {
       for l := 0; l+L <= n; l++ {
           r := l + L - 1
           found := false
           // try all k where 2*k <= L-1
           for k := 1; 2*k <= L-1; k++ {
               for x := l; x+2*k <= r; x++ {
                   b := sBytes[x]
                   if b == sBytes[x+k] && b == sBytes[x+2*k] {
                       found = true
                       break
                   }
               }
               if found {
                   break
               }
           }
           if found {
               smallGood++
           }
       }
   }
   ans := longCount + smallGood
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   fmt.Fprint(w, ans)
}
`

const testcasesRaw = `100
010
11100101
0110010
0
01010011010
01101001011
011110101
0110100
111010110
000
011111010
010110111
1011000
001000010
1010011000101
11111001111
0101010001000
10100111
1
001
01
101100001001010
1010
11101110001
0110101110000
000
11111
10
10010101001010
0011101100
10001000010010
00011111100
10001111110
11
10
1110010010101
10
1010110
000000
0100001
1110010000011
10
000100110001011
111011011
1
0010001111
001
0
110011101
100010
011110100
100110011110001
01000101110
110100110110100
100010010
011110111
00101000010011
1100111
11101
0010111
010101111010
10101100011
1101110101010
010000
1100011
01
00111011
111
010011101000
10110001111000
100011100000
001111000010111
0001110111100
00
0111101000010
011001100000101
00
01011111
00001101101
010100000100011
10
010
110101101010110
10101
00111100101001
1011000
11111
10001
110111111000011
0111111100010
0001001001
00
0010111001011
10101100110
1011001
001111111101
0010101111
000011111111101
010001
011111`

var _ = solutionSource

func expected(s string) int64 {
	n := len(s)
	total := int64(n) * int64(n+1) / 2
	maxL := 8
	if n < maxL {
		maxL = n
	}
	var smallCount int64
	for L := 1; L <= maxL; L++ {
		smallCount += int64(n - L + 1)
	}
	longCount := total - smallCount
	sBytes := []byte(s)
	var smallGood int64
	for L := 1; L <= maxL; L++ {
		for l := 0; l+L <= n; l++ {
			r := l + L - 1
			found := false
			for k := 1; 2*k <= L-1; k++ {
				for x := l; x+2*k <= r; x++ {
					b := sBytes[x]
					if b == sBytes[x+k] && b == sBytes[x+2*k] {
						found = true
						break
					}
				}
				if found {
					break
				}
			}
			if found {
				smallGood++
			}
		}
	}
	return longCount + smallGood
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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
	if !scanner.Scan() {
		fmt.Fprintln(os.Stderr, "empty test file")
		os.Exit(1)
	}
	tline := strings.TrimSpace(scanner.Text())
	t, err := strconv.Atoi(tline)
	if err != nil {
		fmt.Fprintln(os.Stderr, "invalid test count")
		os.Exit(1)
	}
	for i := 0; i < t; i++ {
		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, "missing test case %d\n", i+1)
			os.Exit(1)
		}
		s := strings.TrimSpace(scanner.Text())
		input := s + "\n"
		want := fmt.Sprintf("%d", expected(s))
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %s got %s\n", i+1, want, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", t)
}
