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

// Embedded testcases (same format as original file).
const embeddedTestcases = `100
73 26 173 3941
9 13 71 6416
20 117 127 5320
9 16 47 2260
1 4 5 7809
142 61 151 6952
14 40 53 1595
23 38 79 1672
58 26 109 7799
51 31 59 3124
86 21 109 9875
37 26 43 5747
34 50 107 2739
20 18 23 8032
8 12 19 9271
1 20 61 7309
26 15 43 4741
81 7 163 5904
4 1 5 9375
141 6 179 791
4 5 7 4928
43 83 109 431
98 28 151 6894
51 48 61 716
6 8 79 5719
173 170 181 1942
4 162 173 8615
26 147 179 8859
2 4 5 4582
4 6 17 1993
18 8 29 3791
118 78 127 4259
27 7 101 8554
125 20 139 3744
13 22 23 7508
4 4 11 1941
14 28 29 8043
33 78 103 3146
10 142 151 6398
23 16 29 9386
2 10 23 6330
123 117 163 9772
206 108 223 5441
164 223 229 8733
5 20 47 8761
61 88 137 4977
33 28 53 4463
49 42 67 5243
7 20 43 4904
7 16 41 6695
42 81 149 3488
102 60 103 4011
34 24 61 6431
10 8 11 7209
25 43 193 1975
69 68 139 142
51 26 97 407
42 16 167 6815
43 55 107 6223
3 4 5 9381
4 23 163 581
20 18 29 3877
29 28 31 5443
18 65 107 4783
125 64 227 1224
10 40 61 3172
9 1 53 455
36 39 109 8620
2 2 3 2752
2 1 5 8547
46 12 101 5135
2 1 3 4242
8 56 97 7377
43 56 131 1803
133 25 149 8579
41 36 43 6552
36 30 73 9788
3 6 7 4572
12 10 13 9592
7 30 37 376
46 33 127 7127
3 34 37 2890
14 24 31 3257
69 7 71 896
2 37 47 6055
30 52 59 4623
34 22 61 6135
3 21 23 7225
30 25 149 2516
11 14 31 3554
65 91 149 5799
33 56 59 179
17 108 139 1821
108 118 157 1214
26 185 223 7466
38 98 179 3025
72 5 103 5640
36 37 107 2312
11 10 13 7267
81 49 149 9180`

// Embedded reference solution (from 919E.go).
func solve(a, b, p, x int64) string {
	P := int(p)
	powA := make([]int64, P)
	powA[0] = 1 % p
	for i := 1; i < P; i++ {
		powA[i] = powA[i-1] * a % p
	}
	L := p * (p - 1)
	ans := int64(0)
	for e := 0; e < P-1; e++ {
		r := (b * powA[P-1-e]) % p
		k := int64(int64(e)-r) % (p - 1)
		if k < 0 {
			k += p - 1
		}
		n0 := r + p*k
		if n0 <= x {
			ans += 1 + (x-n0)/L
		}
	}
	return fmt.Sprint(ans)
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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) < 2 || len(os.Args) > 3 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	if bin == "--" && len(os.Args) == 3 {
		bin = os.Args[2]
	}

	scan := bufio.NewScanner(strings.NewReader(embeddedTestcases))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Fprintln(os.Stderr, "invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		vals := make([]int64, 4)
		for i := 0; i < 4; i++ {
			if !scan.Scan() {
				fmt.Fprintln(os.Stderr, "bad file")
				os.Exit(1)
			}
			v, _ := strconv.ParseInt(scan.Text(), 10, 64)
			vals[i] = v
		}
		a, b, p, x := vals[0], vals[1], vals[2], vals[3]
		input := fmt.Sprintf("%d %d %d %d\n", a, b, p, x)
		exp := solve(a, b, p, x)
		got, err := runBin(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\n", caseIdx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:%sexpected:%s\ngot:%s\n", caseIdx+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
