package rete

import (
	"container/list"
	"fmt"
	"strings"
)

type Env map[string]interface{}

type Token struct {
	parent      *Token
	wme         *WME
	node        IReteNode
	children    *list.List
	joinResults *list.List // used in negative nodes
	nccResults  *list.List
	owner       *Token
	binding     Env
}

func (tok *Token) get_wmes() []*WME {
	var ret []*WME
	_ws := list.New()
	_ws.PushFront(tok.wme)
	for tok.parent != nil {
		tok = tok.parent
		_ws.PushFront(tok.wme)
	}
	for e := _ws.Front(); e != nil; e = e.Next() {
		ret = append(ret, e.Value.(*WME))
	}
	return ret
}

func makeToken(node IReteNode, parent *Token, w *WME, b Env) *Token {
	tok := &Token{
		parent:   parent,
		wme:      w,
		node:     node,
		children: list.New(),
		binding:  b,
	}
	if parent != nil {
		parent.children.PushBack(tok)
	}
	if w != nil {
		w.tokens.PushBack(tok)
	}
	return tok
}

func (tok *Token) deleteTokenAndDescendents() {
	for tok.children != nil && tok.children.Len() > 0 {
		e := tok.children.Front()
		child := e.Value.(*Token)
		child.deleteTokenAndDescendents()
		tok.children.Remove(e)
	}
	removeByValue(tok.node.GetItems(), tok)
	if tok.wme != nil {
		removeByValue(tok.wme.tokens, tok)
	}
	if tok.parent != nil {
		removeByValue(tok.parent.children, tok)
	}
}

func (tok *Token) String() string {
	var ret []string
	wmes := tok.get_wmes()

	// Check if the token itself or the wmes slice is nil
	if tok == nil || wmes == nil {
		return "<Token <nil>>"
	}

	for _, v := range wmes {
		if v != nil {
			ret = append(ret, v.String()) // Safely call String() if v is not nil
		} else {
			ret = append(ret, "<nil>") // Explicitly handle nil WMEs
		}
	}

	// Join the string representations of all WMEs (or <nil> placeholders)
	return fmt.Sprintf("<Token %s>", strings.Join(ret, ", "))
}

func (tok *Token) GetBinding(k string) interface{} {
	var v interface{}
	t := tok
	if t.binding != nil {
		v = t.binding[k]
	}
	for v == nil && t.parent != nil {
		t = t.parent
		if t.binding != nil {
			v = t.binding[k]
		}
	}
	return v
}

func (tok *Token) GetRHSParam(k string) interface{} {
	node, ok := tok.node.(*BetaMemory)
	if !ok {
		return nil
	}
	return node.GetExecuteParam(k)
}

func (tok *Token) AllBinding() Env {
	var path []*Token
	t := tok
	for t != nil {
		path = append(path, t)
		t = t.parent
	}
	result := make(Env)
	for _, t := range path {
		for k, v := range t.binding {
			result[k] = v
		}
	}
	return result
}
