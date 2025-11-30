package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcases = `243
312
107
740
407
492
160
94
70
22
413
564
941
298
821
785
62
229
534
551
370
285
800
178
848
110
270
221
967
951
28
850
658
828
268
821
280
200
170
319
298
644
890
751
985
877
870
903
383
90
867
622
347
689
399
520
256
184
255
486
288
93
969
959
839
888
966
562
862
309
9
932
300
588
723
905
321
870
785
522
201
425
435
615
297
443
464
167
240
314
267
834
818
46
85
49
475
643
289
533`

const mod int64 = 998244353

func powMod(a, e int64) int64 {
	a %= mod
	res := int64(1)
	for e > 0 {
		if e&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		e >>= 1
	}
	return res
}

var memo = map[int64][2]int64{}

func xy(n int64) (int64, int64) {
	if n == 1 {
		return 1, 0
	}
	if val, ok := memo[n]; ok {
		return val[0], val[1]
	}
	right := n / 2
	left := n - right
	XL, YL := xy(left)
	XR, YR := xy(right)
	f := powMod(2, n) - powMod(2, left) - powMod(2, right) + 1
	f %= mod
	if f < 0 {
		f += mod
	}
	X := (f + 2*XL + 2*XR) % mod
	Y := (YL + XR + YR) % mod
	memo[n] = [2]int64{X, Y}
	return X, Y
}

func solve(n int64) int64 {
	X, Y := xy(n)
	return (X + Y) % mod
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	lines := strings.Split(testcases, "\n")
	count := 0
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		count++
		n, _ := strconv.ParseInt(line, 10, 64)
		want := solve(n)
		input := fmt.Sprintf("1\n%d\n", n)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		gotVal, err := strconv.ParseInt(strings.TrimSpace(got), 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d output parse error: %v\n", idx+1, err)
			os.Exit(1)
		}
		if gotVal != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\n", idx+1, want, gotVal)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", count)
}
