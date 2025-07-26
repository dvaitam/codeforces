package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func fft(a []complex128, invert bool) {
	n := len(a)
	j := 0
	for i := 1; i < n; i++ {
		bit := n >> 1
		for ; j&bit != 0; bit >>= 1 {
			j ^= bit
		}
		j ^= bit
		if i < j {
			a[i], a[j] = a[j], a[i]
		}
	}
	for length := 2; length <= n; length <<= 1 {
		ang := 2 * math.Pi / float64(length)
		if invert {
			ang = -ang
		}
		wlen := complex(math.Cos(ang), math.Sin(ang))
		for i := 0; i < n; i += length {
			w := complex(1, 0)
			for j := 0; j < length/2; j++ {
				u := a[i+j]
				v := a[i+j+length/2] * w
				a[i+j] = u + v
				a[i+j+length/2] = u - v
				w *= wlen
			}
		}
	}
	if invert {
		invN := complex(1/float64(n), 0)
		for i := range a {
			a[i] *= invN
		}
	}
}

func convolution(a, b []complex128) []complex128 {
	n := 1
	for n < len(a)+len(b) {
		n <<= 1
	}
	fa := make([]complex128, n)
	fb := make([]complex128, n)
	copy(fa, a)
	copy(fb, b)
	fft(fa, false)
	fft(fb, false)
	for i := 0; i < n; i++ {
		fa[i] *= fb[i]
	}
	fft(fa, true)
	return fa
}

func bluestein(a []complex128) []complex128 {
	n := len(a)
	m := 1
	for m < 2*n-1 {
		m <<= 1
	}
	A := make([]complex128, m)
	B := make([]complex128, m)
	wBase := math.Pi / float64(n)
	for j := 0; j < n; j++ {
		angle := -wBase * float64(j*j)
		c := complex(math.Cos(angle), math.Sin(angle))
		A[j] = a[j] * c
	}
	B[n-1] = 1
	for j := 1; j < n; j++ {
		angle := wBase * float64(j*j)
		c := complex(math.Cos(angle), math.Sin(angle))
		B[n-1+j] = c
		B[n-1-j] = c
	}
	conv := convolution(A, B)
	res := make([]complex128, n)
	for k := 0; k < n; k++ {
		angle := -wBase * float64(k*k)
		c := complex(math.Cos(angle), math.Sin(angle))
		res[k] = conv[n-1+k] * c
	}
	return res
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func cmplxAbs(z complex128) float64 {
	return math.Hypot(real(z), imag(z))
}

func solveG(n int, s string) string {
	a := make([]complex128, n)
	for i := 0; i < n; i++ {
		a[i] = complex(float64(s[i]-'0'), 0)
	}
	res := bluestein(a)
	eps := 1e-6 * float64(n)
	for k := 1; k < n; k++ {
		if gcd(k, n) == 1 {
			if cmplxAbs(res[k]) > eps {
				return "NO"
			}
		}
	}
	return "YES"
}

type CaseG struct {
	input    string
	expected string
}

func generateCaseG(rng *rand.Rand) CaseG {
	n := rng.Intn(28) + 3
	b := make([]byte, n)
	for i := range b {
		b[i] = byte('0' + rng.Intn(10))
	}
	s := string(b)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n%s\n", n, s))
	return CaseG{sb.String(), solveG(n, s)}
}

func runCase(exe, input, expected string) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(exe, ".go") {
		cmd = exec.Command("go", "run", exe)
	} else {
		cmd = exec.Command(exe)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	cases := []CaseG{{"3\n000\n", "YES"}}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCaseG(rng))
	}
	for i, tc := range cases {
		if err := runCase(exe, tc.input, tc.expected); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
