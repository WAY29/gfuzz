package filters

import "strconv"

func init() {
	AddFilter("h", &FilterH{})
}

type FilterH struct {
}

func (p *FilterH) AddExpression(exp, key, s string) string {
	if exp != "" {
		exp += " && "
	}
	if key == "x" {
		s = strconv.Quote(s)
		exp += "!(" + key + "=~ " + s + ")"
		return exp
	}
	if s[len(s)-1] != ',' {
		s += ","
	}
	s += "-1"

	exp += "!(" + key + " IN (" + s + "))"
	return exp
}
