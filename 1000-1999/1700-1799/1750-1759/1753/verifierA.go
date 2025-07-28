package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type test struct {
	input string
	n     int
	arr   []int
}

func computePartition(n int, a []int) ([][2]int, bool) {
	sum := 0
	for _, v := range a {
		if v < 0 {
			sum -= v
		} else {
			sum += v
		}
	}
	if sum%2 != 0 {
		return nil, false
	}
	p := []int{1}
	i := 1
	rev := 0
	for {
		for i <= n && a[i-1] == 0 {
			i++
		}
		if i > n {
			break
		}
		l := i
		i++
		for i <= n && a[i-1] == 0 {
			i++
		}
		if i > n {
			break
		}
		r := i
		i++
		if (r-l)%2 == 1 {
			if a[l-1] == a[r-1] {
				continue
			}
			sgnl := (l % 2) ^ rev
			if sgnl != 0 {
				if p[len(p)-1] != r {
					p = append(p, r)
				}
				p = append(p, r+1)
			} else {
				if p[len(p)-1] != l {
					p = append(p, l)
				}
				p = append(p, l+1)
			}
		} else {
			if a[l-1] != a[r-1] {
				continue
			}
			sgnl := (l % 2) ^ rev
			if sgnl != 0 {
				if p[len(p)-1] != r {
					p = append(p, r-1)
				}
				p = append(p, r+1)
				rev ^= 1
			} else {
				if p[len(p)-1] != l {
					p = append(p, l)
				}
				p = append(p, l+1)
			}
		}
	}
	if p[len(p)-1] != n+1 {
		p = append(p, n+1)
	}
	m := len(p) - 1
	segs := make([][2]int, m)
	for j := 0; j < m; j++ {
		segs[j] = [2]int{p[j], p[j+1] - 1}
	}
	return segs, true
}

func verifyOutput(n int, a []int, out string) error {
	fields := strings.Fields(strings.TrimSpace(out))
	if len(fields) == 0 {
		return fmt.Errorf("empty output")
	}
	if fields[0] == "-1" {
		segs, ok := computePartition(n, a)
		if ok && len(segs) > 0 {
			return fmt.Errorf("partition exists but got -1")
		}
		if len(fields) != 1 {
			return fmt.Errorf("extra data after -1")
		}
		return nil
	}
	k, err := strconv.Atoi(fields[0])
	if err != nil {
		return fmt.Errorf("invalid k: %v", err)
	}
	if len(fields) != 1+2*k {
		return fmt.Errorf("expected %d numbers got %d", 1+2*k, len(fields))
	}
	segs := make([][2]int, k)
	idx := 1
	for i := 0; i < k; i++ {
		l, err1 := strconv.Atoi(fields[idx])
		r, err2 := strconv.Atoi(fields[idx+1])
		if err1 != nil || err2 != nil {
			return fmt.Errorf("invalid segment")
		}
		segs[i] = [2]int{l, r}
		idx += 2
	}
	if k == 0 {
		return fmt.Errorf("k=0")
	}
	if segs[0][0] != 1 || segs[k-1][1] != n {
		return fmt.Errorf("segments must cover array")
	}
	for i := 0; i < k; i++ {
		l := segs[i][0]
		r := segs[i][1]
		if l > r || l < 1 || r > n {
			return fmt.Errorf("invalid segment %d %d", l, r)
		}
		if i+1 < k && segs[i+1][0] != r+1 {
			return fmt.Errorf("segments not contiguous")
		}
	}
	total := 0
	for _, seg := range segs {
		sgn := 1
		sum := 0
		for j := seg[0] - 1; j <= seg[1]-1; j++ {
			sum += sgn * a[j]
			sgn *= -1
		}
		total += sum
	}
	if total != 0 {
		return fmt.Errorf("alternating sum not zero")
	}
	return nil
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(42))
	var tests []test
	tests = append(tests, test{input: "1\n1\n1\n", n: 1, arr: []int{1}})
	tests = append(tests, test{input: "1\n2\n1 -1\n", n: 2, arr: []int{1, -1}})
	for len(tests) < 100 {
		n := rng.Intn(6) + 1
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			arr[i] = rng.Intn(3) - 1
		}
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", arr[i]))
		}
		sb.WriteByte('\n')
		tests = append(tests, test{input: sb.String(), n: n, arr: arr})
	}
	return tests
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		out, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := verifyOutput(t.n, t.arr, out); err != nil {
			fmt.Printf("Wrong answer on test %d: %v\nInput:\n%sOutput:%s\n", i+1, err, t.input, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
