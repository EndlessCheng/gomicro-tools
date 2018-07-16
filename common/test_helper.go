package common

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type tNode struct {
	mocks    func()
	children []tNode
	run      func()
}

func (n tNode) runExpects(expects []func()) {
	if n.children == nil {
		for _, expectFunc := range expects {
			if expectFunc != nil {
				expectFunc()
			}
		}
		if n.mocks != nil {
			n.mocks()
		}

		if n.run == nil {
			panic("未填写 run")
		}

		if n.run != nil {
			n.run()
		}
		return

		gomock.InOrder()
	}

	for _, c := range n.children {
		expects = append(expects, n.mocks)
		c.runExpects(expects)
		expects = expects[:len(expects)-1]
	}
}

func (n tNode) runTestNode() {
	n.runExpects([]func(){})
}

//

type equals struct {
	t    *testing.T
	want []interface{}
	got  []interface{}
}

func newEQ(t *testing.T) *equals {
	return &equals{t: t}
}

func (e *equals) w(want ...interface{}) *equals {
	e.want = want
	return e
}

func (e *equals) check(got ...interface{}) {
	e.got = got
	e.assertEquals()
}

func (e *equals) assertEquals() {
	if len(e.want) != len(e.got) {
		fmt.Println(e.want, e.got)
		panic(fmt.Sprintf("期望的结果数目有误，请检查代码: 期望 %d, 实际 %d", len(e.want), len(e.got)))
	}
	for i, result := range e.got {
		// 有时 nil 会被识别成带类型的 nil，此时用 assert.Equal 会显示二者不同
		// 比如：返回类型为 *struct 时返回 nil
		if e.want[i] == nil {
			assert.Nil(e.t, result)
		} else {
			assert.Equal(e.t, e.want[i], result)
		}
	}
}
