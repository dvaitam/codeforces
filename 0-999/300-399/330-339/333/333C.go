package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var k int
   var m int
   if _, err := fmt.Fscan(reader, &k, &m); err != nil {
       return
   }
   ks := strconv.Itoa(k)
   d := len(ks)
   // total digits 8
   P := 8 - d
   results := make([]string, 0, m)
   seen := make(map[string]struct{}, m)
   // suffix-case: prefix zeros with zero-digit in prefix, then + k as suffix
   // prefix of length P
   // generate prefixes in increasing order until m
   limit := 1
   for i := 0; i < P; i++ {
       limit *= 10
   }
   // first, prefix-case: prefix*... + k -> ticket = prefix + ks
   for i := 0; i < limit && len(results) < m; i++ {
       // format i with leading zeros to length P
       pre := fmt.Sprintf("%0*d", P, i)
       // check if contains '0'
       if !containsZero(pre) {
           continue
       }
       ticket := pre + ks
       if _, ok := seen[ticket]; ok {
           continue
       }
       seen[ticket] = struct{}{}
       results = append(results, ticket)
   }
   // second, suffix-case: ks + suffix with zero
   for i := 0; i < limit && len(results) < m; i++ {
       suf := fmt.Sprintf("%0*d", P, i)
       if !containsZero(suf) {
           continue
       }
       ticket := ks + suf
       if _, ok := seen[ticket]; ok {
           continue
       }
       seen[ticket] = struct{}{}
       results = append(results, ticket)
   }
   // output
   for _, t := range results {
       fmt.Fprintln(writer, t)
   }
}

// containsZero reports whether s contains character '0'
func containsZero(s string) bool {
   for i := 0; i < len(s); i++ {
       if s[i] == '0' {
           return true
       }
   }
   return false
}
