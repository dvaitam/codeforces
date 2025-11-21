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

type tc struct {
	n, m, p, q int
}

func run(bin string, input []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func parse(out string, t int) ([]string, error) {
	reader := strings.NewReader(out)
	res := make([]string, 0, t)
	for i := 0; i < t; i++ {
		var tok string
		if _, err := fmt.Fscan(reader, &tok); err != nil {
			return nil, fmt.Errorf("output ended early on case %d: %v", i+1, err)
		}
		up := strings.ToUpper(tok)
		if up != "YES" && up != "NO" {
			return nil, fmt.Errorf("invalid verdict '%s' on case %d", tok, i+1)
		}
		res = append(res, up)
	}
	var extra string
	if _, err := fmt.Fscan(reader, &extra); err == nil {
		return nil, fmt.Errorf("extra output after %d cases", t)
	}
	return res, nil
}

func genCase() tc {
	mode := rand.Intn(6)
	switch mode {
	case 0:
		// n divisible by p with mismatch sum -> NO
		p := rand.Intn(10) + 1
		t := rand.Intn(9) + 1
		n := p * t
		q := rand.Intn(100) + 1
		m := q*t + rand.Intn(5) + 1
		return tc{n, m, p, q}
	case 1:
		// n divisible by p with exact sum -> YES
		p := rand.Intn(10) + 1
		t := rand.Intn(9) + 1
		n := p * t
		q := rand.Intn(100) + 1
		m := q * t
		return tc{n, m, p, q}
	case 2:
		// n not divisible by p -> always YES
		n := rand.Intn(99) + 2
		p := rand.Intn(n-1) + 1
		if n%p == 0 {
			p = p%n + 1
		}
		m := rand.Intn(100) + 1
		q := rand.Intn(100) + 1
		return tc{n, m, p, q}
	case 3:
		// minimal boundaries
		return tc{1, 1, 1, 1}
	case 4:
		// large bounds
		n := 100
		p := rand.Intn(100) + 1
		m := rand.Intn(100) + 1
		q := rand.Intn(100) + 1
		if n%p == 0 && rand.Intn(2) == 0 {
			m = (n / p) * q
		}
		return tc{n, m, p, q}
	default:
		// random
		n := rand.Intn(100) + 1
		p := rand.Intn(n) + 1
		m := rand.Intn(100) + 1
		q := rand.Intn(100) + 1
		return tc{n, m, p, q}
	}
}

func buildInput() ([]byte, []tc) {
	t := rand.Intn(20) + 1
	var sb strings.Builder
	cases := make([]tc, 0, t)
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for i := 0; i < t; i++ {
		c := genCase()
		cases = append(cases, c)
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", c.n, c.m, c.p, c.q))
	}
	return []byte(sb.String()), cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go <binary>")
		os.Exit(1)
	}
	cand := os.Args[1]
	ref := "./refA.bin"
	if err := exec.Command("go", "build", "-o", ref, "2102A.go").Run(); err != nil {
		fmt.Println("failed to build reference solution:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rand.Seed(time.Now().UnixNano())

	for iter := 0; iter < 200; iter++ {
		input, cases := buildInput()
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

		refAns, err := parse(refOut, len(cases))
		if err != nil {
			fmt.Println("failed to parse reference output on iteration", iter+1, ":", err)
			fmt.Println("output:\n", refOut)
			fmt.Println("input:\n", string(input))
			os.Exit(1)
		}
		candAns, err := parse(candOut, len(cases))
		if err != nil {
			fmt.Println("failed to parse candidate output on iteration", iter+1, ":", err)
			fmt.Println("output:\n", candOut)
			fmt.Println("input:\n", string(input))
			os.Exit(1)
		}

		if len(refAns) != len(candAns) {
			fmt.Printf("answer count mismatch on iteration %d\n", iter+1)
			fmt.Println("input:\n", string(input))
			fmt.Println("reference:\n", refOut)
			fmt.Println("candidate:\n", candOut)
			os.Exit(1)
		}
		for i := range refAns {
			if refAns[i] != candAns[i] {
				fmt.Printf("wrong answer on iteration %d case %d\n", iter+1, i+1)
				fmt.Println("input case:", cases[i])
				fmt.Println("reference:", refAns[i])
				fmt.Println("candidate:", candAns[i])
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed.")
}
