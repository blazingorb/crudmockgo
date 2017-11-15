# Mock Storage in Go
This package is intended for mocking a document based datastore.

Basic Usage
```golang
    s := mockstorage.NewMockStorage()
    s.Store("id", 1234)
    
    data := s.Load("id").(int)
    fmt.Println(data) //1234
    
    list := s.List()
    fmt.Println(list) //[{1234}]
    s.Clear()
```
