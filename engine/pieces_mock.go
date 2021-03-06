// Code generated by MockGen. DO NOT EDIT.
// Source: engine/pieces.go

// Package engine is a generated GoMock package.
package engine

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockPiece is a mock of Piece interface.
type MockPiece struct {
	ctrl     *gomock.Controller
	recorder *MockPieceMockRecorder
}

// MockPieceMockRecorder is the mock recorder for MockPiece.
type MockPieceMockRecorder struct {
	mock *MockPiece
}

// NewMockPiece creates a new mock instance.
func NewMockPiece(ctrl *gomock.Controller) *MockPiece {
	mock := &MockPiece{ctrl: ctrl}
	mock.recorder = &MockPieceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPiece) EXPECT() *MockPieceMockRecorder {
	return m.recorder
}

// CanMove mocks base method.
func (m *MockPiece) CanMove(board Board, movements []Movement, from, to Square) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CanMove", board, movements, from, to)
	ret0, _ := ret[0].(bool)
	return ret0
}

// CanMove indicates an expected call of CanMove.
func (mr *MockPieceMockRecorder) CanMove(board, movements, from, to interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CanMove", reflect.TypeOf((*MockPiece)(nil).CanMove), board, movements, from, to)
}

// Color mocks base method.
func (m *MockPiece) Color() Color {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Color")
	ret0, _ := ret[0].(Color)
	return ret0
}

// Color indicates an expected call of Color.
func (mr *MockPieceMockRecorder) Color() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Color", reflect.TypeOf((*MockPiece)(nil).Color))
}

// Identifier mocks base method.
func (m *MockPiece) Identifier() PieceIdentifier {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Identifier")
	ret0, _ := ret[0].(PieceIdentifier)
	return ret0
}

// Identifier indicates an expected call of Identifier.
func (mr *MockPieceMockRecorder) Identifier() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Identifier", reflect.TypeOf((*MockPiece)(nil).Identifier))
}

// String mocks base method.
func (m *MockPiece) String() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "String")
	ret0, _ := ret[0].(string)
	return ret0
}

// String indicates an expected call of String.
func (mr *MockPieceMockRecorder) String() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "String", reflect.TypeOf((*MockPiece)(nil).String))
}
