package index

import (
	"net/http"

	"github.com/Neukod-Academy/neukod-backend/utils"
)

func RetreiveHomepage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
		return
	}
	res := utils.HttpResponseBody{
		Status:  http.StatusOK,
		Message: "Welcome to homepage",
	}
	res.UpdateHttpResponse(w)
}
