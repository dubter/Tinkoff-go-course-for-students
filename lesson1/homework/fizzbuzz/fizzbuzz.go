package fizzbuzz

func FizzBuzz(i int) string {
	if i%15 == 0 {
		return "FizzBuzz"
	}
	if i%5 == 0 {
		return "Buzz"
	}
	return "Fizz"
}
