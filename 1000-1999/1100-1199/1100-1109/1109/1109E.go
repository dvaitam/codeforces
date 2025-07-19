package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

var (
   n, mod, q int
   primes    []int
   st, lz    []int
   fen0      []int
   fenExp    [][]int
   reader    = bufio.NewReader(os.Stdin)
   writer    = bufio.NewWriter(os.Stdout)
)

func mult(a, b int) int {
   return int((int64(a) * int64(b)) % int64(mod))
}

func fpow(a, e int) int {
   res := 1
   a %= mod
   for e > 0 {
       if e&1 == 1 {
           res = mult(res, a)
       }
       a = mult(a, a)
       e >>= 1
   }
   return res
}

// extended gcd
func egcd(a, b int) (x, y int) {
   if b == 0 {
       return 1, 0
   }
   x1, y1 := egcd(b, a%b)
   return y1, x1 - (a/b)*y1
}

func invert(x int) int {
   x %= mod
   if x < 0 {
       x += mod
   }
   inv, _ := egcd(x, mod)
   inv %= mod
   if inv < 0 {
       inv += mod
   }
   return inv
}

func push(p, L, R int) {
   if lz[p] != 1 {
       st[p] = mult(st[p], lz[p])
       if L < R {
           lz[p*2] = mult(lz[p*2], lz[p])
           lz[p*2+1] = mult(lz[p*2+1], lz[p])
       }
       lz[p] = 1
   }
}

// update range [l..r] multiply by v
func updateRange(p, L, R, l, r, v int) {
   push(p, L, R)
   if r < L || R < l {
       return
   }
   if l <= L && R <= r {
       lz[p] = mult(lz[p], v)
       push(p, L, R)
       return
   }
   mid := (L + R) >> 1
   updateRange(p*2, L, mid, l, r, v)
   updateRange(p*2+1, mid+1, R, l, r, v)
   st[p] = st[p*2]
   st[p] += st[p*2+1]
   if st[p] >= mod {
       st[p] -= mod
   }
}

// update point pos to value v
func updatePoint(p, L, R, pos, v int) {
   push(p, L, R)
   if pos < L || R < pos {
       return
   }
   if L == R {
       st[p] = v
       return
   }
   mid := (L + R) >> 1
   updatePoint(p*2, L, mid, pos, v)
   updatePoint(p*2+1, mid+1, R, pos, v)
   st[p] = st[p*2]
   st[p] += st[p*2+1]
   if st[p] >= mod {
       st[p] -= mod
   }
}

// query sum on [l..r]
func queryRange(p, L, R, l, r int) int {
   push(p, L, R)
   if r < L || R < l {
       return 0
   }
   if l <= L && R <= r {
       return st[p]
   }
   mid := (L + R) >> 1
   res := queryRange(p*2, L, mid, l, r) + queryRange(p*2+1, mid+1, R, l, r)
   if res >= mod {
       res -= mod
   }
   return res
}

// fen0: multiplicative BIT, fenExp: exponent BITs
func fen0Update(pos, v int) {
   for i := pos + 1; i <= n; i += i & -i {
       fen0[i] = mult(fen0[i], v)
   }
}

func fen0Query(pos int) int {
   res := 1
   for i := pos + 1; i > 0; i -= i & -i {
       res = mult(res, fen0[i])
   }
   return res
}

func fenExpUpdate(idx, pos, v int) {
   for i := pos + 1; i <= n; i += i & -i {
       fenExp[idx][i] += v
   }
}

func fenExpQuery(idx, pos int) int {
   sum := 0
   for i := pos + 1; i > 0; i -= i & -i {
       sum += fenExp[idx][i]
   }
   return sum
}

func main() {
   defer writer.Flush()
   fmt.Fscan(reader, &n, &mod)
   // factor mod
   t := mod
   cnts := make(map[int]int)
   for i := 2; i*i <= t; i++ {
       for t%i == 0 {
           cnts[i]++
           t /= i
       }
   }
   if t > 1 {
       cnts[t]++
   }
   for p := range cnts {
       primes = append(primes, p)
   }
   sort.Ints(primes)
   // init structures
   st = make([]int, 4*n)
   lz = make([]int, 4*n)
   for i := range lz {
       lz[i] = 1
   }
   fen0 = make([]int, n+2)
   for i := 1; i <= n; i++ {
       fen0[i] = 1
   }
   fenExp = make([][]int, len(primes))
   for i := range fenExp {
       fenExp[i] = make([]int, n+2)
   }
   // read initial array values
   // read initial array
   for i := 0; i < n; i++ {
       var x int
       fmt.Fscan(reader, &x)
       // value x read
       // segment tree point assign
       updatePoint(1, 0, n-1, i, x%mod)
       // factor x
       rem := x
       for idx, p := range primes {
           cnt := 0
           for rem%p == 0 {
               cnt++
               rem /= p
           }
           if cnt > 0 {
               fenExpUpdate(idx, i, cnt)
               if i+1 < n {
                   fenExpUpdate(idx, i+1, -cnt)
               }
           }
       }
       // co-prime part
       cp := rem % mod
       invcp := invert(cp)
       fen0Update(i, cp)
       if i+1 < n {
           fen0Update(i+1, invcp)
       }
   }
   fmt.Fscan(reader, &q)
   for qi := 0; qi < q; qi++ {
       var op int
       fmt.Fscan(reader, &op)
       if op == 1 {
           var l, r, x int
           fmt.Fscan(reader, &l, &r, &x)
           l--
           r--
           v := x % mod
           updateRange(1, 0, n-1, l, r, v)
           rem := x
           for idx, p := range primes {
               cnt := 0
               for rem%p == 0 {
                   cnt++
                   rem /= p
               }
               if cnt > 0 {
                   fenExpUpdate(idx, l, cnt)
                   if r+1 < n {
                       fenExpUpdate(idx, r+1, -cnt)
                   }
               }
           }
           cp := rem % mod
           invcp := invert(cp)
           fen0Update(l, cp)
           if r+1 < n {
               fen0Update(r+1, invcp)
           }
       } else if op == 2 {
           var u, x int
           fmt.Fscan(reader, &u, &x)
           u--
           // g: prime factor part
           g := 1
           rem := x
           for _, p := range primes {
               for rem%p == 0 {
                   g *= p
                   rem /= p
               }
           }
           // r: current co-prime multiplier at u
           rco := fen0Query(u)
           tval := 1
           // adjust prime exponents
           for idx, p := range primes {
               c := fenExpQuery(idx, u)
               total := 0
               for g%p == 0 {
                   c--
                   total++
                   g /= p
               }
               tval = mult(tval, fpow(p, c))
               if total > 0 {
                   fenExpUpdate(idx, u, -total)
               }
               if u+1 < n && total > 0 {
                   fenExpUpdate(idx, u+1, total)
               }
           }
           // co-prime part
           tval = mult(tval, rco)
           invx := invert(rem % mod)
           // update fen0 diff
           fen0Update(u, invx)
           if u+1 < n {
               fen0Update(u+1, rem%mod)
           }
           tval = mult(tval, invx)
           updateRange(1, 0, n-1, u, u, tval)
       } else {
           var l, r int
           fmt.Fscan(reader, &l, &r)
           l--
           r--
           res := queryRange(1, 0, n-1, l, r)
           fmt.Fprintln(writer, res)
       }
   }
}
