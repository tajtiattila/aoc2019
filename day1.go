package main

import "log"

func day1() {
	ch, cherr := daydataInts(1)
	sum1 := 0
	sum2 := 0
	for n := range ch {
		sum1 += fuel(n)
		sum2 += fuelx(n)
	}

	if err := <-cherr; err != nil {
		log.Fatal(err)
	}
	log.Println("day1a:", sum1)
	log.Println("day1b:", sum2)
}

func fuel(n int) int {
	f := (n / 3) - 2
	if f > 0 {
		return f
	}
	return 0
}

func fuelx(n int) int {
	x := 0
	for n > 0 {
		m := fuel(n)
		x += m
		n = m
	}
	return x
}
