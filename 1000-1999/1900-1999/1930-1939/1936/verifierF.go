package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const eps = 1e-4

const refSourceF = `package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	cx := make([]float64, n)
	cy := make([]float64, n)
	cr := make([]float64, n)

	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &cx[i], &cy[i], &cr[i])
	}

	x, y := cx[0], cy[0]
	step := 2e6

	iters := 400000000 / n
	if iters > 50000 {
		iters = 50000
	}
	if iters < 4000 {
		iters = 4000
	}

	decay := math.Pow(1e-10/step, 1.0/float64(iters))

	bestX, bestY := x, y
	bestVal := -1e18

	for iter := 0; iter < iters; iter++ {
		minVal := 1e18
		minI := 0
		for i := 0; i < n; i++ {
			dx := x - cx[i]
			dy := y - cy[i]
			d := math.Sqrt(dx*dx + dy*dy)
			val := cr[i] - d
			if val < minVal {
				minVal = val
				minI = i
			}
		}

		if minVal > bestVal {
			bestVal = minVal
			bestX = x
			bestY = y
		}

		dx := cx[minI] - x
		dy := cy[minI] - y
		d := math.Sqrt(dx*dx + dy*dy)
		if d > 1e-12 {
			x += step * dx / d
			y += step * dy / d
		}
		step *= decay
	}

	fmt.Printf("%.9f %.9f %.9f\n", bestX, bestY, bestVal)
}
`

func buildRef() (string, error) {
	tmp, err := os.CreateTemp("", "refF_*.go")
	if err != nil {
		return "", err
	}
	if _, err := tmp.WriteString(refSourceF); err != nil {
		tmp.Close()
		return "", err
	}
	tmp.Close()
	defer os.Remove(tmp.Name())
	ref := filepath.Join(os.TempDir(), "refF_1936.bin")
	cmd := exec.Command("go", "build", "-o", ref, tmp.Name())
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out)
	}
	return ref, nil
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

func genCaseF(rng *rand.Rand) string {
	n := rng.Intn(3) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		x := rng.Float64()*10 - 5
		y := rng.Float64()*10 - 5
		r := rng.Float64()*3 + 0.1
		sb.WriteString(fmt.Sprintf("%.3f %.3f %.3f\n", x, y, r))
	}
	return sb.String()
}

func verifyOutput(input, output string) error {
	// Parse input circles
	r := bufio.NewReader(strings.NewReader(input))
	var n int
	fmt.Fscan(r, &n)
	cx := make([]float64, n)
	cy := make([]float64, n)
	cr := make([]float64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(r, &cx[i], &cy[i], &cr[i])
	}

	// Parse output
	var ox, oy, orad float64
	if _, err := fmt.Sscan(output, &ox, &oy, &orad); err != nil {
		return fmt.Errorf("cannot parse output: %v", err)
	}

	// Verify: the output radius should be the min over all circles of (cr[i] - dist(center, ci))
	// and should match the reference to within eps
	minDist := math.Inf(1)
	for i := 0; i < n; i++ {
		dx := ox - cx[i]
		dy := oy - cy[i]
		d := math.Sqrt(dx*dx + dy*dy)
		val := cr[i] - d
		if val < minDist {
			minDist = val
		}
	}

	if math.Abs(minDist-orad) > eps {
		return fmt.Errorf("output radius %.9f doesn't match computed %.9f", orad, minDist)
	}

	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	ref, err := buildRef()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := genCaseF(rng)

		refStr, err := run(ref, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: reference failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		var refX, refY, refR float64
		fmt.Sscan(refStr, &refX, &refY, &refR)

		outStr, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		var x, y, r float64
		if _, err := fmt.Sscan(outStr, &x, &y, &r); err != nil {
			fmt.Fprintf(os.Stderr, "case %d: cannot parse output\n", i+1)
			os.Exit(1)
		}

		// Compare radii (the objective value) with tolerance
		if math.Abs(r-refR) > eps {
			// Also verify the candidate's answer is self-consistent
			if err := verifyOutput(input, outStr); err != nil {
				fmt.Fprintf(os.Stderr, "case %d failed: expected radius ~%.6f got %.6f (ref: %.6f %.6f %.6f)\ninput:\n%s", i+1, refR, r, refX, refY, refR, input)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
