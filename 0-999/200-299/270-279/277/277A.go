package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   scanner := bufio.NewScanner(os.Stdin)
   scanner.Split(bufio.ScanWords)
   nextInt := func() int {
       if !scanner.Scan() {
           panic("no more tokens")
       }
       v := 0
       for _, b := range scanner.Bytes() {
           v = v*10 + int(b-'0')
       }
       return v
   }

   n := nextInt()
   m := nextInt()
   empLangs := make([][]int, n)
   langEmps := make([][]int, m+1)
   zerosCount := 0
   totalKnown := 0
   for i := 0; i < n; i++ {
       ki := nextInt()
       if ki == 0 {
           zerosCount++
       }
       totalKnown += ki
       empLangs[i] = make([]int, ki)
       for j := 0; j < ki; j++ {
           lang := nextInt()
           empLangs[i][j] = lang
           langEmps[lang] = append(langEmps[lang], i)
       }
   }

   visited := make([]bool, n)
   var dfs func(int)
   dfs = func(u int) {
       visited[u] = true
       for _, lang := range empLangs[u] {
           for _, v := range langEmps[lang] {
               if !visited[v] {
                   dfs(v)
               }
           }
       }
   }

   compLang := 0
   for i := 0; i < n; i++ {
       if !visited[i] && len(empLangs[i]) > 0 {
           dfs(i)
           compLang++
       }
   }

   if totalKnown == 0 {
       fmt.Println(n)
   } else {
       res := zerosCount + compLang - 1
       fmt.Println(res)
   }
}
