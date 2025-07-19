package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var L, R int
   fmt.Fscan(in, &L, &R)
   // Precompute squares up to R
   v := make([]int, 0)
   for i := 1; i*i <= R; i++ {
       v = append(v, i*i)
   }
   // Sieve small primes up to sqrt(R)
   NN := int(math.Sqrt(float64(R))) + 1
   isP := make([]bool, NN+1)
   primes := make([]int, 0)
   for i := 2; i <= NN; i++ {
       isP[i] = true
   }
   for i := 2; i <= NN; i++ {
       if isP[i] {
           primes = append(primes, i)
           for j := i * 2; j <= NN; j += i {
               isP[j] = false
           }
       }
   }
   // Compute and output result
   res := FInt(L, R+1, v, primes)
   fmt.Println(res)
}

// FInt counts numbers x in [L, R) that are prime and representable as sum of two squares
func FInt(L, R int, v, primes []int) int {
   size := R - L
   ok := make([]bool, size)
   primeArr := make([]bool, size)
   for i := 0; i < size; i++ {
       primeArr[i] = true
   }
   // mark sums of two squares
   V := len(v)
   for _, vi := range v {
       // find first index j where vi + v[j] >= L
       j0 := sort.Search(V, func(j int) bool { return vi+v[j] >= L })
       for j := j0; j < V; j++ {
           x := vi + v[j]
           if x >= R {
               break
           }
           ok[x-L] = true
       }
   }
   // sieve primes in [L, R)
   for _, p := range primes {
       x := (L/p)*p
       if x < L || x <= p {
           x += p
       }
       for x < R {
           primeArr[x-L] = false
           x += p
       }
   }
   // count positions that are prime and sum of squares
   cnt := 0
   for i := 0; i < size; i++ {
       if ok[i] && primeArr[i] {
           cnt++
       }
   }
   return cnt
}
