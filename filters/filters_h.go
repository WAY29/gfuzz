package filters

func init() {
	AddFilter("h", &FilterH{})
}

type FilterH struct {
}

func (p *FilterH) AddExpression(exp, key, s string) string {
	if exp != "" {
		exp += " && "
	}
	if s[len(s)-1] != ',' {
		s += ","
	}
	s += "-1"
	exp += "!(" + key + " IN (" + s + "))"
	return exp
}
