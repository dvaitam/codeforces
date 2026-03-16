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

type pair struct{ n, m int64 }

func expected(x int64) string {
	res := make([]pair, 0)
	for n := int64(1); ; n++ {
		minSq := n * (n + 1) * (2*n + 1) / 6
		if minSq > x {
			break
		}
		denom := n * (n + 1)
		sixx := 6 * x
		if sixx%denom != 0 {
			continue
		}
		t := sixx / denom
		if (t+n-1)%3 != 0 {
			continue
		}
		m := (t + n - 1) / 3
		if m < n {
			continue
		}
		if n*(n+1)*(3*m-n+1)/6 == x {
			res = append(res, pair{int64(n), m})
			if m != n {
				res = append(res, pair{m, int64(n)})
			}
		}
	}
	// sort
	for i := 0; i < len(res); i++ {
		for j := i + 1; j < len(res); j++ {
			if res[j].n < res[i].n || (res[j].n == res[i].n && res[j].m < res[i].m) {
				res[i], res[j] = res[j], res[i]
			}
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(res)))
	for _, p := range res {
		sb.WriteString(fmt.Sprintf("%d %d\n", p.n, p.m))
	}
	return sb.String()
}

func runCase(exe, input, exp string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp = strings.TrimSpace(exp)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

const testcasesDRaw = `100
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	data := []byte(testcasesDRaw)
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		x, _ := strconv.ParseInt(scan.Text(), 10, 64)
		input := fmt.Sprintf("%d\n", x)
		exp := expected(x)
		if err := runCase(exe, input, exp); err != nil {
			fmt.Printf("case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
