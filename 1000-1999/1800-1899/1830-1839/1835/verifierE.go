package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type frac struct {
	num int64
	den int64
}

func makeFrac(num, den int64) frac {
	if den < 0 {
		num, den = -num, -den
	}
	if den == 0 {
		return frac{num: 1, den: 0}
	}
	g := gcd(abs64(num), den)
	return frac{num: num / g, den: den / g}
}

func addFrac(a, b frac) frac {
	return makeFrac(a.num*b.den+b.num*a.den, a.den*b.den)
}

func addInt(f frac, k int64) frac {
	return makeFrac(f.num+k*f.den, f.den)
}

func divInt(f frac, k int64) frac {
	return makeFrac(f.num, f.den*k)
}

func lessFrac(a, b frac) bool {
	return a.num*b.den < b.num*a.den
}

func abs64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

type state struct {
	pos    int
	extra  int
	mask   int
	hasBS  bool
}

// solver encapsulates the brute-force DP for small m, n.
type solver struct {
	m      int
	target []int
	memo   map[state]frac
}

func newSolver(m int, target []int) *solver {
	return &solver{
		m:      m,
		target: target,
		memo:   make(map[state]frac),
	}
}

func (s *solver) solve(pos, extra, mask int, hasBS bool) frac {
	if extra == 0 && pos == len(s.target) {
		return frac{0, 1}
	}
	st := state{pos: pos, extra: extra, mask: mask, hasBS: hasBS}
	if v, ok := s.memo[st]; ok {
		return v
	}
	kKnown := bits.OnesCount(uint(mask))
	uDigits := s.m - kKnown
	uBackspace := 0
	if !hasBS {
		uBackspace = 1
	}
	var options []frac

	unknownAction := func(curPos, curExtra, curMask int, curHasBS bool) frac {
		totalUnknown := uDigits + uBackspace
		if totalUnknown == 0 {
			return frac{num: 1, den: 0} // sentinel for unusable action
		}
		sum := frac{0, 1}
		// undiscovered digits
		for d := 0; d < s.m; d++ {
			if curMask&(1<<d) != 0 {
				continue
			}
			nPos, nExtra := curPos, curExtra
			if nExtra == 0 && nPos < len(s.target) && s.target[nPos] == d {
				nPos++
			} else {
				nExtra++
			}
			sum = addFrac(sum, s.solve(nPos, nExtra, curMask|(1<<d), curHasBS))
		}
		if !curHasBS {
			// discover backspace
			nPos, nExtra := curPos, curExtra
			if nExtra > 0 {
				nExtra--
			} else if nPos > 0 {
				nPos--
			}
			sum = addFrac(sum, s.solve(nPos, nExtra, curMask, true))
		}
		avg := divInt(sum, int64(totalUnknown))
		return addInt(avg, 1)
	}

	if extra > 0 {
		if hasBS {
			options = append(options, addInt(s.solve(pos, extra-1, mask, hasBS), 1))
		}
		if uDigits+uBackspace > 0 {
			options = append(options, unknownAction(pos, extra, mask, hasBS))
		}
	} else { // extra == 0
		if pos < len(s.target) {
			nextDigit := s.target[pos]
			if mask&(1<<nextDigit) != 0 {
				options = append(options, addInt(s.solve(pos+1, 0, mask, hasBS), 1))
			}
			if uDigits+uBackspace > 0 {
				options = append(options, unknownAction(pos, extra, mask, hasBS))
			}
		}
	}

	best := frac{num: 1, den: 0} // infinity sentinel
	for _, opt := range options {
		if opt.den == 0 { // skip invalid
			continue
		}
		if best.den == 0 || lessFrac(opt, best) {
			best = opt
		}
	}
	s.memo[st] = best
	return best
}

func expectedPresses(m int, target []int) frac {
	s := newSolver(m, target)
	return s.solve(0, 0, 0, false)
}

// modular helpers
const mod = int64(1_000_000_007)

func modPow(a, e int64) int64 {
	res := int64(1)
	for e > 0 {
		if e&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		e >>= 1
	}
	return res
}

func modInv(a int64) int64 {
	return modPow((a%mod+mod)%mod, mod-2)
}

func fracToMod(f frac) int64 {
	num := ((f.num%mod)+mod)%mod
	den := ((f.den%mod)+mod)%mod
	return num * modInv(den) % mod
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateCase(rng *rand.Rand) (string, string) {
	m := rng.Intn(5) + 2        // 2..6 to keep brute manageable
	n := rng.Intn(4) + 1        // 1..4
	num := make([]int, n)
	for i := 0; i < n; i++ {
		num[i] = rng.Intn(m)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for i, d := range num {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", d)
	}
	sb.WriteByte('\n')
	input := sb.String()
	expFrac := expectedPresses(m, num)
	expected := fmt.Sprintf("%d", fracToMod(expFrac))
	return input, expected
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 50; i++ {
		inp, exp := generateCase(rng)
		out, err := runCandidate(bin, inp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, inp)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", i+1, exp, out, inp)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
