package gotickreloader

import (
	"errors"
	"runtime"
	"testing"
	"time"
)

var initialGoRoutineCount = runtime.NumGoroutine() + 1

func TestProcess(t *testing.T) {
	var fGetter = func(...interface{}) (interface{}, error) {
		return true, nil
	}
	var tr = NewClient(1*time.Second, fGetter)
	tr.StartTickReload()
	if runtime.NumGoroutine() != initialGoRoutineCount+1 {
		t.Fatalf("unexpected go routine count %d", runtime.NumGoroutine())
	}
	var v, err = tr.Get()
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	if true != v {
		t.Fatalf("unexpected value %v", v)
	}

	tr.StopTickReload()
	checkThatGoRoutineIsClosed(t)
}

func TestGetterOnError(t *testing.T) {
	var experr = errors.New("test")
	var fGetter = func(...interface{}) (interface{}, error) {
		return true, experr
	}
	var tr = NewClient(1*time.Second, fGetter)
	tr.StartTickReload()
	if runtime.NumGoroutine() != initialGoRoutineCount+1 {
		t.Fatalf("unexpected go routine count %d", runtime.NumGoroutine())
	}
	var _, err = tr.Get()
	if err != experr {
		t.Fatalf("unexpected error %v", err)
	}
	tr.StopTickReload()
	checkThatGoRoutineIsClosed(t)
}

func TestTickReload(t *testing.T) {

	var experr = errors.New("test")
	var param = 0
	var fGetter = func(v ...interface{}) (interface{}, error) {
		p, ok := v[0].(*int)
		if ok {
			*p++
		}
		if *p == 3 {
			return nil, experr
		}
		return *p, nil
	}
	var tr = NewClient(5*time.Millisecond, fGetter, &param)
	tr.StartTickReload()
	var v interface{}
	var err error
	for {
		v, err = tr.Get()
		t.Logf("v: %T %v param: %T %v err: %T %v", v, v, param, param, err, err)
		if param == 3 {
			if err == nil {
				t.Fatalf("unexpected nil error %v", err)
			}
			if v != nil {
				t.Fatalf("unexpected value %v", v)
			}
		} else {
			if err != nil {
				t.Fatalf("unexpected error %v", err)
			}
			if v.(int) != param {
				t.Fatalf("unexpected value %v vs %v", v, param)
			}
			if param == 4 {
				break
			}
		}
		time.Sleep(1 * time.Millisecond)
	}
	tr.StopTickReload()
	checkThatGoRoutineIsClosed(t)
}

func checkThatGoRoutineIsClosed(t *testing.T) {
	var i int
	for {

		if initialGoRoutineCount == runtime.NumGoroutine() {
			t.Logf("waiting stop %d == %d", initialGoRoutineCount, runtime.NumGoroutine())
			break
		}

		t.Logf("waiting %d goroutines", runtime.NumGoroutine())
		if i == 10 {
			t.Fatal("all go routines are not closed")
		}
		i++
		time.Sleep(5 * time.Millisecond)
	}
}
