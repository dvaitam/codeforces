package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var m, l, r, k int64
   if _, err := fmt.Fscan(in, &m, &l, &r, &k); err != nil {
       return
   }
   var ans int64 = 1
   t := int64(1)
   lm1 := l - 1
   for t <= r {
       ur := r / t
       ul := lm1 / t
       cnt := ur - ul
       // compute next change point
       next1 := r/ur + 1
       var next2 int64
       if ul > 0 {
           next2 = lm1/ul + 1
       } else {
           next2 = r + 1
       }
       next := next1
       if next2 < next {
           next = next2
       }
       if cnt >= k {
           // any t' in [t, next-1] is valid, take the max
           cand := next - 1
           if cand > ans {
               ans = cand
           }
       }
       t = next
   }
   // compute fibonacci(ans) mod m via fast doubling
   f := func(n int64) int64 {
       var fib func(n int64) (int64, int64)
       fib = func(n int64) (int64, int64) {
           if n == 0 {
               return 0, 1
           }
           a, b := fib(n >> 1)
           // c = F(2k) = F(k)*(2*F(k+1) - F(k))
           t1 := (2*b - a) % m
           if t1 < 0 {
               t1 += m
           }
           c := (a * t1) % m
           // d = F(2k+1) = F(k)^2 + F(k+1)^2
           d := (a*a + b*b) % m
           if n&1 == 0 {
               return c, d
           }
           // for odd: F(n) = d, F(n+1) = c+d
           return d, (c + d) % m
       }
       fn, _ := fib(n)
       return fn
   }
   res := f(ans) % m
   fmt.Println(res)
}
