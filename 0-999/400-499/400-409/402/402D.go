package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   line, _ := in.ReadString('\n')
   parts := splitInts(line)
   n, m := parts[0], parts[1]
   a := make([]int64, n)
   line, _ = in.ReadString('\n')
   vals := splitInts(line)
   for i := 0; i < n; i++ {
       a[i] = int64(vals[i])
   }
   bad := make(map[int64]bool, m)
   line, _ = in.ReadString('\n')
   bvals := splitInts(line)
   for i := 0; i < m; i++ {
       bad[int64(bvals[i])] = true
   }

   primes := sieve(32000)
   // fcache to memoize f(x)
   fcache := make(map[int64]int)
   f := func(x int64) int {
       if v, ok := fcache[x]; ok {
           return v
       }
       orig := x
       score := 0
       for _, p := range primes {
           pp := int64(p)
           if pp*pp > x {
               break
           }
           for x%pp == 0 {
               if bad[pp] {
                   score--
               } else {
                   score++
               }
               x /= pp
           }
       }
       if x > 1 {
           if bad[x] {
               score--
           } else {
               score++
           }
       }
       fcache[orig] = score
       return score
   }

   // initial beauty
   var ans int64
   for _, v := range a {
       ans += int64(f(v))
   }
   // prefix gcds
   prefix := make([]int64, n)
   prefix[0] = a[0]
   for i := 1; i < n; i++ {
       prefix[i] = gcd(prefix[i-1], a[i])
   }
   var rem int64 = 1
   for i := n - 1; i >= 0; i-- {
       // current gcd after removals
       g := prefix[i] / rem
       s := f(g)
       if s < 0 {
           ans -= int64(s) * int64(i+1)
           rem = prefix[i]
       }
   }
   fmt.Fprintln(out, ans)
}

// splitInts splits a line of space-separated ints
func splitInts(s string) []int {
   var res []int
   num := 0
   neg := false
   seen := false
   for i := 0; i < len(s); i++ {
       c := s[i]
       if c == '-' {
           neg = true
       } else if c >= '0' && c <= '9' {
           seen = true
           num = num*10 + int(c-'0')
       } else {
           if seen {
               if neg {
                   num = -num
               }
               res = append(res, num)
               num = 0
               neg = false
               seen = false
           }
       }
   }
   if seen {
       if neg {
           num = -num
       }
       res = append(res, num)
   }
   return res
}

// sieve returns primes up to n inclusive
func sieve(n int) []int {
   isPrime := make([]bool, n+1)
   for i := 2; i <= n; i++ {
       isPrime[i] = true
   }
   for i := 2; i*i <= n; i++ {
       if isPrime[i] {
           for j := i * i; j <= n; j += i {
               isPrime[j] = false
           }
       }
   }
   var primes []int
   for i := 2; i <= n; i++ {
       if isPrime[i] {
           primes = append(primes, i)
       }
   }
   return primes
}

// gcd computes greatest common divisor
func gcd(a, b int64) int64 {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}
