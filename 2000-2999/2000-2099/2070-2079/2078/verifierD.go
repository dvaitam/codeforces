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

type gate struct {
	isMul bool
	val   int
}

func run(bin string, input []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func parseOutputs(out string, t int) ([]int64, error) {
	reader := strings.NewReader(out)
	res := make([]int64, 0, t)
	for i := 0; i < t; i++ {
		var tok string
		if _, err := fmt.Fscan(reader, &tok); err != nil {
			return nil, fmt.Errorf("output ended early on case %d: %v", i+1, err)
		}
		v, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer on case %d: %v", i+1, err)
		}
		res = append(res, v)
	}
	var extra string
	if _, err := fmt.Fscan(reader, &extra); err == nil {
		return nil, fmt.Errorf("extra tokens after %d answers", t)
	}
	return res, nil
}

func genCase() (int, []gate, []gate) {
	n := rand.Intn(30) + 1
	ls := make([]gate, n)
	rs := make([]gate, n)
	for i := 0; i < n; i++ {
		for side := 0; side < 2; side++ {
			isMul := rand.Intn(3) == 0
			if rand.Intn(8) == 0 {
				// bias toward extremes
				isMul = rand.Intn(2) == 0
			}
			val := 0
			if isMul {
				val = 2 + rand.Intn(2) // 2 or 3
			} else {
				if rand.Intn(5) == 0 {
					val = 1000
				} else {
					val = rand.Intn(1000) + 1
				}
			}
			g := gate{isMul: isMul, val: val}
			if side == 0 {
				ls[i] = g
			} else {
				rs[i] = g
			}
		}
	}
	// Occasionally produce deterministic edge cases.
	switch rand.Intn(6) {
	case 0:
		for i := 0; i < n; i++ {
			ls[i] = gate{isMul: true, val: 3}
			rs[i] = gate{isMul: true, val: 3}
		}
	case 1:
		for i := 0; i < n; i++ {
			ls[i] = gate{isMul: false, val: 1000}
			rs[i] = gate{isMul: false, val: 1000}
		}
	}
	return n, ls, rs
}

func buildInput() []byte {
	t := rand.Intn(5) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for i := 0; i < t; i++ {
		n, l, r := genCase()
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for j := 0; j < n; j++ {
			if l[j].isMul {
				sb.WriteString("x ")
			} else {
				sb.WriteString("+ ")
			}
			sb.WriteString(fmt.Sprint(l[j].val))
			sb.WriteByte(' ')
			if r[j].isMul {
				sb.WriteString("x ")
			} else {
				sb.WriteString("+ ")
			}
			sb.WriteString(fmt.Sprint(r[j].val))
			sb.WriteByte('\n')
		}
	}
	return []byte(sb.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go <binary>")
		os.Exit(1)
	}
	cand := os.Args[1]
	ref := "./refD.bin"
	if err := exec.Command("go", "build", "-o", ref, "2078D.go").Run(); err != nil {
		fmt.Println("failed to build reference solution:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rand.Seed(time.Now().UnixNano())

	for iter := 0; iter < 200; iter++ {
		input := buildInput()
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

		var t int
		if _, err := fmt.Fscan(strings.NewReader(string(input)), &t); err != nil {
			fmt.Println("failed to parse generated input:", err)
			os.Exit(1)
		}
		refAns, err := parseOutputs(refOut, t)
		if err != nil {
			fmt.Println("failed to parse reference output on iteration", iter+1, ":", err)
			fmt.Println("output:\n", refOut)
			fmt.Println("input:\n", string(input))
			os.Exit(1)
		}
		candAns, err := parseOutputs(candOut, t)
		if err != nil {
			fmt.Println("failed to parse candidate output on iteration", iter+1, ":", err)
			fmt.Println("output:\n", candOut)
			fmt.Println("input:\n", string(input))
			os.Exit(1)
		}

		match := len(refAns) == len(candAns)
		if match {
			for i := range refAns {
				if refAns[i] != candAns[i] {
					match = false
					break
				}
			}
		}
		if !match {
			fmt.Printf("mismatch on iteration %d\n", iter+1)
			fmt.Println("input:\n", string(input))
			fmt.Println("reference:\n", refOut)
			fmt.Println("candidate:\n", candOut)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
