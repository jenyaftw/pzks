package token

import (
	"regexp"
)

const (
	UnknownType = iota
	NumericLiteralType
	UnaryOperatorType
	OperatorType
	ParanthesesOpenType
	ParanthesesCloseType
	TrigFunction
	OtherFunction
	Symbol
)

type Token struct {
	Type     int
	Text     string
	Value    float64
	StartIdx int
	EndIdx   int
}

func (t Token) String() string {
	o := ""

	switch t.Type {
	case UnknownType:
		o += "[Unknown]"
	case NumericLiteralType:
		o += "[Literal, Numeric]"
	case OperatorType:
		o += "[Operator]"
	case UnaryOperatorType:
		o += "[Operator, Unary]"
	case ParanthesesOpenType:
		o += "[Paranthese, Open]"
	case ParanthesesCloseType:
		o += "[Paranthese, Close]"
	case Symbol:
		o += "[Symbol]"
	case TrigFunction:
		o += "[Function, Trig]"
	case OtherFunction:
		o += "[Function, Other]"
	}

	return o + " : " + t.Text
}

func (t Token) GetPrecedence() int {
	if t.Type != OperatorType && t.Type != UnaryOperatorType {
		return -1
	}

	switch t.Text {
	case "+", "-":
		return 2
	case "*", "/":
		return 3
	case "^":
		return 4
	}

	return -1
}

func (t Token) IsLeftAssociative() bool {
	if t.Type != OperatorType && t.Type != UnaryOperatorType {
		return false
	}

	switch t.Text {
	case "+", "-", "*", "/":
		return true
	case "^":
		return false
	}

	return false
}

func (t Token) IsRightAssociative() bool {
	if t.Type != OperatorType && t.Type != UnaryOperatorType {
		return false
	}

	switch t.Text {
	case "+", "-", "*", "/":
		return false
	case "^":
		return true
	}

	return false
}

type CharacterChecker struct {
	Characters []rune
}

func NewCharacterChecker(characters string) *CharacterChecker {
	return &CharacterChecker{
		Characters: []rune(characters),
	}
}

func (cc *CharacterChecker) Check(char byte) bool {
	for _, c := range cc.Characters {
		if c == rune(char) {
			return true
		}
	}

	return false
}

func (cc *CharacterChecker) CheckRegex(str string) (bool, error) {
	m, err := regexp.MatchString(string(cc.Characters), str)
	if err != nil {
		return false, err
	}

	return m, nil
}

var Operators = NewCharacterChecker("+-*/^")
var NumericDigits = NewCharacterChecker("0123456789")
var Whitespace = NewCharacterChecker(" \t\n\r\v\f")
var RealNumericDigits = NewCharacterChecker(".0123456789")
var Symbols = NewCharacterChecker("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ._")
var TrigRegex = NewCharacterChecker("sin|cos|tan")
var OtherFunctionsRegex = NewCharacterChecker("sqrt|log|ln|exp|abs|ceil|floor|round|pow")
