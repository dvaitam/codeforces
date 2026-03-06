package main

import (
   "bufio"
   "fmt"
   "os"
)

// BIT for 1-indexed array
type BIT struct {
   n    int
   tree []int
}

func NewBIT(n int) *BIT {
   return &BIT{n: n, tree: make([]int, n+1)}
}

// add v at position i
func (b *BIT) Add(i, v int) {
   for x := i; x <= b.n; x += x & -x {
       b.tree[x] += v
   }
}

// Sum returns prefix sum [1..i]
func (b *BIT) Sum(i int) int {
   s := 0
   for x := i; x > 0; x -= x & -x {
       s += b.tree[x]
   }
   return s
}

// Kth returns smallest i such that Sum(i) >= k, assumes k>=1 and k<=Sum(n)
func (b *BIT) Kth(k int) int {
   pos, acc := 0, 0
   // compute max power of two <= n
   maxPow := 1
   for maxPow<<1 <= b.n {
       maxPow <<= 1
   }
   for pw := maxPow; pw > 0; pw >>= 1 {
       np := pos + pw
       if np <= b.n && acc+b.tree[np] < k {
           acc += b.tree[np]
           pos = np
       }
   }
   return pos + 1
}

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n, q int
   fmt.Fscan(in, &n, &q)
   sbytes := make([]byte, n+1)
   var tmp string
   fmt.Fscan(in, &tmp)
   for i := 1; i <= n; i++ {
       sbytes[i] = tmp[i-1]
   }
   // map types: R=0, S=1, P=2
   idx := map[byte]int{'R': 0, 'S': 1, 'P': 2}
   win := [3]int{1, 2, 0}
   lose := [3]int{2, 0, 1}
   bits := make([]*BIT, 3)
   cnt := [3]int{}
   for t := 0; t < 3; t++ {
       bits[t] = NewBIT(n)
   }
   for i := 1; i <= n; i++ {
       t := idx[sbytes[i]]
       bits[t].Add(i, 1)
       cnt[t]++
   }
	// function to compute answer
	calc := func() int {
		res := 0
		for t := 0; t < 3; t++ {
			y := lose[t] // type that beats t (nemesis)
			z := win[t]  // type that t beats (helper that also beats y)
			if cnt[y] == 0 {
				res += cnt[t]
			} else if cnt[z] == 0 {
				// no helper to eliminate nemesis
			} else {
				total := bits[t].Sum(n)
				fy := bits[y].Kth(1)
				fz := bits[z].Kth(1)
				if fz > fy {
					total -= bits[t].Sum(fz) - bits[t].Sum(fy)
				}
				ly := bits[y].Kth(cnt[y])
				lz := bits[z].Kth(cnt[z])
				if ly > lz {
					total -= bits[t].Sum(ly) - bits[t].Sum(lz)
				}
				res += total
			}
		}
		return res
	}
   // initial
   fmt.Fprintln(out, calc())
	for i := 0; i < q; i++ {
		var p int
		var ch string
		fmt.Fscan(in, &p, &ch)
		c := ch[0]
		old := idx[sbytes[p]]
		if c == sbytes[p] {
			fmt.Fprintln(out, calc())
			continue
		}
		// remove old
		bits[old].Add(p, -1)
		cnt[old]--
		// add new
		sbytes[p] = c
		t := idx[c]
		bits[t].Add(p, 1)
		cnt[t]++
		fmt.Fprintln(out, calc())
	}
}
