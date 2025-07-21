package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

var elements = []string{
	"H", "He", "Li", "Be", "B", "C", "N", "O", "F", "Ne",
	"Na", "Mg", "Al", "Si", "P", "S", "Cl", "Ar", "K", "Ca",
	"Sc", "Ti", "V", "Cr", "Mn", "Fe", "Co", "Ni", "Cu", "Zn",
	"Ga", "Ge", "As", "Se", "Br", "Kr", "Rb", "Sr", "Y", "Zr",
	"Nb", "Mo", "Tc", "Ru", "Rh", "Pd", "Ag", "Cd", "In", "Sn",
	"Sb", "Te", "I", "Xe", "Cs", "Ba", "La", "Ce", "Pr", "Nd",
	"Pm", "Sm", "Eu", "Gd", "Tb", "Dy", "Ho", "Er", "Tm", "Yb",
	"Lu", "Hf", "Ta", "W", "Re", "Os", "Ir", "Pt", "Au", "Hg",
	"Tl", "Pb", "Bi", "Po", "At", "Rn", "Fr", "Ra", "Ac", "Th",
	"Pa", "U", "Np", "Pu", "Am", "Cm", "Bk", "Cf", "Es", "Fm",
}

func solveE(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var out bytes.Buffer
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return ""
	}
	seq := make([]int, n)
	tar := make([]int, m)
	elementMap := make(map[string]int, len(elements))
	for i, s := range elements {
		elementMap[s] = i + 1
	}
	var s string
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &s)
		seq[i] = elementMap[s]
	}
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &s)
		tar[i] = elementMap[s]
	}
	lim := (1 << n) - 1
	type data struct{ id, val int }
	f := make([]data, lim+1)
	link := make([]int, lim+1)
	for i := range link {
		link[i] = -1
	}
	for mask := 0; mask <= lim; mask++ {
		cur := f[mask]
		for j := 0; j < n; j++ {
			bit := 1 << j
			if mask&bit != 0 {
				continue
			}
			p := seq[j]
			id0, val0 := cur.id, cur.val
			rem := tar[id0] - val0
			if rem < p {
				continue
			}
			var nd data
			if rem == p {
				nd = data{id: id0 + 1, val: 0}
			} else {
				nd = data{id: id0, val: val0 + p}
			}
			nxt := mask | bit
			prev := f[nxt]
			if prev.id < nd.id || (prev.id == nd.id && prev.val < nd.val) {
				f[nxt] = nd
				link[nxt] = mask
			}
		}
	}
	end := f[lim]
	if end.id == m && end.val == 0 {
		fmt.Fprintln(&out, "YES")
		now := lim
		v := 0
		for now > 0 {
			prev := link[now]
			diff := now ^ prev
			for j := 0; j < n; j++ {
				if diff&(1<<j) != 0 {
					if v > 0 {
						fmt.Fprintf(&out, "+%s", elements[seq[j]-1])
					} else {
						fmt.Fprintf(&out, "%s", elements[seq[j]-1])
					}
					v++
				}
			}
			if f[prev].val == 0 {
				v = 0
				fmt.Fprintf(&out, "->%s\n", elements[tar[f[prev].id]-1])
			}
			now = prev
		}
	} else {
		fmt.Fprintln(&out, "NO")
	}
	return out.String()
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	k := rng.Intn(n) + 1
	start := make([]int, n)
	for i := range start {
		start[i] = rng.Intn(30) + 1
	}
	sum := 0
	for _, v := range start {
		sum += v
	}
	var target []int
	for {
		target = make([]int, k)
		rem := sum
		ok := true
		for i := 0; i < k-1; i++ {
			lo := 1
			hi := rem - (k - 1 - i)
			if hi > 100 {
				hi = 100
			}
			if lo > hi {
				ok = false
				break
			}
			target[i] = rng.Intn(hi-lo+1) + lo
			rem -= target[i]
		}
		if !ok {
			continue
		}
		target[k-1] = rem
		if target[k-1] >= 1 && target[k-1] <= 100 {
			break
		}
	}
	var in strings.Builder
	fmt.Fprintf(&in, "%d %d\n", n, k)
	for i, v := range start {
		if i > 0 {
			in.WriteByte(' ')
		}
		in.WriteString(elements[v-1])
	}
	in.WriteByte('\n')
	for i, v := range target {
		if i > 0 {
			in.WriteByte(' ')
		}
		in.WriteString(elements[v-1])
	}
	in.WriteByte('\n')
	expect := solveE(in.String())
	return in.String(), expect
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, buf.String())
	}
	if strings.TrimSpace(buf.String()) != strings.TrimSpace(expected) {
		return fmt.Errorf("expected\n%s\ngot\n%s", expected, buf.String())
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
