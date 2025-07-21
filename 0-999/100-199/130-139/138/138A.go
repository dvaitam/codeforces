package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   total := 4 * n
   lines := make([]string, total)
   for i := 0; i < total; i++ {
       fmt.Fscan(reader, &lines[i])
   }
   suffixes := make([]string, total)
   for i, s := range lines {
       suffixes[i] = getSuffix(s, k)
   }
   possibleAABB, possibleABAB, possibleABBA := true, true, true
   allAAAA := true
   for i := 0; i < n; i++ {
       a := suffixes[4*i]
       b := suffixes[4*i+1]
       c := suffixes[4*i+2]
       d := suffixes[4*i+3]
       if a == "" || b == "" || c == "" || d == "" {
           fmt.Println("NO")
           return
       }
       if a == b && b == c && c == d {
           continue
       }
       allAAAA = false
       curAABB := (a == b && c == d)
       curABAB := (a == c && b == d)
       curABBA := (a == d && b == c)
       if !curAABB && !curABAB && !curABBA {
           fmt.Println("NO")
           return
       }
       if !curAABB {
           possibleAABB = false
       }
       if !curABAB {
           possibleABAB = false
       }
       if !curABBA {
           possibleABBA = false
       }
   }
   if allAAAA {
       fmt.Println("aaaa")
       return
   }
   count := 0
   var scheme string
   if possibleAABB {
       count++
       scheme = "aabb"
   }
   if possibleABAB {
       count++
       scheme = "abab"
   }
   if possibleABBA {
       count++
       scheme = "abba"
   }
   if count == 1 {
       fmt.Println(scheme)
   } else {
       fmt.Println("NO")
   }
}

// getSuffix returns the suffix of s starting from the k-th vowel from the end
// or empty string if there are fewer than k vowels.
func getSuffix(s string, k int) string {
   count := 0
   for i := len(s) - 1; i >= 0; i-- {
       switch s[i] {
       case 'a', 'e', 'i', 'o', 'u':
           count++
           if count == k {
               return s[i:]
           }
       }
   }
   return ""
}
