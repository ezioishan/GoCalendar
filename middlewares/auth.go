package middleware

import (
	"encoding/json"
	"my-api/utils"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

type ErrorResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

// AuthenticationMiddleware checks if the user has a valid JWT token
func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenString := r.Header.Get("Authorization")
		if len(tokenString) == 0 {
			// customHTTP.NewErrorResponse(w, http.StatusUnauthorized, "Authentication failure")
			error := ErrorResponse{
				true,
				"Authentication failure",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(&error)
			return
		}
		tokenString = tokenString[len("Bearer "):]
		claims, err := utils.VerifyToken(tokenString)
		if err != nil {
			// customHTTP.NewErrorResponse(w, http.StatusUnauthorized, "Error verifying JWT token: "+err.Error())
			error := ErrorResponse{
				true,
				string("Error verifying JWT token: " + err.Error()),
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(&error)
			return
		}

		//pass userId claim to req
		//todo: find a better way to convert the claim to string
		userId := claims.(jwt.MapClaims)["user_id"].(string)
		r.Header.Set("userId", userId)
		next.ServeHTTP(w, r)
	})

}
