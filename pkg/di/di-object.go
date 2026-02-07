package di

import (
	"reflect"
)

// ---------------------------------------------------------------------------
// #region definition

type DIInstanceType int8

const (
	Singleton = iota
)

type DIObject struct {
	token        string
	generator    reflect.Value
	instanceType DIInstanceType
	value        interface{}
}

// #endregion

// ---------------------------------------------------------------------------
// #region constructor

// #endregion

// ---------------------------------------------------------------------------
// #region public

func (x *DIObject) AsSingleton() {
	x.instanceType = Singleton
}

// #endregion

// ---------------------------------------------------------------------------
// #region private

// #endregion
