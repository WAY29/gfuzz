package filters

import (
	"errors"
	"fmt"

	"github.com/Knetic/govaluate"
)

var FiltersArray map[string]Filter = make(map[string]Filter)

// * Filters interface, use go doc filters.FilterBase for more information
type Filter interface {
	AddExpression(exp, key, s string) string
}

/*
* Filter struct, must implement AddExpression method

* To see more information, please run gfuzz filters
 */
type FilterBase struct {
}

func (p *FilterBase) AddExpression(exp, key, s string) string {
	return exp
}

func AddFilter(name string, f Filter) {
	FiltersArray[name] = f
}

func GetFilter(name string) Filter {
	f, ok := FiltersArray[name]
	if !ok {
		return &FilterBase{}
	} else {
		return f
	}
}

// * 
func FinalFilter(expression string, parameters map[string]interface{}) (isPass bool, err error) {
	var isPassI interface{}
	err = errors.New("Filter invaild")
	nerr := errors.New("")
	defer func() {
		if r := recover(); r != nil {
			isPass = false
			fmt.Println("Test", r)
			err = errors.New("Filter invaild")
		}
	}()

	exp, _ := govaluate.NewEvaluableExpression(expression)
	isPassI, nerr = exp.Evaluate(parameters)
	isPass, ok := isPassI.(bool)
	if nerr != nil || !ok {
		if nerr != nil {
			err = nerr
		}
		isPass = false
	} else if nerr == nil {
		err = nerr
	}

	return isPass, err
}
