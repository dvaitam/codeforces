package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

var sm map[int]int
var b1, b2 int
var odds []int

// dfs tries to find two disjoint subsets of odds with equal sum
func dfs(n, le, bitmask, sum int) {
   if b1 != 0 || b2 != 0 {
       return
   }
   if n == 0 {
       if prev, ok := sm[sum]; ok {
           b1 = prev
           b2 = bitmask
       } else {
           sm[sum] = bitmask
       }
       return
   }
   // include element n-1 if we still can pick le elements
   if le > 0 {
       dfs(n-1, le-1, bitmask|(1<<(n-1)), sum+odds[n-1])
   }
   // skip element n-1
   if n > le {
       dfs(n-1, le, bitmask, sum)
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   fmt.Fscan(reader, &n)
   b := make([]int, n)
   in := make(map[int]bool)
   evens := make([]int, 0, n)
   odds = make([]int, 0, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &b[i])
       if !in[b[i]] {
           in[b[i]] = true
       }
       if b[i]%2 == 0 {
           evens = append(evens, b[i])
       } else {
           odds = append(odds, b[i])
       }
   }
   // check duplicates
   if len(in) != n {
       fmt.Fprintln(writer, "YES")
       // unique values
       vals := make([]int, 0, len(in))
       for x := range in {
           vals = append(vals, x)
       }
       sort.Ints(vals)
       for _, x := range vals {
           fmt.Fprint(writer, x, " ")
       }
       // fill zeros
       for i := len(in); i < n; i++ {
           fmt.Fprint(writer, 0, " ")
       }
       return
   }
   // case with easy triple
   if len(evens) >= 3 || (len(evens) >= 1 && len(odds) >= 2) {
       fmt.Fprintln(writer, "YES")
       a := make([]int, 3)
       if len(evens) >= 3 {
           // take last three
           a[0], a[1], a[2] = evens[len(evens)-1], evens[len(evens)-2], evens[len(evens)-3]
           evens = evens[:len(evens)-3]
       } else {
           a[0] = evens[len(evens)-1]
           evens = evens[:len(evens)-1]
           a[1] = odds[len(odds)-1]
           odds = odds[:len(odds)-1]
           a[2] = odds[len(odds)-1]
           odds = odds[:len(odds)-1]
       }
       s := (a[0] + a[1] + a[2]) / 2
       for i := 0; i < 3; i++ {
           a[i] = s - a[i]
       }
       for i := 0; i < 3; i++ {
           fmt.Fprint(writer, a[i], " ")
       }
       // rest odds
       for _, x := range odds {
           fmt.Fprint(writer, x-a[0], " ")
       }
       // rest evens
       for _, x := range evens {
           fmt.Fprint(writer, x-a[0], " ")
       }
       return
   }
   // no solution if too few odds
   if len(odds) <= 1 {
       fmt.Fprintln(writer, "NO")
       return
   }
   // subset sum on odds
   maxn := min(len(odds), 28)
   maxle := min(len(odds)/2, 14)
   // only consider first maxn odds
   odds = odds[:maxn]
   sm = make(map[int]int)
   b1, b2 = 0, 0
   dfs(maxn, maxle, 0, 0)
   if b1 == 0 {
       fmt.Fprintln(writer, "NO")
       return
   }
   // remove common bits
   w := b1 & b2
   b1 ^= w
   b2 ^= w
   fmt.Fprintln(writer, "YES")
   // first zero
   fmt.Fprint(writer, 0, " ")
   la := 0
   i1, i2 := 0, 0
   // alternate picking from b1 and b2
   for b1 != 0 {
       // pick from b1
       for i1 < maxn {
           if (b1 & (1 << i1)) != 0 {
               la = odds[i1] - la
               fmt.Fprint(writer, la, " ")
               odds[i1] = 0
               b1 ^= 1 << i1
               i1++
               break
           }
           i1++
       }
       if b1 == 0 {
           break
       }
       // pick from b2
       for i2 < maxn {
           if (b2 & (1 << i2)) != 0 {
               la = odds[i2] - la
               fmt.Fprint(writer, la, " ")
               odds[i2] = 0
               b2 ^= 1 << i2
               i2++
               break
           }
           i2++
       }
   }
   // one more from b2 if remains
   for i2 < maxn {
       if (b2 & (1 << i2)) != 0 {
           la = odds[i2] - la
           fmt.Fprint(writer, la, " ")
           odds[i2] = 0
           b2 ^= 1 << i2
           break
       }
       i2++
   }
   // rest odds
   for _, x := range odds {
       if x != 0 {
           fmt.Fprint(writer, x, " ")
       }
   }
}

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}
