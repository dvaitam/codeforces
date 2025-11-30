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

const testcasesE = `100
3868
4970
1691
6490
7846
2540
1477
1090
325
6580
9002
4742
965
3637
8526
8793
5903
4534
2829
1740
4289
3513
421
4265
4453
3170
2701
5077
4746
6102
1421
9927
5529
6356
8290
4078
2913
4053
7760
4588
1464
8973
4920
119
4784
9378
5108
8330
3197
6783
6943
9813
4722
7063
7396
2644
3822
4999
4255
709
1329
759
7581
4595
8502
8760
7721
5618
2377
3205
1089
6764
3321
7228
4527
3010
5830
7143
9647
5254
9151
3256
5301
1655
1010
3750
4547
9539
3890
2002
5425
2909
4767
7521
421
702
5851
1354
4681
5360
`

type pair struct{ a, b int }

func prefixRepeat(s string, d int) int {
	var b strings.Builder
	for b.Len() < d {
		b.WriteString(s)
	}
	str := b.String()[:d]
	val, _ := strconv.Atoi(str)
	return val
}

// referencePairs replicates the accepted solution from 1992E.go
// so the verifier does not depend on an external oracle binary.
func referencePairs(n int) []pair {
	ns := strconv.Itoa(n)
	l := len(ns)
	res := make([]pair, 0)
	for a := 1; a <= 10000; a++ {
		for d := 1; d <= 7; d++ {
			b := l*a - d
			if b < 1 || b > a*n || b > 10000 {
				continue
			}
			if d > l*a {
				continue
			}
			pref := prefixRepeat(ns, d)
			if pref == a*n-b {
				res = append(res, pair{a, b})
			}
		}
	}
	return res
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	data := []byte(testcasesE)
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	expected := make([][]pair, t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		expected[i] = referencePairs(n)
	}
	cmd := exec.Command(os.Args[1])
	cmd.Stdin = bytes.NewReader(data)
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("execution failed:", err)
		os.Exit(1)
	}
	outScan := bufio.NewScanner(bytes.NewReader(out))
	outScan.Split(bufio.ScanWords)
	for i := 0; i < t; i++ {
		if !outScan.Scan() {
			fmt.Printf("missing count for test %d\n", i+1)
			os.Exit(1)
		}
		cnt, _ := strconv.Atoi(outScan.Text())
		if cnt != len(expected[i]) {
			fmt.Printf("test %d failed: expected count %d got %d\n", i+1, len(expected[i]), cnt)
			os.Exit(1)
		}
		for j := 0; j < cnt; j++ {
			if !outScan.Scan() {
				fmt.Printf("missing output for test %d pair %d\n", i+1, j+1)
				os.Exit(1)
			}
			a, _ := strconv.Atoi(outScan.Text())
			outScan.Scan()
			b, _ := strconv.Atoi(outScan.Text())
			if j >= len(expected[i]) || a != expected[i][j].a || b != expected[i][j].b {
				fmt.Printf("test %d pair %d mismatch\n", i+1, j+1)
				os.Exit(1)
			}
		}
	}
	if outScan.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed!")
}
