package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   s, _ := reader.ReadString('\n')
   t, _ := reader.ReadString('\n')
   s = strings.TrimSpace(s)
   t = strings.TrimSpace(t)
   n := len(s)
   m := len(t)
   // prefix matches: pref[i] = length of t matched in s[0..i-1]
   pref := make([]int, n+1)
   for i := 1; i <= n; i++ {
       pref[i] = pref[i-1]
       if pref[i] < m && s[i-1] == t[pref[i]] {
           pref[i]++
       }
   }
   // suffix matches: suff[i] = length of suffix of t matched in s[i-1..n-1]
   suff := make([]int, n+2)
   for i := n; i >= 1; i-- {
       suff[i] = suff[i+1]
       if suff[i+1] < m && s[i-1] == t[m-1-suff[i+1]] {
           suff[i]++
       }
   }
   // positions of each character in t (1-based)
   pos := make([][]int, 26)
   for j := 0; j < m; j++ {
       c := t[j] - 'a'
       pos[c] = append(pos[c], j+1)
   }
   // check each position in s
   ok := true
   for i := 1; i <= n; i++ {
       c := s[i-1] - 'a'
       // allowable j range: B <= j <= A
       A := pref[i-1] + 1
       B := m - suff[i+1]
       found := false
       if B <= A {
           arr := pos[c]
           idx := sort.Search(len(arr), func(k int) bool { return arr[k] >= B })
           if idx < len(arr) && arr[idx] <= A {
               found = true
           }
       }
       if !found {
           ok = false
           break
       }
   }
   if ok {
       fmt.Println("Yes")
   } else {
       fmt.Println("No")
   }
}
