package mkpub

import (
	"html/template"
)

var TemplateFuncs = template.FuncMap{
	"add": add,
	"sub": sub,
	"mul": mul,
	"div": div,
	"inc": inc,
	"dec": dec,
}

func add(nums ...float64) float64 {
	var result float64
	for i, n := range nums {
		if i == 0 {
			result = n
			continue
		}
		result += n
	}
	return result
}

func sub(nums ...float64) float64 {
	var result float64
	for i, n := range nums {
		if i == 0 {
			result = n
			continue
		}
		result -= n
	}
	return result
}

func mul(nums ...float64) float64 {
	var result float64
	for i, n := range nums {
		if i == 0 {
			result = n
			continue
		}
		result *= n
	}
	return result
}

func div(nums ...float64) float64 {
	var result float64
	for i, n := range nums {
		if i == 0 {
			result = n
			continue
		}
		result /= n
	}
	return result
}

func inc(n float64) float64 {
	return n + 1
}

func dec(n float64) float64 {
	return n - 1
}
