package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m, s int
   if _, err := fmt.Fscan(in, &n, &m, &s); err != nil {
       return
   }
   amax := (n - 1) / 2
   bmax := (m - 1) / 2
   // Precompute smallest prime factors up to s
   maxVal := s
   if maxVal < 1 {
       maxVal = 1
   }
   spf := make([]int, maxVal+1)
   for i := 2; i <= maxVal; i++ {
       if spf[i] == 0 {
           for j := i; j <= maxVal; j += i {
               if spf[j] == 0 {
                   spf[j] = i
               }
           }
       }
   }
   // Precompute odd divisors of s
   oddDivS := oddDivisors(s, spf)
   var ans int64
   // Iterate over parameters a,b
   for a := 0; a <= amax; a++ {
       A := 2*a + 1
       xa := int64(n - 2*a)
       if xa <= 0 {
           continue
       }
       for b := 0; b <= bmax; b++ {
           B := 2*b + 1
           yb := int64(m - 2*b)
           if yb <= 0 {
               continue
           }
           AB := A * B
           t := s - AB
           // case1: c<=a, d<=b, t==0
           if t == 0 {
               cnt := int64(a+1) * int64(b+1)
               ans += cnt * xa * yb
           }
           // case2 and case3: t>0 and even
           if t > 0 && t%2 == 0 {
               T := t / 2
               // get odd divisors of T
               oddDivT := oddDivisors(T, spf)
               // case2: c<=a, d>b
               for _, C := range oddDivT {
                   c := (C - 1) / 2
                   if c > a {
                       continue
                   }
                   K := T / C
                   d := b + K
                   if d > bmax {
                       continue
                   }
                   yd := int64(m - 2*d)
                   if yd > 0 {
                       ans += xa * yd
                   }
               }
               // case3: c>a, d<=b
               for _, D := range oddDivT {
                   dIdx := (D - 1) / 2
                   if dIdx > b {
                       continue
                   }
                   K := T / D
                   c := a + K
                   if c > amax {
                       continue
                   }
                   xd := int64(n - 2*c)
                   if xd > 0 {
                       ans += xd * yb
                   }
               }
           }
           // case4: c>a, d>b, C*D = s
           for _, C := range oddDivS {
               if C <= A {
                   continue
               }
               if s%C != 0 {
                   continue
               }
               D := s / C
               if D%2 == 0 {
                   continue
               }
               c := (C - 1) / 2
               d := (D - 1) / 2
               if c <= a || d <= b || c > amax || d > bmax {
                   continue
               }
               xd := int64(n - 2*c)
               yd := int64(m - 2*d)
               if xd > 0 && yd > 0 {
                   ans += xd * yd
               }
           }
       }
   }
   fmt.Println(ans)
}

// oddDivisors returns all odd divisors of x using spf sieve
func oddDivisors(x int, spf []int) []int {
   if x <= 0 {
       return nil
   }
   // remove factors of 2
   for x%2 == 0 {
       x /= 2
   }
   // factor x
   m := x
   primes := make([]int, 0)
   counts := make([]int, 0)
   for m > 1 {
       p := spf[m]
       if p == 0 {
           p = m
       }
       cnt := 0
       for m%p == 0 {
           m /= p
           cnt++
       }
       primes = append(primes, p)
       counts = append(counts, cnt)
   }
   // generate divisors
   divs := []int{1}
   for i, p := range primes {
       pow := 1
       var next []int
       for e := 0; e < counts[i]; e++ {
           pow *= p
           for _, d := range divs {
               next = append(next, d*pow)
           }
       }
       divs = append(divs, next...)
   }
   return divs
}
