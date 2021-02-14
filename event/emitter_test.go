package event

import (
	"testing"
	"time"
)

var testEvent Type = "test"

func TestEmit(t *testing.T) {
	type args struct {
		eventType Type
		param1    string
		param2    string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test 1",
			args: args{
				eventType: testEvent,
				param1:    "First param from unit test",
				param2:    "Second param from unit test",
			},
		},
	}

	var receivedParam1 string
	var receivedParam2 string
	Subscribe(testEvent, func(args ...interface{}) () {
		p1 := args[0].(string)
		p2 := args[1].(string)
		receivedParam1 = p1
		receivedParam2 = p2
	})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Emit(testEvent, tt.args.param1, tt.args.param2)
			time.Sleep(1 * time.Second)
			if receivedParam1 != tt.args.param1 || receivedParam2 != tt.args.param2 {
				t.Errorf("Emitted event not received by the handler")
			}
		})
	}
}
