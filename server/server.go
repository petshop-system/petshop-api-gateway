package server

import (
	"go.uber.org/zap"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type ServeReverseProxyPass struct {
	LoggerSugar  *zap.SugaredLogger
	RouterConfig map[string]map[string]string
}

func NewServerPass(loggerSugar *zap.SugaredLogger, fileName string) ServeReverseProxyPass {

	config := LoadRouterConfig(fileName)
	return ServeReverseProxyPass{
		LoggerSugar:  loggerSugar,
		RouterConfig: config,
	}
}

func (h *ServeReverseProxyPass) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	partsPath := strings.Split(r.RequestURI, "/")
	initialPath, hostRedirect, appContext := h.getRouterConfigInfo(partsPath)

	h.LoggerSugar.Infow("request server pass received",
		"host", r.Host, "request_uri", r.RequestURI, "request_url_path", r.URL.Path,
		"initialPath", initialPath)

	reverseProxy, err := h.buildReverseProxy(hostRedirect, appContext, r)
	if err != nil {
		h.LoggerSugar.Errorw("error to process url destination",
			"host_redirect", hostRedirect)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	reverseProxy.ServeHTTP(w, r)
	h.LoggerSugar.Infow("server pass successfully",
		"new_host", hostRedirect, "new_request_uri", r.RequestURI)

}

func (h *ServeReverseProxyPass) buildReverseProxy(hostRedirect, appContext string, r *http.Request) (*httputil.ReverseProxy, error) {

	destination, err := url.Parse(hostRedirect)
	if err != nil {
		return nil, err
	}

	rp := httputil.NewSingleHostReverseProxy(destination)
	rp.Director = func(req *http.Request) {
		req.Host = hostRedirect
		req.URL.Scheme = destination.Scheme
		req.URL.Host = destination.Host
		//req.URL.Path = singleJoiningSlash(u.Path, req.URL.Path)
		//if targetQuery == "" || req.URL.RawQuery == "" {
		//	req.URL.RawQuery = targetQuery + req.URL.RawQuery
		//} else {
		//	req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		//}
		req.RequestURI = strings.ReplaceAll(r.RequestURI, "petshop-system", appContext)
		req.URL.Path = strings.ReplaceAll(r.URL.Path, "petshop-system", appContext)
		req.Header.Set("api-gateway", "true")
	}

	rp.ModifyResponse = func(w *http.Response) error {
		w.Header.Set("api-gateway", "true")
		return nil
	}

	return rp, nil
}

func (h *ServeReverseProxyPass) getRouterConfigInfo(partsPath []string) (string, string, string) {

	initialPath := partsPath[2]
	hostRedirect := h.RouterConfig[initialPath]["host"]
	appContext := h.RouterConfig[initialPath]["app-context"]

	return initialPath, hostRedirect, appContext
}
