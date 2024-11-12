package tokenizer

import (
	"strconv"

	"github.com/jenyaftw/lab1/token"
)

type TokenizerError struct {
	Message  string
	StartIdx int
	EndIdx   int
}

const (
	NewToken = iota
	NumericLiteral
	SymbolName
	ParanthesesOpen
	ParanthesesClose
	Operator
	CompleteToken
)

type Tokenizer struct{}

func NewTokenizer() *Tokenizer {
	return &Tokenizer{}
}

func (p *Tokenizer) Tokenize(expression string) ([]token.Token, []TokenizerError) {
	errors := []TokenizerError{}
	if expression == "" {
		errors = append(errors, TokenizerError{Message: "empty expression", StartIdx: 0, EndIdx: 0})
		return nil, errors
	}

	tokens := []token.Token{}

	stateNow := NewToken
	stateNext := NewToken
	currentTokenStr := ""
	currentToken := token.Token{}

	decimalPointCount := 0
	parantheseCount := 0

	idxNow := 0
	for idxNow < len(expression) {
		charNow := expression[idxNow]

		switch stateNow {
		case NewToken:
			currentTokenStr = ""
			currentToken = token.Token{
				Type:     token.UnknownType,
				Text:     "",
				StartIdx: idxNow,
			}
			decimalPointCount = 0

			if token.Whitespace.Check(charNow) { // Check for whitespace
				stateNext = NewToken
				idxNow += 1
			} else if token.RealNumericDigits.Check(charNow) { // Check for numeric literal
				currentTokenStr += string(charNow)

				stateNext = NumericLiteral
				idxNow += 1
			} else {
				if idxNow == 0 && !token.Symbols.Check(charNow) && string(charNow) != "(" {
					errors = append(errors, TokenizerError{Message: "неочікуваний на початку виразу символ " + string(charNow), StartIdx: idxNow, EndIdx: idxNow})
				}

				if token.Operators.Check(charNow) { // Check for operator
					stateNext = Operator
				} else if string(charNow) == "(" { // Check for parantheses open
					stateNext = ParanthesesOpen
				} else if string(charNow) == ")" { // Check for parantheses close
					stateNext = ParanthesesClose
				} else {
					currentTokenStr += string(charNow)
					idxNow += 1
					stateNext = SymbolName
				}
			}
		case NumericLiteral:
			if token.RealNumericDigits.Check(charNow) {
				if string(charNow) == "." {
					if decimalPointCount > 0 {
						errors = append(errors, TokenizerError{Message: "друга десяткова крапка в числі", StartIdx: idxNow, EndIdx: idxNow})
					} else {
						decimalPointCount += 1
					}
				}

				currentTokenStr += string(charNow)
				idxNow += 1
				stateNext = NumericLiteral
			} else {
				if len(tokens) > 0 && tokens[len(tokens)-1].Text == "-" {
					tokens[len(tokens)-1].Type = token.UnaryOperatorType
					if len(errors) > 0 {
						errors = errors[:len(errors)-1]
					}
				}

				stateNext = CompleteToken
				currentToken.Type = token.NumericLiteralType
				currentToken.Text = currentTokenStr
				currentToken.EndIdx = idxNow - 1
				currentToken.Value, _ = strconv.ParseFloat(currentTokenStr, 32)
			}
		case Operator:
			if token.Operators.Check(charNow) {
				if len(tokens) > 0 {
					if tokens[len(tokens)-1].Type == token.ParanthesesOpenType {
						errors = append(errors, TokenizerError{Message: "оператор після відкриваючої дужки, очікувалась змінна", StartIdx: idxNow, EndIdx: idxNow})
					} else if tokens[len(tokens)-1].Type == token.OperatorType {
						errors = append(errors, TokenizerError{Message: "два оператори підряд", StartIdx: idxNow, EndIdx: idxNow})
					}
				}

				currentToken.Type = token.OperatorType
				currentToken.Text = string(charNow)
				currentToken.StartIdx = idxNow
				currentToken.EndIdx = idxNow

				stateNext = CompleteToken
				idxNow += 1
			} else {
				errors = append(errors, TokenizerError{Message: "invalid operator", StartIdx: idxNow, EndIdx: idxNow})
			}
		case CompleteToken:
			tokens = append(tokens, currentToken)
			stateNext = NewToken
		case ParanthesesOpen:
			currentTokenStr += string(charNow)
			idxNow += 1
			parantheseCount += 1

			currentToken.Type = token.ParanthesesOpenType
			currentToken.Text = currentTokenStr
			currentToken.EndIdx = idxNow - 1
			stateNext = CompleteToken
		case ParanthesesClose:
			if len(tokens) > 0 && (tokens[len(tokens)-1].Type == token.OperatorType || tokens[len(tokens)-1].Type == token.TrigFunction || tokens[len(tokens)-1].Type == token.OtherFunction) {
				errors = append(errors, TokenizerError{Message: "закриваюча дужка після оператора, очікувалась змінна", StartIdx: idxNow, EndIdx: idxNow})
			}

			if len(tokens) > 0 && tokens[len(tokens)-1].Type == token.ParanthesesOpenType {
				errors = append(errors, TokenizerError{Message: "порожня дужка", StartIdx: idxNow, EndIdx: idxNow})
			}

			currentTokenStr += string(charNow)
			idxNow += 1
			parantheseCount -= 1

			currentToken.Type = token.ParanthesesCloseType
			currentToken.Text = currentTokenStr
			currentToken.EndIdx = idxNow - 1
			stateNext = CompleteToken
		case SymbolName:
			if token.Symbols.Check(charNow) {
				currentTokenStr += string(charNow)
				idxNow += 1
				stateNext = SymbolName
			} else {
				currentToken.Type = token.Symbol

				if m, _ := token.TrigRegex.CheckRegex(currentTokenStr); m {
					currentToken.Type = token.TrigFunction
				}

				if m, _ := token.OtherFunctionsRegex.CheckRegex(currentTokenStr); m {
					currentToken.Type = token.OtherFunction
				}

				currentToken.Text = currentTokenStr
				currentToken.EndIdx = idxNow - 1
				stateNext = CompleteToken
			}
		}

		stateNow = stateNext
	}

	if stateNow == NumericLiteral {
		currentToken = token.Token{
			Type: token.NumericLiteralType,
			Text: currentTokenStr,
		}
		currentToken.Value, _ = strconv.ParseFloat(currentTokenStr, 32)
		tokens = append(tokens, currentToken)
	} else if stateNow == Operator {
		if token.Operators.Check(currentTokenStr[0]) {
			currentToken = token.Token{
				Type:     token.OperatorType,
				Text:     currentTokenStr,
				EndIdx:   idxNow - 1,
				StartIdx: idxNow - 1,
			}
			tokens = append(tokens, currentToken)
		} else {
			errors = append(errors, TokenizerError{Message: "невалідний оператор", StartIdx: idxNow, EndIdx: idxNow})
		}
	} else if stateNow == CompleteToken {
		tokens = append(tokens, currentToken)
	} else if stateNow == SymbolName {
		currentToken = token.Token{
			Type: token.Symbol,
			Text: currentTokenStr,
		}
		tokens = append(tokens, currentToken)
	}

	if len(tokens) > 0 && tokens[len(tokens)-1].Type == token.OperatorType {
		errors = append(errors, TokenizerError{Message: "кінець виразу після оператора, очікувалась змінна", StartIdx: tokens[len(tokens)-1].StartIdx, EndIdx: tokens[len(tokens)-1].EndIdx})
	}

	parantheseCount = 0
	for i := 0; i < len(tokens); i++ {
		if tokens[i].Type == token.ParanthesesOpenType {
			parantheseCount += 1
			foundClosed := false
			localParantheseCount := 1
			for j := i + 1; j < len(tokens); j++ {
				if tokens[j].Type == token.ParanthesesOpenType {
					localParantheseCount += 1
				} else if tokens[j].Type == token.ParanthesesCloseType {
					localParantheseCount -= 1
					if localParantheseCount == 0 {
						foundClosed = true
					}
					break
				}
			}
			if !foundClosed {
				errors = append(errors, TokenizerError{Message: "не закрита дужка", StartIdx: tokens[i].StartIdx, EndIdx: tokens[i].EndIdx})
			}
		} else if tokens[i].Type == token.ParanthesesCloseType {
			parantheseCount -= 1
			if parantheseCount < 0 {
				errors = append(errors, TokenizerError{Message: "зайва дужка", StartIdx: tokens[i].StartIdx, EndIdx: tokens[i].EndIdx})
			}
		}
	}

	return tokens, errors
}
