package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

const M = 60

func nextPow2(n int) int {
	p := 1
	for p < n {
		p <<= 1
	}
	return p
}

func solveCase(n int, a []int, queries []string) string {
	size := nextPow2(n)
	total := 2 * size * M
	f := make([]int32, total)
	for i := 0; i < size; i++ {
		node := (i + size) * M
		if i < n {
			ai := a[i]
			for t := 0; t < M; t++ {
				if t%ai == 0 {
					f[node+t] = 2
				} else {
					f[node+t] = 1
				}
			}
		} else {
			for t := 0; t < M; t++ {
				f[node+t] = 0
			}
		}
	}
	for i := size - 1; i >= 1; i-- {
		left := 2 * i
		right := left + 1
		base := i * M
		lbase := left * M
		rbase := right * M
		for t := 0; t < M; t++ {
			tL := f[lbase+t]
			tR := f[rbase+((t+int(tL))%M)]
			f[base+t] = tL + tR
		}
	}
	var sb strings.Builder
	for _, q := range queries {
		parts := strings.Fields(q)
		if parts[0] == "A" {
			x := atoi(parts[1])
			y := atoi(parts[2])
			l := x - 1 + size
			r := y - 2 + size
			tmp := make([]int32, M)
			leftArr := make([]int32, M)
			rightNodes := make([]int, 0)
			for t := 0; t < M; t++ {
				leftArr[t] = 0
			}
			for l <= r {
				if l&1 == 1 {
					base := l * M
					for t := 0; t < M; t++ {
						tL := leftArr[t]
						tmp[t] = tL + f[base+((t+int(tL))%M)]
					}
					copy(leftArr, tmp)
					l++
				}
				if r&1 == 0 {
					rightNodes = append(rightNodes, r)
					r--
				}
				l >>= 1
				r >>= 1
			}
			for i := len(rightNodes) - 1; i >= 0; i-- {
				node := rightNodes[i]
				base := node * M
				for t := 0; t < M; t++ {
					tL := leftArr[t]
					tmp[t] = tL + f[base+((t+int(tL))%M)]
				}
				copy(leftArr, tmp)
			}
			sb.WriteString(fmt.Sprintf("%d\n", leftArr[0]))
		} else {
			pos := atoi(parts[1]) - 1
			val := atoi(parts[2])
			a[pos] = val
			nodeIdx := pos + size
			base := nodeIdx * M
			for t := 0; t < M; t++ {
				if t%val == 0 {
					f[base+t] = 2
				} else {
					f[base+t] = 1
				}
			}
			for nodeIdx >>= 1; nodeIdx >= 1; nodeIdx >>= 1 {
				left := 2 * nodeIdx
				right := left + 1
				base := nodeIdx * M
				lbase := left * M
				rbase := right * M
				for t := 0; t < M; t++ {
					tL := f[lbase+t]
					tR := f[rbase+((t+int(tL))%M)]
					f[base+t] = tL + tR
				}
			}
		}
	}
	return sb.String()
}

func atoi(s string) int {
	var v int
	fmt.Sscanf(s, "%d", &v)
	return v
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

func generateCase(rng *rand.Rand) (string, int, []int, []string) {
	n := rng.Intn(10) + 1
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = rng.Intn(5) + 2
	}
	q := rng.Intn(20) + 1
	queries := make([]string, q)
	for i := 0; i < q; i++ {
		if rng.Intn(2) == 0 {
			x := rng.Intn(n) + 1
			y := rng.Intn(n-x+1) + x + 1
			queries[i] = fmt.Sprintf("A %d %d", x, y)
		} else {
			x := rng.Intn(n) + 1
			y := rng.Intn(5) + 2
			queries[i] = fmt.Sprintf("C %d %d", x, y)
			a[x-1] = y
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d ", a[i]))
	}
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("%d\n", q))
	for _, qq := range queries {
		sb.WriteString(fmt.Sprintf("%s\n", qq))
	}
	return sb.String(), n, a, queries
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	inputs := make([]string, 0, 101)
	ns := make([]int, 0, 101)
	arrays := make([][]int, 0, 101)
	qlists := make([][]string, 0, 101)
	in, n, arr, ql := generateCase(rng)
	inputs = append(inputs, in)
	ns = append(ns, n)
	arrays = append(arrays, arr)
	qlists = append(qlists, ql)
	for i := 0; i < 100; i++ {
		in, n, arr, ql := generateCase(rng)
		inputs = append(inputs, in)
		ns = append(ns, n)
		arrays = append(arrays, arr)
		qlists = append(qlists, ql)
	}
	for i := range inputs {
		out, err := runBinary(bin, inputs[i])
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		exp := solveCase(ns[i], arrays[i], qlists[i])
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected:\n%s\nfound:\n%s\n", i+1, inputs[i], exp, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
