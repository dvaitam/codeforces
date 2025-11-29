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

// Embedded source for the reference solution (was 1138B.go).
const solutionSource = `package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   _, err := fmt.Fscan(reader, &n)
   if err != nil {
       return
   }
   // Read strings a and b
   var aStr, bStr string
   // Consume next tokens as strings (could be without spaces)
   fmt.Fscan(reader, &aStr)
   fmt.Fscan(reader, &bStr)
   a := make([]int, n)
   b := make([]int, n)
   v := make([]int, n)
   a0, a1, a2 := 0, 0, 0
   sumb := 0
   for i := 0; i < n; i++ {
       if i < len(aStr) {
           a[i] = int(aStr[i] - '0')
       }
       if i < len(bStr) {
           b[i] = int(bStr[i] - '0')
       }
       if b[i] == 1 {
           sumb++
       }
       v[i] = a[i] + b[i]
       switch v[i] {
       case 0:
           a0++
       case 1:
           a1++
       case 2:
           a2++
       }
   }
   half := n / 2
   // Try choose cnt2=i, cnt1=j
   for i2 := 0; i2 <= a2; i2++ {
       for j1 := 0; j1 <= a1; j1++ {
           if i2+j1 > half {
               continue
           }
           if 2*i2+j1 == sumb {
               cnt2, cnt1 := i2, j1
               cnt0 := half - i2 - j1
               if cnt0 > a0 {
                   continue
               }
               out := bufio.NewWriter(os.Stdout)
               defer out.Flush()
               for k := 0; k < n; k++ {
                   if v[k] == 1 && cnt1 > 0 {
                       fmt.Fprintf(out, "%d ", k+1)
                       cnt1--
                   } else if v[k] == 2 && cnt2 > 0 {
                       fmt.Fprintf(out, "%d ", k+1)
                       cnt2--
                   } else if v[k] == 0 && cnt0 > 0 {
                       fmt.Fprintf(out, "%d ", k+1)
                       cnt0--
                   }
               }
               return
           }
       }
   }
   // No solution
   fmt.Print("-1")
}
`

const testcasesRaw = `10 1101110111 1001000001
14 01010100010110 01111011010010
16 0111010100000001 1011111001110100
18 001100001011111101 101110111111000011
14 11011110111111 01101111011011
10 1000010001 1001011110
16 1101001111010001 0001011011101001
12 101101001111 110001110000
16 0111100110101000 0110110000001100
16 1010011000010011 1011100011000111
10 0000111101 0000011011
6 110100 100001
12 101111100110 011111011000
18 011100111000010001 000111001011110001
2 01 10
18 110001110010111011 000000101110110110
16 0010100001011011 0100000110000101
10 0100000000 0111011000
6 110001 001000
18 000101100000001001 100111101000111101
2 00 01
16 1001011010101101 1010101001011110
10 1000011111 1010101011
2 01 11
18 101001110101001101 101101100001000101
4 1101 0010
6 010000 100011
18 111100100100100111 110000010001110001
10 1011000100 0110101011
18 010100010111001100 110111100101000011
8 00001100 00111010
2 10 10
14 11101111011010 11001110011110
18 010010110010110110 000010011000010001
6 000001 111001
16 1110101000011110 1101101000111000
16 0010110000100010 0001011001110000
6 111110 100100
6 011110 111011
12 010111101101 110010000010
10 1101101101 0011101100
16 0101100000101010 0111011101001010
18 011000101001000100 101100010111000110
2 10 11
18 010001011011000010 101001001101111011
6 110110 000111
18 110111111100011001 010001010001111010
18 101001001101010110 010001011101101001
6 001101 001101
8 01111110 10100001
12 110101110101 001111110000
4 0010 1001
4 1011 0011
14 11100111000110 10110100010101
16 0100001111000011 0011001110011111
2 10 10
14 00010011100101 01111100000001
6 110110 011111
12 001100011000 010000100001
6 011001 000001
12 011001011100 001101010100
8 01000111 00001111
10 1001110011 0100110010
2 00 10
6 110000 101010
10 1011101011 0011111001
12 001101000100 100101010110
10 0001011110 1110111010
2 01 01
10 0101001101 1011101001
2 10 01
8 00000111 10111111
10 0111011111 0111010000
10 0111001100 0011000011
2 11 11
4 1111 0011
6 010000 000111
2 00 10
16 0101010001111011 1000100110010110
16 1011011011010110 0010100010000111
12 110010001011 101100010100
12 010010011100 111110000110
14 10010000010100 10101010110010
2 10 11
12 111010011100 000110011101
6 010110 110110
18 111011000111111111 000111101000011011
14 01101010001011 11100110101001
4 0101 0011
16 1110111001110101 1001110111011001
6 111001 100111
12 101010000010 110011110001
6 101101 111101
12 111010011011 110000000010
18 111100110010101110 111110010110111011
6 001101 101111
18 001110011111100001 001001101100111101
10 1001111101 0000111110
18 000011101001011111 001011000010100000
16 0100000101101111 1001100101101110`

