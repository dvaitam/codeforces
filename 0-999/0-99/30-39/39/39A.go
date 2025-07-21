package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type Task struct {
   w    int   // weight = sign * coefficient
   c    int   // coefficient
   pre  bool  // true for ++a, false for a++
   sign int   // +1 or -1
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var a0 int
   var expr string
   if _, err := fmt.Fscan(reader, &a0, &expr); err != nil {
       return
   }
   // parse expression into tasks
   tasks := make([]Task, 0, len(expr))
   i := 0
   n := len(expr)
   sign := 1
   for i < n {
       // parse coefficient if present
       coef := 1
       start := i
       tmp := 0
       for i < n && expr[i] >= '0' && expr[i] <= '9' {
           tmp = tmp*10 + int(expr[i]-'0')
           i++
       }
       if i < n && expr[i] == '*' && i > start {
           coef = tmp
           i++
       } else {
           // no valid coefficient, reset
           i = start
           coef = 1
       }
       // parse increment
       pre := false
       if i+1 < n && expr[i] == '+' && expr[i+1] == '+' {
           pre = true
           i += 2
           // expect 'a'
           i++
       } else if i < n && expr[i] == 'a' {
           // postfix, expect 'a++'
           i++
           if i+1 < n && expr[i] == '+' && expr[i+1] == '+' {
               pre = false
               i += 2
           }
       }
       tasks = append(tasks, Task{w: sign * coef, c: coef, pre: pre, sign: sign})
       // parse operator
       if i < n {
           if expr[i] == '+' {
               sign = 1
           } else if expr[i] == '-' {
               sign = -1
           }
           i++
       }
   }
   // sort by weight ascending
   sort.Slice(tasks, func(i, j int) bool {
       return tasks[i].w < tasks[j].w
   })
   // simulate
   a := a0
   var total int64
   for _, t := range tasks {
       if t.pre {
           a++
           total += int64(t.sign) * int64(t.c) * int64(a)
       } else {
           total += int64(t.sign) * int64(t.c) * int64(a)
           a++
       }
   }
   fmt.Println(total)
}
