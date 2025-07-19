package main

import (
   "bufio"
   "fmt"
   "math/bits"
   "os"
   "sort"
   "strconv"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   sort.Ints(a)
   // remove duplicates
   b := make([]int, 0, n)
   for i, v := range a {
       if i == 0 || v != a[i-1] {
           b = append(b, v)
       }
   }
   a = b
   n = len(a)
   // split into low and high parts
   aLow := make([]int, n)
   aHigh := make([]int, n)
   for i, v := range a {
       aLow[i] = v & 0x7fff
       aHigh[i] = v >> 15
   }
   mp := make(map[string]int, 1<<15)
   d := make([]int, n)
   var sb strings.Builder
   // first half masks
   for mask := 0; mask < (1 << 15); mask++ {
       for i := 0; i < n; i++ {
           d[i] = bits.OnesCount(uint(aLow[i]^mask))
       }
       base := d[0]
       for i := 1; i < n; i++ {
           d[i] -= base
       }
       d[0] = 0
       sb.Reset()
       for i := 0; i < n; i++ {
           sb.WriteString(strconv.Itoa(d[i]))
           sb.WriteByte(',')
       }
       mp[sb.String()] = mask
   }
   // second half masks
   for mask := 0; mask < (1 << 15); mask++ {
       for i := 0; i < n; i++ {
           cnt := bits.OnesCount(uint(aHigh[i] ^ mask))
           d[i] = 30 - cnt
       }
       base := d[0]
       for i := 1; i < n; i++ {
           d[i] -= base
       }
       d[0] = 0
       sb.Reset()
       for i := 0; i < n; i++ {
           sb.WriteString(strconv.Itoa(d[i]))
           sb.WriteByte(',')
       }
       key := sb.String()
       if v, ok := mp[key]; ok {
           res := ((mask << 15) ^ v)
           fmt.Println(res)
           return
       }
   }
   fmt.Println(-1)
}
