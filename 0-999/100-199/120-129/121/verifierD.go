package main

import (
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

const solution121DSource = `package main

import (
   "bufio"
   "fmt"
   "math/big"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   var k int64
   fmt.Fscan(reader, &n, &k)
   li := make([]int64, n)
   ri := make([]int64, n)
   Lmin := int64(1<<62)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &li[i], &ri[i])
       if ri[i]-li[i] < Lmin {
           Lmin = ri[i] - li[i]
       }
   }
   // generate lucky numbers up to 1e18
   luckies := genLuckies()
   sort.Slice(luckies, func(i, j int) bool { return luckies[i] < luckies[j] })

   // sort segment endpoints
   lis := make([]int64, n)
   ris := make([]int64, n)
   copy(lis, li)
   copy(ris, ri)
   sort.Slice(lis, func(i, j int) bool { return lis[i] < lis[j] })
   sort.Slice(ris, func(i, j int) bool { return ris[i] < ris[j] })

   // two pointers for segments
   p1, p2 := 0, 0
   c1, c2 := 0, 0
   sumRi := big.NewInt(0)
   sumLi := big.NewInt(0)
   // initial X
   X := luckies[0]
   // initial p2 and sumLi for li > X
   for p2 < n && lis[p2] <= X {
       p2++
   }
   c2 = n - p2
   for i := p2; i < n; i++ {
       sumLi.Add(sumLi, big.NewInt(lis[i]))
   }

   best := 0
   l := 0
   bigX := big.NewInt(X)
   bigY := big.NewInt(0)
   cost1 := big.NewInt(0)
   cost2 := big.NewInt(0)
   tmp := big.NewInt(0)
   sumCost := big.NewInt(0)
   bigK := big.NewInt(k)

   for r := 0; r < len(luckies); r++ {
       Y := luckies[r]
       bigY.SetInt64(Y)
       // include new S1 segments with ri < Y
       for p1 < n && ris[p1] < Y {
           sumRi.Add(sumRi, big.NewInt(ris[p1]))
           p1++
           c1++
       }
       // adjust l to satisfy window length and cost
       for l <= r {
           // current X and Y
           X = luckies[l]
           // window length check
           if Y-X > Lmin {
               // move l
           } else {
               // compute cost1 = c1*Y - sumRi
               tmp.SetInt64(int64(c1))
               tmp.Mul(tmp, bigY)
               cost1.Sub(tmp, sumRi)
               // compute cost2 = sumLi - c2*X
               tmp.SetInt64(int64(c2))
               tmp.Mul(tmp, bigX)
               cost2.Sub(sumLi, tmp)
               // total cost
               sumCost.Add(cost1, cost2)
               if sumCost.Cmp(bigK) <= 0 {
                   break
               }
           }
           // increment l
           oldX := X
           l++
           if l >= len(luckies) {
               break
           }
           X = luckies[l]
           bigX.SetInt64(X)
           // remove from S2 those with li <= X
           for p2 < n && lis[p2] <= X {
               sumLi.Sub(sumLi, big.NewInt(lis[p2]))
               p2++
               c2--
           }
           // adjust S1 and counts when moving X forward
           for p1 > 0 && ris[p1-1] < oldX {
               p1--
               c1--
               sumRi.Sub(sumRi, big.NewInt(ris[p1]))
           }
       }
       // update best
       if l <= r {
           cnt := r - l + 1
           if cnt > best {
               best = cnt
           }
       }
   }
   fmt.Fprint(writer, best)
}

// generate all lucky numbers up to 1e18
func genLuckies() []int64 {
   var res []int64
   queue := []int64{4, 7}
   for i := 0; i < len(queue); i++ {
       v := queue[i]
       res = append(res, v)
       if v <= 1000000000000000000/10 {
           queue = append(queue, v*10+4, v*10+7)
       }
   }
   return res
}
`

// Keep the embedded reference solution reachable so it is preserved in the binary.
var _ = solution121DSource

type segment struct {
	l int
	r int
}

type testCase struct {
	n        int
	k        int64
	segments []segment
}

var testcases = []testCase{
	{n: 2, k: 75, segments: []segment{{l: 70, r: 72}, {l: 78, r: 80}}},
	{n: 5, k: 1, segments: []segment{{l: 61, r: 69}, {l: 30, r: 35}, {l: 61, r: 78}, {l: 61, r: 71}, {l: 20, r: 22}}},
	{n: 5, k: 49, segments: []segment{{l: 95, r: 95}, {l: 21, r: 22}, {l: 39, r: 39}, {l: 61, r: 73}, {l: 92, r: 98}}},
	{n: 5, k: 56, segments: []segment{{l: 18, r: 19}, {l: 5, r: 8}, {l: 28, r: 34}, {l: 100, r: 109}, {l: 54, r: 66}}},
	{n: 5, k: 44, segments: []segment{{l: 69, r: 82}, {l: 75, r: 80}, {l: 88, r: 88}, {l: 78, r: 83}, {l: 42, r: 45}}},
	{n: 2, k: 81, segments: []segment{{l: 74, r: 78}, {l: 16, r: 17}}},
	{n: 4, k: 11, segments: []segment{{l: 45, r: 46}, {l: 20, r: 20}, {l: 55, r: 68}, {l: 16, r: 16}}},
	{n: 4, k: 91, segments: []segment{{l: 76, r: 84}, {l: 36, r: 43}, {l: 5, r: 5}, {l: 10, r: 10}}},
	{n: 2, k: 52, segments: []segment{{l: 38, r: 46}, {l: 20, r: 21}}},
	{n: 3, k: 46, segments: []segment{{l: 18, r: 24}, {l: 59, r: 71}, {l: 83, r: 100}}},
	{n: 1, k: 79, segments: []segment{{l: 65, r: 71}}},
	{n: 2, k: 38, segments: []segment{{l: 56, r: 64}, {l: 39, r: 49}}},
	{n: 1, k: 100, segments: []segment{{l: 54, r: 64}}},
	{n: 1, k: 48, segments: []segment{{l: 79, r: 83}}},
	{n: 1, k: 81, segments: []segment{{l: 81, r: 88}}},
	{n: 3, k: 86, segments: []segment{{l: 46, r: 54}, {l: 95, r: 95}, {l: 76, r: 76}}},
	{n: 3, k: 32, segments: []segment{{l: 81, r: 85}, {l: 76, r: 86}, {l: 23, r: 25}}},
	{n: 3, k: 97, segments: []segment{{l: 48, r: 56}, {l: 39, r: 40}, {l: 99, r: 99}}},
	{n: 3, k: 64, segments: []segment{{l: 29, r: 37}, {l: 31, r: 33}, {l: 87, r: 97}}},
	{n: 1, k: 13, segments: []segment{{l: 77, r: 82}}},
	{n: 2, k: 56, segments: []segment{{l: 22, r: 23}, {l: 95, r: 101}}},
	{n: 5, k: 57, segments: []segment{{l: 35, r: 36}, {l: 5, r: 11}, {l: 41, r: 46}, {l: 36, r: 46}, {l: 11, r: 22}}},
	{n: 5, k: 16, segments: []segment{{l: 54, r: 62}, {l: 35, r: 40}, {l: 82, r: 86}, {l: 54, r: 67}, {l: 5, r: 7}}},
	{n: 2, k: 0, segments: []segment{{l: 62, r: 78}, {l: 56, r: 63}}},
	{n: 1, k: 95, segments: []segment{{l: 59, r: 68}}},
	{n: 5, k: 43, segments: []segment{{l: 30, r: 32}, {l: 37, r: 38}, {l: 6, r: 6}, {l: 56, r: 57}, {l: 2, r: 5}}},
	{n: 2, k: 64, segments: []segment{{l: 39, r: 39}, {l: 68, r: 81}}},
	{n: 1, k: 78, segments: []segment{{l: 15, r: 17}}},
	{n: 3, k: 69, segments: []segment{{l: 62, r: 63}, {l: 29, r: 29}, {l: 69, r: 70}}},
	{n: 2, k: 35, segments: []segment{{l: 17, r: 17}, {l: 81, r: 93}}},
	{n: 1, k: 96, segments: []segment{{l: 35, r: 39}}},
	{n: 5, k: 67, segments: []segment{{l: 67, r: 67}, {l: 61, r: 61}, {l: 8, r: 8}, {l: 16, r: 16}, {l: 62, r: 62}}},
	{n: 5, k: 64, segments: []segment{{l: 63, r: 65}, {l: 41, r: 42}, {l: 50, r: 62}, {l: 76, r: 81}, {l: 34, r: 36}}},
	{n: 4, k: 15, segments: []segment{{l: 17, r: 17}, {l: 92, r: 104}, {l: 11, r: 16}, {l: 6, r: 13}}},
	{n: 5, k: 83, segments: []segment{{l: 70, r: 80}, {l: 6, r: 19}, {l: 7, r: 17}, {l: 64, r: 70}, {l: 89, r: 96}}},
	{n: 1, k: 31, segments: []segment{{l: 28, r: 36}}},
	{n: 5, k: 9, segments: []segment{{l: 55, r: 61}, {l: 17, r: 17}, {l: 48, r: 56}, {l: 16, r: 27}, {l: 16, r: 28}}},
	{n: 1, k: 93, segments: []segment{{l: 41, r: 58}}},
	{n: 1, k: 75, segments: []segment{{l: 92, r: 92}}},
	{n: 2, k: 30, segments: []segment{{l: 100, r: 100}, {l: 68, r: 70}}},
	{n: 1, k: 84, segments: []segment{{l: 49, r: 49}}},
	{n: 3, k: 15, segments: []segment{{l: 4, r: 7}, {l: 90, r: 99}, {l: 39, r: 39}}},
	{n: 5, k: 65, segments: []segment{{l: 68, r: 69}, {l: 71, r: 71}, {l: 71, r: 80}, {l: 24, r: 24}, {l: 24, r: 31}}},
	{n: 4, k: 78, segments: []segment{{l: 90, r: 94}, {l: 48, r: 60}, {l: 45, r: 58}, {l: 11, r: 19}}},
	{n: 2, k: 52, segments: []segment{{l: 96, r: 99}, {l: 89, r: 107}}},
	{n: 5, k: 87, segments: []segment{{l: 62, r: 65}, {l: 20, r: 20}, {l: 64, r: 78}, {l: 76, r: 77}, {l: 35, r: 36}}},
	{n: 5, k: 65, segments: []segment{{l: 41, r: 45}, {l: 86, r: 95}, {l: 75, r: 83}, {l: 28, r: 28}, {l: 35, r: 47}}},
	{n: 2, k: 22, segments: []segment{{l: 73, r: 76}, {l: 42, r: 46}}},
	{n: 4, k: 89, segments: []segment{{l: 62, r: 68}, {l: 60, r: 77}, {l: 4, r: 6}, {l: 52, r: 53}}},
	{n: 2, k: 30, segments: []segment{{l: 83, r: 83}, {l: 33, r: 36}}},
	{n: 3, k: 17, segments: []segment{{l: 24, r: 25}, {l: 33, r: 33}, {l: 41, r: 44}}},
	{n: 1, k: 93, segments: []segment{{l: 11, r: 11}}},
	{n: 3, k: 37, segments: []segment{{l: 5, r: 12}, {l: 75, r: 75}, {l: 4, r: 9}}},
	{n: 4, k: 48, segments: []segment{{l: 63, r: 63}, {l: 83, r: 98}, {l: 51, r: 55}, {l: 41, r: 43}}},
	{n: 1, k: 85, segments: []segment{{l: 56, r: 59}}},
	{n: 5, k: 32, segments: []segment{{l: 13, r: 24}, {l: 87, r: 94}, {l: 38, r: 46}, {l: 14, r: 24}, {l: 73, r: 89}}},
	{n: 1, k: 85, segments: []segment{{l: 64, r: 75}}},
	{n: 1, k: 91, segments: []segment{{l: 38, r: 43}}},
	{n: 2, k: 22, segments: []segment{{l: 48, r: 62}, {l: 16, r: 17}}},
	{n: 3, k: 82, segments: []segment{{l: 93, r: 112}, {l: 54, r: 63}, {l: 83, r: 86}}},
	{n: 4, k: 39, segments: []segment{{l: 23, r: 23}, {l: 92, r: 96}, {l: 70, r: 82}, {l: 46, r: 48}}},
	{n: 3, k: 49, segments: []segment{{l: 7, r: 7}, {l: 62, r: 70}, {l: 32, r: 43}}},
	{n: 3, k: 51, segments: []segment{{l: 58, r: 60}, {l: 46, r: 49}, {l: 20, r: 21}}},
	{n: 1, k: 72, segments: []segment{{l: 100, r: 101}}},
	{n: 2, k: 72, segments: []segment{{l: 54, r: 65}, {l: 17, r: 21}}},
	{n: 4, k: 24, segments: []segment{{l: 70, r: 75}, {l: 73, r: 74}, {l: 33, r: 37}, {l: 4, r: 18}}},
	{n: 4, k: 49, segments: []segment{{l: 41, r: 50}, {l: 82, r: 91}, {l: 86, r: 86}, {l: 77, r: 82}}},
	{n: 1, k: 13, segments: []segment{{l: 99, r: 106}}},
	{n: 2, k: 67, segments: []segment{{l: 81, r: 84}, {l: 25, r: 31}}},
	{n: 1, k: 64, segments: []segment{{l: 83, r: 84}}},
	{n: 5, k: 36, segments: []segment{{l: 85, r: 86}, {l: 60, r: 62}, {l: 7, r: 7}, {l: 80, r: 81}, {l: 64, r: 64}}},
	{n: 3, k: 41, segments: []segment{{l: 43, r: 54}, {l: 89, r: 89}, {l: 77, r: 77}}},
	{n: 3, k: 26, segments: []segment{{l: 9, r: 12}, {l: 90, r: 97}, {l: 41, r: 41}}},
	{n: 4, k: 9, segments: []segment{{l: 26, r: 29}, {l: 64, r: 66}, {l: 69, r: 72}, {l: 84, r: 93}}},
	{n: 1, k: 59, segments: []segment{{l: 59, r: 66}}},
	{n: 2, k: 58, segments: []segment{{l: 5, r: 10}, {l: 48, r: 56}}},
	{n: 3, k: 76, segments: []segment{{l: 52, r: 52}, {l: 27, r: 32}, {l: 19, r: 27}}},
	{n: 2, k: 20, segments: []segment{{l: 27, r: 27}, {l: 75, r: 83}}},
	{n: 2, k: 81, segments: []segment{{l: 4, r: 4}, {l: 78, r: 81}}},
	{n: 4, k: 23, segments: []segment{{l: 8, r: 8}, {l: 58, r: 64}, {l: 5, r: 5}, {l: 52, r: 53}}},
	{n: 4, k: 3, segments: []segment{{l: 29, r: 30}, {l: 50, r: 56}, {l: 22, r: 31}, {l: 15, r: 16}}},
	{n: 5, k: 6, segments: []segment{{l: 94, r: 98}, {l: 60, r: 67}, {l: 32, r: 40}, {l: 4, r: 14}, {l: 45, r: 46}}},
	{n: 1, k: 87, segments: []segment{{l: 56, r: 58}}},
	{n: 5, k: 0, segments: []segment{{l: 14, r: 14}, {l: 3, r: 7}, {l: 5, r: 6}, {l: 25, r: 41}, {l: 43, r: 49}}},
	{n: 4, k: 43, segments: []segment{{l: 62, r: 72}, {l: 5, r: 9}, {l: 97, r: 109}, {l: 12, r: 14}}},
	{n: 4, k: 14, segments: []segment{{l: 65, r: 73}, {l: 43, r: 55}, {l: 23, r: 31}, {l: 46, r: 48}}},
	{n: 4, k: 56, segments: []segment{{l: 30, r: 42}, {l: 90, r: 101}, {l: 35, r: 39}, {l: 93, r: 105}}},
	{n: 4, k: 5, segments: []segment{{l: 20, r: 25}, {l: 97, r: 97}, {l: 12, r: 14}, {l: 31, r: 32}}},
	{n: 5, k: 6, segments: []segment{{l: 58, r: 72}, {l: 94, r: 104}, {l: 48, r: 48}, {l: 25, r: 37}, {l: 14, r: 23}}},
	{n: 3, k: 14, segments: []segment{{l: 58, r: 60}, {l: 27, r: 27}, {l: 20, r: 24}}},
	{n: 5, k: 1, segments: []segment{{l: 15, r: 19}, {l: 27, r: 33}, {l: 65, r: 75}, {l: 69, r: 72}, {l: 23, r: 25}}},
	{n: 1, k: 14, segments: []segment{{l: 77, r: 77}}},
	{n: 2, k: 32, segments: []segment{{l: 11, r: 14}, {l: 52, r: 53}}},
	{n: 4, k: 97, segments: []segment{{l: 86, r: 92}, {l: 78, r: 90}, {l: 15, r: 24}, {l: 57, r: 60}}},
	{n: 1, k: 69, segments: []segment{{l: 1, r: 5}}},
	{n: 1, k: 43, segments: []segment{{l: 45, r: 48}}},
	{n: 1, k: 70, segments: []segment{{l: 88, r: 94}}},
	{n: 1, k: 77, segments: []segment{{l: 67, r: 68}}},
	{n: 3, k: 7, segments: []segment{{l: 43, r: 49}, {l: 57, r: 58}, {l: 28, r: 30}}},
	{n: 2, k: 92, segments: []segment{{l: 28, r: 42}, {l: 92, r: 105}}},
}

func genLuckies() []int64 {
	var res []int64
	queue := []int64{4, 7}
	for i := 0; i < len(queue); i++ {
		v := queue[i]
		res = append(res, v)
		if v <= 1000000000000000000/10 {
			queue = append(queue, v*10+4, v*10+7)
		}
	}
	return res
}

func solveExpected(tc testCase) int {
	n := tc.n
	k := tc.k
	li := make([]int64, n)
	ri := make([]int64, n)
	Lmin := int64(1 << 62)
	for i, s := range tc.segments {
		li[i] = int64(s.l)
		ri[i] = int64(s.r)
		if ri[i]-li[i] < Lmin {
			Lmin = ri[i] - li[i]
		}
	}

	luckies := genLuckies()
	sort.Slice(luckies, func(i, j int) bool { return luckies[i] < luckies[j] })

	lis := append([]int64(nil), li...)
	ris := append([]int64(nil), ri...)
	sort.Slice(lis, func(i, j int) bool { return lis[i] < lis[j] })
	sort.Slice(ris, func(i, j int) bool { return ris[i] < ris[j] })

	p1, p2 := 0, 0
	c1, c2 := 0, 0
	sumRi := big.NewInt(0)
	sumLi := big.NewInt(0)

	X := luckies[0]
	for p2 < n && lis[p2] <= X {
		p2++
	}
	c2 = n - p2
	for i := p2; i < n; i++ {
		sumLi.Add(sumLi, big.NewInt(lis[i]))
	}

	best := 0
	l := 0
	bigX := big.NewInt(X)
	bigY := big.NewInt(0)
	cost1 := big.NewInt(0)
	cost2 := big.NewInt(0)
	tmp := big.NewInt(0)
	sumCost := big.NewInt(0)
	bigK := big.NewInt(k)

	for r := 0; r < len(luckies); r++ {
		Y := luckies[r]
		bigY.SetInt64(Y)
		for p1 < n && ris[p1] < Y {
			sumRi.Add(sumRi, big.NewInt(ris[p1]))
			p1++
			c1++
		}
		for l <= r {
			X = luckies[l]
			if Y-X > Lmin {
				// move l forward
			} else {
				tmp.SetInt64(int64(c1))
				tmp.Mul(tmp, bigY)
				cost1.Sub(tmp, sumRi)

				tmp.SetInt64(int64(c2))
				tmp.Mul(tmp, bigX)
				cost2.Sub(sumLi, tmp)

				sumCost.Add(cost1, cost2)
				if sumCost.Cmp(bigK) <= 0 {
					break
				}
			}
			l++
			if l >= len(luckies) {
				break
			}
			X = luckies[l]
			bigX.SetInt64(X)
			for p2 < n && lis[p2] <= X {
				sumLi.Sub(sumLi, big.NewInt(lis[p2]))
				p2++
				c2--
			}
		}
		if l <= r {
			if r-l+1 > best {
				best = r - l + 1
			}
		}
	}
	return best
}

func runCase(bin string, idx int, tc testCase) error {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.k)
	for _, s := range tc.segments {
		fmt.Fprintf(&sb, "%d %d\n", s.l, s.r)
	}
	input := sb.String()
	expect := solveExpected(tc)

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("case %d failed: %v\nstderr: %s\ninput:\n%s", idx, err, string(out), input)
	}
	gotStr := strings.TrimSpace(string(out))
	got, err := strconv.Atoi(gotStr)
	if err != nil {
		return fmt.Errorf("case %d failed: invalid output %q\ninput:\n%s", idx, gotStr, input)
	}
	if got != expect {
		return fmt.Errorf("case %d failed: expected %d got %d\ninput:\n%s", idx, expect, got, input)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for i, tc := range testcases {
		if err := runCase(bin, i+1, tc); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
