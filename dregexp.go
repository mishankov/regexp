package dregexp

import (
	"fmt"
	"strconv"
	"strings"
)

func digits() string {
	result := ""
	for i := 0; i < 10; i++ {
		result += strconv.Itoa(i)
	}
	return result
}

func alphaNumeric() string {
	result := ""

	for i := 'a'; i < 'z'; i++ {
		result += string(i)
	}

	for i := 'A'; i < 'Z'; i++ {
		result += string(i)
	}

	result += digits() + "_"

	return result
}

type TokenQuantifier struct {
	isOneOrMore bool
	isZeroOrOne bool
}

type Token struct {
	isLineStart     bool
	isLineEnd       bool
	isLiteral       bool
	literal         rune
	isDigit         bool
	isAlphaNumeric  bool
	isAny           bool
	isPositiveGroup bool
	isNegativeGroup bool
	isAlternation   bool
	alternations    [][]Token
	group           []Token
	quantifier      TokenQuantifier
}

func (t Token) String() string {
	result := ""
	switch {
	case t.isLineStart:
		result += "Line_start"
	case t.isLineEnd:
		result += "Line_end"
	case t.isLiteral:
		result += fmt.Sprintf("%q", t.literal)
	case t.isDigit:
		result += "Digit"
	case t.isAlphaNumeric:
		result += "Alpha_numeric"
	case t.isAny:
		result += "Any"
	case t.isNegativeGroup:
		result += fmt.Sprintf("Negative_group(%v)", t.group)
	case t.isPositiveGroup:
		result += fmt.Sprintf("Positive_group(%v)", t.group)
	case t.isAlternation:
		result += fmt.Sprintf("Alternation(%v)", t.alternations)
	default:
		result += "Ivalid_token"
	}

	switch {
	case t.quantifier.isOneOrMore:
		result += "[1+]"
	case t.quantifier.isZeroOrOne:
		result += "[?]"
	}

	return result
}

func (t Token) matches(s rune) bool {
	switch {
	case t.isLiteral:
		return s == t.literal
	case t.isAlphaNumeric:
		return strings.ContainsAny(string(s), alphaNumeric())
	case t.isDigit:
		return strings.ContainsAny(string(s), digits())
	case t.isAny:
		return true
	case t.isNegativeGroup:
		for _, token := range t.group {
			if token.matches(s) {
				return false
			}
		}
		return true
	case t.isPositiveGroup:
		for _, token := range t.group {
			if token.matches(s) {
				return true
			}
		}
	case t.isAlternation:
		// for _, alternation := range t.alternations {
		// 	ok, _ := matchTokens()
		// }
	}
	return false
}

func parseTokens(pattern string) []Token {
	tokens := []Token{}

	i := 0
	for s := pattern[0]; i < len(pattern); s = pattern[i] {
		token := Token{}
		switch s {
		case '^':
			token = Token{isLineStart: true}
			i++
		case '$':
			token = Token{isLineEnd: true}
			i++
		case '.':
			token = Token{isAny: true}
			i++
		case '\\':
			switch pattern[i+1] {
			case 'd':
				token = Token{isDigit: true}
				i += 2
			case 'w':
				token = Token{isAlphaNumeric: true}
				i += 2
			default:
				token = Token{isLiteral: true, literal: '\\'}
				i += 2
			}
		case '[':
			if pattern[i+1] == '^' {
				token.isNegativeGroup = true
			} else {
				token.isPositiveGroup = true
			}

			for groupEndIndex := i + 1; groupEndIndex < len(pattern); groupEndIndex++ {
				if pattern[groupEndIndex] == ']' {
					token.group = parseTokens(pattern[i+1 : groupEndIndex])
					i = groupEndIndex + 1
					break
				}
			}
		case '(':
			token.isAlternation = true
			alternations := [][]Token{}
			for index := i + 1; index < len(pattern); index++ {
				if pattern[index] == '|' {
					alternation := parseTokens(pattern[i+1 : index])
					alternations = append(alternations, alternation)
					i = index
				} else if pattern[index] == ')' {
					alternation := parseTokens(pattern[i+1 : index])
					alternations = append(alternations, alternation)
					token.alternations = alternations
					i = index + 1
					break
				}
			}
		default:
			token = Token{isLiteral: true, literal: rune(s)}
			i++
		}

		if i < len(pattern) {
			next_s := pattern[i]

			switch next_s {
			case '+':
				token.quantifier = TokenQuantifier{isOneOrMore: true}
				i++
			case '?':
				token.quantifier = TokenQuantifier{isZeroOrOne: true}
				i++
			}
		}

		tokens = append(tokens, token)

		if i >= len(pattern) {
			break
		}
	}

	return tokens
}

