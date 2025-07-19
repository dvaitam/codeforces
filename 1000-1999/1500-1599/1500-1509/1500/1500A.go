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

   var n int
   fmt.Fscan(reader, &n)
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // count occurrences and record last positions
   cnt := make(map[int]int)
   num := make(map[int]int)
   // fresh pairs from direct duplicates
   fresh := make(map[int][2]int)
   for i, v := range a {
       if cnt[v] == 1 {
           fresh[2*v] = [2]int{num[v], i}
       }
       cnt[v]++
       num[v] = i
   }
   // sorted unique values
   keys := make([]int, 0, len(cnt))
   for k := range cnt {
       keys = append(keys, k)
   }
   sort.Ints(keys)
   // check for value with >=4 occurrences or two pairs
   c := 0
   mr := make([]int, 0, 4)
   for _, v := range keys {
       c += cnt[v] / 2
       repeats := (cnt[v] / 2) * 2
       for j := 0; j < repeats && len(mr) < 4; j++ {
           mr = append(mr, v)
       }
   }
   if c >= 2 {
       fmt.Fprintln(writer, "Yes")
       need := make(map[int]int)
       for _, v := range mr {
           need[v]++
       }
       res := make([]int, 0, 4)
       for i, v := range a {
           if need[v] > 0 {
               res = append(res, i)
               need[v]--
               if len(res) == 4 {
                   break
               }
           }
       }
       // adjust order if first two are same
       if a[res[0]] == a[res[1]] {
           res[0], res[3] = res[3], res[0]
       }
       fmt.Fprintf(writer, "%d %d %d %d\n", res[0]+1, res[1]+1, res[2]+1, res[3]+1)
       return
   }
   // search for two distinct pairs with equal sums
   for i1 := 0; i1 < len(keys); i1++ {
       v1 := keys[i1]
       for _, v2 := range keys[i1+1:] {
           sum := v1 + v2
           if old, exist := fresh[sum]; exist {
               fmt.Fprintln(writer, "Yes")
               fmt.Fprintf(writer, "%d %d %d %d\n", old[0]+1, old[1]+1, num[v1]+1, num[v2]+1)
               return
           }
           fresh[sum] = [2]int{num[v1], num[v2]}
       }
   }
   fmt.Fprintln(writer, "No")
}
