package commands

import (
	"reflect"
	"testing"

	"github.com/helewud/redis-clone/resp"
)

func TestPing(t *testing.T) {
	tests := []struct {
		name string
		args []resp.Value
		want resp.Value
	}{
		{
			name: "simple ping",
			args: []resp.Value{},
			want: resp.Value{T: resp.RespTString, String: "PONG"},
		},
		{
			name: "ping with argument",
			args: []resp.Value{{T: resp.RespTBulk, Bulk: "hello"}},
			want: resp.Value{T: resp.RespTString, String: "hello"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ping(tt.args)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ping() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetAndGet(t *testing.T) {
	// Clear the map before testing
	SETs = map[string]string{}

	tests := []struct {
		name     string
		setArgs  []resp.Value
		getArgs  []resp.Value
		wantSet  resp.Value
		wantGet  resp.Value
		scenario string
	}{
		{
			name: "simple set and get",
			setArgs: []resp.Value{
				{T: resp.RespTBulk, Bulk: "key1"},
				{T: resp.RespTBulk, Bulk: "resp.value1"},
			},
			getArgs: []resp.Value{
				{T: resp.RespTBulk, Bulk: "key1"},
			},
			wantSet:  resp.Value{T: resp.RespTString, String: "OK"},
			wantGet:  resp.Value{T: resp.RespTBulk, Bulk: "resp.value1"},
			scenario: "success",
		},
		{
			name: "get non-existent key",
			getArgs: []resp.Value{
				{T: resp.RespTBulk, Bulk: "nonexistent"},
			},
			wantGet:  resp.Value{T: resp.RespTNull},
			scenario: "get_only",
		},
		{
			name: "set with wrong args",
			setArgs: []resp.Value{
				{T: resp.RespTBulk, Bulk: "key1"},
			},
			wantSet: resp.Value{
				T:      resp.RespTError,
				String: "ERR wrong number of arguments for 'SET' command",
			},
			scenario: "set_error",
		},
		{
			name: "get with wrong args",
			getArgs: []resp.Value{
				{T: resp.RespTBulk, Bulk: "key1"},
				{T: resp.RespTBulk, Bulk: "extra"},
			},
			wantGet: resp.Value{
				T:      resp.RespTError,
				String: "ERR wrong number of arguments for 'GET' command",
			},
			scenario: "get_error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.scenario {
			case "success":
				got := set(tt.setArgs)
				if !reflect.DeepEqual(got, tt.wantSet) {
					t.Errorf("set() = %v, want %v", got, tt.wantSet)
				}
				got = get(tt.getArgs)
				if !reflect.DeepEqual(got, tt.wantGet) {
					t.Errorf("get() = %v, want %v", got, tt.wantGet)
				}
			case "get_only":
				got := get(tt.getArgs)
				if !reflect.DeepEqual(got, tt.wantGet) {
					t.Errorf("get() = %v, want %v", got, tt.wantGet)
				}
			case "set_error":
				got := set(tt.setArgs)
				if !reflect.DeepEqual(got, tt.wantSet) {
					t.Errorf("set() = %v, want %v", got, tt.wantSet)
				}
			case "get_error":
				got := get(tt.getArgs)
				if !reflect.DeepEqual(got, tt.wantGet) {
					t.Errorf("get() = %v, want %v", got, tt.wantGet)
				}
			}
		})
	}
}

func TestSetAndDel(t *testing.T) {

	tests := []struct {
		name     string
		scenario string

		setArgs []resp.Value
		wantSet resp.Value

		delArgs []resp.Value
		wantDel resp.Value

		getArgs []resp.Value
		wantGet resp.Value
	}{
		{
			name:     "simple set and del",
			scenario: "success",

			setArgs: []resp.Value{
				{T: resp.RespTBulk, Bulk: "key1"},
				{T: resp.RespTBulk, Bulk: "resp.value1"},
			},
			wantSet: resp.Value{T: resp.RespTString, String: "OK"},

			delArgs: []resp.Value{
				{T: resp.RespTBulk, Bulk: "key1"},
			},
			wantDel: resp.Value{T: resp.RespTString, String: "OK"},
		},
		{
			name:     "simple set, del and get",
			scenario: "del_and_get",

			setArgs: []resp.Value{
				{T: resp.RespTBulk, Bulk: "key1"},
				{T: resp.RespTBulk, Bulk: "resp.value1"},
			},
			wantSet: resp.Value{T: resp.RespTString, String: "OK"},

			delArgs: []resp.Value{
				{T: resp.RespTBulk, Bulk: "key1"},
			},
			wantDel: resp.Value{T: resp.RespTString, String: "OK"},

			getArgs: []resp.Value{
				{T: resp.RespTBulk, Bulk: "key1"},
			},
			wantGet: resp.Value{T: resp.RespTNull},
		},
		{
			name:     "del non-existent key",
			scenario: "del_only",

			delArgs: []resp.Value{
				{T: resp.RespTBulk, Bulk: "nonexistent"},
			},
			wantDel: resp.Value{T: resp.RespTNull},
		},
		{
			name:     "del with wrong args",
			scenario: "del_error",

			delArgs: []resp.Value{
				{T: resp.RespTBulk, Bulk: "key1"},
				{T: resp.RespTBulk, Bulk: "extra"},
			},

			wantDel: resp.Value{
				T:      resp.RespTError,
				String: "ERR wrong number of arguments for 'DEL' command",
			},
		},
	}

	for _, tt := range tests {
		// Clear the map before testing
		SETs = map[string]string{}

		t.Run(tt.name, func(t *testing.T) {
			switch tt.scenario {
			case "success":
				got := set(tt.setArgs)
				if !reflect.DeepEqual(got, tt.wantSet) {
					t.Errorf("set() = %v, want %v", got, tt.wantSet)
				}
				got = del(tt.delArgs)
				if !reflect.DeepEqual(got, tt.wantDel) {
					t.Errorf("del() = %v, want %v", got, tt.wantDel)
				}
			case "del_only":
				got := del(tt.delArgs)
				if !reflect.DeepEqual(got, tt.wantDel) {
					t.Errorf("del() = %v, want %v", got, tt.wantDel)
				}
			case "del_and_get":
				got := set(tt.setArgs)
				if !reflect.DeepEqual(got, tt.wantSet) {
					t.Errorf("set() = %v, want %v", got, tt.wantSet)
				}
				got = del(tt.delArgs)
				if !reflect.DeepEqual(got, tt.wantDel) {
					t.Errorf("del() = %v, want %v", got, tt.wantDel)
				}
				got = get(tt.delArgs)
				if !reflect.DeepEqual(got, tt.wantGet) {
					t.Errorf("del() = %v, want %v", got, tt.wantGet)
				}
			case "set_error":
				got := set(tt.setArgs)
				if !reflect.DeepEqual(got, tt.wantSet) {
					t.Errorf("set() = %v, want %v", got, tt.wantSet)
				}
			case "del_error":
				got := del(tt.delArgs)
				if !reflect.DeepEqual(got, tt.wantDel) {
					t.Errorf("del() = %v, want %v", got, tt.wantDel)
				}
			}
		})
	}
}
