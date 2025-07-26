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
)

func buildRef() (string, func(), error) {
	_, cur, _, _ := runtime.Caller(0)
	dir := filepath.Dir(cur)
	src := filepath.Join(dir, "1284E.go")
	tmpDir, err := os.MkdirTemp("", "ref1284E")
	if err != nil {
		return "", nil, err
	}
	bin := filepath.Join(tmpDir, "refbin")
	cmd := exec.Command("go", "build", "-o", bin, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	cleanup := func() { os.RemoveAll(tmpDir) }
	return bin, cleanup, nil
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
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func collinear(x1, y1, x2, y2, x3, y3 int) bool {
	return (x2-x1)*(y3-y1)-(y2-y1)*(x3-x1) == 0
}

func genCase(r *rand.Rand) string {
	n := r.Intn(8) + 5
	type pt struct{ x, y int }
	pts := make([]pt, 0, n)
	for len(pts) < n {
		x := r.Intn(200) - 100
		y := r.Intn(200) - 100
		ok := true
		for j := 0; j < len(pts); j++ {
			if pts[j].x == x && pts[j].y == y {
				ok = false
				break
			}
		}
		if !ok {
			continue
		}
		col := false
		for j := 0; j < len(pts) && !col; j++ {
			for k := j + 1; k < len(pts); k++ {
				if collinear(pts[j].x, pts[j].y, pts[k].x, pts[k].y, x, y) {
					col = true
					break
				}
			}
		}
		if col {
			continue
		}
		pts = append(pts, pt{x, y})
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for _, p := range pts {
		sb.WriteString(fmt.Sprintf("%d %d\n", p.x, p.y))
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}
	candidate := os.Args[1]
	ref, cleanup, err := buildRef()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer cleanup()
	r := rand.New(rand.NewSource(1))
	for tc := 1; tc <= 100; tc++ {
		in := genCase(r)
		want, err := run(ref, in)
		if err != nil {
			fmt.Printf("reference failed on case %d: %v\n", tc, err)
			os.Exit(1)
		}
		got, err := run(candidate, in)
		if err != nil {
			fmt.Printf("case %d: %v\n", tc, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(want) {
			fmt.Printf("case %d failed\ninput:\n%sexpected: %q\ngot: %q\n", tc, in, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All 100 tests passed")
}
