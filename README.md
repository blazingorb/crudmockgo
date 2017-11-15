[![Go Report Card](https://goreportcard.com/badge/github.com/blazingorb/mockstoragego)](https://goreportcard.com/report/github.com/blazingorb/mockstoragego)

# Mock JSON Storage in Go
This package is intended for rapid mocking of front-end applications that requires some persistent data.

```sh
go install
mockstoragego -p 8000
```

The HTTP listen port can be set with `-p` or "PORT" environment variable.

----
# Available Endpoints

----
# Endpoints

**Write Data**
----
  Write Any JSON Data To the Store

* **URL**

  /write

* **Method:**

  `POST`
  
*  **URL Params**

   **Required:**
 
   `id=[integer]`

* **Data Params**

  { id : "ANY_STRING", data : "ANY_OBJECT" }

* **Success Response:**

  * **Code:** 200 <br />

 
* **Error Response:**

  * **Code:** 4** Response Ranges <br />
  

* **Sample Call:**
  ```javascript
    let payload = { some: "data", num: 1 }
    $.ajax({
      url: "/write",
      type : "POST",
      body: payload,
      headers: {
        'Content-Type': 'application/json',
      },
      success : function(r) {
        console.log(r);
      }
    });
  ```

**Read Data**
----
  Returns json data based on the ID parameter sent

* **URL**

  /read

* **Method:**

  `GET`
  
*  **URL Params**

   **Required:**
 
   `id=[string]`

* **Success Response:**

  * **Code:** 200 <br />
  `{ JSON OBJECT}`

 
* **Error Response:**

  * **Code:** 4** Response Ranges <br />

* **Sample Call:**
  ```javascript
    let id = "1234"
    $.ajax({
      url: "/read?id=" + id,
      type : "GET",
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded',
      },
      success : function(r) {
        console.log(r);
      }
    });
  ```
