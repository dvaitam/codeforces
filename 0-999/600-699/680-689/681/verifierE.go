package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type segment struct{ l, r float64 }

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solve(x0, y0, v, T int64, circs [][3]int) string {
	R := float64(v) * float64(T)
	if R > 6e9 {
		R = 6e9
	}
	segs := make([]segment, 0, len(circs))
	inside := false
	for _, c := range circs {
		dx := float64(int64(c[0]) - x0)
		dy := float64(int64(c[1]) - y0)
		dist2 := dx*dx + dy*dy
		r := float64(c[2])
		if dist2 <= r*r {
			inside = true
			break
		}
		dist := math.Sqrt(dist2)
		if r == 0 || R+r <= dist {
			continue
		}
		angl := math.Atan2(dx, dy)
		al := math.Asin(r / dist)
		if R*R+r*r < dist*dist {
			cosv := (R*R + dist*dist - r*r) / (2 * R * dist)
			if cosv > 1 {
				cosv = 1
			} else if cosv < -1 {
				cosv = -1
			}
			al = math.Acos(cosv)
		}
		ang1 := angl - al
		ang2 := angl + al
		if ang1 < -math.Pi {
			segs = append(segs, segment{-math.Pi, ang2})
			segs = append(segs, segment{ang1 + 2*math.Pi, math.Pi})
		} else if ang2 > math.Pi {
			segs = append(segs, segment{ang1, math.Pi})
			segs = append(segs, segment{-math.Pi, ang2 - 2*math.Pi})
		} else {
			segs = append(segs, segment{ang1, ang2})
		}
	}
	if inside {
		return fmt.Sprintf("%.9f", 1.0)
	}
	if len(segs) == 0 {
		return fmt.Sprintf("%.9f", 0.0)
	}
	sort.Slice(segs, func(i, j int) bool {
		if segs[i].l != segs[j].l {
			return segs[i].l < segs[j].l
		}
		return segs[i].r < segs[j].r
	})
	var total, currLen, end float64
	end = -math.Pi
	for _, s := range segs {
		if end >= s.l {
			if end < s.r {
				currLen += s.r - end
				end = s.r
			}
		} else {
			total += currLen
			currLen = s.r - s.l
			end = s.r
		}
	}
	total += currLen
	p := total / (2 * math.Pi)
	return fmt.Sprintf("%.9f", p)
}

func generateCase(rng *rand.Rand) (string, string) {
	x0 := int64(rng.Intn(21) - 10)
	y0 := int64(rng.Intn(21) - 10)
	v := int64(rng.Intn(5) + 1)
	T := int64(rng.Intn(5) + 1)
	n := rng.Intn(3) + 1
	circs := make([][3]int, n)
	for i := 0; i < n; i++ {
		circs[i] = [3]int{rng.Intn(21) - 10, rng.Intn(21) - 10, rng.Intn(10) + 1}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d %d\n", x0, y0, v, T))
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for _, c := range circs {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", c[0], c[1], c[2]))
	}
	return sb.String(), solve(x0, y0, v, T, circs)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, exp := generateCase(rng)
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
