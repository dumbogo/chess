// Code generated by MockGen. DO NOT EDIT.
// Source: engine/board.go

// Package engine is a generated GoMock package.
package engine

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockBoard is a mock of Board interface.
type MockBoard struct {
	ctrl     *gomock.Controller
	recorder *MockBoardMockRecorder
}

// MockBoardMockRecorder is the mock recorder for MockBoard.
type MockBoardMockRecorder struct {
	mock *MockBoard
}

// NewMockBoard creates a new mock instance.
func NewMockBoard(ctrl *gomock.Controller) *MockBoard {
	mock := &MockBoard{ctrl: ctrl}
	mock.recorder = &MockBoardMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBoard) EXPECT() *MockBoardMockRecorder {
	return m.recorder
}

// BlackPlayer mocks base method.
func (m *MockBoard) BlackPlayer() *Player {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BlackPlayer")
	ret0, _ := ret[0].(*Player)
	return ret0
}

// BlackPlayer indicates an expected call of BlackPlayer.
func (mr *MockBoardMockRecorder) BlackPlayer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BlackPlayer", reflect.TypeOf((*MockBoard)(nil).BlackPlayer))
}

// EatPiece mocks base method.
func (m *MockBoard) EatPiece(loc SquareIdentifier) Piece {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EatPiece", loc)
	ret0, _ := ret[0].(Piece)
	return ret0
}

// EatPiece indicates an expected call of EatPiece.
func (mr *MockBoardMockRecorder) EatPiece(loc interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EatPiece", reflect.TypeOf((*MockBoard)(nil).EatPiece), loc)
}

// Squares mocks base method.
func (m *MockBoard) Squares() Squares {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Squares")
	ret0, _ := ret[0].(Squares)
	return ret0
}

// Squares indicates an expected call of Squares.
func (mr *MockBoardMockRecorder) Squares() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Squares", reflect.TypeOf((*MockBoard)(nil).Squares))
}

// String mocks base method.
func (m *MockBoard) String() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "String")
	ret0, _ := ret[0].(string)
	return ret0
}

// String indicates an expected call of String.
func (mr *MockBoardMockRecorder) String() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "String", reflect.TypeOf((*MockBoard)(nil).String))
}

// WhitePlayer mocks base method.
func (m *MockBoard) WhitePlayer() *Player {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WhitePlayer")
	ret0, _ := ret[0].(*Player)
	return ret0
}

// WhitePlayer indicates an expected call of WhitePlayer.
func (mr *MockBoardMockRecorder) WhitePlayer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WhitePlayer", reflect.TypeOf((*MockBoard)(nil).WhitePlayer))
}