package ast

type stack []interface{}

var lastPop interface{}

func (s *stack) Push(v interface{}) {
	*s = append(*s, v)
}

func (s *stack) Pop() interface{} {
	self := *s
	object := self[len(self)-1]
	*s = self[:len(self)-1]
	lastPop = object
	return object
}

func (s *stack) PopN(n int) []interface{} {
	self := *s
	target := self[len(self)-n:]
	*s = self[:len(self)-n]
	lastPop = target
	return target
}
