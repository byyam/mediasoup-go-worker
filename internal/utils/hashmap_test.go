package utils

import (
	"encoding/json"
	"fmt"
	"testing"
)

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

func TestHashMapGet(t *testing.T) {
	hashmap := NewHashMap()
	// case1
	hashmap.Store("p", "c1", &Person{
		Age:  10,
		Name: "hello",
	})
	hashmap.Store("p", "c2", &Person{
		Age:  12,
		Name: "world",
	})
	value, ok := hashmap.Get("p")
	if !ok {
		t.Fatal("get failed")
	}
	r, ok := value.(map[interface{}]interface{})
	if !ok {
		t.Fatal("get value failed")
	}
	for k, v := range r {
		t.Logf("k:%s, v:%+v", k.(string), v.(*Person))
	}
}

func TestJson(t *testing.T) {
	type Person struct {
		Age int64 `json:"age"`
	}
	p := &Person{Age: 18528326348443648}
	data, _ := json.Marshal(p)
	t.Logf("p:%s", string(data))

	var q map[string]interface{}
	// wrong
	_ = json.Unmarshal(data, &q)
	// correct
	//d := json.NewDecoder(bytes.NewReader(data))
	//d.UseNumber()
	//_ = d.Decode(&q)
	t.Logf("q:%+v", q["age"])

	var ret Person
	_ = InterfaceToStruct(q, &ret)
	t.Logf("ret:%+v", ret)
}

func TestJsonString(t *testing.T) {
	var test interface{}
	// str := `{"id":18502511176978432, "name":"golang"}`
	str := `{"id":18528326348443648, "name":"golang"}`
	err := json.Unmarshal([]byte(str), &test)
	if err != nil {
		fmt.Println(err)
	}
	m := test.(map[string]interface{})
	fmt.Printf("type:%T, value:%v\n", m["id"], m["id"])
}

func InterfaceToStruct(i interface{}, dst interface{}) error {
	data, err := json.Marshal(i)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, dst)
	return err
}
