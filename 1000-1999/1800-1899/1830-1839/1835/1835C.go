package main

import (
   "bufio"
   "fmt"
   "math/rand"
   "os"
   "sort"
)

func solve(reader *bufio.Reader, writer *bufio.Writer) {
   var T, k int
   if _, err := fmt.Fscan(reader, &T); err != nil {
       return
   }
   // random generator seeded as in C++ mt19937(9)
   R := rand.New(rand.NewSource(9))
   for T > 0 {
       T--
       fmt.Fscan(reader, &k)
       n := 1 << (k + 1)
       // prefix xor array s[0..n]
       s := make([]uint64, n+1)
       s[0] = 0
       for i := 1; i <= n; i++ {
           var g uint64
           fmt.Fscan(reader, &g)
           s[i] = s[i-1] ^ g
       }
       mp := make(map[uint64]uint64)
       for {
           l := R.Intn(n + 1)
           r := R.Intn(n + 1)
           if l == r {
               continue
           }
           if l > r {
               l, r = r, l
           }
           vi := (uint64(l) << 32) | uint64(r)
           vv := s[l] ^ s[r]
           if ci, ok := mp[vv]; !ok {
               mp[vv] = vi
           } else {
               // found matching trait
               l1 := int(ci >> 32)
               r1 := int(uint32(ci))
               l2 := l
               r2 := r
               // ensure not same interval
               if l1 == l2 && r1 == r2 {
                   mp[vv] = vi
                   continue
               }
               ans := []int{l1, r1, l2, r2}
               sort.Ints(ans)
               // print as [a,b] [c,d]
               // converting to 1-based for left bounds
               a := ans[0] + 1
               b := ans[1]
               c := ans[2] + 1
               d := ans[3]
               fmt.Fprintf(writer, "%d %d %d %d\n", a, b, c, d)
               break
           }
       }
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   solve(reader, writer)
}
