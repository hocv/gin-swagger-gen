# gin-swagger-gen

generate gin swagger comment

## install

```go
go install github.com/hocv/gin-swagger-gen
```

## params

| param         | short | default | desc                                 |
| ------------- | ----- | ------- | ------------------------------------ |
| dir           | d     | ./      | project root dir                     |
| func.name     | f     | -       | specify the funcion to add comment   |
| info.disable  | -     | false   | add comment add main function        |
| info.all      | -     | false   | add all comment to main function     |
| info.license  | -     | false   | add license comment to main function |
| info.tags     | -     | false   | add tags comment to main function    |
| info.security | -     | false   | add security coment to main function |

## features

1. add comment to main function
2. add comment to gin handler function
   1. route, method
   2. params in path, query, form
   3. produce, status code
   4. accept

## example

```go
func route() {
    g := gin.Default()
    setRoute(g)
    _ = g.Run(":9090")
}

func setRoute(a *gin.Engine) {
    a.GET("/api/:id", normalHandle)
}

// @Summary normalHandle
// @Description normalHandle
// @Accept json,multipart/form-data
// @Produce string
// @Param id path string true "id"
// @Param lg body login true "lg"
// @Param q1 query string true "q1"
// @Param q2 query string true "q2 default 0"
// @Param f1 formData string true "f1"
// @Failure 400 {string} string
// @Success 200 {string} string
// @Router /api/{id} [get]
func normalHandle(c *gin.Context) {
    lg := &login{}
    q := c.Query("q1")
    b := c.DefaultQuery("q2", "0")
    f, _ := c.GetPostForm("f1")
    fmt.Println(q, b, f)
    if err := c.BindJSON(lg); err != nil {
        c.String(http.StatusBadRequest, "f")
        return
    }
    resp(c)
}

func resp(d *gin.Context) {
    d.String(200, "f")
}
```

### bugs

1. not supported `BingQuery` for now
2. Other unknown
