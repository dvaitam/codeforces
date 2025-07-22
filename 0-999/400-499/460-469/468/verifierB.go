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
	"time"
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
	return out.String(), nil
}

func hasSolution(arr []int, a, b int) bool {
	if b < a {
		a, b = b, a
	}
	m := make(map[int]int)
	for i, v := range arr {
		m[v] = i
	}
	for _, x := range arr {
		if _, ok := m[x]; !ok {
			continue
		}
		if _, ok := m[b-x]; ok {
			delete(m, x)
			delete(m, b-x)
		} else if _, ok := m[a-x]; ok {
			delete(m, x)
			delete(m, a-x)
		} else {
			return false
		}
	}
	return true
}

func checkCase(output string, n, a, b int, arr []int) error {
	scan := bufio.NewScanner(strings.NewReader(output))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		return fmt.Errorf("missing answer")
	}
	ans := strings.ToUpper(scan.Text())
	if ans == "NO" {
		if scan.Scan() {
			return fmt.Errorf("extra output after NO")
		}
		if hasSolution(arr, a, b) {
			return fmt.Errorf("solution exists but output NO")
		}
		return nil
	}
	if ans != "YES" {
		return fmt.Errorf("first token not YES or NO")
	}
	assign := make([]int, n)
	for i := 0; i < n; i++ {
		if !scan.Scan() {
			return fmt.Errorf("missing assignment value")
		}
		v, err := strconv.Atoi(scan.Text())
		if err != nil || (v != 0 && v != 1) {
			return fmt.Errorf("invalid assignment value")
		}
		assign[i] = v
	}
	if scan.Scan() {
		return fmt.Errorf("extra output")
	}
	pos := make(map[int]int)
	for i, v := range arr {
		pos[v] = i
	}
	for i, val := range arr {
		if assign[i] == 0 {
			j, ok := pos[a-val]
			if !ok || assign[j] != 0 {
				return fmt.Errorf("pair for %d in set A missing", val)
			}
		} else {
			j, ok := pos[b-val]
			if !ok || assign[j] != 1 {
				return fmt.Errorf("pair for %d in set B missing", val)
			}
		}
	}
	return nil
}

func generateCase(rng *rand.Rand) (int, int, int, []int) {
	n := rng.Intn(6) + 2
	a := rng.Intn(100) + 1
	b := rng.Intn(100) + 1
	vals := make([]int, 0, n)
	used := map[int]bool{}
	for len(vals) < n {
		v := rng.Intn(100) + 1
		if !used[v] {
			used[v] = true
			vals = append(vals, v)
		}
	}
	return n, a, b, vals
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, a, b, arr := generateCase(rng)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, a, b))
		for j, v := range arr {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteString("\n")
		out, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := checkCase(out, n, a, b, arr); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput: %s", i+1, err, sb.String())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
