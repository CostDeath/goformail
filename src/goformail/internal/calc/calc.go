package calc

func Add[T int | float32 | float64](a, b T) T {
	return a + b
}

func Sub[T int | float32 | float64](a, b T) T {
	return a - b
}

func Mul[T int | float32 | float64](a, b T) T {
	return a * b
}

func Div[T int | float32 | float64](a, b T) T {
	return a / b
}
