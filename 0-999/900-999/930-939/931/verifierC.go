package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type result struct {
	minEq  int
	counts [3]int
	minV   int
}

func compute(arr []int) result {
	n := len(arr)
	if n == 0 {
		return result{}
	}
	minV, maxV := arr[0], arr[0]
	for _, v := range arr[1:] {
		if v < minV {
			minV = v
		}
		if v > maxV {
			maxV = v
		}
	}
	cnt0, cnt1, cnt2 := 0, 0, 0
	for _, v := range arr {
		switch v - minV {
		case 0:
			cnt0++
		case 1:
			cnt1++
		case 2:
			cnt2++
		}
	}
	if maxV-minV <= 1 {
		return result{minEq: n, counts: [3]int{cnt0, cnt1, cnt2}, minV: minV}
	}
	best := -1
	best0, best1, best2 := 0, 0, 0
	for i := -n; i <= n; i++ {
		if cnt1-2*i < 0 || cnt0+i < 0 || cnt2+i < 0 {
			continue
		}
		now := 0
		if cnt1-2*i > cnt1 {
			now += (cnt1 - 2*i) - cnt1
		}
		if cnt0+i > cnt0 {
			now += (cnt0 + i) - cnt0
		}
		if cnt2+i > cnt2 {
			now += (cnt2 + i) - cnt2
		}
		if now > best {
			best = now
			best0 = cnt0 + i
			best1 = cnt1 - 2*i
			best2 = cnt2 + i
		}
	}
	return result{minEq: n - best, counts: [3]int{best0, best1, best2}, minV: minV}
}

