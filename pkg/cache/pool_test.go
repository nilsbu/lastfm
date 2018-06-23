package cache

import (
	"errors"
	"reflect"
	"testing"
)

type Sum struct{}

func (Sum) Do(job Job) (data interface{}, err error) {
	nums, ok := job.([]int)
	if !ok {
		return nil, errors.New("wrong job")
	}

	sum := 0
	for _, x := range nums {
		sum += x
	}

	return sum, nil
}

type String struct{}

func (String) Do(job Job) (data interface{}, err error) {
	return "i am a string", nil
}

func TestWorkSum(t *testing.T) {
	cases := []struct {
		workers  []Worker
		job      interface{}
		data     int
		ctorOK   bool
		workerOK bool
		resultOK bool
	}{
		{[]Worker{}, []int{1}, 1, false, false, false}, // no workers
		{[]Worker{Sum{}}, []int{1, 2, 3, 4}, 10, true, true, true},
		{[]Worker{Sum{}}, "wrong", -1, true, false, false},
		{[]Worker{String{}}, []int{1, 2}, 3, true, true, false},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			p, err := NewPool(c.workers)
			if err != nil {
				if c.ctorOK {
					t.Error("unexpected error in constructor:", err)
				}
				return
			}

			resultChan := p.Work(c.job)
			result := <-resultChan
			if result.Err != nil && c.workerOK {
				t.Error("unexpected error in worker:", result.Err)
			} else if result.Err == nil && !c.workerOK {
				t.Error("expected error but none occurred")
			}
			if err == nil {
				data, ok := result.Data.(int)
				if ok && !c.resultOK {
					t.Error("expected result type other than 'int' but got 'int'")
				} else if !ok && c.resultOK {
					t.Errorf("result type is '%v', expected 'int'",
						reflect.TypeOf(result.Data))
				}
				if ok {
					if data != c.data {
						t.Errorf("wrong result: got '%v', expected '%v'", data, c.data)
					}
				}
			}
		})
	}
}
