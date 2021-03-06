package svg

import (
	"fmt"
	"strconv"

	mt "github.com/rustyoz/Mtransform"
	gl "github.com/rustyoz/genericlexer"
)

func parseNumber(i gl.Item) (float64, error) {
	var n float64
	var ok error
	if i.Type == gl.ItemNumber {
		n, ok = strconv.ParseFloat(i.Value, 64)
		if ok != nil {
			return n, fmt.Errorf("Error passing number %s", ok)
		}
	}
	return n, nil
}

func parseTuple(l *gl.Lexer) (Tuple, error) {
	t := Tuple{}

	l.ConsumeWhiteSpace()

	ni := l.NextItem()
	if ni.Type == gl.ItemNumber {
		n, ok := strconv.ParseFloat(ni.Value, 64)
		if ok != nil {
			return t, fmt.Errorf("Error passing number %s", ok)
		}
		t[0] = n
	} else {
		return t, fmt.Errorf("Error passing Tuple expected Number got %v", ni)
	}

	if l.PeekItem().Type == gl.ItemWSP || l.PeekItem().Type == gl.ItemComma {
		l.NextItem()
	}
	ni = l.NextItem()
	if ni.Type == gl.ItemNumber {
		n, ok := strconv.ParseFloat(ni.Value, 64)
		if ok != nil {
			return t, fmt.Errorf("Error passing Number %s", ok)
		}
		t[1] = n
	} else {
		return t, fmt.Errorf("Error passing Tuple expected Number got: %v", ni)
	}

	return t, nil
}

func parseTransform(tstring string) (mt.Transform, error) {
	var tm mt.Transform
	lexer, _ := gl.Lex("tlexer", tstring)
	for {
		i := lexer.NextItem()
		switch i.Type {
		case gl.ItemEOS:
			break
		case gl.ItemWord:
			switch i.Value {
			case "matrix":
				err := parseMatrix(lexer, &tm)
				return tm, err
				// case "scale":
				// case "rotate":

			}
		}
	}
}

func parseMatrix(l *gl.Lexer, t *mt.Transform) error {
	i := l.NextItem()
	if i.Type != gl.ItemParan {
		return fmt.Errorf("Error Parsing Transform Matrix: Expected Opening Parantheses")
	}
	var ncount int
	for {
		if ncount > 0 {
			for l.PeekItem().Type == gl.ItemComma || l.PeekItem().Type == gl.ItemWSP {
				l.NextItem()
			}
		}
		if l.PeekItem().Type != gl.ItemNumber {
			return fmt.Errorf("Error Parsing Transform Matrix: Expected Number got %v", l.PeekItem().String())
		}
		n, err := parseNumber(l.NextItem())
		if err != nil {
			return err
		}
		t[ncount%2][ncount/3] = n
		ncount++
		if ncount > 5 {
			i = l.PeekItem()
			if i.Type != gl.ItemParan {
				return fmt.Errorf("Error Parsing Transform Matrix: Expected Closing Parantheses")
			}
			l.NextItem() // consume Parantheses
			return nil
		}
	}
}
