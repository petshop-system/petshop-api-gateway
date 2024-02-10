package server

import (
	"fmt"
	"github.com/google/uuid"
	database "github.com/petshop-system/petshop-api-gateway/server/db"
	"go.uber.org/zap"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

type ServeReverseProxyPass struct {
	LoggerSugar  *zap.SugaredLogger
	RouterConfig map[string]map[string]string
}

func NewServerPass(loggerSugar *zap.SugaredLogger, gatewayDB *database.GatewayDB) ServeReverseProxyPass {
	serveReverseProxyPass := ServeReverseProxyPass{
		LoggerSugar: loggerSugar,
	}

	serveReverseProxyPass.LoadRouterConfig(loggerSugar, gatewayDB)

	return serveReverseProxyPass
}

func (h *ServeReverseProxyPass) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	partsPath := strings.Split(r.URL.Path, "/")
	initialPath, hostRedirect, appContext := h.getRouterConfigInfo(partsPath)

	random, _ := uuid.NewRandom()
	requestID := fmt.Sprintf("%s.%d", random.String(), time.Now().UnixNano())
	logger := h.LoggerSugar.With("host", r.Host, "request_uri", r.RequestURI, "request_url_path", r.URL.Path,
		"initialPath", initialPath, "request_id", requestID)
	logger.Infow("request server pass received")

	reverseProxy, newRequestURI, newURLPath, err := h.buildReverseProxy(hostRedirect, appContext, requestID, r)
	if err != nil {
		h.LoggerSugar.Errorw("error to process url destination",
			"host_redirect", hostRedirect)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	reverseProxy.ServeHTTP(w, r)
	logger.Infow("server pass done", "new_host", hostRedirect,
		"new_request_uri", newRequestURI, "new_url_path", newURLPath)
}

func (h *ServeReverseProxyPass) buildReverseProxy(hostRedirect, appContext, requestID string, r *http.Request) (*httputil.ReverseProxy, string, string, error) {

	newRequestURI := strings.ReplaceAll(r.RequestURI, "petshop-system", appContext)
	newURLPath := strings.ReplaceAll(r.URL.Path, "petshop-system", appContext)

	destinationTo := fmt.Sprintf("%s%s", hostRedirect, newRequestURI)
	destination, err := url.Parse(destinationTo)
	if err != nil {
		return nil, "", "", err
	}

	rp := httputil.NewSingleHostReverseProxy(destination)
	rp.Director = func(req *http.Request) {
		req.Host = destination.Host
		req.URL.Scheme = destination.Scheme
		req.URL.Host = destination.Host
		//req.URL.Path = singleJoiningSlash(u.Path, req.URL.Path)
		//if targetQuery == "" || req.URL.RawQuery == "" {
		//	req.URL.RawQuery = targetQuery + req.URL.RawQuery
		//} else {
		//	req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		//}

		req.RequestURI = newRequestURI
		req.URL.Path = destination.Path
		req.Header.Set("request_id", requestID)
	}

	rp.ModifyResponse = func(w *http.Response) error {
		w.Header.Set("request_id", requestID)
		return nil
	}

	return rp, newRequestURI, newURLPath, nil
}

func (h *ServeReverseProxyPass) getRouterConfigInfo(partsPath []string) (string, string, string) {

	initialPath := partsPath[2]
	hostRedirect := h.RouterConfig[initialPath]["host"]
	appContext := h.RouterConfig[initialPath]["app-context"]

	return initialPath, hostRedirect, appContext
}

func (h *ServeReverseProxyPass) LoadRouterConfig(loggerSugar *zap.SugaredLogger, gatewayDDB *database.GatewayDB) {

	loadRoutersFunc := func() map[string]map[string]string {

		routersDB := gatewayDDB.GetAllRouter()
		routers := make(map[string]map[string]string, len(routersDB))

		for _, router := range routersDB {

			host, _ := GetMapValueFromJsonRawMessage[string](router.Configuration, "host")
			appContext, _ := GetMapValueFromJsonRawMessage[string](router.Configuration, "app-context")
			routers[router.Router] = map[string]string{
				"host":        host.(string),
				"app-context": appContext.(string),
			}

		}

		loggerSugar.Infow("loaded routers", "routers", routers)
		return routers
	}

	h.RouterConfig = loadRoutersFunc()

	go func() {
		for range time.Tick(1 * time.Minute) {
			h.RouterConfig = loadRoutersFunc()
		}
	}()

}
