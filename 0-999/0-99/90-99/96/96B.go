package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   s, _ := reader.ReadString('\n')
   s = strings.TrimSpace(s)
   nLen := len(s)
   // If length is odd, next even length minimal super lucky
   if nLen%2 == 1 {
       k := (nLen + 1) / 2
       fmt.Println(strings.Repeat("4", k) + strings.Repeat("7", k))
       return
   }
   // even length: try same length
   L := nLen
   k := L / 2
   found := false
   var out string
   var dfs func(pos, c4, c7 int, cur []byte)
   dfs = func(pos, c4, c7 int, cur []byte) {
       if found {
           return
       }
       if pos == L {
           cand := string(cur)
           if strings.Compare(cand, s) >= 0 {
               out = cand
               found = true
           }
           return
       }
       if c4 < k {
           cur[pos] = '4'
           dfs(pos+1, c4+1, c7, cur)
           if found {
               return
           }
       }
       if c7 < k {
           cur[pos] = '7'
           dfs(pos+1, c4, c7+1, cur)
           if found {
               return
           }
       }
   }
   cur := make([]byte, L)
   dfs(0, 0, 0, cur)
   if found {
       fmt.Println(out)
       return
   }
   // no candidate of this length, go to next even length
   k++
   fmt.Println(strings.Repeat("4", k) + strings.Repeat("7", k))
}
