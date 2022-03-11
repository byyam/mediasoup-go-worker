package utils

import "testing"

type Person struct {
	Age  int
	Name string
}

func TestNewHashMap(t *testing.T) {
	hashmap := NewHashMap()
	// case1
	hashmap.Store("p1", "c1", nil)
	_, ok := hashmap.Load("p1", "c1")
	if ok {
		t.Fatal("load failed")
	}

	// case2
	hashmap.Store("p1", "c1", &Person{
		Age:  10,
		Name: "hello",
	})
	value, ok := hashmap.Load("p1", 2)
	if ok {
		t.Fatal("id2 should not be load")
	}

	// case3
	value, ok = hashmap.Load("p1", "c1")
	if !ok {
		t.Fatal("load failed")
	}
	person, ok := value.(*Person)
	if !ok {
		t.Fatal("load person failed")
	}
	t.Logf("person:%+v", person)

	// case4
	hashmap.Delete("p1", "c1")
	value, ok = hashmap.Load("p1", "c1")
	if ok {
		t.Fatal("id2 should not be load")
	}

	// case5
	hashmap.Store("p1", "c1", &Person{
		Age:  10,
		Name: "hello",
	})
	value, ok = hashmap.Load("p1", "c1")
	if !ok {
		t.Fatal("load failed")
	}
	person, ok = value.(*Person)
	if !ok {
		t.Fatal("load person failed")
	}
	hashmap.Erase("p1")
	value, ok = hashmap.Load("p1", "c1")
	if ok {
		t.Fatal("id2 should not be load")
	}
}
