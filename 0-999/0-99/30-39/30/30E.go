package main

import (
   "fmt"
)

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

// palindrom_extend computes odd-length palindrome radii: extend[i] = max radius around i
func palindrom_extend(s []int, extend []int) {
   n := len(s)
   idx, reach := 0, 0
   extend[0] = 0
   for i := 1; i < n; i++ {
       l := 0
       if reach > i {
           j := 2*idx - i
           if j >= 0 && j < n {
               l = min(reach-i, extend[j])
           }
       }
       // expand around i
       for i-l-1 >= 0 && i+l+1 < n && s[i-l-1] == s[i+l+1] {
           l++
       }
       extend[i] = l
       if i+l > reach {
           reach = i + l
           idx = i
       }
   }
}

func main() {
   var line string
   if _, err := fmt.Scan(&line); err != nil {
       return
   }
   n := len(line)
   s := make([]int, n)
   for i, ch := range line {
       s[i] = int(ch)
   }
   // compute palindrome radii
   extend := make([]int, n)
   palindrom_extend(s, extend)
   // reverse s
   r := make([]int, n)
   for i := 0; i < n; i++ {
       r[i] = s[n-1-i]
   }
   // build KMP fail for r
   fail := make([]int, n)
   fail[0] = -1
   p := -1
   for i := 1; i < n; i++ {
       for p != -1 && r[p+1] != r[i] {
           p = fail[p]
       }
       if r[p+1] == r[i] {
           p++
       }
       fail[i] = p
   }
   // compute match and pos
   match := make([]int, n)
   pos := make([]int, n)
   p = -1
   for i := 0; i < n; i++ {
       for p != -1 && r[p+1] != s[i] {
           p = fail[p]
       }
       if r[p+1] == s[i] {
           p++
       }
       // ensure r prefix doesn't exceed suffix
       for p+1+i+1 > n {
           p = fail[p]
       }
       match[i] = p + 1
       pos[i] = i
       if i > 0 && match[i] < match[i-1] {
           match[i] = match[i-1]
           pos[i] = pos[i-1]
       }
   }
   // search best decomposition
   best := 0
   x1, l1, x2, l2, x3, l3 := 0, 0, 0, 0, 0, 0
   for i := 0; i < n; i++ {
       rad := extend[i]
       ll := i - rad - 1
       rr := i + rad + 1
       pl, sl := 0, 0
       pstart, sstart := 0, 0
       if ll >= 0 && rr < n {
           ln := match[ll]
           if ln > n-rr {
               ln = n - rr
           }
           pstart = pos[ll] - ln + 1
           sstart = n - ln
           pl = ln
           sl = ln
       }
       cur := pl + sl + 2*rad + 1
       if cur > best {
           best = cur
           x1, l1 = pstart, pl
           x2, l2 = ll+1, 2*rad+1
           x3, l3 = sstart, sl
       }
   }
   // output
   ans := 0
   if l1 > 0 {
       ans++
   }
   if l2 > 0 {
       ans++
   }
   if l3 > 0 {
       ans++
   }
   fmt.Println(ans)
   if l1 > 0 {
       fmt.Printf("%d %d\n", x1+1, l1)
   }
   if l2 > 0 {
       fmt.Printf("%d %d\n", x2+1, l2)
   }
   if l3 > 0 {
       fmt.Printf("%d %d\n", x3+1, l3)
   }
}
