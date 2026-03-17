package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

func parseTest(in []byte) (int, int, [][]int) {
	sc := bufio.NewScanner(bytes.NewReader(in))
	sc.Split(bufio.ScanWords)
	readInt := func() int {
		sc.Scan()
		v, _ := strconv.Atoi(sc.Text())
		return v
	}
	n := readInt()
	m := readInt()
	acts := make([][]int, n)
	for i := 0; i < n; i++ {
		k := readInt()
		a := make([]int, k)
		for j := 0; j < k; j++ {
			a[j] = readInt()
		}
		sort.Ints(a)
		acts[i] = a
	}
	return n, m, acts
}

func validateAnswer(in []byte, expected, got string) bool {
	expLines := strings.Fields(expected)
	gotLines := strings.Fields(got)

	if len(expLines) == 0 || len(gotLines) == 0 {
		return false
	}

	expYes := strings.ToUpper(expLines[0]) == "YES"
	gotYes := strings.ToUpper(gotLines[0]) == "YES"

	if !expYes && !gotYes {
		return true // both NO
	}
	if expYes != gotYes {
		return false // disagree on YES/NO
	}

	// Both YES, validate candidate's pair
	if len(gotLines) < 3 {
		return false
	}
	u, err1 := strconv.Atoi(gotLines[1])
	v, err2 := strconv.Atoi(gotLines[2])
	if err1 != nil || err2 != nil {
		return false
	}

	_, _, acts := parseTest(in)
	u--
	v--
	if u < 0 || u >= len(acts) || v < 0 || v >= len(acts) || u == v {
		return false
	}

	// Check goodPair: shared element, unique in both
	a, b := acts[u], acts[v]
	i, j := 0, 0
	shared, diffA, diffB := false, false, false
	for i < len(a) && j < len(b) {
		if a[i] == b[j] {
			shared = true
			i++
			j++
		} else if a[i] < b[j] {
			diffA = true
			i++
		} else {
			diffB = true
			j++
		}
	}
	if i < len(a) { diffA = true }
	if j < len(b) { diffB = true }
	return shared && diffA && diffB
}

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		if out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput(); err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, string(out))
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func runProg(exe string, input []byte) (string, error) {
	cmd := exec.Command(exe)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return out.String() + errBuf.String(), err
	}
	return strings.TrimSpace(out.String()), nil
}

func buildRef() (string, error) {
	src := os.Getenv("REFERENCE_SOURCE_PATH")
	if src == "" {
		return "", fmt.Errorf("REFERENCE_SOURCE_PATH not set")
	}
	content, err := os.ReadFile(src)
	if err != nil {
		return "", fmt.Errorf("read reference: %v", err)
	}
	refBin := "/tmp/1949F_ref.bin"
	if strings.Contains(string(content), "#include") {
		cppSrc := "/tmp/1949F_ref.cpp"
		if err := os.WriteFile(cppSrc, content, 0644); err != nil {
			return "", err
		}
		cmd := exec.Command("g++", "-O2", "-o", refBin, cppSrc)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", fmt.Errorf("build reference failed: %v\n%s", err, string(out))
		}
	} else {
		cmd := exec.Command("go", "build", "-o", refBin, src)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", fmt.Errorf("build reference failed: %v\n%s", err, string(out))
		}
	}
	return refBin, nil
}

func genTest() []byte {
	n := rand.Intn(5) + 2
	m := rand.Intn(8) + n
	if m < 4 {
		m = 4
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		maxK := 4
		if m < maxK {
			maxK = m
		}
		k := rand.Intn(maxK) + 1
		// Generate k unique activities
		used := make(map[int]bool)
		acts := make([]int, 0, k)
		for len(acts) < k {
			val := rand.Intn(m) + 1
			if !used[val] {
				used[val] = true
				acts = append(acts, val)
			}
		}
		sb.WriteString(fmt.Sprintf("%d", k))
		for _, a := range acts {
			sb.WriteString(fmt.Sprintf(" %d", a))
		}
		sb.WriteByte('\n')
	}
	return []byte(sb.String())
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		return
	}
	exe, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer cleanup()
	ref, err := buildRef()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer os.Remove(ref)
	for i := 1; i <= 100; i++ {
		in := genTest()
		expected, err := runProg(ref, in)
		if err != nil {
			fmt.Printf("reference failed on test %d: %v\n%s", i, err, expected)
			os.Exit(1)
		}
		got, err := runProg(exe, in)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n%s", i, err, got)
			os.Exit(1)
		}
		if !validateAnswer(in, expected, got) {
			fmt.Printf("wrong answer on test %d\ninput:\n%sexpected:%s\ngot:%s\n", i, string(in), expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
