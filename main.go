package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type NodeType int

const (
	NodeNumber NodeType = iota
	NodeBinOp
)

type Node interface {
	nodeType() NodeType
}

type NumberNode struct {
	Value int
}

func (n *NumberNode) nodeType() NodeType {
	return NodeNumber
}

type BinOpNode struct {
	Left     Node
	Right    Node
	Operator string
}

func (n *BinOpNode) nodeType() NodeType {
	return NodeBinOp
}

type TokenType int

const (
	TokenNumber TokenType = iota
	TokenOperator
	TokenLeftParen
	TokenRightParen
	TokenError
	TokenEOF
)

type Token struct {
	Type  TokenType
	Value string
}

type Lexer struct {
	input        string
	position     int
	currentToken Token
}

type Parser struct {
	lexer        *Lexer
	currentToken Token
}

func NewParser(lexer *Lexer) *Parser {
	parser := &Parser{
		lexer: lexer,
	}
	parser.consumeToken() // Initialize currentToken
	return parser
}

func (l *Lexer) readNumber() Token {
	startPosition := l.position
	for l.position < len(l.input) && unicode.IsDigit(rune(l.input[l.position])) {
		l.position++
	}
	return Token{Type: TokenNumber, Value: l.input[startPosition:l.position]}
}

func (l *Lexer) readOperator() Token {
	currentChar := l.input[l.position]
	l.position++
	return Token{Type: TokenOperator, Value: string(currentChar)}
}

func (l *Lexer) skipWhitespace() {
	for l.position < len(l.input) && unicode.IsSpace(rune(l.input[l.position])) {
		l.position++
	}
}

func NewLexer(input string) *Lexer {
	return &Lexer{
		input:    input,
		position: 0,
	}
}

func (l *Lexer) NextToken() Token {
	l.skipWhitespace()

	if l.position >= len(l.input) {
		return Token{Type: TokenEOF, Value: ""}
	}

	currentChar := l.input[l.position]

	switch {
	case unicode.IsDigit(rune(currentChar)):
		return l.readNumber()
	case strings.ContainsRune("+-*/", rune(currentChar)):
		return l.readOperator()
	case currentChar == '(':
		l.position++
		return Token{Type: TokenLeftParen, Value: string(currentChar)}
	case currentChar == ')':
		l.position++
		return Token{Type: TokenRightParen, Value: string(currentChar)}
	default:
		return Token{Type: TokenError, Value: string(currentChar)}
	}
}

func (p *Parser) consumeToken() {
	p.currentToken = p.lexer.NextToken()
}

func (p *Parser) parseExpression() Node {
	node := p.parseTerm()

	for p.currentToken.Type == TokenOperator && (p.currentToken.Value == "+" || p.currentToken.Value == "-") {
		operator := p.currentToken.Value
		p.consumeToken()
		right := p.parseTerm()
		node = &BinOpNode{Left: node, Right: right, Operator: operator}
	}

	return node
}

func (p *Parser) parseTerm() Node {
	node := p.parseFactor()

	for p.currentToken.Type == TokenOperator && (p.currentToken.Value == "*" || p.currentToken.Value == "/") {
		operator := p.currentToken.Value
		p.consumeToken()
		right := p.parseFactor()
		node = &BinOpNode{Left: node, Right: right, Operator: operator}
	}

	return node
}

func (p *Parser) parseFactor() Node {
	token := p.currentToken

	switch token.Type {
	case TokenNumber:
		p.consumeToken()
		value, _ := strconv.Atoi(token.Value)
		return &NumberNode{Value: value}
	case TokenLeftParen:
		p.consumeToken()
		node := p.parseExpression()
		if p.currentToken.Type != TokenRightParen {
			panic("Expected ')' after expression")
		}
		p.consumeToken() // Consume the ')'
		return node
	default:
		panic("Unexpected token")
	}
}

func evaluate(node Node) int {
	switch n := node.(type) {
	case *NumberNode:
		return n.Value
	case *BinOpNode:
		left := evaluate(n.Left)
		right := evaluate(n.Right)
		switch n.Operator {
		case "+":
			return left + right
		case "-":
			return left - right
		case "*":
			return left * right
		case "/":
			if right == 0 {
				panic("division by zero")
			}
			return left / right
		default:
			panic("unknown operator")
		}
	default:
		panic("unknown node type")
	}
}

func main() {
	for {
		fmt.Print("$ ")
		var input string
		fmt.Scanln(&input)

		if input == "exit" {
			break
		}

		lexer := NewLexer(input)
		parser := NewParser(lexer)
		ast := parser.parseExpression()

		result := evaluate(ast)
		fmt.Printf("Result: %d\n", result)
	}
}