var _ = solutionSource

func hasSolution(n int, s, t string) bool {
	a0, a1, a2 := 0, 0, 0
	sumB := 0
	v := make([]int, n)
	for i := 0; i < n; i++ {
		a := int(s[i] - '0')
		b := int(t[i] - '0')
		if b == 1 {
			sumB++
		}
		val := a + b
		v[i] = val
		switch val {
		case 0:
			a0++
		case 1:
			a1++
		case 2:
			a2++
		}
	}
	half := n / 2
	for i2 := 0; i2 <= a2; i2++ {
		for j1 := 0; j1 <= a1; j1++ {
			if i2+j1 > half {
				continue
			}
			if 2*i2+j1 == sumB && half-i2-j1 <= a0 {
				return true
			}
		}
	}
	return false
}

func isValidOutput(n int, s, t string, out string) bool {
	if len(s) != n || len(t) != n {
		return false
	}
	out = strings.TrimSpace(out)
	if out == "-1" {
		return false
	}
	fields := strings.Fields(out)
	if len(fields) != n/2 {
		return false
	}
	seen := make(map[int]bool)
	firstClown := 0
	secondAcrobat := 0
	for _, f := range fields {
		idx, err := strconv.Atoi(f)
		if err != nil || idx < 1 || idx > n || seen[idx] {
			return false
		}
		seen[idx] = true
		if s[idx-1] == '1' {
			firstClown++
		}
	}
	for i := 0; i < n; i++ {
		if !seen[i+1] && t[i] == '1' {
			secondAcrobat++
		}
	}
	return firstClown == secondAcrobat
}

func runProg(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	var _ = solutionSource
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		fields := strings.Fields(line)
		if len(fields) != 3 {
			fmt.Fprintf(os.Stderr, "test %d: invalid line\n", idx)
			os.Exit(1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: bad n\n", idx)
			os.Exit(1)
		}
		aStr, bStr := fields[1], fields[2]
		if len(aStr) != n || len(bStr) != n {
			fmt.Fprintf(os.Stderr, "test %d: length mismatch\n", idx)
			os.Exit(1)
		}
		expectExists := hasSolution(n, aStr, bStr)
		input := fmt.Sprintf("%d\n%s\n%s\n", n, aStr, bStr)
		res, err := runProg(bin, input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", idx, err)
			fmt.Printf("input:\n%s", input)
			os.Exit(1)
		}
		if expectExists {
			if res == "-1" || !isValidOutput(n, aStr, bStr, res) {
				fmt.Printf("test %d failed: invalid answer\n", idx)
				os.Exit(1)
			}
		} else {
			if res != "-1" {
				fmt.Printf("test %d failed: expected -1 got %s\n", idx, res)
				os.Exit(1)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
