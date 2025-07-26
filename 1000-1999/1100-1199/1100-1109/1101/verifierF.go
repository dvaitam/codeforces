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

type caseF struct {
	n, m   int
	pos    []int
	trucks [][4]int
}

func readCases(path string) ([]caseF, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		return nil, fmt.Errorf("bad file")
	}
	t, _ := strconv.Atoi(scan.Text())
	cases := make([]caseF, t)
	for i := 0; i < t; i++ {
		scan.Scan()
		n, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		m, _ := strconv.Atoi(scan.Text())
		pos := make([]int, n)
		for j := 0; j < n; j++ {
			scan.Scan()
			pos[j], _ = strconv.Atoi(scan.Text())
		}
		trucks := make([][4]int, m)
		for j := 0; j < m; j++ {
			scan.Scan()
			s, _ := strconv.Atoi(scan.Text())
			scan.Scan()
			f, _ := strconv.Atoi(scan.Text())
			scan.Scan()
			c, _ := strconv.Atoi(scan.Text())
			scan.Scan()
			r, _ := strconv.Atoi(scan.Text())
			trucks[j] = [4]int{s, f, c, r}
		}
		cases[i] = caseF{n, m, pos, trucks}
	}
	return cases, nil
}

func solveCase(cf caseF) int64 {
	n := cf.n
	A := cf.pos
	N := n
	NN := N * N
	size := NN * N
	dp := make([]int32, size)
	const inf32 = int32(1000000005)
	for i := 0; i < N; i++ {
		iNN := i * NN
		for j := i; j < N; j++ {
			ij := iNN + j*N
			dp[ij] = int32(A[j] - A[i])
			s := i
			for k := 1; k <= j-i; k++ {
				for s < j-1 && dp[iNN+s*N+(k-1)] < int32(A[j]-A[s]) {
					s++
				}
				v := inf32
				if s != i {
					v1 := dp[iNN+(s-1)*N+(k-1)]
					d1 := int32(A[j] - A[s-1])
					if v1 > d1 {
						v = v1
					} else {
						v = d1
					}
				}
				v2 := dp[iNN+s*N+(k-1)]
				d2 := int32(A[j] - A[s])
				t32 := v2
				if d2 > t32 {
					t32 = d2
				}
				if v > t32 {
					v = t32
				}
				dp[ij+k] = v
			}
		}
	}
	var ans int64
	for _, tr := range cf.trucks {
		s := tr[0] - 1
		t := tr[1] - 1
		c := tr[2]
		r := tr[3]
		maxRef := t - s
		if r > maxRef {
			r = maxRef
		}
		val := int64(dp[s*NN+t*N+r]) * int64(c)
		if val > ans {
			ans = val
		}
	}
	return ans
}

func run(bin string, input string) (string, error) {
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
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	if bin == "--" && len(os.Args) > 2 {
		bin = os.Args[2]
	}
	cases, err := readCases("testcasesF.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for idx, cs := range cases {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", cs.n, cs.m))
		for i, v := range cs.pos {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		for _, tr := range cs.trucks {
			sb.WriteString(fmt.Sprintf("%d %d %d %d\n", tr[0], tr[1], tr[2], tr[3]))
		}
		input := sb.String()
		expected := solveCase(cs)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		got, err := strconv.ParseInt(strings.TrimSpace(out), 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d bad output: %v", idx+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\ninput:\n%s", idx+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
