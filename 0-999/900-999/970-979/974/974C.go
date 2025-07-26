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

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for ; t > 0; t-- {
       var n int
       fmt.Fscan(reader, &n)
       a := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &a[i])
       }
       if n < 3 {
           fmt.Fprintln(writer, 0)
           continue
       }
       pairCount := make(map[uint64]int)
       tripleCount := make(map[[3]int]int)
       // process triples
       for i := 0; i+2 < n; i++ {
           x, y, z := a[i], a[i+1], a[i+2]
           // count unique pairs of values per triple
           freq := make(map[int]int)
           freq[x]++
           freq[y]++
           freq[z]++
           // count same-value pairs (u,u)
           for u, c := range freq {
               if c >= 2 {
                   key := (uint64(u) << 32) | uint64(u)
                   pairCount[key]++
               }
           }
           // count distinct-value pairs (u,v)
           // collect unique values
           uvals := make([]int, 0, len(freq))
           for u := range freq {
               uvals = append(uvals, u)
           }
           sort.Ints(uvals)
           for ii := 0; ii < len(uvals); ii++ {
               for jj := ii + 1; jj < len(uvals); jj++ {
                   uu, vv := uvals[ii], uvals[jj]
                   key := (uint64(uu) << 32) | uint64(vv)
                   pairCount[key]++
               }
           }
           // count triple multiset
           t3 := []int{x, y, z}
           sort.Ints(t3)
           key3 := [3]int{t3[0], t3[1], t3[2]}
           tripleCount[key3]++
       }
       // sum pairs sharing at least two values
       var sumPairs int64
       for _, cnt := range pairCount {
           if cnt > 1 {
               sumPairs += int64(cnt) * int64(cnt-1) / 2
           }
       }
       // subtract over-count of triples sharing all three values
       var sumTriples int64
       for _, cnt := range tripleCount {
           if cnt > 1 {
               sumTriples += int64(cnt) * int64(cnt-1) / 2
           }
       }
       ans := sumPairs - 2*sumTriples
       fmt.Fprintln(writer, ans)
   }
}
