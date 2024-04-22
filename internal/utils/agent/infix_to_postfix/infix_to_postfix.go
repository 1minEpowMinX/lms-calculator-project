package infix_to_postfix

import (
	"strings"
)

// Певращение обычных выражений в постфиксную запись
// s string = "2+ 2 *2"
// return "2 2 2 * +"
func ToPostfix(s string) string {
	// объявления стека, постфиксной строки
	var stack Stack
	postfix := ""
	length := len(s)

	for i := 0; i < length; i++ {
		char := string(s[i])

		if char == " " {
			continue
		}

		if char == "(" {
			stack.Push(char)
		} else if char == ")" {
			for !stack.Empty() {
				str, _ := stack.TopFunc().(string)
				if str == "(" {
					break
				}
				postfix += " " + str
				stack.Pop()
			}
			stack.Pop()
		} else if !IsOperator(s[i]) {
			j := i
			number := ""

			for ; j < length && IsOperand(s[j]); j++ {
				number += string(s[j])
			}
			postfix += " " + number
			i = j - 1
		} else {
			for !stack.Empty() {
				top, _ := stack.TopFunc().(string)
				if top == "(" || !HasHigherPrecedence(top, char) {
					break
				}
				postfix += " " + top
				stack.Pop()
			}
			stack.Push(char)
		}
	}
	for !stack.Empty() { // дозапись оставшихся элементов
		str, _ := stack.Pop().(string)
		postfix += " " + str
	}
	return strings.TrimSpace(postfix)
}

func IsOperator(c uint8) bool {
	return strings.ContainsAny(string(c), "+ & - & * & /")
}

func IsOperand(c uint8) bool {
	return c >= '0' && c <= '9'
}

// Определение порядка выполнения арифметических действий
func Precedence(op string) int {
	switch op {
	case "+", "-":
		return 1
	case "*", "/":
		return 2
	}
	return -1
}

// Определение какая из операций имеет высший приоритет
func HasHigherPrecedence(op1, op2 string) bool {
	op1Prec := Precedence(op1)
	op2Prec := Precedence(op2)
	return op1Prec >= op2Prec
}
