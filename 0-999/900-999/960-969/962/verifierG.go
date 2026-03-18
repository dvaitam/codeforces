package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Point struct{ x, y int }

const embeddedRefGo = `package main

import (
	"bufio"
	"fmt"
	"io"
	"math/bits"
	"os"
)

type FastScanner struct {
	data []byte
	idx  int
}

func NewFastScanner() *FastScanner {
	data, _ := io.ReadAll(os.Stdin)
	return &FastScanner{data: data}
}

func (fs *FastScanner) NextInt() int {
	n := len(fs.data)
	for fs.idx < n {
		c := fs.data[fs.idx]
		if c == '-' || (c >= '0' && c <= '9') {
			break
		}
		fs.idx++
	}
	sign := 1
	if fs.data[fs.idx] == '-' {
		sign = -1
		fs.idx++
	}
	val := 0
	for fs.idx < n {
		c := fs.data[fs.idx]
		if c < '0' || c > '9' {
			break
		}
		val = val*10 + int(c-'0')
		fs.idx++
	}
	return sign * val
}

func main() {
	in := NewFastScanner()

	wx1 := in.NextInt()
	wy1 := in.NextInt()
	wx2 := in.NextInt()
	wy2 := in.NextInt()

	if wx1 >= wx2 || wy2 >= wy1 {
		out := bufio.NewWriterSize(os.Stdout, 1<<20)
		fmt.Fprint(out, 0)
		out.Flush()
		return
	}

	n := in.NextInt()
	events := make([][]int, 15002)

	fx := in.NextInt()
	fy := in.NextInt()
	px, py := fx, fy

	for i := 1; i < n; i++ {
		x := in.NextInt()
		y := in.NextInt()
		if x == px {
			a, b := py, y
			if a > b {
				a, b = b, a
			}
			events[a] = append(events[a], x)
			events[b] = append(events[b], x)
		}
		px, py = x, y
	}

	if px == fx {
		a, b := py, fy
		if a > b {
			a, b = b, a
		}
		events[a] = append(events[a], px)
		events[b] = append(events[b], px)
	}

	width := wx2 - wx1
	words := (width + 63) >> 6
	curr := make([]uint64, words)
	prev := make([]uint64, words)

	all := ^uint64(0)
	lastMask := all
	if rem := width & 63; rem != 0 {
		lastMask = (uint64(1) << uint(rem)) - 1
	}

	var aCount, bCount, cCount, dCount int64
	havePrev := false

	for y := 0; y < wy1; y++ {
		for _, x := range events[y] {
			if x >= wx2 {
				continue
			}
			if x <= wx1 {
				for i := 0; i < words; i++ {
					curr[i] ^= all
				}
			} else {
				start := x - wx1
				wi := start >> 6
				bi := uint(start & 63)
				curr[wi] ^= all << bi
				for i := wi + 1; i < words; i++ {
					curr[i] ^= all
				}
			}
		}

		curr[words-1] &= lastMask

		if y < wy2 {
			continue
		}

		var ones, hadj, vadj, blocks int64
		var carry1, carry2 uint64

		if havePrev {
			for i := 0; i < words; i++ {
				w := curr[i]
				ones += int64(bits.OnesCount64(w))
				hadj += int64(bits.OnesCount64(w & ((w << 1) | carry1)))
				carry1 = w >> 63

				c := w & prev[i]
				vadj += int64(bits.OnesCount64(c))
				blocks += int64(bits.OnesCount64(c & ((c << 1) | carry2)))
				carry2 = c >> 63
			}
			cCount += vadj
			dCount += blocks
		} else {
			for i := 0; i < words; i++ {
				w := curr[i]
				ones += int64(bits.OnesCount64(w))
				hadj += int64(bits.OnesCount64(w & ((w << 1) | carry1)))
				carry1 = w >> 63
			}
			havePrev = true
		}

		aCount += ones
		bCount += hadj
		copy(prev, curr)
	}

	ans := aCount - bCount - cCount + dCount

	out := bufio.NewWriterSize(os.Stdout, 1<<20)
	fmt.Fprint(out, ans)
	out.Flush()
}
`

func buildRef() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("getwd: %v", err)
	}
	ref := filepath.Join(wd, "refG.bin")
	goPath := filepath.Join(wd, "refG.go")
	if err := os.WriteFile(goPath, []byte(embeddedRefGo), 0644); err != nil {
		return "", fmt.Errorf("write go: %v", err)
	}
	cmd := exec.Command("go", "build", "-o", ref, goPath)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference go: %v: %s", err, string(out))
	}
	return ref, nil
}

func genCase(rng *rand.Rand) string {
	x1 := rng.Intn(20)
	x2 := x1 + rng.Intn(5) + 1
	y2 := rng.Intn(20)
	y1 := y2 + rng.Intn(5) + 1
	// polygon as rectangle
	px1 := rng.Intn(20)
	py2 := rng.Intn(20)
	px2 := px1 + rng.Intn(5) + 1
	py1 := py2 + rng.Intn(5) + 1
	poly := []Point{{px1, py1}, {px2, py1}, {px2, py2}, {px1, py2}}
	n := len(poly)
	input := fmt.Sprintf("%d %d %d %d\n%d\n", x1, y1, x2, y2, n)
	for _, p := range poly {
		input += fmt.Sprintf("%d %d\n", p.x, p.y)
	}
	return input
}

func runBin(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierG /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer func() {
		os.Remove(ref)
		wd, _ := os.Getwd()
		os.Remove(filepath.Join(wd, "refG.go"))
	}()

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := genCase(rng)
		exp, err := runBin(ref, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: reference error: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		got, err := runBin(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: candidate error: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
