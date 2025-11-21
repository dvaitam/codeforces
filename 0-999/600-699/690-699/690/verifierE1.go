package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const refSource = "0-999/600-699/690-699/690/690E1.go"

type testCase struct {
	input string
	q     int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE1.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for i, tc := range tests {
		wantRaw, err := runProgram(refBin, tc.input)
		if err != nil {
			fail("reference runtime error on test %d: %v\ninput:\n%s", i+1, err, tc.input)
		}
		gotRaw, err := runProgram(candidate, tc.input)
		if err != nil {
			fail("candidate runtime error on test %d: %v\ninput:\n%s", i+1, err, tc.input)
		}
		want, err := parseAnswers(wantRaw, tc.q)
		if err != nil {
			fail("could not parse reference output on test %d: %v\noutput:\n%s", i+1, err, wantRaw)
		}
		got, err := parseAnswers(gotRaw, tc.q)
		if err != nil {
			fail("could not parse candidate output on test %d: %v\noutput:\n%s", i+1, err, gotRaw)
		}
		for idx := 0; idx < tc.q; idx++ {
			if got[idx] != want[idx] {
				fail("wrong answer on test %d, query %d: expected %s got %s\ninput:\n%s", i+1, idx+1, want[idx], got[idx], tc.input)
			}
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "690E1-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("build reference failed: %v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func runProgram(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", filepath.Clean(bin))
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseAnswers(out string, q int) ([]string, error) {
	sc := bufio.NewScanner(strings.NewReader(out))
	sc.Split(bufio.ScanWords)
	sc.Buffer(make([]byte, 1024), 1<<20)
	var ans []string
	for sc.Scan() {
		token := strings.ToUpper(sc.Text())
		if token == "" {
			continue
		}
		if token != "YES" && token != "NO" {
			return nil, fmt.Errorf("unexpected token %q", token)
		}
		ans = append(ans, token)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	if len(ans) != q {
		return nil, fmt.Errorf("expected %d answers, got %d", q, len(ans))
	}
	return ans, nil
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	tests = append(tests, makeUniformTest(1, 2, 1, 0, false))
	tests = append(tests, makeUniformTest(1, 2, 1, 255, true))
	tests = append(tests, makeUniformTest(3, 4, 3, 128, false))
	tests = append(tests, makeUniformTest(2, 6, 6, 64, true))
	for i := 0; i < 30; i++ {
		tests = append(tests, randomTest(rng, rng.Intn(5)+1, 12, 15))
	}
	for i := 0; i < 10; i++ {
		tests = append(tests, randomTest(rng, rng.Intn(3)+1, 200, 200))
	}
	tests = append(tests, randomTest(rng, 2, 600, 600))
	return tests
}

func randomTest(rng *rand.Rand, q, maxH, maxW int) testCase {
	if q <= 0 {
		q = 1
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", q)
	for i := 0; i < q; i++ {
		h := rng.Intn(maxH/2)*2 + 2
		if h > maxH {
			h = maxH
			if h%2 == 1 {
				h--
			}
			if h < 2 {
				h = 2
			}
		}
		w := rng.Intn(maxW) + 1
		fmt.Fprintf(&sb, "%d %d\n", h, w)
		img := randomImage(rng, h, w)
		if rng.Intn(2) == 1 {
			img = swapHalf(img)
		}
		for _, row := range img {
			for j, val := range row {
				if j > 0 {
					sb.WriteByte(' ')
				}
				sb.WriteString(fmt.Sprint(val))
			}
			sb.WriteByte('\n')
		}
	}
	return testCase{input: sb.String(), q: q}
}

func randomImage(rng *rand.Rand, h, w int) [][]int {
	img := make([][]int, h)
	for i := 0; i < h; i++ {
		row := make([]int, w)
		for j := 0; j < w; j++ {
			row[j] = rng.Intn(256)
		}
		img[i] = row
	}
	return img
}

func swapHalf(img [][]int) [][]int {
	h := len(img)
	if h == 0 {
		return img
	}
	mid := h / 2
	res := make([][]int, 0, h)
	res = append(res, img[mid:]...)
	res = append(res, img[:mid]...)
	return res
}

func makeUniformTest(q, h, w int, val int, swapped bool) testCase {
	if h%2 == 1 {
		h++
	}
	img := make([][]int, h)
	for i := range img {
		row := make([]int, w)
		for j := range row {
			row[j] = val
		}
		img[i] = row
	}
	if swapped {
		img = swapHalf(cloneImage(img))
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", q)
	for i := 0; i < q; i++ {
		fmt.Fprintf(&sb, "%d %d\n", h, w)
		for _, row := range img {
			for j, v := range row {
				if j > 0 {
					sb.WriteByte(' ')
				}
				sb.WriteString(fmt.Sprint(v))
			}
			sb.WriteByte('\n')
		}
	}
	return testCase{input: sb.String(), q: q}
}

func cloneImage(img [][]int) [][]int {
	res := make([][]int, len(img))
	for i := range img {
		row := make([]int, len(img[i]))
		copy(row, img[i])
		res[i] = row
	}
	return res
}
