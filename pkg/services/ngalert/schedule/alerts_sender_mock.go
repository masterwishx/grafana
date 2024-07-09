// Code generated by mockery v2.34.2. DO NOT EDIT.

package schedule

import (
	context "context"

	definitions "github.com/grafana/grafana/pkg/services/ngalert/api/tooling/definitions"
	mock "github.com/stretchr/testify/mock"

	models "github.com/grafana/grafana/pkg/services/ngalert/models"
)

// AlertsSenderMock is an autogenerated mock type for the AlertsSender type
type AlertsSenderMock struct {
	mock.Mock
}

type AlertsSenderMock_Expecter struct {
	mock *mock.Mock
}

func (_m *AlertsSenderMock) EXPECT() *AlertsSenderMock_Expecter {
	return &AlertsSenderMock_Expecter{mock: &_m.Mock}
}

// Send provides a mock function with given fields: ctx, key, alerts
func (_m *AlertsSenderMock) Send(ctx context.Context, key models.AlertRuleKey, alerts definitions.PostableAlerts) {
	_m.Called(ctx, key, alerts)
}

// AlertsSenderMock_Send_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Send'
type AlertsSenderMock_Send_Call struct {
	*mock.Call
}

// Send is a helper method to define mock.On call
//   - ctx context.Context
//   - key models.AlertRuleKey
//   - alerts definitions.PostableAlerts
func (_e *AlertsSenderMock_Expecter) Send(ctx interface{}, key interface{}, alerts interface{}) *AlertsSenderMock_Send_Call {
	return &AlertsSenderMock_Send_Call{Call: _e.mock.On("Send", ctx, key, alerts)}
}

func (_c *AlertsSenderMock_Send_Call) Run(run func(ctx context.Context, key models.AlertRuleKey, alerts definitions.PostableAlerts)) *AlertsSenderMock_Send_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(models.AlertRuleKey), args[2].(definitions.PostableAlerts))
	})
	return _c
}

func (_c *AlertsSenderMock_Send_Call) Return() *AlertsSenderMock_Send_Call {
	_c.Call.Return()
	return _c
}

func (_c *AlertsSenderMock_Send_Call) RunAndReturn(run func(context.Context, models.AlertRuleKey, definitions.PostableAlerts)) *AlertsSenderMock_Send_Call {
	_c.Call.Return(run)
	return _c
}

// NewAlertsSenderMock creates a new instance of AlertsSenderMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAlertsSenderMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *AlertsSenderMock {
	mock := &AlertsSenderMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}