func intersectionCount(a, b []int) int {
	cntA := map[int]int{}
	for _, v := range a {
		cntA[v]++
	}
	res := 0
	for _, v := range b {
		if cntA[v] > 0 {
			res++
			cntA[v]--
		}
	}
	return res
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	const testcasesRaw = `100
28
-5 -5 -4 -5 -3 -3 -4 -4 -3 -5 -3 -5 -3 -3 -5 -4 -3 -4 -3 -3 -4 -3 -4 -3 -4 -5 -5 -4
15
1 1 2 0 2 0 0 0 0 0 1 0 0 2 2
12
5 5 3 4 4 5 5 4 5 4 4 4
6
3 3 2 3 3 1
16
0 1 1 0 1 0 0 0 1 1 1 1 0 0 1 -1
11
-1 -2 -2 -2 -2 -1 -1 -1 -1 -1 -1
20
5 5 6 4 5 6 5 6 6 4 5 6 4 4 6 4 4 6 6 4
9
4 6 4 6 4 5 4 4 4
14
-5 -4 -4 -5 -5 -3 -5 -5 -5 -5 -5 -5 -3 -5
12
-1 -1 1 -1 1 1 -1 0 1 -1 -1 -1
2
-4 -3
21
-3 -3 -3 -4 -3 -3 -2 -2 -2 -4 -3 -3 -2 -2 -4 -3 -4 -4 -2 -2 -3
27
-4 -3 -4 -2 -2 -3 -3 -2 -3 -4 -3 -3 -3 -2 -3 -2 -4 -2 -2 -4 -2 -4 -3 -4 -4 -4 -4
4
4 2 4 4
30
-5 -5 -3 -4 -5 -4 -5 -3 -5 -3 -3 -3 -4 -4 -3 -4 -4 -3 -5 -5 -5 -4 -4 -5 -5 -3 -3 -5 -5 -5
4
-5 -5 -5 -5
1
5
15
3 4 4 3 2 4 2 4 3 3 4 2 4 4 2
29
3 3 1 1 3 2 2 1 3 1 3 2 2 3 2 2 1 3 2 1 1 2 1 3 1 2 1 2 3
16
2 4 4 2 2 3 2 3 3 4 2 2 3 2 2 4
12
3 2 1 2 2 1 2 1 1 1 3 2
21
1 3 1 1 3 3 2 1 3 3 2 2 2 2 1 2 2 2 3 2 3
24
2 3 2 2 1 1 2 3 2 3 2 3 3 3 1 3 3 3 1 1 2 1 3 1
26
1 1 3 3 1 1 2 2 1 3 3 3 2 2 1 3 2 1 1 3 2 1 2 3 1 3
17
-1 -1 0 1 -1 0 0 1 1 -1 0 0 1 -1 1 0 1
6
3 1 2 2 2 2
23
3 3 2 2 3 2 3 3 3 1 3 1 3 3 1 2 3 2 3 1 2 2 3
23
4 2 2 2 3 2 4 4 4 3 2 2 3 4 2 3 4 4 4 2 3 4 2
1
6
14
3 4 4 2 3 4 2 2 3 2 3 4 4 4
3
5 6 5
15
1 0 -1 -1 1 0 -1 0 0 1 1 1 -1 -1 -1
16
2 3 2 1 1 2 3 2 1 2 1 3 2 3 2 1
28
6 6 7 5 6 7 7 7 6 5 6 6 6 6 7 6 7 7 6 7 5 7 7 5 6 7 7 5
26
5 3 4 3 3 4 5 3 5 4 3 3 3 5 4 4 3 4 4 3 4 4 4 3 4 4
28
0 -1 1 0 1 -1 0 0 -1 -1 -1 0 -1 1 0 -1 0 -1 -1 -1 0 1 1 0 0 1 0 0
29
6 6 6 6 5 6 5 6 6 6 6 5 6 6 6 6 4 5 5 5 5 4 5 5 5 5 4 4 5
13
-3 -3 -2 -3 -2 -3 -1 -1 -2 -3 -1 -2 -3
27
5 6 5 5 6 4 5 5 4 5 4 4 4 5 5 4 5 5 5 4 5 4 4 6 4 6 6
2
1 2
25
4 4 4 4 5 5 5 4 5 6 6 4 5 6 6 6 5 4 4 5 5 5 6 5 6
27
6 4 6 6 6 6 5 5 6 5 4 5 5 6 5 4 4 5 6 4 4 6 6 4 4 5 5
15
-4 -5 -3 -4 -5 -5 -5 -5 -4 -3 -5 -4 -3 -4 -5
28
2 2 3 2 4 2 2 4 4 4 2 2 2 4 2 4 3 4 4 2 2 3 3 2 2 4 2 2
4
2 2 4 2
21
2 3 3 2 3 3 3 3 1 3 1 1 3 1 1 2 1 2 2 2 1
14
-2 -2 -2 -3 -1 -3 -2 -2 -2 -1 -2 -2 -3 -2
23
-1 -1 -1 -3 -3 -1 -1 -3 -3 -3 -2 -1 -3 -1 -3 -1 -3 -1 -2 -1 -1 -2 -1
2
-2 0
11
-3 -5 -5 -5 -3 -3 -3 -5 -5 -3 -4
21
-1 -1 0 1 0 1 1 -1 0 0 1 1 -1 0 0 1 -1 1 1 1 0
30
-5 -4 -4 -4 -3 -3 -4 -3 -3 -3 -5 -4 -5 -3 -5 -4 -4 -4 -4 -4 -3 -4 -3 -5 -3 -4 -4 -5 -4 -3
11
-1 -1 0 1 1 0 -1 1 1 1 1
29
-1 -2 -3 -1 -2 -3 -3 -3 -3 -2 -3 -2 -1 -2 -1 -1 -2 -2 -3 -2 -2 -1 -3 -2 -1 -1 -3 -3 -3
6
5 6 4 6 6 5
30
1 -1 0 0 -1 -1 1 0 -1 -1 1 1 0 1 1 1 0 1 0 -1 1 -1 1 -1 0 1 -1 -1 1 1
13
1 2 2 2 2 2 2 1 1 1 1 3 3
27
-2 -2 -2 -4 -4 -3 -2 -3 -3 -4 -2 -3 -4 -3 -2 -3 -2 -3 -2 -4 -2 -3 -2 -2 -4 -3 -4
1
2
5
5 5 6 5 7
2
4 6
10
-2 -2 -2 -4 -2 -2 -4 -3 -2 -4
2
1 1
9
2 4 3 3 4 2 4 3 2
20
1 2 1 1 1 2 2 0 0 0 1 1 2 0 1 2 2 2 0 2
3
5 5 4
14
1 1 1 3 1 1 3 3 2 3 1 2 1 1
8
4 5 4 5 4 6 4 4
8
3 3 3 5 4 3 5 4
19
5 5 7 7 7 6 6 5 5 6 5 6 7 6 5 7 5 7 6
14
6 5 4 5 6 4 5 4 4 5 6 4 4 6
14
6 5 5 5 5 4 5 5 5 4 5 6 4 5
20
-4 -4 -5 -3 -4 -3 -4 -3 -4 -5 -4 -3 -4 -4 -3 -3 -4 -3 -5 -3
28
3 4 4 4 4 3 5 5 3 5 5 5 4 3 4 5 4 5 5 4 3 5 4 4 5 5 5 5
12
4 3 3 3 5 5 3 4 4 4 5 5
26
-2 -2 -2 -1 -2 -2 -1 -1 -2 -1 -2 -1 -3 -2 -1 -3 -3 -2 -3 -3 -3 -2 -2 -1 -2 -1
13
-1 -2 -3 -1 -1 -2 -2 -1 -3 -2 -3 -1 -3
12
-4 -4 -5 -3 -3 -5 -5 -4 -4 -4 -3 -3
1
4
30
4 4 6 6 5 4 5 5 5 6 6 6 5 5 5 4 6 4 6 5 4 5 5 4 4 6 6 6 6 4
7
5 4 5 5 6 4 6
28
-3 -2 -2 -1 -2 -1 -3 -1 -3 -1 -1 -1 -3 -2 -1 -3 -1 -3 -3 -1 -3 -3 -1 -3 -2 -1 -3 -1
6
0 1 -1 -1 0 -1
10
-3 -4 -4 -4 -5 -5 -4 -5 -5 -4
30
5 6 5 6 7 5 7 7 6 6 5 7 7 5 5 6 5 5 7 6 6 5 5 7 6 6 5 5 5 7
12
6 7 7 5 5 5 5 5 6 5 6 5
12
-2 -2 -2 -4 -3 -4 -3 -2 -2 -4 -3 -2
25
-5 -4 -3 -5 -3 -5 -3 -5 -5 -3 -4 -4 -3 -4 -4 -3 -4 -4 -5 -3 -5 -3 -5 -5 -5
29
-2 -1 -3 -2 -1 -1 -1 -3 -3 -3 -1 -1 -1 -3 -2 -3 -3 -1 -2 -3 -2 -1 -1 -3 -2 -2 -2 -1 -1
12
-1 -3 -1 -3 -3 -1 -2 -3 -2 -1 -3 -1
9
6 6 4 4 6 6 6 6 4
7
7 7 7 5 7 6 5
17
-2 0 -1 0 -1 -1 -2 -1 -2 -2 -1 -2 -1 -2 0 -2 0
26
0 1 -1 0 0 -1 -1 -1 -1 -1 1 -1 0 -1 -1 1 0 0 0 -1 -1 0 -1 1 -1 -1
7
1 1 2 2 1 1 2
28
0 1 -1 1 0 0 1 0 0 0 0 0 0 -1 1 0 0 1 0 1 0 1 1 0 1 -1 1 -1
9
3 2 3 1 3 3 2 3 3
9
-3 -2 -3 -1 -3 -3 -2 -3 -1
29
0 -1 -2 -2 -1 -1 -2 0 -2 0 -1 0 -2 -2 0 -2 -1 0 -2 -1 0 -1 -2 -2 0 -1 -2 -2 0
27
0 -1 0 -1 0 -1 -2 -2 0 0 -1 0 0 -2 -2 -1 -2 0 -1 -1 0 0 -2 -1 -1 -2 -1`

	scan := bufio.NewScanner(strings.NewReader(testcasesRaw))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("empty test file")
		os.Exit(1)
	}
	var t int
	fmt.Sscan(scan.Text(), &t)
	for caseNum := 1; caseNum <= t; caseNum++ {
		if !scan.Scan() {
			fmt.Printf("case %d missing n\n", caseNum)
			os.Exit(1)
		}
		var n int
		fmt.Sscan(scan.Text(), &n)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			if !scan.Scan() {
				fmt.Printf("case %d missing element %d\n", caseNum, i+1)
				os.Exit(1)
			}
			fmt.Sscan(scan.Text(), &arr[i])
		}
		exp := compute(arr)
		var input bytes.Buffer
		fmt.Fprintf(&input, "%d\n", n)
		for i, v := range arr {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(fmt.Sprintf("%d", v))
		}
		input.WriteByte('\n')
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewReader(input.Bytes())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n", caseNum, err)
			os.Exit(1)
		}
		lines := strings.Split(strings.TrimSpace(string(out)), "\n")
		if len(lines) < 2 {
			fmt.Printf("case %d: expected two lines of output\n", caseNum)
			os.Exit(1)
		}
		var gotMin int
		if _, err := fmt.Sscan(strings.TrimSpace(lines[0]), &gotMin); err != nil {
			fmt.Printf("case %d: cannot parse first line: %v\n", caseNum, err)
			os.Exit(1)
		}
		numStrs := strings.Fields(lines[1])
		if len(numStrs) != n {
			fmt.Printf("case %d: expected %d numbers got %d\n", caseNum, n, len(numStrs))
			os.Exit(1)
		}
		outArr := make([]int, n)
		for i, s := range numStrs {
			if _, err := fmt.Sscan(s, &outArr[i]); err != nil {
				fmt.Printf("case %d: bad number in output: %v\n", caseNum, err)
				os.Exit(1)
			}
		}
		// check bounds
		minV := exp.minV
		maxV := minV + 2
		for _, v := range outArr {
			if v < minV || v > maxV {
				fmt.Printf("case %d: output value out of range\n", caseNum)
				os.Exit(1)
			}
		}
		// check average
		sumIn := 0
		for _, v := range arr {
			sumIn += v
		}
		sumOut := 0
		for _, v := range outArr {
			sumOut += v
		}
		if sumIn != sumOut {
			fmt.Printf("case %d: averages differ\n", caseNum)
			os.Exit(1)
		}
		// check minimality
		if gotMin != exp.minEq {
			fmt.Printf("case %d: expected minimal equal %d got %d\n", caseNum, exp.minEq, gotMin)
			os.Exit(1)
		}
		if intersectionCount(arr, outArr) != gotMin {
			fmt.Printf("case %d: reported %d equal but got %d\n", caseNum, gotMin, intersectionCount(arr, outArr))
			os.Exit(1)
		}
		// check counts match expected counts
		counts := map[int]int{}
		for _, v := range outArr {
			counts[v]++
		}
		if counts[minV] != exp.counts[0] || counts[minV+1] != exp.counts[1] || counts[minV+2] != exp.counts[2] {
			fmt.Printf("case %d: counts do not match expected distribution\n", caseNum)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
