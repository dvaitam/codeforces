package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   fmt.Fscan(reader, &n, &m)
   a := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   b := make([]int, m+1)
   for i := 1; i <= m; i++ {
       fmt.Fscan(reader, &b[i])
   }
   // Build P: value->position in b
   P := make([]int, m+2)
   for i := 1; i <= m; i++ {
       P[b[i]] = i
   }
   // Text T is P[1..m]
   T := make([]int, m+1)
   for v := 1; v <= m; v++ {
       T[v] = P[v]
   }
   // Compute prev smaller/larger for pattern a
   plPat := make([]int, n+1)
   glPat := make([]int, n+1)
   // prev smaller
   stack := make([]int, 0, n)
   for i := 1; i <= n; i++ {
       for len(stack) > 0 && a[stack[len(stack)-1]] >= a[i] {
           stack = stack[:len(stack)-1]
       }
       if len(stack) == 0 {
           plPat[i] = 0
       } else {
           plPat[i] = stack[len(stack)-1]
       }
       stack = append(stack, i)
   }
   // prev larger
   stack = stack[:0]
   for i := 1; i <= n; i++ {
       for len(stack) > 0 && a[stack[len(stack)-1]] <= a[i] {
           stack = stack[:len(stack)-1]
       }
       if len(stack) == 0 {
           glPat[i] = 0
       } else {
           glPat[i] = stack[len(stack)-1]
       }
       stack = append(stack, i)
   }
   // Compute prev smaller/larger for text T
   plTxt := make([]int, m+1)
   glTxt := make([]int, m+1)
   // prev smaller in T
   stack = stack[:0]
   for i := 1; i <= m; i++ {
       for len(stack) > 0 && T[stack[len(stack)-1]] >= T[i] {
           stack = stack[:len(stack)-1]
       }
       if len(stack) == 0 {
           plTxt[i] = 0
       } else {
           plTxt[i] = stack[len(stack)-1]
       }
       stack = append(stack, i)
   }
   // prev larger in T
   stack = stack[:0]
   for i := 1; i <= m; i++ {
       for len(stack) > 0 && T[stack[len(stack)-1]] <= T[i] {
           stack = stack[:len(stack)-1]
       }
       if len(stack) == 0 {
           glTxt[i] = 0
       } else {
           glTxt[i] = stack[len(stack)-1]
       }
       stack = append(stack, i)
   }
   // Build pi for pattern
   pi := make([]int, n+1)
   pi[1] = 0
   for i := 2; i <= n; i++ {
       j := pi[i-1]
       for j > 0 && !matchPat(plPat, glPat, plPat, glPat, i, j+1) {
           j = pi[j]
       }
       if matchPat(plPat, glPat, plPat, glPat, i, j+1) {
           j++
       }
       pi[i] = j
   }
   // KMP over text
   q := 0
   count := 0
   for i := 1; i <= m; i++ {
       for q > 0 && !matchTxt(plPat, glPat, plTxt, glTxt, i, q+1) {
           q = pi[q]
       }
       if matchTxt(plPat, glPat, plTxt, glTxt, i, q+1) {
           q++
       }
       if q == n {
           count++
           q = pi[q]
       }
   }
   // Output count
   w := bufio.NewWriter(os.Stdout)
   fmt.Fprintln(w, count)
   w.Flush()
}

// match pattern against pattern for pi building
func matchPat(plBase, glBase, plTxt, glTxt []int, i, j int) bool {
   // here plBase/glBase are plPat/glPat, and plTxt/glTxt also plPat/glPat
   // pattern window start in pattern: s = i - j
   s := i - j
   // smaller
   if plBase[j] == 0 {
       if plTxt[i] > s {
           return false
       }
   } else {
       if plTxt[i] != s + (j - plBase[j]) {
           return false
       }
   }
   // larger
   if glBase[j] == 0 {
       if glTxt[i] > s {
           return false
       }
   } else {
       if glTxt[i] != s + (j - glBase[j]) {
           return false
       }
   }
   return true
}

// match pattern against text
func matchTxt(plPat, glPat, plTxt, glTxt []int, i, j int) bool {
   s := i - j
   // smaller
   if plPat[j] == 0 {
       if plTxt[i] > s {
           return false
       }
   } else {
       if plTxt[i] != s + (j - plPat[j]) {
           return false
       }
   }
   // larger
   if glPat[j] == 0 {
       if glTxt[i] > s {
           return false
       }
   } else {
       if glTxt[i] != s + (j - glPat[j]) {
           return false
       }
   }
   return true
}
