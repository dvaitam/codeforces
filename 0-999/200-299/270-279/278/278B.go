package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   titles := make([]string, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &titles[i])
   }
   // find shortest original title
   for length := 1; ; length++ {
       subs := make(map[string]struct{})
       // collect all substrings of given length
       for _, t := range titles {
           if len(t) < length {
               continue
           }
           for j := 0; j+length <= len(t); j++ {
               subs[t[j:j+length]] = struct{}{}
           }
       }
       // search lex smallest string of this length not in subs
       buf := make([]byte, length)
       var found string
       var dfs func(pos int) bool
       dfs = func(pos int) bool {
           if pos == length {
               s := string(buf)
               if _, exists := subs[s]; !exists {
                   found = s
                   return true
               }
               return false
           }
           for c := byte('a'); c <= 'z'; c++ {
               buf[pos] = c
               if dfs(pos + 1) {
                   return true
               }
           }
           return false
       }
       if dfs(0) {
           fmt.Println(found)
           return
       }
   }
}
