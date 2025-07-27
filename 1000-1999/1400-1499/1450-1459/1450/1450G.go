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

   var n, a, b int
   if _, err := fmt.Fscan(reader, &n, &a, &b); err != nil {
       return
   }
   var s string
   fmt.Fscan(reader, &s)

   // positions and counts for each letter
   const alpha = 26
   fpos := make([]int, alpha)
   lpos := make([]int, alpha)
   cnt := make([]int, alpha)
   for i := 0; i < alpha; i++ {
       fpos[i] = n + 1
   }
   for i, ch := range s {
       idx := int(ch - 'a')
       cnt[idx]++
       // positions are 1-based
       pos := i + 1
       if pos < fpos[idx] {
           fpos[idx] = pos
       }
       if pos > lpos[idx] {
           lpos[idx] = pos
       }
   }

   // check for each letter x: if it can be converted (density condition)
   good := make([]bool, alpha)
   for i := 0; i < alpha; i++ {
       if cnt[i] == 0 {
           continue
       }
       // length of interval covering all i
       length := lpos[i] - fpos[i] + 1
       // condition: a * length <= b * cnt[i]
       if a*length <= b*cnt[i] {
           good[i] = true
       }
   }

   // find obtainable target letters
   var ans []rune
   for i := 0; i < alpha; i++ {
       if cnt[i] == 0 {
           continue
       }
       ok := true
       // all other letters must be convertible to this target
       for j := 0; j < alpha; j++ {
           if j == i || cnt[j] == 0 {
               continue
           }
           if !good[j] {
               ok = false
               break
           }
       }
       if ok {
           ans = append(ans, rune('a'+i))
       }
   }

   // output
   fmt.Fprint(writer, len(ans))
   if len(ans) > 0 {
       fmt.Fprint(writer, " ")
       for k, ch := range ans {
           if k > 0 {
               fmt.Fprint(writer, " ")
           }
           fmt.Fprint(writer, string(ch))
       }
   }
   fmt.Fprintln(writer)
}
