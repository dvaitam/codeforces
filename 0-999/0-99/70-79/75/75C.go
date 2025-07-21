package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var a, b int64
   fmt.Fscan(reader, &a, &b)
   g := gcd(a, b)
   // factorize g and generate divisors
   primes := factorize(g)
   divs := make([]int64, 0, 64)
   genDivs(primes, 0, 1, &divs)
   sort.Slice(divs, func(i, j int) bool { return divs[i] < divs[j] })

   var n int
   fmt.Fscan(reader, &n)
   for i := 0; i < n; i++ {
       var low, high int64
       fmt.Fscan(reader, &low, &high)
       // find rightmost divisor <= high
       idx := sort.Search(len(divs), func(i int) bool { return divs[i] > high }) - 1
       if idx >= 0 && divs[idx] >= low {
           fmt.Fprintln(writer, divs[idx])
       } else {
           fmt.Fprintln(writer, -1)
       }
   }
}

// gcd returns the greatest common divisor of x and y
func gcd(x, y int64) int64 {
   for y != 0 {
       x, y = y, x%y
   }
   return x
}

// factorize returns the prime factorization of n as a slice of (prime, exponent) pairs
func factorize(n int64) [][2]int64 {
   res := make([][2]int64, 0)
   for p := int64(2); p*p <= n; p++ {
       if n%p == 0 {
           cnt := int64(0)
           for n%p == 0 {
               n /= p
               cnt++
           }
           res = append(res, [2]int64{p, cnt})
       }
   }
   if n > 1 {
       res = append(res, [2]int64{n, 1})
   }
   return res
}

// genDivs recursively generates divisors from prime factors
func genDivs(primes [][2]int64, idx int, curr int64, divs *[]int64) {
   if idx == len(primes) {
       *divs = append(*divs, curr)
       return
   }
   p := primes[idx][0]
   cnt := primes[idx][1]
   // for exponent from 0 to cnt
   for e := int64(0); e <= cnt; e++ {
       genDivs(primes, idx+1, curr, divs)
       curr *= p
   }
}
