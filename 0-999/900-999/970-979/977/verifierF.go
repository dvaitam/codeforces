package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

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

func genCase(rng *rand.Rand) string {
	n := rng.Intn(20) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", rng.Intn(30)))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func maxLen(a []int) int {
	m := make(map[int]int)
	best := 0
	for _, v := range a {
		m[v] = m[v-1] + 1
		if m[v] > best {
			best = m[v]
		}
	}
	return best
}

func check(input, output string) error {
	scan := bufio.NewScanner(strings.NewReader(input))
	scan.Split(bufio.ScanWords)
	scan.Scan()
	n, _ := strconv.Atoi(scan.Text())
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		scan.Scan()
		arr[i], _ = strconv.Atoi(scan.Text())
	}
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) < 1 {
		return fmt.Errorf("empty output")
	}
	wantLen, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		return fmt.Errorf("invalid length")
	}
	idx := []int{}
	if wantLen > 0 {
		if len(lines) < 2 {
			return fmt.Errorf("missing index line")
		}
		for _, f := range strings.Fields(lines[1]) {
			v, err := strconv.Atoi(f)
			if err != nil {
				return fmt.Errorf("bad index")
			}
			idx = append(idx, v)
		}
		if len(idx) != wantLen {
			return fmt.Errorf("length mismatch")
		}
		prev := 0
		for _, v := range idx {
			if v < 1 || v > n {
				return fmt.Errorf("index out of range")
			}
			if v <= prev {
				return fmt.Errorf("indices not strictly increasing")
			}
			prev = v
		}
		for i := 1; i < wantLen; i++ {
			if arr[idx[i]-1] != arr[idx[i-1]-1]+1 {
				return fmt.Errorf("values not consecutive")
			}
		}
	}
	if wantLen != maxLen(arr) {
		return fmt.Errorf("reported length incorrect")
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierF.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(6))
	for i := 0; i < 100; i++ {
		in := genCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := check(in, out); err != nil {
			fmt.Printf("case %d failed: %v\ninput:%soutput:%s\n", i+1, err, in, out)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
