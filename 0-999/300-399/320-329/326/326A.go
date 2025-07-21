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

   var s string
   var n int
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }

   // count frequencies in s
   freq := make([]int, 26)
   for _, ch := range s {
       freq[ch-'a']++
   }
   // collect non-zero frequencies
   var freqs []int
   for _, f := range freq {
       if f > 0 {
           freqs = append(freqs, f)
       }
   }
   distinct := len(freqs)
   if distinct > n {
       fmt.Fprintln(writer, -1)
       return
   }
   // find minimal k
   maxf := 0
   for _, f := range freqs {
       if f > maxf {
           maxf = f
       }
   }
   ansK := -1
   for k := 1; k <= maxf; k++ {
       need := 0
       for _, f := range freqs {
           need += (f + k - 1) / k
       }
       if need <= n {
           ansK = k
           break
       }
   }
   if ansK == -1 {
       fmt.Fprintln(writer, -1)
       return
   }
   // build sheet string
   t := make([]byte, 0, n)
   used := 0
   for i, f := range freq {
       if f > 0 {
           cnt := (f + ansK - 1) / ansK
           for j := 0; j < cnt; j++ {
               t = append(t, byte('a'+i))
           }
           used += cnt
       }
   }
   // fill remaining slots
   for used < n {
       t = append(t, 'a')
       used++
   }

   fmt.Fprintln(writer, ansK)
   fmt.Fprintln(writer, string(t))
}
