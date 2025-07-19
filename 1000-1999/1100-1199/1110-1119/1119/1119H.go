package main

import (
   "bufio"
   "fmt"
   "os"
)

const md = 998244353

func add(x, y int) int {
   x += y
   if x >= md {
       x -= md
   }
   return x
}

func sub(x, y int) int {
   x -= y
   if x < 0 {
       x += md
   }
   return x
}

func mul(x, y int) int {
   return int((int64(x) * int64(y)) % md)
}

func power(x, y int) int {
   res := 1
   xx := x
   for y > 0 {
       if y&1 == 1 {
           res = mul(res, xx)
       }
       xx = mul(xx, xx)
       y >>= 1
   }
   return res
}

func fwt(a []int) {
   n := len(a)
   for l := 1; l < n; l <<= 1 {
       for i := 0; i < n; i += l << 1 {
           for j := 0; j < l; j++ {
               u := a[i+j]
               v := a[i+j+l]
               a[i+j] = add(u, v)
               a[i+j+l] = sub(u, v)
           }
       }
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m, x, y, z int
   if _, err := fmt.Fscan(reader, &n, &m, &x, &y, &z); err != nil {
       return
   }
   N := 1 << m
   foo := make([]int, N)
   bar := make([]int, N)
   baz := make([]int, N)
   xorAll := 0
   for i := 0; i < n; i++ {
       var a, b, c int
       fmt.Fscan(reader, &a, &b, &c)
       a ^= c
       b ^= c
       xorAll ^= c
       foo[a]++
       bar[b]++
       baz[a^b]++
   }
   fwt(foo)
   fwt(bar)
   fwt(baz)
   ans := make([]int, N)
   inv2 := (md + 1) / 2
   // mod inputs
   x %= md; if x < 0 { x += md }
   y %= md; if y < 0 { y += md }
   z %= md; if z < 0 { z += md }
   // precompute bases
   b1 := add((x+y)%md, z)
   b2 := add((x-y+md)%md, z)
   b3 := add((-x+y+md)%md, z)
   b4 := add(((-x-y+2*md)%md), z)
   for i := 0; i < N; i++ {
       // adjust transforms
       foo[i] = add(foo[i], n)
       bar[i] = add(bar[i], n)
       baz[i] = add(baz[i], n)
       a1 := mul(foo[i], inv2)
       b1t := mul(bar[i], inv2)
       c1 := mul(baz[i], inv2)
       d := (a1 + b1t + c1 - n) / 2
       e := a1 - d
       f := b1t - d
       g := c1 - d
       // combine
       res := 1
       res = mul(res, power(b1, d))
       res = mul(res, power(b2, e))
       res = mul(res, power(b3, f))
       res = mul(res, power(b4, g))
       ans[i] = res
   }
   fwt(ans)
   invN := power(N, md-2)
   for i := 0; i < N; i++ {
       if i > 0 {
           writer.WriteByte(' ')
       }
       idx := xorAll ^ i
       writer.WriteString(fmt.Sprintf("%d", mul(ans[idx], invN)))
   }
   writer.WriteByte('\n')
}
