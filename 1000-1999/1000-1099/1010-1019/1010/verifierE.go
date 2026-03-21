package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

const embeddedSolverE = `package main

import (
	"bufio"
	"io"
	"os"
	"sort"
)

type Point2 struct {
	u, v int
}

type Check2 struct {
	u, v int
	idx  int
}

type Point3 struct {
	u, v, w int
}

type Check3 struct {
	u, v, w int
	idx     int
}

const INF = int(1 << 30)

func status(val, l, r int) int {
	if val < l {
		return 1
	}
	if val > r {
		return 2
	}
	return 0
}

func transform(mx, side, val int) int {
	if side == 1 {
		return mx - val + 1
	}
	return val
}

func bitUpdate(bit []int, i, val int) {
	for i < len(bit) {
		if val < bit[i] {
			bit[i] = val
		}
		i += i & -i
	}
}

func bitQuery(bit []int, i int) int {
	res := INF
	for i > 0 {
		if bit[i] < res {
			res = bit[i]
		}
		i -= i & -i
	}
	return res
}

func main() {
	data, _ := io.ReadAll(os.Stdin)
	ptr := 0
	nextInt := func() int {
		for ptr < len(data) && (data[ptr] < '0' || data[ptr] > '9') {
			ptr++
		}
		val := 0
		for ptr < len(data) && data[ptr] >= '0' && data[ptr] <= '9' {
			val = val*10 + int(data[ptr]-'0')
			ptr++
		}
		return val
	}

	xMax := nextInt()
	yMax := nextInt()
	zMax := nextInt()
	n := nextInt()
	m := nextInt()
	k := nextInt()

	maxes := [3]int{xMax, yMax, zMax}
	L := [3]int{INF, INF, INF}
	R := [3]int{0, 0, 0}

	for i := 0; i < n; i++ {
		x := nextInt()
		y := nextInt()
		z := nextInt()
		if x < L[0] {
			L[0] = x
		}
		if x > R[0] {
			R[0] = x
		}
		if y < L[1] {
			L[1] = y
		}
		if y > R[1] {
			R[1] = y
		}
		if z < L[2] {
			L[2] = z
		}
		if z > R[2] {
			R[2] = z
		}
	}

	var statusCode [27][3]int
	var outerDims [27][3]int
	var outerCnt [27]int
	for code := 0; code < 27; code++ {
		tmp := code
		cnt := 0
		for d := 0; d < 3; d++ {
			s := tmp % 3
			tmp /= 3
			statusCode[code][d] = s
			if s != 0 {
				outerDims[code][cnt] = d
				cnt++
			}
		}
		outerCnt[code] = cnt
	}

	var min1 [27]int
	for i := 0; i < 27; i++ {
		min1[i] = INF
	}

	var points2 [27][]Point2
	var points3 [27][]Point3

	for i := 0; i < m; i++ {
		x := nextInt()
		y := nextInt()
		z := nextInt()

		sx := status(x, L[0], R[0])
		sy := status(y, L[1], R[1])
		sz := status(z, L[2], R[2])

		if sx == 0 && sy == 0 && sz == 0 {
			w := bufio.NewWriterSize(os.Stdout, 1<<20)
			w.WriteString("INCORRECT\n")
			w.Flush()
			return
		}

		code := sx + 3*sy + 9*sz
		cnt := outerCnt[code]
		coords := [3]int{x, y, z}

		if cnt == 1 {
			d := outerDims[code][0]
			tv := transform(maxes[d], statusCode[code][d], coords[d])
			if tv < min1[code] {
				min1[code] = tv
			}
		} else if cnt == 2 {
			d0 := outerDims[code][0]
			d1 := outerDims[code][1]
			u := transform(maxes[d0], statusCode[code][d0], coords[d0])
			v := transform(maxes[d1], statusCode[code][d1], coords[d1])
			points2[code] = append(points2[code], Point2{u, v})
		} else {
			u := transform(maxes[0], sx, x)
			v := transform(maxes[1], sy, y)
			w := transform(maxes[2], sz, z)
			points3[code] = append(points3[code], Point3{u, v, w})
		}
	}

	ans := make([]byte, k)
	var checks2 [27][]Check2
	var checks3 [27][]Check3

	for i := 0; i < k; i++ {
		x := nextInt()
		y := nextInt()
		z := nextInt()
		coords := [3]int{x, y, z}
		s := [3]int{
			status(x, L[0], R[0]),
			status(y, L[1], R[1]),
			status(z, L[2], R[2]),
		}

		if s[0] == 0 && s[1] == 0 && s[2] == 0 {
			ans[i] = 1
			continue
		}

		var od [3]int
		t := 0
		for d := 0; d < 3; d++ {
			if s[d] != 0 {
				od[t] = d
				t++
			}
		}

		for mask := 1; mask < (1 << t); mask++ {
			c0, c1, c2 := 0, 0, 0
			subCnt := 0

			if mask&1 != 0 {
				d := od[0]
				if d == 0 {
					c0 = s[0]
				} else if d == 1 {
					c1 = s[1]
				} else {
					c2 = s[2]
				}
				subCnt++
			}
			if t > 1 && mask&2 != 0 {
				d := od[1]
				if d == 0 {
					c0 = s[0]
				} else if d == 1 {
					c1 = s[1]
				} else {
					c2 = s[2]
				}
				subCnt++
			}
			if t > 2 && mask&4 != 0 {
				d := od[2]
				if d == 0 {
					c0 = s[0]
				} else if d == 1 {
					c1 = s[1]
				} else {
					c2 = s[2]
				}
				subCnt++
			}

			code := c0 + 3*c1 + 9*c2

			if subCnt == 1 {
				d := outerDims[code][0]
				thr := transform(maxes[d], statusCode[code][d], coords[d])
				if min1[code] <= thr {
					ans[i] = 2
				}
			} else if subCnt == 2 {
				d0 := outerDims[code][0]
				d1 := outerDims[code][1]
				u := transform(maxes[d0], statusCode[code][d0], coords[d0])
				v := transform(maxes[d1], statusCode[code][d1], coords[d1])
				checks2[code] = append(checks2[code], Check2{u, v, i})
			} else {
				u := transform(maxes[0], c0, x)
				v := transform(maxes[1], c1, y)
				w := transform(maxes[2], c2, z)
				checks3[code] = append(checks3[code], Check3{u, v, w, i})
			}
		}
	}

	for code := 0; code < 27; code++ {
		if outerCnt[code] != 2 || len(points2[code]) == 0 || len(checks2[code]) == 0 {
			continue
		}

		pts := points2[code]
		qs := checks2[code]

		sort.Slice(pts, func(i, j int) bool {
			if pts[i].u == pts[j].u {
				return pts[i].v < pts[j].v
			}
			return pts[i].u < pts[j].u
		})
		sort.Slice(qs, func(i, j int) bool {
			if qs[i].u == qs[j].u {
				return qs[i].v < qs[j].v
			}
			return qs[i].u < qs[j].u
		})

		minV := INF
		p := 0
		for _, q := range qs {
			for p < len(pts) && pts[p].u <= q.u {
				if pts[p].v < minV {
					minV = pts[p].v
				}
				p++
			}
			if minV <= q.v {
				ans[q.idx] = 2
			}
		}
	}

	for code := 0; code < 27; code++ {
		if outerCnt[code] != 3 || len(points3[code]) == 0 || len(checks3[code]) == 0 {
			continue
		}

		pts := points3[code]
		qs := checks3[code]

		sort.Slice(pts, func(i, j int) bool {
			if pts[i].u != pts[j].u {
				return pts[i].u < pts[j].u
			}
			if pts[i].v != pts[j].v {
				return pts[i].v < pts[j].v
			}
			return pts[i].w < pts[j].w
		})
		sort.Slice(qs, func(i, j int) bool {
			if qs[i].u != qs[j].u {
				return qs[i].u < qs[j].u
			}
			if qs[i].v != qs[j].v {
				return qs[i].v < qs[j].v
			}
			return qs[i].w < qs[j].w
		})

		bit := make([]int, yMax+2)
		for i := 1; i < len(bit); i++ {
			bit[i] = INF
		}

		p := 0
		for _, q := range qs {
			for p < len(pts) && pts[p].u <= q.u {
				bitUpdate(bit, pts[p].v, pts[p].w)
				p++
			}
			if bitQuery(bit, q.v) <= q.w {
				ans[q.idx] = 2
			}
		}
	}

	w := bufio.NewWriterSize(os.Stdout, 1<<20)
	w.WriteString("CORRECT\n")
	for i := 0; i < k; i++ {
		if ans[i] == 1 {
			w.WriteString("OPEN\n")
		} else if ans[i] == 2 {
			w.WriteString("CLOSED\n")
		} else {
			w.WriteString("UNKNOWN\n")
		}
	}
	w.Flush()
}
`

