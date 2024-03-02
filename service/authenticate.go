package service

import (
	"go.uber.org/zap"
	"net/http"
)

type Authenticate struct {
	LoggerSugar      *zap.SugaredLogger
	ClientConnection *http.Client
}

func NewAuthenticate(loggerSugar *zap.SugaredLogger, clientConnection *http.Client) Authenticate {
	return Authenticate{
		LoggerSugar:      loggerSugar,
		ClientConnection: clientConnection,
	}
}

func (authenticate *Authenticate) Execute(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		request, _ := http.NewRequest(http.MethodPost, "", r.Body)
		cookie, _ := r.Cookie("petshop-authenticate")
		if cookie == nil || len(cookie.Value) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		resp, err := authenticate.ClientConnection.Do(request)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		if resp.StatusCode != http.StatusOK {
			w.WriteHeader(418)
			w.Write([]byte("I'm a teapot!!!!"))
			return
		}

		next(w, r)
	}
}
