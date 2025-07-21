package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   var x int
   if _, err := fmt.Fscan(in, &n, &x); err != nil {
       return
   }
   a := make([]int, 0, n)
   hasOne := false
   for i := 0; i < n; i++ {
       var v int
       fmt.Fscan(in, &v)
       if v == 1 {
           hasOne = true
       }
       a = append(a, v)
   }
   // If only one possible k
   if x <= 2 {
       fmt.Println(0)
       return
   }
   // Mode with a_i == 1 gives exact k
   if hasOne {
       fmt.Println(1)
       return
   }
   L := x - 1
   const maxL = 2000000
   if L > maxL {
       fmt.Println(-1)
       return
   }
   // filter a_i <= L and unique
   uniq := make([]bool, L+1)
   for _, v := range a {
       if v <= L {
           uniq[v] = true
       }
   }
   vals := make([]int, 0, len(uniq))
   for v, ok := range uniq {
       if ok && v >= 2 {
           vals = append(vals, v)
       }
   }
   if len(vals) == 0 {
       fmt.Println(-1)
       return
   }
   sort.Ints(vals)
   // sieve smallest prime factor up to L
   spf := make([]int, L+1)
   for i := 2; i <= L; i++ {
       if spf[i] == 0 {
           for j := i; j <= L; j += i {
               if spf[j] == 0 {
                   spf[j] = i
               }
           }
       }
   }
   // build B of minimal elements
   hasB := make([]bool, L+1)
   B := make([]int, 0, len(vals))
   for _, v := range vals {
       // check if any b in B divides v
       ok := true
       // generate divisors via factorization
       // use recursive generation
       divs := []int{1}
       vv := v
       for vv > 1 {
           p := spf[vv]
           cnt := 0
           for vv%p == 0 {
               vv /= p
               cnt++
           }
           // multiply existing divs by p^e
           base := len(divs)
           mul := 1
           for e := 1; e <= cnt; e++ {
               mul *= p
               for i := 0; i < base; i++ {
                   divs = append(divs, divs[i]*mul)
               }
           }
       }
       for _, d := range divs {
           if d == v {
               continue
           }
           if hasB[d] {
               ok = false
               break
           }
       }
       if ok {
           hasB[v] = true
           B = append(B, v)
       }
   }
   // coverage check
   covered := make([]bool, L+1)
   for _, b := range B {
       for m := b; m <= L; m += b {
           covered[m] = true
       }
   }
   for k := 2; k <= L; k++ {
       if !covered[k] {
           fmt.Println(-1)
           return
       }
   }
   fmt.Println(len(B))
}
