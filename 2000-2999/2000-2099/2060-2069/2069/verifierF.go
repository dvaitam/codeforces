package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type query struct {
	g    byte
	u, v int
}

func run(bin string, input []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func parseOutputs(out string, q int) ([]int, error) {
	reader := strings.NewReader(out)
	res := make([]int, 0, q)
	for i := 0; i < q; i++ {
		var s string
		if _, err := fmt.Fscan(reader, &s); err != nil {
			return nil, fmt.Errorf("output ended early at answer %d: %v", i+1, err)
		}
		v, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("non-integer output at answer %d: %v", i+1, err)
		}
		res = append(res, v)
	}
	var extra string
	if _, err := fmt.Fscan(reader, &extra); err == nil {
		return nil, fmt.Errorf("extra output after %d answers", q)
	}
	return res, nil
}

func genQueries(n, q int) []query {
	qs := make([]query, 0, q)
	for i := 0; i < q; i++ {
		g := byte('A')
		if rand.Intn(2) == 1 {
			g = 'B'
		}
		u := rand.Intn(n) + 1
		v := rand.Intn(n-1) + 1
		if v >= u {
			v++
		}
		qs = append(qs, query{g: g, u: u, v: v})
	}
	return qs
}

func genTest() []byte {
	var sb strings.Builder
	n := rand.Intn(500) + 2
	q := rand.Intn(400) + 1
	// Occasionally push closer to limits.
	if rand.Intn(8) == 0 {
		n = rand.Intn(10000) + 3000
		q = rand.Intn(5000) + 3000
	}
	if rand.Intn(25) == 0 {
		n = 400000
		q = 400000
	}
	sb.WriteString(fmt.Sprintf("%d %d\n", n, q))
	for _, qu := range genQueries(n, q) {
		sb.WriteString(fmt.Sprintf("%c %d %d\n", qu.g, qu.u, qu.v))
	}
	return []byte(sb.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go <binary>")
		os.Exit(1)
	}
	cand := os.Args[1]
	ref := "./refF.bin"
	if err := exec.Command("go", "build", "-o", ref, "2069F.go").Run(); err != nil {
		fmt.Println("failed to build reference solution:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rand.Seed(time.Now().UnixNano())

	for iter := 0; iter < 200; iter++ {
		input := genTest()
		refOut, err := run(ref, input)
		if err != nil {
			fmt.Println("reference failed on iteration", iter+1, ":", err)
			fmt.Println("input:\n", string(input))
			os.Exit(1)
		}
		candOut, err := run(cand, input)
		if err != nil {
			fmt.Printf("candidate runtime error on iteration %d: %v\n", iter+1, err)
			fmt.Println("input:\n", string(input))
			os.Exit(1)
		}

		var q int
		fmt.Fscan(strings.NewReader(string(input)), new(int), &q) // skip n, read q

		refParsed, err := parseOutputs(refOut, q)
		if err != nil {
			fmt.Println("failed to parse reference output on iteration", iter+1, ":", err)
			fmt.Println("output:\n", refOut)
			fmt.Println("input:\n", string(input))
			os.Exit(1)
		}
		candParsed, err := parseOutputs(candOut, q)
		if err != nil {
			fmt.Println("failed to parse candidate output on iteration", iter+1, ":", err)
			fmt.Println("output:\n", candOut)
			fmt.Println("input:\n", string(input))
			os.Exit(1)
		}

		equal := len(refParsed) == len(candParsed)
		if equal {
			for i := range refParsed {
				if refParsed[i] != candParsed[i] {
					equal = false
					break
				}
			}
		}
		if !equal {
			fmt.Printf("mismatch on iteration %d\n", iter+1)
			fmt.Println("input:\n", string(input))
			fmt.Println("reference:\n", refOut)
			fmt.Println("candidate:\n", candOut)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
