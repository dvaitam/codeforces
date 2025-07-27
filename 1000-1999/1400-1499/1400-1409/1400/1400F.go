package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s string
   var x int
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   fmt.Fscan(reader, &x)
   n := len(s)
   // prefix sums
   ps := make([]int, n+1)
   for i := 0; i < n; i++ {
       ps[i+1] = ps[i] + int(s[i]-'0')
   }
   type interval struct{ l, r int }
   var ivs []interval
   // find x-prime substrings
   for l := 0; l < n; l++ {
       for r := l; r < n; r++ {
           sum := ps[r+1] - ps[l]
           if sum > x {
               break
           }
           if sum < x {
               continue
           }
           // sum == x, check minimality: no sub-interval (strict) with sum divides x
           ok := true
           for i := l; i <= r && ok; i++ {
               for j := i; j <= r; j++ {
                   if i == l && j == r {
                       continue
                   }
                   sub := ps[j+1] - ps[i]
                   if sub != 0 && x%sub == 0 {
                       ok = false
                       break
                   }
               }
           }
           if ok {
               ivs = append(ivs, interval{l, r})
           }
           break
       }
   }
   if len(ivs) == 0 {
       fmt.Println(0)
       return
   }
   // greedy hitting set: pick minimal points
   sort.Slice(ivs, func(i, j int) bool {
       if ivs[i].r != ivs[j].r {
           return ivs[i].r < ivs[j].r
       }
       return ivs[i].l < ivs[j].l
   })
   ans := 0
   last := -1
   for _, it := range ivs {
       if it.l > last {
           ans++
           last = it.r
       }
   }
   fmt.Println(ans)
}
