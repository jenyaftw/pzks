package tokenizer

import "testing"

func TestValidExpressions(t *testing.T) {
	expressions := []string{
		"3.15 + 4545.313",
		"3 + 5",
		"x * (y + 2)",
		"sin(x) + cos(y)",
		"2 * (3 + 4) / 5",
		"sqrt(a) - log(b)",
		"(x + y) * (z - w)",
		"3 + 5 * (2 - 4)",
		"x / (y - z) * a",
		"abs(-5) + round(4.67)",
		"log(10) + sqrt(16)",
	}

	for _, expression := range expressions {
		tokenizer := NewTokenizer()
		_, errors := tokenizer.Tokenize(expression)
		if len(errors) > 0 {
			t.Errorf("Expected no errors, got %v", errors)
		}
	}
}

func TestStartingWithInvalidCharacter(t *testing.T) {
	expressions := []string{
		"* 3 + 5",
		"/ 2 - 1",
		") x + 5",
		"+ 7 * 8",
		"- (3 + 5)",
	}

	for _, expression := range expressions {
		tokenizer := NewTokenizer()
		_, errors := tokenizer.Tokenize(expression)
		if len(errors) == 0 {
			t.Errorf("Expected error, got %v", errors)
		}
	}
}

func TestEndsWithInvalidCharacter(t *testing.T) {
	expressions := []string{
		"3 + 5 *",
		"(x - y /",
		"sqrt(16) +",
		"(a + b -",
	}

	for _, expression := range expressions {
		tokenizer := NewTokenizer()
		_, errors := tokenizer.Tokenize(expression)
		if len(errors) == 0 {
			t.Errorf("Expected error, got %v", errors)
		}
	}
}

func TestMismatchedBrackets(t *testing.T) {
	expressions := []string{
		"3 + (4 - 5",
		"(2 + 3)) * 7",
		"(x * (y + z)",
		"((a - b) + c",
		"log(10 + sqrt(25)) + abs(x",
	}

	for _, expression := range expressions {
		tokenizer := NewTokenizer()
		_, errors := tokenizer.Tokenize(expression)
		if len(errors) == 0 {
			t.Errorf("Expected error, got %v", errors)
		}
	}
}

func TestEmptyBrackets(t *testing.T) {
	expressions := []string{
		"x * () + 5",
		"sin() + cos(y)",
		"(3 + ) * 4",
		"log( + b)",
	}

	for _, expression := range expressions {
		tokenizer := NewTokenizer()
		_, errors := tokenizer.Tokenize(expression)
		if len(errors) == 0 {
			t.Errorf("Expected error, got %v", errors)
		}
	}
}
