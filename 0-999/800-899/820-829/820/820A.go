package main

import "fmt"

func main() {
	var c, v0, v1, a, l int
	if _, err := fmt.Scan(&c, &v0, &v1, &a, &l); err != nil {
		return
	}
	days := 1
	read := v0
	speed := v0
	for read < c {
		days++
		speed += a
		if speed > v1 {
			speed = v1
		}
		read += speed - l
	}
	fmt.Println(days)
}
