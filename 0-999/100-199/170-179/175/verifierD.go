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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type params struct {
	hp, dt, l, r, p int
}

func solveCase(a, b params) string {
	hp := [2]int{a.hp, b.hp}
	dtv := [2]int{a.dt, b.dt}
	l := [2]int{a.l, b.l}
	r := [2]int{a.r, b.r}
	p := [2]int{a.p, b.p}
	z := [2]int{r[0] - l[0] + 1, r[1] - l[1] + 1}
	var state [2][2][205]float64
	w := 0
	if p[0] == 100 || p[1] == 100 {
		if p[0] == 100 {
			return fmt.Sprintf("%.6f", 0.0)
		}
		return fmt.Sprintf("%.6f", 1.0)
	}
	state[w][0][hp[0]] = 1.0
	state[w][1][hp[1]] = 1.0
	var result float64
	goStep := func(f int) {
		wp := 1 - w
		for j := 0; j <= hp[f]; j++ {
			state[wp][f][j] = float64(p[1-f]) * 0.01 * state[w][f][j]
		}
		var s float64
		for j := hp[f]; j >= 0; j-- {
			if j+l[1-f] <= hp[f] {
				s += state[w][f][j+l[1-f]]
			}
			if j+r[1-f]+1 <= hp[f] {
				s -= state[w][f][j+r[1-f]+1]
			}
			state[wp][f][j] += 0.01 * float64(100-p[1-f]) / float64(z[1-f]) * s
		}
		for j := 0; j <= hp[f]; j++ {
			if j >= r[1-f] {
				continue
			}
			a := j - r[1-f]
			b := min(j-l[1-f], -1)
			pq := float64(b-a+1) / float64(z[1-f]) * float64(100-p[1-f]) * 0.01
			state[wp][f][0] += pq * state[w][f][j]
		}
		for j := 0; j <= hp[1-f]; j++ {
			state[wp][1-f][j] = state[w][1-f][j]
		}
		w = wp
	}
outer:
	for t := 0; ; t++ {
		if t%dtv[0] == 0 {
			result -= state[w][1][0] * (1.0 - state[w][0][0])
			goStep(1)
			result += state[w][1][0] * (1.0 - state[w][0][0])
			q := state[w][0][0] + state[w][1][0] - state[w][0][0]*state[w][1][0]
			if q+1e-7 > 1.0 {
				break outer
			}
		}
		if t%dtv[1] == 0 {
			goStep(0)
		}
	}
	return fmt.Sprintf("%.6f", result)
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	data, err := os.ReadFile("testcasesD.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot read testcasesD.txt: %v\n", err)
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
	for i := 0; i < t; i++ {
		var a, b params
		scan.Scan()
		a.hp, _ = strconv.Atoi(scan.Text())
		scan.Scan()
		a.dt, _ = strconv.Atoi(scan.Text())
		scan.Scan()
		a.l, _ = strconv.Atoi(scan.Text())
		scan.Scan()
		a.r, _ = strconv.Atoi(scan.Text())
		scan.Scan()
		a.p, _ = strconv.Atoi(scan.Text())
		scan.Scan()
		b.hp, _ = strconv.Atoi(scan.Text())
		scan.Scan()
		b.dt, _ = strconv.Atoi(scan.Text())
		scan.Scan()
		b.l, _ = strconv.Atoi(scan.Text())
		scan.Scan()
		b.r, _ = strconv.Atoi(scan.Text())
		scan.Scan()
		b.p, _ = strconv.Atoi(scan.Text())
		expectedOut[i] = solveCase(a, b)
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
