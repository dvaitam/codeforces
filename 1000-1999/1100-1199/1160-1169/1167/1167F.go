package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

const MOD = 1000000007

// Fenwick implements a Fenwick tree (BIT) for prefix sums modulo MOD.
type Fenwick struct {
   n int
   t []int64
}

// NewFenwick creates a Fenwick tree of size n.
func NewFenwick(n int) *Fenwick {
   return &Fenwick{n: n, t: make([]int64, n+1)}
}

// Update adds v at index i (1-based).
func (f *Fenwick) Update(i int, v int64) {
   for ; i <= f.n; i += i & -i {
       f.t[i] = (f.t[i] + v) % MOD
   }
}

// Query returns the prefix sum up to i (1-based).
func (f *Fenwick) Query(i int) int64 {
   var s int64
   for ; i > 0; i -= i & -i {
       s = (s + f.t[i]) % MOD
   }
   return s
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   // Pair values with original indices
   arr := make([]struct{ val, idx int }, n)
   for i := 0; i < n; i++ {
       arr[i] = struct{ val, idx int }{a[i], i + 1}
   }
   // Sort by value ascending
   sort.Slice(arr, func(i, j int) bool {
       return arr[i].val < arr[j].val
   })
   ftPos := NewFenwick(n)   // sums of positions
   ftSuf := NewFenwick(n)   // sums of (n-pos+1)
   var ans int64
   // Process in increasing value order
   for _, p := range arr {
       v := int64(p.val % MOD)
       j := p.idx
       // sum of positions <= j
       x1 := ftPos.Query(j)
       // sum of (n-pos+1) for positions > j
       totalSuf := ftSuf.Query(n)
       x2 := (totalSuf - ftSuf.Query(j) + MOD) % MOD
       // number of choices for right endpoint weight
       nj1 := int64(n - j + 1)
       // contributions
       t1 := x1 * nj1 % MOD
       t2 := int64(j) * x2 % MOD
       self := int64(j) * nj1 % MOD
       Sj := (t1 + t2 + self) % MOD
       ans = (ans + v*Sj) % MOD
       // insert current position
       ftPos.Update(j, int64(j))
       ftSuf.Update(j, nj1)
   }
   fmt.Println(ans)
}
