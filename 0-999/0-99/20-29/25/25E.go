package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

// compute prefix function (pi) for KMP
func prefixFunc(p []byte) []int {
   n := len(p)
   pi := make([]int, n)
   for i := 1; i < n; i++ {
       j := pi[i-1]
       for j > 0 && p[i] != p[j] {
           j = pi[j-1]
       }
       if p[i] == p[j] {
           j++
       }
       pi[i] = j
   }
   return pi
}

// check if pattern b is substring of text a
func contains(a, b string) bool {
   la, lb := len(a), len(b)
   if lb > la {
       return false
   }
   pat := []byte(b)
   pi := prefixFunc(pat)
   j := 0
   for i := 0; i < la; i++ {
       for j > 0 && a[i] != b[j] {
           j = pi[j-1]
       }
       if a[i] == b[j] {
           j++
           if j == lb {
               return true
           }
       }
   }
   return false
}

// overlap length: longest suffix of a matching prefix of b
func overlap(a, b string) int {
   la, lb := len(a), len(b)
   // only need suffix of a of length <= lb
   start := 0
   if la > lb {
       start = la - lb
   }
   pat := []byte(b)
   pi := prefixFunc(pat)
   j := 0
   for i := start; i < la; i++ {
       for j > 0 && a[i] != b[j] {
           j = pi[j-1]
       }
       if a[i] == b[j] {
           j++
           if j == lb {
               // full b matched; can't get longer
               break
           }
       }
   }
   return j
}

// merge two strings with maximal overlap
func merge(a, b string) string {
   k := overlap(a, b)
   return a + b[k:]
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   strs := make([]string, 3)
   for i := 0; i < 3; i++ {
       line, err := reader.ReadString('\n')
       if err != nil && err != os.EOF {
           fmt.Fprintln(os.Stderr, "read error:", err)
           return
       }
       strs[i] = strings.TrimSpace(line)
   }
   // filter out strings contained in others or duplicates
   keep := make([]string, 0, 3)
   for i, s := range strs {
       skip := false
       for j, t := range strs {
           if i == j {
               continue
           }
           if contains(t, s) {
               // skip if t longer or equal but earlier index
               if len(t) > len(s) || (len(t) == len(s) && j < i) {
                   skip = true
                   break
               }
           }
       }
       if !skip {
           keep = append(keep, s)
       }
   }
   var result int
   n := len(keep)
   if n == 0 {
       result = len(strs[0])
   } else if n == 1 {
       result = len(keep[0])
   } else if n == 2 {
       a, b := keep[0], keep[1]
       // two orders
       l1 := len(a) + len(b) - overlap(a, b)
       l2 := len(a) + len(b) - overlap(b, a)
       if l1 < l2 {
           result = l1
       } else {
           result = l2
       }
   } else {
       // 3 strings, try all permutations
       perm := []int{0, 1, 2}
       result = -1
       // generate permutations of 0,1,2
       var gen func(int)
       gen = func(idx int) {
           if idx == 3 {
               a, b, c := keep[perm[0]], keep[perm[1]], keep[perm[2]]
               t := merge(a, b)
               t = merge(t, c)
               l := len(t)
               if result < 0 || l < result {
                   result = l
               }
               return
           }
           for i := idx; i < 3; i++ {
               perm[idx], perm[i] = perm[i], perm[idx]
               gen(idx + 1)
               perm[idx], perm[i] = perm[i], perm[idx]
           }
       }
       gen(0)
   }
   fmt.Println(result)
}
