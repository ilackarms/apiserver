package filters

import (
	"github.com/emicklei/go-restful"
	"github.com/golang/glog"
	"net/http"
	"strings"
)

func WithCompression(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		//don't compress watches
		if req.URL.Query().Get("watch") == "true" {
			handler.ServeHTTP(w, req)
			return
		}
		wantsCompression, encoding := wantsCompressedResponse(req)
		if wantsCompression {
			compressionWriter, err := restful.NewCompressingResponseWriter(w, encoding)
			if err != nil {
				glog.V(0).Infof("Error: failed to install compression response writer")
				handler.ServeHTTP(w, req)
				return
			}
			handler.ServeHTTP(compressionWriter, req)
		}
	})
}

// WantsCompressedResponse reads the Accept-Encoding header to see if and which encoding is requested.
func wantsCompressedResponse(httpRequest *http.Request) (bool, string) {
	header := httpRequest.Header.Get(restful.HEADER_AcceptEncoding)
	gi := strings.Index(header, restful.ENCODING_GZIP)
	zi := strings.Index(header, restful.ENCODING_DEFLATE)
	// use in order of appearance
	if gi == -1 {
		return zi != -1, restful.ENCODING_DEFLATE
	} else if zi == -1 {
		return gi != -1, restful.ENCODING_GZIP
	} else {
		if gi < zi {
			return true, restful.ENCODING_GZIP
		}
		return true, restful.ENCODING_DEFLATE
	}
}
