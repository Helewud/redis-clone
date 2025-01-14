package commands

import (
	"reflect"
	"testing"

	"github.com/helewud/redis-clone/resp"
)

func TestHash(t *testing.T) {
	// Clear the map before testing
	HSETs = map[string]map[string]string{}

	tests := []struct {
		name     string
		scenario string

		hsetArgs []resp.Value
		wantHset resp.Value

		hgetArgs []resp.Value
		wantHget resp.Value

		hgetallArgs []resp.Value
		wantHgetall resp.Value
	}{
		{
			name:     "simple hash operations",
			scenario: "success",

			hsetArgs: []resp.Value{
				{T: resp.RespTBulk, Bulk: "hash1"},
				{T: resp.RespTBulk, Bulk: "field1"},
				{T: resp.RespTBulk, Bulk: "value1"},
			},
			wantHset: resp.Value{T: resp.RespTString, String: "OK"},

			hgetArgs: []resp.Value{
				{T: resp.RespTBulk, Bulk: "hash1"},
				{T: resp.RespTBulk, Bulk: "field1"},
			},
			wantHget: resp.Value{T: resp.RespTBulk, Bulk: "value1"},

			hgetallArgs: []resp.Value{
				{T: resp.RespTBulk, Bulk: "hash1"},
			},
			wantHgetall: resp.Value{
				T: "array",
				Array: []resp.Value{
					{T: resp.RespTBulk, Bulk: "field1"},
					{T: resp.RespTBulk, Bulk: "value1"},
				},
			},
		},
		{
			name:     "hget non-existent field",
			scenario: "hget_only",

			hgetArgs: []resp.Value{
				{T: resp.RespTBulk, Bulk: "hash1"},
				{T: resp.RespTBulk, Bulk: "nonexistent"},
			},
			wantHget: resp.Value{T: resp.RespTNull},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.scenario {
			case "success":
				got := hset(tt.hsetArgs)
				if !reflect.DeepEqual(got, tt.wantHset) {
					t.Errorf("hset() = %v, want %v", got, tt.wantHset)
				}
				got = hget(tt.hgetArgs)
				if !reflect.DeepEqual(got, tt.wantHget) {
					t.Errorf("hget() = %v, want %v", got, tt.wantHget)
				}
				got = hgetall(tt.hgetallArgs)
				if !reflect.DeepEqual(got, tt.wantHgetall) {
					t.Errorf("hgetall() = %v, want %v", got, tt.wantHgetall)
				}

			case "hget_only":
				got := hget(tt.hgetArgs)
				if !reflect.DeepEqual(got, tt.wantHget) {
					t.Errorf("hget() = %v, want %v", got, tt.wantHget)
				}
			}
		})
	}
}

// TestConcurrentAccess tests thread safety of the operations
// func TestConcurrentAccess(t *testing.T) {
// 	// Clear maps before testing
// 	SETs = map[string]string{}
// 	HSETs = map[string]map[string]string{}

// 	var wg sync.WaitGroup
// 	numGoroutines := 100

// 	// Test concurrent SET/GET
// 	wg.Add(numGoroutines)
// 	for i := 0; i < numGoroutines; i++ {
// 		go func(i int) {
// 			defer wg.Done()
// 			key := "key" + string(rune(i))
// 			value := "value" + string(rune(i))

// 			// Test SET
// 			setArgs := []resp.Value{
// 				{T: resp.RespTBulk, Bulk: key},
// 				{T: resp.RespTBulk, Bulk: value},
// 			}
// 			setResult := set(setArgs)
// 			if setResult.String != "OK" {
// 				t.Errorf("concurrent set failed: %v", setResult)
// 			}

// 			// Test GET
// 			getArgs := []resp.Value{{T: resp.RespTBulk, Bulk: key}}
// 			getResult := get(getArgs)
// 			if getResult.Bulk != value {
// 				t.Errorf("concurrent get failed: got %v, want %v", getResult.Bulk, value)
// 			}
// 		}(i)
// 	}

// 	// Test concurrent HSET/HGET
// 	wg.Add(numGoroutines)
// 	for i := 0; i < numGoroutines; i++ {
// 		go func(i int) {
// 			defer wg.Done()
// 			hash := "hash" + string(rune(i))
// 			field := "field" + string(rune(i))
// 			value := "value" + string(rune(i))

// 			// Test HSET
// 			hsetArgs := []resp.Value{
// 				{T: resp.RespTBulk, Bulk: hash},
// 				{T: resp.RespTBulk, Bulk: field},
// 				{T: resp.RespTBulk, Bulk: value},
// 			}
// 			hsetResult := hset(hsetArgs)
// 			if hsetResult.String != "OK" {
// 				t.Errorf("concurrent hset failed: %v", hsetResult)
// 			}

// 			// Test HGET
// 			hgetArgs := []resp.Value{
// 				{T: resp.RespTBulk, Bulk: hash},
// 				{T: resp.RespTBulk, Bulk: field},
// 			}
// 			hgetResult := hget(hgetArgs)
// 			if hgetResult.Bulk != value {
// 				t.Errorf("concurrent hget failed: got %v, want %v", hgetResult.Bulk, value)
// 			}
// 		}(i)
// 	}

// 	wg.Wait()
// }
