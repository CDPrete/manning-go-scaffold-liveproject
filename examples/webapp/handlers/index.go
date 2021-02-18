package handlers

import (
	"fmt"
	"net/http"
)

func HandleIndex(w http.ResponseWriter, _ *http.Request)  {
	var html = `
<html>
	<body>
		<h1>Hello World from Website</h1>
	</body>
</html>
`
	_, _ = fmt.Fprint(w, html)
}
