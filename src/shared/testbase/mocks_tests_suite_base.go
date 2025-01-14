package testbase

import (
	mocks "github.com/otterize/intents-operator/src/operator/controllers/intents_reconcilers/mocks"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"k8s.io/client-go/tools/record"
)

const (
	fakeRecorderBufferSize = 100
)

type MocksSuiteBase struct {
	suite.Suite
	Controller *gomock.Controller
	Recorder   *record.FakeRecorder
	Client     *mocks.MockClient
}

func (s *MocksSuiteBase) SetupTest() {
	s.Controller = gomock.NewController(s.T())
	s.Client = mocks.NewMockClient(s.Controller)
	s.Recorder = record.NewFakeRecorder(fakeRecorderBufferSize)
}

func (s *MocksSuiteBase) TearDownTest() {
	s.ExpectNoEvent()
	s.Recorder = nil
	s.Client = nil
	s.Controller.Finish()
}

func (s *MocksSuiteBase) ExpectEvent(expectedEventReason string) {
	select {
	case event := <-s.Recorder.Events:
		s.Require().Contains(event, expectedEventReason)
	default:
		s.Fail("Expected event not found")
	}
}

func (s *MocksSuiteBase) ExpectNoEvent() {
	select {
	case event := <-s.Recorder.Events:
		s.Fail("Unexpected event found", event)
	default:
		// Amazing, no events left behind!
	}
}
