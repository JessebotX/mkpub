package mkpub

import (
	"errors"
	"html/template"
	"strconv"
)

var (
	ErrNotFloat = errors.New("arguments must be rational numbers")
)

var TemplateFuncs = template.FuncMap{
	"add": add,
	"sub": sub,
	"mul": mul,
	"div": div,
	"inc": inc,
	"dec": dec,
}

func convFloat(num any) (float64, error) {
	numStr, ok := num.(string)
	if !ok {
		return -1, ErrNotFloat
	}

	n, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return -1, ErrNotFloat
	}

	return n, nil
}

func add(nums ...any) (float64, error) {
	var result float64
	for i, num := range nums {
		n, err := convFloat(num)
		if err != nil {
			return -1, err
		}

		if i == 0 {
			result = n
			continue
		}
		result += n
	}
	return result, nil
}

func sub(nums ...any) (float64, error) {
	var result float64
	for i, num := range nums {
		n, err := convFloat(num)
		if err != nil {
			return -1, err
		}

		if i == 0 {
			result = n
			continue
		}
		result -= n
	}
	return result, nil
}

func mul(nums ...any) (float64, error) {
	var result float64
	for i, num := range nums {
		n, err := convFloat(num)
		if err != nil {
			return -1, err
		}

		if i == 0 {
			result = n
			continue
		}
		result *= n
	}
	return result, nil
}

func div(nums ...any) (float64, error) {
	var result float64
	for i, num := range nums {
		n, err := convFloat(num)
		if err != nil {
			return -1, err
		}

		if i == 0 {
			result = n
			continue
		}
		result /= n
	}
	return result, nil
}

func inc(num any) (float64, error) {
	n, err := convFloat(num)
	if err != nil {
		return -1, err
	}

	return n + 1, nil
}

func dec(num any) (float64, error) {
	n, err := convFloat(num)
	if err != nil {
		return -1, err
	}

	return n - 1, nil
}
