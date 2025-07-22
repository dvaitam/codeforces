package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type fig struct{ cnt, cost int64 }

func expectedCase(figs []fig, p []int64) int64 {
	sort.Slice(figs, func(i, j int) bool { return figs[i].cost < figs[j].cost })
	var total int64
	for _, f := range figs {
		total += f.cnt
	}
	seg := make([]int64, 0, len(p)+1)
	var prev int64
	for _, v := range p {
		if prev >= total {
			break
		}
		up := v
		if up > total {
			up = total
		}
		cur := up - prev
		if cur > 0 {
			seg = append(seg, cur)
			prev += cur
		}
	}
	if prev < total {
		seg = append(seg, total-prev)
	}
	var ans int64
	si := 0
	var rem int64
	if len(seg) > 0 {
		rem = seg[0]
	}
	for _, f := range figs {
		cnt := f.cnt
		for cnt > 0 && si < len(seg) {
			take := cnt
			if take > rem {
				take = rem
			}
			ans += take * f.cost * int64(si+1)
			cnt -= take
			rem -= take
			if rem == 0 {
				si++
				if si < len(seg) {
					rem = seg[si]
				}
			}
		}
	}
	return ans
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	data, err := os.ReadFile("testcasesC.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot read testcasesC.txt: %v\n", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Fprintln(os.Stderr, "bad file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	expectedOut := make([]string, t)
	inputs := make([]string, t)
	for caseNum := 0; caseNum < t; caseNum++ {
		scan.Scan()
		n, _ := strconv.Atoi(scan.Text())
		figs := make([]fig, n)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			scan.Scan()
			k, _ := strconv.Atoi(scan.Text())
			scan.Scan()
			c, _ := strconv.Atoi(scan.Text())
			figs[i] = fig{int64(k), int64(c)}
			sb.WriteString(fmt.Sprintf("%d %d\n", k, c))
		}
		scan.Scan()
		t2, _ := strconv.Atoi(scan.Text())
		sb.WriteString(fmt.Sprintf("%d\n", t2))
		p := make([]int64, t2)
		for i := 0; i < t2; i++ {
			scan.Scan()
			v, _ := strconv.Atoi(scan.Text())
			p[i] = int64(v)
			sb.WriteString(fmt.Sprintf("%d", v))
			if i+1 < t2 {
				sb.WriteByte(' ')
			} else {
				sb.WriteByte('\n')
			}
		}
		inputs[caseNum] = sb.String()
		ans := expectedCase(figs, p)
		expectedOut[caseNum] = fmt.Sprintf("%d", ans)
	}

	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(data)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "runtime error: %v\n%s", err, errBuf.String())
		os.Exit(1)
	}
	outScan := bufio.NewScanner(bytes.NewReader(out.Bytes()))
	outScan.Split(bufio.ScanWords)
	for i := 0; i < t; i++ {
		if !outScan.Scan() {
			fmt.Fprintf(os.Stderr, "missing output for case %d\n", i+1)
			os.Exit(1)
		}
		got := outScan.Text()
		if got != expectedOut[i] {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", i+1, expectedOut[i], got)
			os.Exit(1)
		}
	}
	if outScan.Scan() {
		fmt.Fprintln(os.Stderr, "extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
