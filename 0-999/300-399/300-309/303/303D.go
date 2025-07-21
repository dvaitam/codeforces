package main

import (
   "bufio"
   "fmt"
   "os"
)

func isPrime(n int) bool {
   if n < 2 {
       return false
   }
   if n%2 == 0 {
       return n == 2
   }
   for i := 3; i*i <= n; i += 2 {
       if n%i == 0 {
           return false
       }
   }
   return true
}

func powMod(a, e, mod int) int {
   res := 1
   a %= mod
   for e > 0 {
       if e&1 != 0 {
           res = int((int64(res) * int64(a)) % int64(mod))
       }
       a = int((int64(a) * int64(a)) % int64(mod))
       e >>= 1
   }
   return res
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   var x int
   if _, err := fmt.Fscan(in, &n, &x); err != nil {
       return
   }
   // special case n == 1
   if n == 1 {
       if x > 2 {
           fmt.Fprintln(out, x-1)
       } else {
           fmt.Fprintln(out, -1)
       }
       return
   }
   p := n + 1
   if !isPrime(p) {
       fmt.Fprintln(out, -1)
       return
   }
   // factor p-1 = n
   m := n
   var primes []int
   for i := 2; i*i <= m; i++ {
       if m%i == 0 {
           primes = append(primes, i)
           for m%i == 0 {
               m /= i
           }
       }
   }
   if m > 1 {
       primes = append(primes, m)
   }
   // prepare division
   M := x - 1
   kMax := M / p
   rem := M - kMax*p
   // check residues
   // function to test primitive root
   isRoot := func(r int) bool {
       if r%p == 0 {
           return false
       }
       for _, q := range primes {
           if powMod(r, n/q, p) == 1 {
               return false
           }
       }
       return true
   }
   // search best
   var ans int = -1
   if rem > 0 {
       for r := rem; r >= 1; r-- {
           if isRoot(r) {
               ans = r + kMax*p
               break
           }
       }
       if ans >= 2 {
           fmt.Fprintln(out, ans)
           return
       }
   }
   // fallback kMax-1
   k := kMax - 1
   if k >= 0 {
       for r := p - 1; r >= 1; r-- {
           if isRoot(r) {
               ans = r + k*p
               break
           }
       }
   }
   if ans >= 2 && ans < x {
       fmt.Fprintln(out, ans)
   } else {
       fmt.Fprintln(out, -1)
   }
}
