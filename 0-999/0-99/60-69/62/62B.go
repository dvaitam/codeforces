package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func abs(x int) int {
   if x < 0 {
       return -x
   }
   return x
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   var s string
   fmt.Fscan(reader, &s)
   // positions of each letter in s (1-based)
   pos := make([][]int, 26)
   for i, ch := range s {
       idx := ch - 'a'
       pos[idx] = append(pos[idx], i+1)
   }

   // process each potential address
   for qi := 0; qi < n; qi++ {
       var c string
       fmt.Fscan(reader, &c)
       clen := len(c)
       var f int64
       for i, ch := range c {
           lst := pos[ch-'a']
           if len(lst) == 0 {
               f += int64(clen)
           } else {
               // desired position is i+1
               target := i + 1
               idx := sort.Search(len(lst), func(j int) bool { return lst[j] >= target })
               // find minimal distance
               d := int(1e9)
               if idx < len(lst) {
                   d0 := abs(lst[idx] - target)
                   if d0 < d {
                       d = d0
                   }
               }
               if idx > 0 {
                   d0 := abs(lst[idx-1] - target)
                   if d0 < d {
                       d = d0
                   }
               }
               f += int64(d)
           }
       }
       fmt.Fprintln(writer, f)
   }
}
