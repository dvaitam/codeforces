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

// Embedded reference solver for 446D.
const embeddedSolver446D = `package main

import (
	"fmt"
	"io"
	"math"
	"os"
)

func matMul(a, b []float64, t int) []float64 {
	c := make([]float64, t*t)
	for i := 0; i < t; i++ {
		ai := i * t
		ci := i * t
		for k := 0; k < t; k++ {
			av := a[ai+k]
			if av == 0 {
				continue
			}
			bk := k * t
			for j := 0; j < t; j++ {
				c[ci+j] += av * b[bk+j]
			}
		}
	}
	return c
}

func vecMulMat(v, m []float64, t int) []float64 {
	res := make([]float64, t)
	for i := 0; i < t; i++ {
		vi := v[i]
		if vi == 0 {
			continue
		}
		mi := i * t
		for j := 0; j < t; j++ {
			res[j] += vi * m[mi+j]
		}
	}
	return res
}

func main() {
	data, _ := io.ReadAll(os.Stdin)
	p := 0
	nextInt := func() int {
		for p < len(data) && (data[p] < '0' || data[p] > '9') {
			p++
		}
		v := 0
		for p < len(data) && data[p] >= '0' && data[p] <= '9' {
			v = v*10 + int(data[p]-'0')
			p++
		}
		return v
	}

	n := nextInt()
	m := nextInt()
	k := int64(nextInt())

	trapFlag := make([]int, n)
	for i := 0; i < n; i++ {
		trapFlag[i] = nextInt()
	}

	cnt := make([][]int, n)
	for i := 0; i < n; i++ {
		cnt[i] = make([]int, n)
	}
	deg := make([]int, n)

	for i := 0; i < m; i++ {
		u := nextInt() - 1
		v := nextInt() - 1
		cnt[u][v]++
		cnt[v][u]++
		deg[u]++
		deg[v]++
	}

	trapIdx := make([]int, n)
	nonIdx := make([]int, n)
	for i := 0; i < n; i++ {
		trapIdx[i] = -1
		nonIdx[i] = -1
	}

	traps := make([]int, 0)
	nons := make([]int, 0)
	for i := 0; i < n; i++ {
		if trapFlag[i] == 1 {
			trapIdx[i] = len(traps)
			traps = append(traps, i)
		} else {
			nonIdx[i] = len(nons)
			nons = append(nons, i)
		}
	}

	t := len(traps)
	r := len(nons)
	idxN := trapIdx[n-1]

	tot := r + t
	mat := make([][]float64, r)
	for i := 0; i < r; i++ {
		v := nons[i]
		row := make([]float64, tot)
		row[i] = float64(deg[v])
		for u := 0; u < n; u++ {
			c := cnt[v][u]
			if c == 0 {
				continue
			}
			if trapIdx[u] >= 0 {
				row[r+trapIdx[u]] += float64(c)
			} else if u != v {
				row[nonIdx[u]] -= float64(c)
			}
		}
		mat[i] = row
	}

	for col := 0; col < r; col++ {
		pivot := col
		best := math.Abs(mat[col][col])
		for i := col + 1; i < r; i++ {
			v := math.Abs(mat[i][col])
			if v > best {
				best = v
				pivot = i
			}
		}
		if pivot != col {
			mat[col], mat[pivot] = mat[pivot], mat[col]
		}
		pv := mat[col][col]
		prow := mat[col]
		for i := col + 1; i < r; i++ {
			f := mat[i][col] / pv
			if f == 0 {
				continue
			}
			row := mat[i]
			row[col] = 0
			for j := col + 1; j < tot; j++ {
				row[j] -= f * prow[j]
			}
		}
	}

	X := make([]float64, r*t)
	for i := r - 1; i >= 0; i-- {
		diag := mat[i][i]
		base := i * t
		for j := 0; j < t; j++ {
			sum := mat[i][r+j]
			for c := i + 1; c < r; c++ {
				sum -= mat[i][c] * X[c*t+j]
			}
			X[base+j] = sum / diag
		}
	}

	pi := make([]float64, t)
	startBase := nonIdx[0] * t
	copy(pi, X[startBase:startBase+t])
	sumPi := 0.0
	for j := 0; j < t; j++ {
		sumPi += pi[j]
	}
	if sumPi != 0 {
		for j := 0; j < t; j++ {
			pi[j] /= sumPi
		}
	}

	Q := make([]float64, t*t)
	for i := 0; i < t; i++ {
		v := traps[i]
		rb := i * t
		for u := 0; u < n; u++ {
			c := cnt[v][u]
			if c == 0 {
				continue
			}
			cf := float64(c)
			if trapIdx[u] >= 0 {
				Q[rb+trapIdx[u]] += cf
			} else {
				xb := nonIdx[u] * t
				for j := 0; j < t; j++ {
					Q[rb+j] += cf * X[xb+j]
				}
			}
		}
		d := float64(deg[v])
		sumRow := 0.0
		for j := 0; j < t; j++ {
			Q[rb+j] /= d
			sumRow += Q[rb+j]
		}
		if sumRow != 0 {
			for j := 0; j < t; j++ {
				Q[rb+j] /= sumRow
			}
		}
	}

	exp := k - 2
	vec := pi
	if exp > 0 {
		power := Q
		for exp > 0 {
			if exp&1 == 1 {
				vec = vecMulMat(vec, power, t)
			}
			exp >>= 1
			if exp > 0 {
				power = matMul(power, power, t)
			}
		}
	}

	ans := vec[idxN]
	if ans < 0 {
		ans = 0
	}
	if ans > 1 {
		ans = 1
	}
	fmt.Printf("%.10f\n", ans)
}
`

