package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(in, &n)
   var s string
   fmt.Fscan(in, &s)
   // convert to byte array of 0/1
   y := make([]byte, n)
   for i := 0; i < n; i++ {
       if s[i] == '1' {
           y[i] = 1
       } else {
           y[i] = 0
       }
   }
   // factor n
   nn := n
   primes := make([]int, 0)
   for p := 2; p*p <= nn; p++ {
       if nn%p == 0 {
           primes = append(primes, p)
           for nn%p == 0 {
               nn /= p
           }
       }
   }
   if nn > 1 {
       primes = append(primes, nn)
   }
   // get divisors of n
   divs := make([]int, 0)
   for i := 1; i*i <= n; i++ {
       if n%i == 0 {
           divs = append(divs, i)
           if i != n/i {
               divs = append(divs, n/i)
           }
       }
   }
   // compute answer
   var ans int64 = 0
   for _, d := range divs {
       // check if for all classes mod d, xor is 0
       ok := true
       for r := 0; r < d; r++ {
           x := byte(0)
           for j := r; j < n; j += d {
               x ^= y[j]
           }
           if x != 0 {
               ok = false
               break
           }
       }
       if !ok {
           continue
       }
       // count k with gcd(n,k)==d: that's phi(n/d)
       m := n / d
       phi := int64(m)
       for _, p := range primes {
           if m%p == 0 {
               phi = phi / int64(p) * int64(p-1)
           }
       }
       ans += phi
   }
   // output result
   fmt.Println(ans)
}
