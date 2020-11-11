package filters

func init() {
	AddFilter("s", &FilterS{})
}

type FilterS struct {
}

func (p *FilterS) AddExpression(exp, key, s string) string {
	if exp != "" {
		exp += " && "
	}
	if s[len(s)-1] != ',' {
		s += ","
	}
	s += "-1"
	exp += key + " IN (" + s + ")"
	return exp
}