type Matcher struct {
	tokens            []Token
	currentTokenIndex int
	line              string
	currentLineIndex  *int
	currentToken      Token
	usingQuantifier   bool
}

func newMatcher(tokens []Token, index *int) *Matcher {
	return &Matcher{tokens: tokens, currentTokenIndex: 0, currentLineIndex: index, currentToken: Token{}, usingQuantifier: false}
}

func (m *Matcher) matchRune(r rune) {
	if m.currentToken.matches(r) {
		// Handling quantifiers
		if m.currentToken.quantifier != (TokenQuantifier{}) {
			switch {
			case m.currentToken.quantifier.isOneOrMore:
				// If rune matches token with "1+" quantifier, set flag and keep current tocken
				m.usingQuantifier = true
			case m.currentToken.quantifier.isZeroOrOne:
				// If rune matches token with "zero or one" quantifier, increaze token index
				m.currentTokenIndex++
			}
		} else {
			// Increase current token index if current rune matches current token
			m.currentTokenIndex++
		}
		// Increase rune index if current rune matches current token
		*m.currentLineIndex++
	} else {
		if m.usingQuantifier {
			// If rune does not match token, but current token with "1+" quantifier wath mathed earlier,
			// reset quantifier flag and move to next token without incrementing rune index
			m.usingQuantifier = false
			m.currentTokenIndex++
		} else if m.currentToken.quantifier != (TokenQuantifier{}) && m.currentToken.quantifier.isZeroOrOne {
			// If rune does not match token, but token has "zero or one" quantifier, move to next token without incrementing rune index
			m.currentTokenIndex++
		} else {
			// If rune does not match token, reset token index and increase rune index
			m.currentTokenIndex = 0
			*m.currentLineIndex++
		}
	}
}

func (m Matcher) matchLine(line string) (bool, error) {
	fmt.Println("Tokens: ", m.tokens, line[*m.currentLineIndex:])
	m.line = line

	for s := line[*m.currentLineIndex]; *m.currentLineIndex < len(line); s = line[*m.currentLineIndex] {
		r := rune(s)
		m.currentToken = m.tokens[m.currentTokenIndex]

		// Handling line start logic
		if m.currentTokenIndex == 0 && m.currentToken.isLineStart && *m.currentLineIndex != 0 {
			return false, nil
		} else if m.currentTokenIndex == 0 && m.currentToken.isLineStart {
			m.currentTokenIndex = 1
			m.currentToken = m.tokens[m.currentTokenIndex]
		}

		// Handling alternations
		matchedAlt := false
		if m.currentToken.isAlternation {
			currentLineIndex := *m.currentLineIndex
			for _, alteration := range m.currentToken.alternations {
				alterationMatcher := newMatcher(alteration, m.currentLineIndex)
				ok, _ := alterationMatcher.matchLine(line)

				if ok {
					fmt.Println("ok", *m.currentLineIndex)
					matchedAlt = ok
					m.currentTokenIndex++
					break
				}

				*m.currentLineIndex = currentLineIndex
			}
		}

		// Skip matching if alteration matched
		if !matchedAlt {
			// Match current token to current rune
			m.matchRune(r)
		}

		if m.currentTokenIndex == len(m.tokens) {
			// If it got to last token, it means that all token matched
			return true, nil
		}

		if *m.currentLineIndex >= len(line) {
			breakOuter := false
			for _, token := range m.tokens[m.currentTokenIndex:] {
				if token.quantifier == (TokenQuantifier{}) || !token.quantifier.isZeroOrOne {
					// If it is end of line and only tokens with "zero or one" quantifier left, it means that string matches pattern
					breakOuter = true
					break
				}
			}
			if breakOuter {
				break
			}

			return true, nil
		}
	}

	if m.currentTokenIndex == len(m.tokens)-1 && m.tokens[m.currentTokenIndex].isLineEnd {
		// Handle line end
		return true, nil
	}

	return false, nil
}

func MatchLine(line string, pattern string) (bool, error) {
	tokens := parseTokens(pattern)
	index := 0
	matcher := newMatcher(tokens, &index)

	return matcher.matchLine(line)
}
