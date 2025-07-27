package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 1000000007

type Fenwick struct {
   n    int
   tree []int64
}

func NewFenwick(n int) *Fenwick {
   f := &Fenwick{n: n, tree: make([]int64, n+2)}
   for i := range f.tree {
       f.tree[i] = 1
   }
   return f
}

func (f *Fenwick) mul(i int, v int64) {
   for ; i <= f.n; i += i & -i {
       f.tree[i] = f.tree[i] * v % MOD
   }
}

// Range multiply [l..r] by v, invV = inverse of v modulo MOD
func (f *Fenwick) rangeMul(l, r int, v, invV int64) {
   if l > r {
       return
   }
   f.mul(l, v)
   if r+1 <= f.n {
       f.mul(r+1, invV)
   }
}

// prefix product at x
func (f *Fenwick) query(x int) int64 {
   res := int64(1)
   for i := x; i > 0; i -= i & -i {
       res = res * f.tree[i] % MOD
   }
   return res
}

// mod exponentiation
func modPow(a int64, e int) int64 {
   res := int64(1)
   for e > 0 {
       if e&1 != 0 {
           res = res * a % MOD
       }
       a = a * a % MOD
       e >>= 1
   }
   return res
}

// modular inverse via Fermat's little theorem
func modInv(a int64) int64 {
   return modPow(a, MOD-2)
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   fmt.Fscan(in, &n)
   a := make([]int, n+1)
   maxA := 0
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &a[i])
       if a[i] > maxA {
           maxA = a[i]
       }
   }
   // sieve for spf
   spf := make([]int, maxA+1)
   for i := 2; i <= maxA; i++ {
       if spf[i] == 0 {
           for j := i; j <= maxA; j += i {
               if spf[j] == 0 {
                   spf[j] = i
               }
           }
       }
   }
   // stacks for each prime: slice of (pos, exp)
   type pe struct{ pos, exp int }
   stacks := make([][]pe, maxA+1)
   for p := 2; p <= maxA; p++ {
       if spf[p] == p { // only primes, but we'll init all for simplicity
           stacks[p] = []pe{{0, 0}}
       }
   }
   bit := NewFenwick(n)
   // process array
   for i := 1; i <= n; i++ {
       // factorize a[i]
       x := a[i]
       for x > 1 {
           p := spf[x]
           cnt := 0
           for x%p == 0 {
               x /= p
               cnt++
           }
           // ensure stack initialized
           if stacks[p] == nil {
               stacks[p] = []pe{{0, 0}}
           }
           stk := stacks[p]
           // remove smaller or equal exponents
           // pop until previous exponent > cnt, keep sentinel at bottom
           for len(stk) > 1 && stk[len(stk)-1].exp <= cnt {
               last := stk[len(stk)-1]
               stk = stk[:len(stk)-1]
               prev := stk[len(stk)-1]
               // remove contribution of last.exp on (prev.pos+1 .. last.pos)
               pe_k := modPow(int64(p), last.exp)
               inv_pe_k := modInv(pe_k)
               l := prev.pos + 1
               r := last.pos
               bit.rangeMul(l, r, inv_pe_k, pe_k)
           }
           // add new contribution
           prev := stk[len(stk)-1]
           pe_v := modPow(int64(p), cnt)
           inv_pe_v := modInv(pe_v)
           l := prev.pos + 1
           r := i
           bit.rangeMul(l, r, pe_v, inv_pe_v)
           // push current
           stk = append(stk, pe{i, cnt})
           stacks[p] = stk
       }
   }
   // queries
   var q int
   fmt.Fscan(in, &q)
   last := int64(0)
   for qi := 0; qi < q; qi++ {
       var x, y int
       fmt.Fscan(in, &x, &y)
       l := (last + int64(x))%int64(n) + 1
       r := (last + int64(y))%int64(n) + 1
       if l > r {
           l, r = r, l
       }
       prer := bit.query(int(r))
       prel := bit.query(int(l) - 1)
       invPrel := modInv(prel)
       ans := prer * invPrel % MOD
       fmt.Fprint(out, ans)
       if qi+1 < q {
           fmt.Fprint(out, ' ')
       }
       last = ans
   }
   fmt.Fprintln(out)
}
