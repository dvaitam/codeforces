package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
)

type point struct{ x, y, t int }

func sanitize(input []byte) ([]byte, error) {
	var n, m, k, q int
	r := bytes.NewReader(input)
	if _, err := fmt.Fscan(r, &n, &m, &k, &q); err != nil {
		return nil, err
	}
	mp := make(map[[2]int]int)
	for i := 0; i < q; i++ {
		var x, y, t int
		if _, err := fmt.Fscan(r, &x, &y, &t); err != nil {
			return nil, err
		}
		key := [2]int{x, y}
		if old, ok := mp[key]; !ok || t < old {
			mp[key] = t
		}
	}
	arr := make([]point, 0, len(mp))
	for k, v := range mp {
		arr = append(arr, point{k[0], k[1], v})
	}
	sort.Slice(arr, func(i, j int) bool {
		if arr[i].x != arr[j].x {
			return arr[i].x < arr[j].x
		}
		return arr[i].y < arr[j].y
	})
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%d %d %d %d\n", n, m, k, len(arr))
	for _, p := range arr {
		fmt.Fprintf(&buf, "%d %d %d\n", p.x, p.y, p.t)
	}
	return buf.Bytes(), nil
}

func runTests(dir, binary string) error {
	files, err := filepath.Glob(filepath.Join(dir, "*.in"))
	if err != nil {
		return err
	}
	sort.Strings(files)
	for _, inFile := range files {
		outFile := inFile[:len(inFile)-3] + ".out"
		input, err := os.ReadFile(inFile)
		if err != nil {
			return err
		}
		sanitized, err := sanitize(input)
		if err != nil {
			return err
		}
		expected, err := os.ReadFile(outFile)
		if err != nil {
			return err
		}
		cmd := exec.Command(binary)
		cmd.Stdin = bytes.NewReader(sanitized)
		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("%s: %v", filepath.Base(inFile), err)
		}
		if string(out) != string(expected) {
			return fmt.Errorf("%s: expected\n%sbut got\n%s", filepath.Base(inFile), expected, out)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	_, file, _, _ := runtime.Caller(0)
	base := filepath.Dir(file)
	testDir := filepath.Join(base, "tests", "D")
	if err := runTests(testDir, binary); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("All tests passed!")
}
