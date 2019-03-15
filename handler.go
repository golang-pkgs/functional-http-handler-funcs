package fph

import (
	"net/http"
)

// Handler is typed to a nextable  http.HandlerFunc
type Handler func(res http.ResponseWriter, req *http.Request) Signal

// HandlerCondition returns a bool result by http response and request
type HandlerCondition func(res http.ResponseWriter, req *http.Request) bool

// Signal presents the continuation type of control flow
type Signal int

const (
	// Next present continue
	Next Signal = iota
	// Complete present done
	Complete
	// Error present panic
	Error
)

// EmptyHandlerFunc do nothing
func EmptyHandlerFunc(res http.ResponseWriter, req *http.Request) {}

// NextIt returns Handler which executes http.HandlerFunc and then call next automatically
func NextIt(hf http.HandlerFunc) Handler {
	return func(res http.ResponseWriter, req *http.Request) Signal {
		hf(res, req)
		return Next
	}
}

// CompleteIt returns Handler which executes http.HandlerFunc and then call next automatically
func CompleteIt(hf http.HandlerFunc) Handler {
	return func(res http.ResponseWriter, req *http.Request) Signal {
		hf(res, req)
		return Complete
	}
}

// ErrorIt returns Handler which executes http.HandlerFunc and then call err automatically
func ErrorIt(hf http.HandlerFunc) Handler {
	return func(res http.ResponseWriter, req *http.Request) Signal {
		hf(res, req)
		return Error
	}
}

// IfElse call a pair Handler by condition
func IfElse(condFunc HandlerCondition, handlerFunc1, handlerFunc2 Handler) Handler {
	return func(res http.ResponseWriter, req *http.Request) Signal {
		if condFunc(res, req) {
			return handlerFunc1(res, req)
		}

		return handlerFunc2(res, req)
	}
}

// Compose compose multiple Handlers
func Compose(hs ...Handler) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		for _, h := range hs {

			if sig := h(res, req); sig == Next {
				continue
			} else if sig == Complete {
				break
			}

			panic("reject")
		}
	}
}
