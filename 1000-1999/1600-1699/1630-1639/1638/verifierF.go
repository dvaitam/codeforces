package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type testCase struct {
	n int
	h []int64
}

// Embedded correct solver for 1638F
func solve1638F(tc testCase) string {
	n := tc.n
	h := make([]int64, n+2)
	for i := 1; i <= n; i++ {
		h[i] = tc.h[i-1]
	}

	L := make([]int, n+2)
	R := make([]int, n+2)
	for i := 1; i <= n; i++ {
		L[i] = i
		for L[i] > 1 && h[L[i]-1] >= h[i] {
			L[i] = L[L[i]-1]
		}
	}
	for i := n; i >= 1; i-- {
		R[i] = i
		for R[i] < n && h[R[i]+1] >= h[i] {
			R[i] = R[R[i]+1]
		}
	}

	var ans int64

	pref := make([]int64, n+2)
	for m := 1; m <= n; m++ {
		pref[m] = pref[m-1]
		for i := 1; i <= m; i++ {
			rightBound := R[i]
			if rightBound > m {
				rightBound = m
			}
			area := h[i] * int64(rightBound-L[i]+1)
			if area > pref[m] {
				pref[m] = area
			}
		}
	}

	suff := make([]int64, n+2)
	for m := n; m >= 1; m-- {
		suff[m] = suff[m+1]
		for i := m; i <= n; i++ {
			leftBound := L[i]
			if leftBound < m {
				leftBound = m
			}
			area := h[i] * int64(R[i]-leftBound+1)
			if area > suff[m] {
				suff[m] = area
			}
		}
	}

	for m := 1; m <= n-1; m++ {
		if pref[m]+suff[m+1] > ans {
			ans = pref[m] + suff[m+1]
		}
	}
	if pref[n] > ans {
		ans = pref[n]
	}

	H_left := make([]int64, n+2)
	H_right := make([]int64, n+2)
	C := make([]int64, 0, 2*n+5)

	for k := 1; k <= n; k++ {
		minH := h[k]
		for x := k; x >= 1; x-- {
			if h[x] < minH {
				minH = h[x]
			}
			H_left[x] = minH
		}

		minH = h[k]
		for y := k; y <= n; y++ {
			if h[y] < minH {
				minH = h[y]
			}
			H_right[y] = minH
		}

		C = C[:0]
		C = append(C, 0, h[k])
		for x := 1; x <= k; x++ {
			C = append(C, H_left[x])
		}
		for y := k; y <= n; y++ {
			C = append(C, h[k]-H_right[y])
		}

		sort.Slice(C, func(i, j int) bool { return C[i] < C[j] })

		leftPtr := 1
		rightPtr := k

		for i := 0; i < len(C); i++ {
			if i > 0 && C[i] == C[i-1] {
				continue
			}
			HA := C[i]
			if HA < 0 || HA > h[k] {
				continue
			}
			HB := h[k] - HA

			for leftPtr <= k && H_left[leftPtr] < HA {
				leftPtr++
			}
			for rightPtr+1 <= n && H_right[rightPtr+1] >= HB {
				rightPtr++
			}

			area := HA*int64(R[k]-leftPtr+1) + HB*int64(rightPtr-L[k]+1)
			if area > ans {
				ans = area
			}
		}
	}

	return fmt.Sprint(ans)
}

func generateCases() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]testCase, 100)
	for i := range cases {
		n := rng.Intn(50) + 1
		h := make([]int64, n)
		for j := range h {
			h[j] = rng.Int63n(1_000_000_000) + 1
		}
		cases[i] = testCase{n: n, h: h}
	}
	return cases
}

func runCase(bin string, tc testCase) (string, string, error) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", tc.n))
	for i, v := range tc.h {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	input := sb.String()

	expected := solve1638F(tc)

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out strings.Builder
	var stderr strings.Builder
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return expected, "", fmt.Errorf("runtime error: %v\nstderr: %s", err, stderr.String())
	}
	got := strings.TrimSpace(out.String())
	return expected, got, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases := generateCases()
	for i, c := range cases {
		expected, got, err := runCase(bin, c)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
