package newaward

import (
	"crypto/rand"
	"encoding/binary"
	"io"
)

const (
	SIZE = 36
)

type ChildAction struct {
	doOne    func()
	doTwenty func()
	doEighty func()
}

func NewDummyChildAction() *ChildAction {
	return &ChildAction{
		doOne:    func() {},
		doTwenty: func() {},
		doEighty: func() {},
	}
}

type Award struct {
	left             int
	buf              []byte
	oneAction        func()
	noAwardAction    func()
	leftChildAction  *ChildAction
	rightChildAction *ChildAction
}

type Options func(*Award)

func WithLeftChildAction(one, twenty, eighty func()) Options {
	return func(a *Award) {
		a.leftChildAction.doOne = one
		a.leftChildAction.doTwenty = twenty
		a.leftChildAction.doEighty = eighty
	}
}

func WithRightChildAction(one, twenty, eighty func()) Options {
	return func(a *Award) {
		a.rightChildAction.doOne = one
		a.rightChildAction.doTwenty = twenty
		a.rightChildAction.doEighty = eighty
	}
}

func WithOneAction(one func()) Options {
	return func(a *Award) {
		a.oneAction = one
	}
}

func WithNoAwardAction(one func()) Options {
	return func(a *Award) {
		a.noAwardAction = one
	}
}

func NewAward(opts ...Options) *Award {
	a := &Award{
		buf:              make([]byte, SIZE),
		oneAction:        func() {},
		noAwardAction:    func() {},
		leftChildAction:  NewDummyChildAction(),
		rightChildAction: NewDummyChildAction(),
	}
	for _, o := range opts {
		o(a)
	}
	a.checkOverhead(0)
	return a
}

func (a *Award) checkOverhead(required int) {
	if a.left <= required {
		io.ReadFull(rand.Reader, a.buf)
		a.left = SIZE
	}
}

func (a *Award) FiftyPercentage() bool {
	a.checkOverhead(1)
	base := a.buf[SIZE-a.left:]
	expect := base[0] % 2
	a.left--
	return expect == 0
}

func (a *Award) OnePercentage() bool {
	a.checkOverhead(8)
	base := a.buf[SIZE-a.left:]
	expect := binary.LittleEndian.Uint32(base[0:4]) % 100
	r := binary.LittleEndian.Uint32(base[4:8]) % 100
	a.left -= 8
	return r == expect
}

func (a *Award) TwentyPercentage() bool {
	a.checkOverhead(8)
	base := a.buf[SIZE-a.left:]
	expect := binary.LittleEndian.Uint32(base[0:4]) % 5
	r := binary.LittleEndian.Uint32(base[4:8]) % 5
	a.left -= 8
	return r == expect
}

func (a *Award) PickLeftChild() {
	if a.OnePercentage() {
		a.leftChildAction.doOne()
	} else if a.TwentyPercentage() {
		a.leftChildAction.doTwenty()
	} else {
		a.leftChildAction.doEighty()
	}
}

func (a *Award) PickRightChild() {
	if a.OnePercentage() {
		a.rightChildAction.doOne()
	} else if a.TwentyPercentage() {
		a.rightChildAction.doTwenty()
	} else {
		a.rightChildAction.doEighty()
	}
}

func (a *Award) Pick() {
	if a.TwentyPercentage() {
		a.noAwardAction()
		return
	}

	if a.OnePercentage() {
		a.oneAction()
	} else if a.FiftyPercentage() {
		a.PickLeftChild()
	} else {
		a.PickRightChild()
	}
}
