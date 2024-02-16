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
	RouterConfig map[string]Router
}

type Router struct {
	Host                 string
	ReplaceOldAppContext string
	ReplaceNewAppContext string
}

func NewServerPass(loggerSugar *zap.SugaredLogger, gatewayDB *database.GatewayDB, tickerReloadRouters time.Duration) ServeReverseProxyPass {
	serveReverseProxyPass := ServeReverseProxyPass{
		LoggerSugar: loggerSugar,
	}

	serveReverseProxyPass.LoadRouterConfig(loggerSugar, gatewayDB, tickerReloadRouters)

	return serveReverseProxyPass
}

func (h *ServeReverseProxyPass) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	partsPath := strings.Split(r.URL.Path, "/")
	initialPath, router := h.getRouterConfigInfo(partsPath)

	random, _ := uuid.NewRandom()
	requestID := fmt.Sprintf("%s.%d", random.String(), time.Now().UnixNano())
	logger := h.LoggerSugar.With("host", r.Host, "request_uri", r.RequestURI, "request_url_path", r.URL.Path,
		"initialPath", initialPath, "request_id", requestID)
	logger.Infow("request server pass received")

	reverseProxy, newRequestURI, newURLPath, err := h.buildReverseProxy(requestID, r, router)
	if err != nil {
		h.LoggerSugar.Errorw("error to process url destination",
			"host_redirect", router.Host)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	reverseProxy.ServeHTTP(w, r)
	logger.Infow("server pass done", "new_host", router.Host,
		"new_request_uri", newRequestURI, "new_url_path", newURLPath)
}

func (h *ServeReverseProxyPass) buildReverseProxy(requestID string, r *http.Request, router Router) (*httputil.ReverseProxy, string, string, error) {

	newRequestURI := strings.ReplaceAll(r.RequestURI, router.ReplaceOldAppContext, router.ReplaceNewAppContext)
	newURLPath := strings.ReplaceAll(r.URL.Path, router.ReplaceOldAppContext, router.ReplaceNewAppContext)

	destinationTo := fmt.Sprintf("%s%s", router.Host, newRequestURI)
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

func (h *ServeReverseProxyPass) getRouterConfigInfo(partsPath []string) (string, Router) {
	initialPath := partsPath[2]
	router := h.RouterConfig[initialPath]
	return initialPath, router
}

func (h *ServeReverseProxyPass) LoadRouterConfig(loggerSugar *zap.SugaredLogger, gatewayDDB *database.GatewayDB, tickerReloadRouters time.Duration) {

	loadRoutersFunc := func() map[string]Router {

		routersDB := gatewayDDB.GetAllRouter()
		routers := make(map[string]Router, len(routersDB))

		for _, router := range routersDB {

			host, _ := GetMapValueFromJsonRawMessage[string](router.Configuration, "host")
			//appContext, _ := GetMapValueFromJsonRawMessage[string](router.Configuration, "app-context")
			replaceOldAppContext, _ := GetMapValueFromJsonRawMessage[string](router.Configuration, "replace-old-app-context")
			replaceNewAppContext, _ := GetMapValueFromJsonRawMessage[string](router.Configuration, "replace-new-app-context")

			routers[router.Router] = Router{
				Host:                 host.(string),
				ReplaceOldAppContext: replaceOldAppContext.(string),
				ReplaceNewAppContext: replaceNewAppContext.(string),
			}

		}

		loggerSugar.Infow("loaded routers", "routers", routers)
		return routers
	}

	h.RouterConfig = loadRoutersFunc()

	go func() {
		for range time.Tick(tickerReloadRouters) {
			h.RouterConfig = loadRoutersFunc()
		}
	}()

}