func buildReference() (string, error) {
	tmpSrc, err := os.CreateTemp("", "446D-src-*.go")
	if err != nil {
		return "", err
	}
	srcPath := tmpSrc.Name()
	if _, err := tmpSrc.WriteString(embeddedSolver446D); err != nil {
		tmpSrc.Close()
		os.Remove(srcPath)
		return "", err
	}
	tmpSrc.Close()
	defer os.Remove(srcPath)

	tmp, err := os.CreateTemp("", "446D-ref-*")
	if err != nil {
		return "", err
	}
	binPath := tmp.Name()
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", binPath, srcPath)
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(binPath)
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return binPath, nil
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(4) + 2
	maxM := n * (n - 1) / 2
	m := n - 1 + rng.Intn(maxM-(n-1)+1)
	// create connected graph via tree
	edges := make([][2]int, 0, m)
	for i := 2; i <= n; i++ {
		u := rng.Intn(i-1) + 1
		edges = append(edges, [2]int{u, i})
	}
	used := make(map[[2]int]bool)
	for _, e := range edges {
		if e[0] > e[1] {
			e[0], e[1] = e[1], e[0]
		}
		used[[2]int{e[0], e[1]}] = true
	}
	for len(edges) < m {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v {
			continue
		}
		a := u
		b := v
		if a > b {
			a, b = b, a
		}
		key := [2]int{a, b}
		if used[key] {
			continue
		}
		used[key] = true
		edges = append(edges, [2]int{u, v})
	}
	k := int64(rng.Intn(4) + 2)
	a := make([]int, n+1)
	trapCnt := 1
	for i := 2; i < n; i++ {
		if rng.Intn(2) == 0 && trapCnt < 3 {
			a[i] = 1
			trapCnt++
		}
	}
	a[n] = 1
	a[1] = 0
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, len(edges), k))
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(a[i]))
	}
	sb.WriteByte('\n')
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	return sb.String()
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

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in := generateCase(rng)
		refOut, err := run(refBin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on case %d: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		var exp float64
		if _, err := fmt.Sscan(refOut, &exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d: cannot parse reference output %q: %v\n", i+1, refOut, err)
			os.Exit(1)
		}
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		var got float64
		if _, err := fmt.Sscan(out, &got); err != nil {
			fmt.Fprintf(os.Stderr, "case %d: cannot parse output: %v\n", i+1, err)
			os.Exit(1)
		}
		if diff := got - exp; diff < -1e-4 || diff > 1e-4 {
			fmt.Fprintf(os.Stderr, "case %d: expected %.6f got %.6f\ninput:\n%s", i+1, exp, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
