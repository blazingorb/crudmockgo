package mockstorage

import "testing"

func TestMockStorage(t *testing.T) {
	type Data struct {
		Name string
	}
	id1 := "1"
	data1 := &Data{"Denis"}
	id2 := "2"
	data2 := &Data{"Tiffany"}

	s := NewMockStorage()
	s.Store(id1, data1)
	s.Store(id2, data2)
	if len(s.store) != 2 {
		t.Errorf("Expected Store Length %d but got %d", 2, len(s.store))
	}

	data1Load := s.Load(id1).(*Data)
	if data1Load.Name != data1.Name {
		t.Errorf("Expected Data Name %s but got %s", data1.Name, data1Load.Name)
	}
	data2Load := s.Load(id2).(*Data)
	if data2Load.Name != data2.Name {
		t.Errorf("Expected Data Name %s but got %s", data2.Name, data2Load.Name)
	}
	_ = s.Load("")
	list := s.List()
	if len(list) != 2 {
		t.Errorf("Expected List Length %d but got %d", 2, len(list))
	}

	s.Clear()
	if len(s.store) != 0 {
		t.Errorf("Expected Store Length %d but got %d", 0, len(s.store))
	}
}
