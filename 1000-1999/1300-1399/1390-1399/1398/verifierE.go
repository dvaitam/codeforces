package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type spell struct {
	tp int
	d  int64
}

func maxDamage(fire []int64, light []int64) int64 {
	spells := make([]spell, 0, len(fire)+len(light))
	for _, v := range fire {
		spells = append(spells, spell{0, v})
	}
	for _, v := range light {
		spells = append(spells, spell{1, v})
	}
	n := len(spells)
	used := make([]bool, n)
	var best int64
	var dfs func(int, int64, int64)
	dfs = func(k int, mult int64, dmg int64) {
		if k == n {
			if dmg > best {
				best = dmg
			}
			return
		}
		for i := 0; i < n; i++ {
			if !used[i] {
				used[i] = true
				sp := spells[i]
				if sp.tp == 0 {
					dfs(k+1, mult, dmg+mult*sp.d)
				} else {
					dfs(k+1, mult*2, dmg+mult*sp.d)
				}
				used[i] = false
			}
		}
	}
	dfs(0, 1, 0)
	return best
}

func solveCase(input string) string {
	sc := bufio.NewScanner(strings.NewReader(input))
	sc.Split(bufio.ScanWords)
	sc.Scan()
	n, _ := strconv.Atoi(sc.Text())
	fire := make(map[int64]bool)
	light := make(map[int64]bool)
	var outputs []string
	for i := 0; i < n; i++ {
		sc.Scan()
		tp, _ := strconv.Atoi(sc.Text())
		sc.Scan()
		d64, _ := strconv.ParseInt(sc.Text(), 10, 64)
		if d64 > 0 {
			if tp == 0 {
				fire[d64] = true
			} else {
				light[d64] = true
			}
		} else {
			if tp == 0 {
				delete(fire, -d64)
			} else {
				delete(light, -d64)
			}
		}
		f := make([]int64, 0, len(fire))
		l := make([]int64, 0, len(light))
		for v := range fire {
			f = append(f, v)
		}
		for v := range light {
			l = append(l, v)
		}
		outputs = append(outputs, fmt.Sprint(maxDamage(f, l)))
	}
	return strings.Join(outputs, "\n")
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(5) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintln(n))
	fire := make(map[int64]bool)
	light := make(map[int64]bool)
	for i := 0; i < n; i++ {
		tp := rng.Intn(2)
		add := rng.Intn(2) == 0 || len(fire)+len(light) == 0
		var d int64
		if add {
			// ensure unique
			for {
				d = int64(rng.Intn(9) + 1)
				if tp == 0 {
					if !fire[d] {
						break
					}
				} else {
					if !light[d] {
						break
					}
				}
			}
			if tp == 0 {
				fire[d] = true
			} else {
				light[d] = true
			}
			sb.WriteString(fmt.Sprintf("%d %d\n", tp, d))
		} else {
			// remove existing
			if tp == 0 && len(fire) > 0 {
				for v := range fire {
					d = v
					break
				}
				delete(fire, d)
				sb.WriteString(fmt.Sprintf("%d %d\n", tp, -d))
			} else if tp == 1 && len(light) > 0 {
				for v := range light {
					d = v
					break
				}
				delete(light, d)
				sb.WriteString(fmt.Sprintf("%d %d\n", tp, -d))
			} else {
				// if chosen type empty, add instead
				d = int64(rng.Intn(9) + 1)
				if tp == 0 {
					fire[d] = true
				} else {
					light[d] = true
				}
				sb.WriteString(fmt.Sprintf("%d %d\n", tp, d))
			}
		}
	}
	return sb.String()
}

func runCase(bin, input, expected string) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expected = strings.TrimSpace(expected)
	if got != expected {
		return fmt.Errorf("expected %q got %q", expected, got)
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
		in := generateCase(rng)
		exp := solveCase(in)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
