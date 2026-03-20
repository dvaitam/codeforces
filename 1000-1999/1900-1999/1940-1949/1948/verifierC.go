package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const refSource = `package main

import (
	"bufio"
	"io"
	"os"
)

type FastScanner struct {
	data []byte
	idx  int
	n    int
}

func NewFastScanner() *FastScanner {
	data, _ := io.ReadAll(os.Stdin)
	return &FastScanner{data: data, n: len(data)}
}

func (fs *FastScanner) skip() {
	for fs.idx < fs.n {
		b := fs.data[fs.idx]
		if b > ' ' {
			break
		}
		fs.idx++
	}
}

func (fs *FastScanner) NextInt() int {
	fs.skip()
	val := 0
	for fs.idx < fs.n {
		b := fs.data[fs.idx]
		if b < '0' || b > '9' {
			break
		}
		val = val*10 + int(b-'0')
		fs.idx++
	}
	return val
}

func (fs *FastScanner) NextBytes() []byte {
	fs.skip()
	start := fs.idx
	for fs.idx < fs.n && fs.data[fs.idx] > ' ' {
		fs.idx++
	}
	return fs.data[start:fs.idx]
}

func solve(n int, a, b []byte) bool {
	rows := [2][]byte{a, b}
	off := 2 * n
	total := 4 * n

	vis := make([]bool, total)
	q := make([]int, total)
	head, tail := 0, 0

	vis[0] = true
	q[tail] = 0
	tail++

	for head < tail {
		id := q[head]
		head++

		phase := 0
		if id >= off {
			phase = 1
			id -= off
		}

		row := id / n
		col := id % n

		if row == 1 && col == n-1 {
			return true
		}

		if phase == 0 {
			if col > 0 {
				nid := off + row*n + col - 1
				if !vis[nid] {
					vis[nid] = true
					q[tail] = nid
					tail++
				}
			}
			if col+1 < n {
				nid := off + row*n + col + 1
				if !vis[nid] {
					vis[nid] = true
					q[tail] = nid
					tail++
				}
			}
			nid := off + (1-row)*n + col
			if !vis[nid] {
				vis[nid] = true
				q[tail] = nid
				tail++
			}
		} else {
			nc := col - 1
			if rows[row][col] == '>' {
				nc = col + 1
			}
			nid := row*n + nc
			if !vis[nid] {
				vis[nid] = true
				q[tail] = nid
				tail++
			}
		}
	}

	return false
}

func main() {
	fs := NewFastScanner()
	out := bufio.NewWriterSize(os.Stdout, 1<<20)
	defer out.Flush()

	t := fs.NextInt()
	for ; t > 0; t-- {
		n := fs.NextInt()
		a := fs.NextBytes()
		b := fs.NextBytes()
		if solve(n, a, b) {
			out.WriteString("YES\n")
		} else {
			out.WriteString("NO\n")
		}
	}
}
`

func buildRef() (string, func()) {
	tmpDir, err := os.MkdirTemp("", "ref1948C")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	srcPath := filepath.Join(tmpDir, "ref.go")
	os.WriteFile(srcPath, []byte(refSource), 0644)
	binPath := filepath.Join(tmpDir, "ref")
	cmd := exec.Command("go", "build", "-o", binPath, srcPath)
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to build ref: %v\n%s\n", err, string(out))
		os.Exit(1)
	}
	return binPath, func() { os.RemoveAll(tmpDir) }
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		return
	}
	bin := os.Args[1]

	refBin, cleanup := buildRef()
	defer cleanup()

	rand.Seed(42)
	const t = 100
	ns := make([]int, t)
	rows1 := make([]string, t)
	rows2 := make([]string, t)
	for i := 0; i < t; i++ {
		n := rand.Intn(19) + 2 // 2..20
		ns[i] = n
		b1 := make([]byte, n)
		b2 := make([]byte, n)
		for j := 0; j < n; j++ {
			if rand.Intn(2) == 0 {
				b1[j] = '<'
			} else {
				b1[j] = '>'
			}
			if rand.Intn(2) == 0 {
				b2[j] = '<'
			} else {
				b2[j] = '>'
			}
		}
		rows1[i] = string(b1)
		rows2[i] = string(b2)
	}

	var input bytes.Buffer
	fmt.Fprintf(&input, "%d\n", t)
	for i := 0; i < t; i++ {
		fmt.Fprintf(&input, "%d\n%s\n%s\n", ns[i], rows1[i], rows2[i])
	}

	// Run reference
	refCmd := exec.Command(refBin)
	refCmd.Stdin = bytes.NewReader(input.Bytes())
	var refOut bytes.Buffer
	refCmd.Stdout = &refOut
	if err := refCmd.Run(); err != nil {
		fmt.Println("failed to run reference:", err)
		os.Exit(1)
	}

	// Run candidate
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		fmt.Println("failed to run binary:", err)
		os.Exit(1)
	}

	refScanner := bufio.NewScanner(bytes.NewReader(refOut.Bytes()))
	scanner := bufio.NewScanner(bytes.NewReader(out.Bytes()))
	for i := 0; i < t; i++ {
		if !refScanner.Scan() {
			fmt.Printf("reference missing output for case %d\n", i+1)
			os.Exit(1)
		}
		if !scanner.Scan() {
			fmt.Printf("missing output for case %d\n", i+1)
			os.Exit(1)
		}
		want := strings.TrimSpace(refScanner.Text())
		got := strings.TrimSpace(scanner.Text())
		if got != want {
			fmt.Printf("case %d: expected %s, got %s\n", i+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
