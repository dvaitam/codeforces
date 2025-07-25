package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 998244353

// Fenwick tree for sum
type Fenwick struct {
   n    int
   tree []int
}

func NewFenwick(n int) *Fenwick {
   return &Fenwick{n, make([]int, n+1)}
}
func (f *Fenwick) Add(i, v int) {
   for ; i <= f.n; i += i & -i {
       f.tree[i] += v
       if f.tree[i] >= MOD {
           f.tree[i] -= MOD
       }
   }
}
func (f *Fenwick) Sum(i int) int {
   s := 0
   for ; i > 0; i -= i & -i {
       s += f.tree[i]
       if s >= MOD {
           s -= MOD
       }
   }
   return s
}
func (f *Fenwick) RangeSum(l, r int) int {
   if r < l {
       return 0
   }
   s := f.Sum(r) - f.Sum(l-1)
   if s < 0 {
       s += MOD
   }
   return s
}

func modexp(a, e int) int {
   res := 1
   for e > 0 {
       if e&1 != 0 {
           res = int64(res) * a % MOD
       }
       a = int(int64(a) * a % MOD)
       e >>= 1
   }
   return res
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(in, &n)
   a := make([][]int, n)
   for i := 0; i < n; i++ {
       a[i] = make([]int, n)
       for j := 0; j < n; j++ {
           fmt.Fscan(in, &a[i][j])
       }
   }
   // precompute factorials and derangements
   fact := make([]int, n+1)
   fact[0] = 1
   for i := 1; i <= n; i++ {
       fact[i] = int64(fact[i-1]) * i % MOD
   }
   g := make([]int, n+1)
   g[0] = 1
   if n >= 1 {
       g[1] = 0
   }
   for i := 2; i <= n; i++ {
       g[i] = int64(i-1) * (int64(g[i-1]) + int64(g[i-2]) % MOD) % MOD
   }
   // precompute Gpows
   G := g[n]
   gpows := make([]int, n+1)
   gpows[0] = 1
   for i := 1; i <= n; i++ {
       gpows[i] = int64(gpows[i-1]) * int64(G) % MOD
   }

   ans := 0
   // first row
   bit := NewFenwick(n)
   for i := 1; i <= n; i++ {
       bit.Add(i, 1)
   }
   rank := 0
   for j := 0; j < n; j++ {
       x := a[0][j]
       less := bit.RangeSum(1, x-1)
       rank = (rank + int64(less)*fact[n-j-1]%MOD) % MOD
       bit.Add(x, MOD-1)
   }
   if n > 1 {
       ans = (ans + int64(rank)*int64(gpows[n-1])%MOD) % MOD
   }
   // subsequent rows
   for i := 1; i < n; i++ {
       prev := a[i-1]
       curr := a[i]
       // init bits
       U := NewFenwick(n)
       B := NewFenwick(n)
       for j := 1; j <= n; j++ {
           U.Add(j, 1)
       }
       for j := 0; j < n; j++ {
           B.Add(prev[j], 1)
       }
       rank = 0
       for j := 0; j < n; j++ {
           // remove p[j]
           B.Add(prev[j], MOD-1)
           x := curr[j]
           less := U.RangeSum(1, x-1)
           bad := B.RangeSum(1, x-1)
           good := less - bad
           if good < 0 {
               good += MOD
           }
           rank = (rank + int64(good)*int64(g[n-j-1])%MOD) % MOD
           // remove x
           if B.RangeSum(x, x) > 0 {
               B.Add(x, MOD-1)
           }
           U.Add(x, MOD-1)
       }
       rem := n - 1 - i
       if rem >= 0 {
           ans = (ans + int64(rank)*int64(gpows[rem])%MOD) % MOD
       }
   }
   fmt.Println(ans)
}
