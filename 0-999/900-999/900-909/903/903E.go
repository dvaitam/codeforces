package main

import (
   "bufio"
   "fmt"
   "os"
)

func possible(target, s string, canStay bool) bool {
   var diffs []int
   for i := 0; i < len(target); i++ {
       if target[i] != s[i] {
           diffs = append(diffs, i)
       }
   }
   if canStay && len(diffs) == 0 {
       return true
   }
   if len(diffs) != 2 {
       return false
   }
   if target[diffs[0]] != s[diffs[1]] || target[diffs[1]] != s[diffs[0]] {
       return false
   }
   return true
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var k, n int
   if _, err := fmt.Fscan(in, &k, &n); err != nil {
       return
   }
   uniq := make(map[string]struct{}, k)
   for i := 0; i < k; i++ {
       var s string
       fmt.Fscan(in, &s)
       uniq[s] = struct{}{}
   }
   if len(uniq) == 1 {
       // one unique string: swap first two chars
       var res string
       for s := range uniq {
           res = s
           break
       }
       if n >= 2 {
           b := []byte(res)
           b[0], b[1] = b[1], b[0]
           res = string(b)
       }
       fmt.Print(res)
       return
   }
   // compute char counts and canStay
   var canStay bool
   first := ""
   for s := range uniq {
       first = s
       break
   }
   cntRef := make([]int, 26)
   for i := 0; i < len(first); i++ {
       idx := first[i] - 'a'
       cntRef[idx]++
       if cntRef[idx] > 1 {
           canStay = true
       }
   }
   // verify same multiset
   for s := range uniq {
       cnt := make([]int, 26)
       for i := 0; i < len(s); i++ {
           cnt[s[i]-'a']++
       }
       for j := 0; j < 26; j++ {
           if cnt[j] != cntRef[j] {
               fmt.Print(-1)
               return
           }
       }
   }
   // build vector of unique strings
   strVec := make([]string, 0, len(uniq))
   for s := range uniq {
       strVec = append(strVec, s)
   }
   // find diffs between first two
   a := strVec[0]
   bstr := strVec[1]
   var diffs []int
   for i := 0; i < len(a); i++ {
       if a[i] != bstr[i] {
           diffs = append(diffs, i)
       }
   }
   if len(diffs) < 2 || len(diffs) > 4 {
       fmt.Print(-1)
       return
   }
   // generate targets
   targets := make(map[string]struct{})
   if canStay {
       targets[a] = struct{}{}
   }
   // swap any two diff positions
   bts := []byte(a)
   for i := 0; i < len(diffs); i++ {
       for j := i + 1; j < len(diffs); j++ {
           x, y := diffs[i], diffs[j]
           bts[x], bts[y] = bts[y], bts[x]
           targets[string(bts)] = struct{}{}
           bts[x], bts[y] = bts[y], bts[x]
       }
   }
   // filter targets
   for _, s := range strVec[1:] {
       for t := range targets {
           if !possible(t, s, canStay) {
               delete(targets, t)
           }
       }
   }
   for t := range targets {
       fmt.Print(t)
       return
   }
   fmt.Print(-1)
}
