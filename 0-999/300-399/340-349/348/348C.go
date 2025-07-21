package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n, m, q int
   fmt.Fscan(reader, &n, &m, &q)
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   sets := make([][]int, m)
   const B = 320
   heavyIndex := make([]int, m)
   for i := 0; i < m; i++ {
       heavyIndex[i] = -1
       var sz int
       fmt.Fscan(reader, &sz)
       s := make([]int, sz)
       for j := 0; j < sz; j++ {
           fmt.Fscan(reader, &s[j])
           s[j]--
       }
       sets[i] = s
   }
   // identify heavy sets
   heavies := make([]int, 0, m)
   for i := 0; i < m; i++ {
       if len(sets[i]) > B {
           heavyIndex[i] = len(heavies)
           heavies = append(heavies, i)
       }
   }
   hc := len(heavies)
   // for each element, which heavies contain it
   heaviesContaining := make([][]int, n)
   for h, si := range heavies {
       for _, j := range sets[si] {
           heaviesContaining[j] = append(heaviesContaining[j], h)
       }
   }
   // compute interHeavy
   interHeavy := make([][]int, hc)
   for h := 0; h < hc; h++ {
       interHeavy[h] = make([]int, hc)
   }
   for h1, si := range heavies {
       for _, j := range sets[si] {
           for _, h2 := range heaviesContaining[j] {
               interHeavy[h1][h2]++
           }
       }
   }
   // heavySum: sum of a[j] for heavy sets
   heavySum := make([]int64, hc)
   for h, si := range heavies {
       var sum int64
       for _, j := range sets[si] {
           sum += a[j]
       }
       heavySum[h] = sum
   }
   tag := make([]int64, hc)

   // process queries
   for qi := 0; qi < q; qi++ {
       var op string
       fmt.Fscan(reader, &op)
       if op == "?" {
           var k int
           fmt.Fscan(reader, &k)
           k--
           var res int64
           hi := heavyIndex[k]
           if hi >= 0 {
               res = heavySum[hi]
           } else {
               for _, j := range sets[k] {
                   res += a[j]
               }
           }
           // add contributions from heavy tags
           if hi >= 0 {
               // heavy query: use precomputed intersections
               for h := 0; h < hc; h++ {
                   cnt := interHeavy[h][hi]
                   if cnt != 0 {
                       res += tag[h] * int64(cnt)
                   }
               }
           } else {
               // small query: accumulate tag per common element
               for _, j := range sets[k] {
                   for _, h := range heaviesContaining[j] {
                       res += tag[h]
                   }
               }
           }
           fmt.Fprintln(writer, res)
       } else if op == "+" {
           var k int
           var x int64
           fmt.Fscan(reader, &k, &x)
           k--
           hi := heavyIndex[k]
           if hi >= 0 {
               tag[hi] += x
           } else {
               // small update
               for _, j := range sets[k] {
                   a[j] += x
                   // update heavySum for heavies containing j
                   for _, h := range heaviesContaining[j] {
                       heavySum[h] += x
                   }
               }
           }
       }
   }
}