func buildEmbeddedOracle() (string, func(), error) {
	tmpSrc, err := os.CreateTemp("", "oracle1010E-*.go")
	if err != nil {
		return "", nil, err
	}
	if _, err := tmpSrc.WriteString(embeddedSolverE); err != nil {
		tmpSrc.Close()
		os.Remove(tmpSrc.Name())
		return "", nil, err
	}
	tmpSrc.Close()

	tmpBin, err := os.CreateTemp("", "oracle1010E-bin-*")
	if err != nil {
		os.Remove(tmpSrc.Name())
		return "", nil, err
	}
	tmpBin.Close()

	cmd := exec.Command("go", "build", "-o", tmpBin.Name(), tmpSrc.Name())
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(tmpSrc.Name())
		os.Remove(tmpBin.Name())
		return "", nil, fmt.Errorf("build oracle: %v\n%s", err, out)
	}
	os.Remove(tmpSrc.Name())
	return tmpBin.Name(), func() { os.Remove(tmpBin.Name()) }, nil
}

func runBinary(bin string, input []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func runTests(dir, binary, oracle string) error {
	files, err := filepath.Glob(filepath.Join(dir, "*.in"))
	if err != nil {
		return err
	}
	sort.Strings(files)
	for _, inFile := range files {
		input, err := os.ReadFile(inFile)
		if err != nil {
			return err
		}
		expected, err := runBinary(oracle, input)
		if err != nil {
			return fmt.Errorf("%s: oracle error: %v", filepath.Base(inFile), err)
		}
		got, err := runBinary(binary, input)
		if err != nil {
			return fmt.Errorf("%s: %v", filepath.Base(inFile), err)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expected) {
			return fmt.Errorf("%s: expected\n%sgot\n%s", filepath.Base(inFile), expected, got)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	_, file, _, _ := runtime.Caller(0)
	base := filepath.Dir(file)
	testDir := filepath.Join(base, "tests", "E")

	oracle, cleanup, err := buildEmbeddedOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	if err := runTests(testDir, binary, oracle); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("All tests passed!")
}
