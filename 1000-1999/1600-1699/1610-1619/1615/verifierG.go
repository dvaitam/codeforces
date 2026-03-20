package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const refSourceG = `package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func lis(a []int) int {
	d := make([]int, len(a))
	l := 0
	for _, x := range a {
		i := sort.Search(l, func(i int) bool { return d[i] >= x })
		if i == l {
			d[l] = x
			l++
		} else {
			d[i] = x
		}
	}
	return l
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	if _, err := fmt.Fscan(in, &n, &q); err != nil {
		return
	}
	m := 2*n + 1
	p := make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &p[i])
	}

	for ; q > 0; q-- {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u--
		v--
		p[u], p[v] = p[v], p[u]

		ans := -1
		for k := 0; k < m; k++ {
			b := make([]int, m)
			for i := 0; i < m; i++ {
				b[i] = p[(k+i)%m]
			}
			if lis(b) <= n {
				ans = k
				break
			}
		}
		fmt.Fprintln(out, ans)
	}
}
`

type Test struct {
	input string
}

func runExe(path, input string) (string, error) {
	cmd := exec.Command(path)
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		err = fmt.Errorf("%v: %s", err, errBuf.String())
	}
	return out.String(), err
}

func genTests() []Test {
	rand.Seed(6)
	tests := make([]Test, 0, 20)
	for t := 0; t < 19; t++ {
		n := rand.Intn(5) + 1
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			sb.WriteString(fmt.Sprintf("%d ", rand.Intn(4)))
		}
		sb.WriteByte('\n')
		tests = append(tests, Test{sb.String()})
	}
	tests = append(tests, Test{"1\n0\n"})
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		return
	}
	bin := os.Args[1]

	tmp, err := os.CreateTemp("", "refG_*.go")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create temp file: %v\n", err)
		os.Exit(1)
	}
	if _, err := tmp.WriteString(refSourceG); err != nil {
		tmp.Close()
		fmt.Fprintf(os.Stderr, "failed to write temp file: %v\n", err)
		os.Exit(1)
	}
	tmp.Close()
	defer os.Remove(tmp.Name())

	ref := filepath.Join(os.TempDir(), "refG_1615.bin")
	cmd := exec.Command("go", "build", "-o", ref, tmp.Name())
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("build reference: %v: %s", err, string(out)))
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := genTests()
	for i, tc := range tests {
		want, err := runExe(ref, tc.input)
		if err != nil {
			fmt.Printf("reference runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runExe(bin, tc.input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(want) != strings.TrimSpace(got) {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%sGot:\n%s\n", i+1, tc.input, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
