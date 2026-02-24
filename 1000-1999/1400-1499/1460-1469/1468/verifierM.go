package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"strconv"
)

func buildRef() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	src := filepath.Join(dir, "1468M.go")
	bin := filepath.Join(os.TempDir(), "ref1468M.bin")
	cmd := exec.Command("go", "build", "-o", bin, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\\n%s", err, out)
	}
	return bin, nil
}

func runBinary(bin string, input []byte) ([]byte, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("runtime error: %v\\n%s", err, errBuf.String())
	}
	return out.Bytes(), nil
}

type Case struct{ input []byte }

func genCases() []Case {
	rng := rand.New(rand.NewSource(1468))
	cases := make([]Case, 100)
	for i := range cases {
		n := rng.Intn(4) + 1
		var buf strings.Builder
		fmt.Fprintf(&buf, "1\\n%d\\n", n)
		for j := 0; j < n; j++ {
			m := rng.Intn(4) + 2
			fmt.Fprintf(&buf, "%d", m)
			used := make(map[int]bool)
			for k := 0; k < m; {
				val := rng.Intn(10) + 1
				if !used[val] {
					used[val] = true
					fmt.Fprintf(&buf, " %d", val)
					k++
				}
			}
			buf.WriteByte('\\')
			buf.WriteByte('n')
		}
		
		inputStr := buf.String()
		inputStr = strings.ReplaceAll(inputStr, "\\n", "\n")
		cases[i] = Case{[]byte(inputStr)}
	}
	return cases
}

func parseSets(input string) [][]int {
	tokens := strings.Fields(input)
	if len(tokens) < 2 {
		return nil
	}
	t, _ := strconv.Atoi(tokens[0])
	if t != 1 {
		return nil
	}
	n, _ := strconv.Atoi(tokens[1])
	sets := make([][]int, n)
	idx := 2
	for i := 0; i < n; i++ {
		if idx >= len(tokens) {
			break
		}
		m, _ := strconv.Atoi(tokens[idx])
		idx++
		sets[i] = make([]int, m)
		for j := 0; j < m; j++ {
			if idx < len(tokens) {
				sets[i][j], _ = strconv.Atoi(tokens[idx])
				idx++
			}
		}
	}
	return sets
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierM.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	cases := genCases()
	for i, c := range cases {
		expBytes, err := runBinary(ref, c.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference error on case %d: %v\\n", i+1, err)
			os.Exit(1)
		}
		outBytes, err := runBinary(bin, c.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on case %d: %v\\n", i+1, err)
			os.Exit(1)
		}
		exp := strings.TrimSpace(string(expBytes))
		got := strings.TrimSpace(string(outBytes))

		if exp == "-1" {
			if got != "-1" {
				fmt.Printf("wrong answer on case %d\\ninput:\\n%sexpected:\\n%s\\ngot:\\n%s\\n", i+1, string(c.input), exp, got)
				os.Exit(1)
			}
		} else {
			if got == "-1" {
				fmt.Printf("wrong answer on case %d\\ninput:\\n%sexpected:\\n%s\\ngot:\\n%s\\n", i+1, string(c.input), exp, got)
				os.Exit(1)
			}
			parts := strings.Fields(got)
			if len(parts) != 2 {
				fmt.Printf("wrong answer on case %d (invalid format)\\ninput:\\n%sexpected:\\n%s\\ngot:\\n%s\\n", i+1, string(c.input), exp, got)
				os.Exit(1)
			}
			u, err1 := strconv.Atoi(parts[0])
			v, err2 := strconv.Atoi(parts[1])
			if err1 != nil || err2 != nil {
				fmt.Printf("wrong answer on case %d (non-integer indices)\\ninput:\\n%sexpected:\\n%s\\ngot:\\n%s\\n", i+1, string(c.input), exp, got)
				os.Exit(1)
			}
			sets := parseSets(string(c.input))
			if u < 1 || u > len(sets) || v < 1 || v > len(sets) || u == v {
				fmt.Printf("wrong answer on case %d (invalid indices)\\ninput:\\n%sexpected:\\n%s\\ngot:\\n%s\\n", i+1, string(c.input), exp, got)
				os.Exit(1)
			}
			
			setU := sets[u-1]
			setV := sets[v-1]
			common := 0
			mapU := make(map[int]bool)
			for _, val := range setU {
				mapU[val] = true
			}
			for _, val := range setV {
				if mapU[val] {
					common++
				}
			}
			if common < 2 {
				fmt.Printf("wrong answer on case %d (intersection < 2)\\ninput:\\n%sexpected:\\n%s\\ngot:\\n%s\\n", i+1, string(c.input), exp, got)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
