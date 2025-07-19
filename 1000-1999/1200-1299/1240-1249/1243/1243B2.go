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

   var T int
   if _, err := fmt.Fscan(reader, &T); err != nil {
       return
   }
   for T > 0 {
       T--
       var n int
       fmt.Fscan(reader, &n)
       var sStr, tStr string
       fmt.Fscan(reader, &sStr)
       fmt.Fscan(reader, &tStr)
       s := []byte(sStr)
       t := []byte(tStr)

       // check total character counts
       cnt := make([]int, 26)
       for i := 0; i < n; i++ {
           cnt[s[i]-'a']++
           cnt[t[i]-'a']++
       }
       possible := true
       for i := 0; i < 26; i++ {
           if cnt[i]%2 != 0 {
               possible = false
               break
           }
       }
       if !possible {
           fmt.Fprintln(writer, "No")
           continue
       }
       fmt.Fprintln(writer, "Yes")
       // perform swaps
       type pair struct{ x, y int }
       ops := make([]pair, 0, 2*n)
       for i := 0; i < n; i++ {
           if s[i] != t[i] {
               // try to find in s
               found := false
               for j := i + 1; j < n; j++ {
                   if s[j] == s[i] {
                       // swap s[j] with t[i]
                       ops = append(ops, pair{j, i})
                       s[j], t[i] = t[i], s[j]
                       found = true
                       break
                   }
               }
               if found {
                   continue
               }
               // find in t
               for j := i + 1; j < n; j++ {
                   if t[j] == s[i] {
                       // swap s[j] with t[j]
                       ops = append(ops, pair{j, j})
                       s[j], t[j] = t[j], s[j]
                       // swap s[j] with t[i]
                       ops = append(ops, pair{j, i})
                       s[j], t[i] = t[i], s[j]
                       break
                   }
               }
           }
       }
       // output operations
       fmt.Fprintln(writer, len(ops))
       for _, op := range ops {
           // convert to 1-based
           fmt.Fprintln(writer, op.x+1, op.y+1)
       }
   }
}
