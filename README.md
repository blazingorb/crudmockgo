<<<<<<< HEAD
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
=======
[![Go Report Card](https://goreportcard.com/badge/github.com/blazingorb/mockstoragego)](https://goreportcard.com/report/github.com/blazingorb/mockstoragego)

# Mock JSON Storage in Go
This package is intended for rapid mocking of front-end applications that requires some persistent data.

```sh
go install
mockstoragego -p 8000
>>>>>>> b7b8e2007a24a1e4588c9298e39dd1116514af4d
```
