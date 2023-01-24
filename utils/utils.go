package utils

type httpMethod string

const (
	GET    httpMethod = "GET"
	POST   httpMethod = "POST"
	DELETE httpMethod = "DELETE"
	PATCH  httpMethod = "PATCH"
)

func Checkmethod(method string, checkmethod httpMethod) bool {
	return method == string(checkmethod)
}

// func ResponseWriter(c *gin.Context, status int, data interface{}, message string) {
// 	c.JSON(status, gin.H{
// 		"code":    status,
// 		"message": message,
// 		"data":    data})

// }
