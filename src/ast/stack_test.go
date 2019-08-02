package ast

import "testing"

func TestStackPush(t *testing.T) {
	testStack := createStack()
	testStack.Push(1)
	testStack.Push(2)
	testStack.Push(3)
	if !checkArraysMatch(testStack, []interface{}{1, 2, 3}) {
		t.Error("Stack doesn't match", testStack)
	}
}

func TestStackPop(t *testing.T) {
	testStack := createStack()
	testStack.Push(1)
	testStack.Push(2)
	testStack.Push(3)
	// Stack = [1,2,3]
	a := testStack.Pop()
	if a != 3 {
		t.Error("Popping failed. Returned incorrect value", a)
	}
	// Stack = [1,2]
	testStack.Push(4)
	testStack.Push(5)
	// Stack = [1,2,4,5]
	b := testStack.Pop()
	if b != 5 {
		t.Error("Popping failed. Returned incorrect value", b)
	}
	// Stack = [1,2,4]
	c := testStack.Pop()
	if c != 4 {
		t.Error("Popping failed. Returned incorrect value", c)
	}
	// Stack = [1,2]
	if !checkArraysMatch(testStack, []interface{}{1, 2}) {
		t.Error("Stack doesn't match", testStack)
	}
}

func TestStackPopN(t *testing.T) {
	testStack := createStack()
	testStack.Push(1)
	testStack.Push(2)
	testStack.Push(3)
	testStack.Push(4)
	// Stack = [1,2,3,4]
	a := testStack.PopN(2)
	if !checkArraysMatch(a, []interface{}{3, 4}) {
		t.Error("Popping failed. Returned incorrect value", a)
	}
	// Stack = [1,2]
	testStack.Push(5)
	// Stack = [1,2,5]
	b := testStack.PopN(3)
	if !checkArraysMatch(b, []interface{}{1, 2, 5}) {
		t.Error("Popping failed. Returned incorrect value", b)
	}
	// Stack = []
	if !checkArraysMatch(testStack, []interface{}{}) {
		t.Error("Stack doesn't match", testStack)
	}
}

func TestStackDebugging(t *testing.T) {
	testStack := createStack()
	testStack.Push(1)
	testStack.Push(2)
	testStack.Push(3)
	testStack.Pop() //3 should be popped
	if lastPop != 3 {
		t.Error("Last pop doesn't match", lastPop)
	}
	testStack.PopN(2) //[1, 2] should be popped
	if !checkArraysMatch(lastPop.(stack), []interface{}{1, 2}) {
		t.Error("Last pop doesn't match", lastPop)
	}

	//Check that stack is clean while we're at it
	if !checkArraysMatch(testStack, []interface{}{}) {
		t.Error("Stack doesn't match", testStack)
	}
}

func createStack() stack {
	return make([]interface{}, 0)
}

func checkArraysMatch(a, b []interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
