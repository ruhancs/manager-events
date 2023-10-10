package events

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TestEvent struct {
	Name    string
	Payload interface{}
}

func (event *TestEvent) GetName() string {
	return event.Name
}

func (event *TestEvent) GetDateTime() time.Time {
	return time.Now()
}

func (event *TestEvent) GetPayload() any {
	return event.Payload
}

type TestEventHandler struct {
	ID int
}

func (ed *TestEventHandler) Handle(event EventInterface, wg *sync.WaitGroup) {
	wg.Done()
}

type EventDispatcherTestSuite struct {
	suite.Suite
	event           TestEvent
	event2          TestEvent
	handler         TestEventHandler
	handler2        TestEventHandler
	handler3        TestEventHandler
	eventDispatcher *EventDispatcher
}

func (suite *EventDispatcherTestSuite) SetupTest() {
	suite.eventDispatcher = NewEventDispatcher()
	suite.handler = TestEventHandler{ID: 1}
	suite.handler2 = TestEventHandler{ID: 2}
	suite.handler3 = TestEventHandler{ID: 3}
	suite.event = TestEvent{Name: "test1", Payload: "test1"}
	suite.event2 = TestEvent{Name: "test2", Payload: "test2"}
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Register() {
	err := suite.eventDispatcher.Register(suite.event.GetName(),&suite.handler)
	suite.Nil(err)
	suite.Len(suite.eventDispatcher.handlers[suite.event.GetName()],1)
	
	err = suite.eventDispatcher.Register(suite.event.GetName(),&suite.handler2)
	suite.Nil(err)
	suite.Len(suite.eventDispatcher.handlers[suite.event.GetName()],2)

	suite.Equal(&suite.handler, suite.eventDispatcher.handlers[suite.event.GetName()][0])
	suite.Equal(&suite.handler2, suite.eventDispatcher.handlers[suite.event.GetName()][1])
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Register_WithSameHandler() {
	err := suite.eventDispatcher.Register(suite.event.GetName(),&suite.handler)
	suite.Nil(err)
	suite.Len(suite.eventDispatcher.handlers[suite.event.GetName()],1)
	
	err = suite.eventDispatcher.Register(suite.event.GetName(),&suite.handler)
	suite.NotNil(err)
	suite.Equal(errHandlerAlreadyRegistered,err)
	suite.Len(suite.eventDispatcher.handlers[suite.event.GetName()],1)
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Clear() {
	err := suite.eventDispatcher.Register(suite.event.GetName(),&suite.handler)
	suite.Nil(err)
	suite.Len(suite.eventDispatcher.handlers[suite.event.GetName()],1)
	
	err = suite.eventDispatcher.Register(suite.event.GetName(),&suite.handler2)
	suite.Nil(err)
	suite.Len(suite.eventDispatcher.handlers[suite.event.GetName()],2)
	
	err = suite.eventDispatcher.Register(suite.event2.GetName(),&suite.handler3)
	suite.Nil(err)
	suite.Len(suite.eventDispatcher.handlers[suite.event2.GetName()],1)
	
	err = suite.eventDispatcher.Clear()
	suite.Nil(err)
	suite.Len(suite.eventDispatcher.handlers,0)
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Has() {
	err := suite.eventDispatcher.Register(suite.event.GetName(),&suite.handler)
	suite.Nil(err)
	suite.Len(suite.eventDispatcher.handlers[suite.event.GetName()],1)
	
	err = suite.eventDispatcher.Register(suite.event.GetName(),&suite.handler2)
	suite.Nil(err)
	suite.Len(suite.eventDispatcher.handlers[suite.event.GetName()],2)

	suite.True(suite.eventDispatcher.Has(suite.event.GetName(), &suite.handler))
	suite.True(suite.eventDispatcher.Has(suite.event.GetName(), &suite.handler2))
	suite.False(suite.eventDispatcher.Has(suite.event.GetName(), &suite.handler3))
}

type MockHandler struct {
	mock.Mock
}
func(mock *MockHandler) Handle(event EventInterface, wg *sync.WaitGroup) {
	mock.Called(event)
	wg.Done()
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Dispatch() {
	eventHandler := &MockHandler{}
	eventHandler.On("Handle", &suite.event)
	suite.eventDispatcher.Register(suite.event.GetName(), eventHandler)
	suite.eventDispatcher.Dispatch(&suite.event)
	
	eventHandler.AssertExpectations(suite.T())
	eventHandler.AssertNumberOfCalls(suite.T(),"Handle",1)
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Remove() {
	err := suite.eventDispatcher.Register(suite.event.GetName(),&suite.handler)
	suite.Nil(err)
	suite.Len(suite.eventDispatcher.handlers[suite.event.GetName()],1)
	
	err = suite.eventDispatcher.Register(suite.event.GetName(),&suite.handler2)
	suite.Nil(err)
	suite.Len(suite.eventDispatcher.handlers[suite.event.GetName()],2)
	
	suite.eventDispatcher.Remove(suite.event.GetName(), &suite.handler)
	suite.Len(suite.eventDispatcher.handlers[suite.event.GetName()],1)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(EventDispatcherTestSuite))
}
