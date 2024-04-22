package infix_to_postfix

type Stack struct {
	Items []interface{}
}

func (s *Stack) Empty() bool {
	return len(s.Items) == 0
}

func (s *Stack) TopFunc() interface{} {
	if !s.Empty() {
		return s.Items[len(s.Items)-1]
	}
	return nil
}

func (s *Stack) Push(value interface{}) {
	s.Items = append(s.Items, value)
}

func (s *Stack) Pop() (value interface{}) {
	if !s.Empty() {
		top := s.Items[len(s.Items)-1]
		s.Items = s.Items[:len(s.Items)-1]
		return top
	}
	return nil
}
