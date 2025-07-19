package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   v := make([]int, 0, n)
   var prev, x int
   // read first element
   fmt.Fscan(in, &prev)
   v = append(v, prev)
   for i := 1; i < n; i++ {
       fmt.Fscan(in, &x)
       if x != prev {
           v = append(v, x)
       }
       prev = x
   }
   // count frequency of each block value
   freq := make(map[int]int)
   var maxFreq int
   for _, val := range v {
       freq[val]++
       if freq[val] > maxFreq {
           maxFreq = freq[val]
       }
   }
   blocks := len(v)
   // minimum operations is total blocks minus most frequent block count
   ans := blocks - maxFreq
   fmt.Fprint(out, ans)
}
