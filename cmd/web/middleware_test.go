package main

import (
	"net/http"
	"testing"
)

//func TestNoSurf(t *testing.T) {
//	var myHandler *MyHandler
//	// ?为什么这里必须要有指针类型的handler
//	switch v := NoSurf(myHandler).(type) {
//	case http.Handler:
//	// do nothing
//	default:
//		t.Errorf("v is not http.Handler but is %T", v)
//	}
//}

func TestNoSurf(t *testing.T) {
	var myH MyHandler
	h := NoSurf(&myH)

	switch v := h.(type) {
	case http.Handler:
		// do nothing
	default:
		t.Errorf("type is not http.Handler but is %T", v)
	}
}

func TestSession(t *testing.T) {
	var myH MyHandler
	switch v := SessionLoad(&myH).(type) {
	case http.Handler:
		// do nothing
	default:
		t.Errorf("type is not http.Handler but is %T", v)
	}
}
