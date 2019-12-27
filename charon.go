package charon

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/charon/errors"
	logr "github.com/charon/logger"
)

// type ServerHandler interface {
// 	ServeHTTP(http.ResponseWriter, *http.Request)
// }

type charonServerHandler struct {
	logger       *logr.Logger
	respHandler  ResponseHandler
	pathHandlers map[PathDetail]RouteHandler
}

//ServeRequest method serves all incoming requests
func (serverHandler *charonServerHandler) ServeRequest(resp http.ResponseWriter, req *http.Request) {

	isServed := false
	found := false

	path := req.URL.Path
	method := req.Method
	header := req.Header

	rDetails := RouteDetails{method: method, path: path, headers: header, log: strings.Builder{}}

	serverHandler.logger.LogInfo(fmt.Sprint("Incoming Request  ", method, " : ", path), nil, &rDetails)

	defer func() {
		if r := recover(); r != nil {
			//serverHandler.logger.LogSevere(string(debug.Stack()), nil, &rDetails)
			serverHandler.logger.LogPanic(r.(string), nil, &rDetails)
			err := errors.InternalError{Err: "Unknown server error"}
			handleLog(rDetails, serverHandler.logger)
			handleResponse(resp, nil, err, serverHandler.respHandler)
		}
	}()

	body := make(map[string]interface{})

	// fmt.Println("Incoming Request  ", method, ":", path, "  ", time.Now())
	if method == "GET" {
		// body = req.URL.Query()
		for k, v := range req.URL.Query() {
			body[k] = v
		}
	} else {
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&body)
		if err != nil {
			if err.Error() == "EOF" {
				body = nil
			} else {
				// fmt.Println("Error:  ", err.Error(), "  ", time.Now())
				serverHandler.logger.LogSevere(fmt.Sprint("Error:  ", err.Error()), nil, &rDetails)
				handleLog(rDetails, serverHandler.logger)
				handleResponse(resp, nil, errors.InternalError{Err: err.Error()}, serverHandler.respHandler)
				isServed = true
			}
		}
	}

	rDetails.body = body
	// var rDetails = RouteDetails{method, path, header, body, req.Context(), strings.Builder{}}

	if !isServed {
		for pDetail, handler := range serverHandler.pathHandlers {

			if pDetail.PathRegex == rDetails.Path() && rDetails.Method() == pDetail.Method {
				found = true
				handledResp, err := HandleRequest(handler, &rDetails, req)
				if err != nil {
					serverHandler.logger.LogSevere(fmt.Sprint("Error:  ", err.Error()), nil, &rDetails)
				}
				handleLog(rDetails, serverHandler.logger)
				handleResponse(resp, handledResp, err, serverHandler.respHandler)
				break
			}
		}
		if !found {
			// fmt.Println("Error:  Path not found", "  ", time.Now())
			serverHandler.logger.LogSevere("Path not found", nil, &rDetails)
			err := errors.AuthenticationError{Err: "Path not found", Mess: "Path not found"}
			handleLog(rDetails, serverHandler.logger)
			handleResponse(resp, nil, err, serverHandler.respHandler)
		}
	}
}

// RouteDetails - route details for the incoming request
type RouteDetails struct {
	method  string
	path    string
	headers http.Header
	body    map[string]interface{}
	ctx     context.Context
	log     strings.Builder
}

//Method returns the http method for the incoming http request
func (detail RouteDetails) Method() string {
	return detail.method
}

//Path returns the path(uri) for the incoming http request
func (detail RouteDetails) Path() string {
	return detail.path
}

//Headers returns the http headers for the incoming http request
func (detail RouteDetails) Headers() http.Header {
	return detail.headers
}

//Body returns the body(if any) of the incoming http request
func (detail RouteDetails) Body() map[string]interface{} {
	return detail.body
}

//Context returns the context of the incoming http request
func (detail RouteDetails) Context() context.Context {
	return detail.ctx
}

//WriteLog implementation of logger.Writer, to be used along with charon logger
func (detail *RouteDetails) WriteLog(logStr string) {
	detail.log.WriteString(logStr)
}

//GetLog implementation
func (detail RouteDetails) GetLog() string {
	return detail.log.String()
}

//NewRouteDetail created and returns a new RouteDetail Obj
func NewRouteDetail(ctx context.Context, method, path string, header http.Header, body map[string]interface{}, logBldr strings.Builder) *RouteDetails {
	return &RouteDetails{
		method, path, header, body, ctx, logBldr,
	}
}

// RouteHandler type to be registered with each regex url with cerberus
type RouteHandler interface {
	IsAuthenticated(context.Context, http.Header) (context.Context, *url.Userinfo, errors.Error)
	IsValidInput(RouteDetails) errors.Error
	//IsValidInput(RouteDetails) Error
	HandleCall(*RouteDetails) ([]byte, errors.Error)
}

// PathDetail type to be registered with each regex url with cerberus
type PathDetail struct {
	Method    string
	PathRegex string
}

//ResponseHandler function does response handling in the format specified by the user
type ResponseHandler func(http.ResponseWriter, []byte, errors.Error)

//AuthenticateAndSetContext does authentication and setting of context adn userinfo
// takes in header, and request context, returns a new context, UserInfo and Cerberus error if any
// type AuthenticateAndSetContext func(context.Context, http.Header, string, string) (context.Context, *url.Userinfo, errors.Error)

// RegisterValidatedRoutes function registers all the given handlers against the given path and method combo
// TODO :- paths can also be regexes
func RegisterValidatedRoutes(handlers map[PathDetail]RouteHandler, respHandler ResponseHandler,
	logger *logr.Logger) {

	serverHandler := &charonServerHandler{
		logger:       logger,
		respHandler:  respHandler,
		pathHandlers: handlers,
	}
	http.HandleFunc("/", serverHandler.ServeRequest)
}

//HandleRequest handle incoming requests
func HandleRequest(handler RouteHandler, rDetails *RouteDetails, req *http.Request) ([]byte, errors.Error) {
	ctx, userInfo, auErr := handler.IsAuthenticated(req.Context(), req.Header)
	if auErr != nil {
		return nil, auErr
	} else {
		if ctx != nil {
			req = req.WithContext(ctx)
		}
		if userInfo != nil {
			req.URL.User = userInfo
		}
	}
	rDetails.ctx = req.Context()

	authError := handler.IsValidInput(*rDetails)
	if authError != nil {
		return nil, authError
	}

	resp, err := handler.HandleCall(rDetails)
	return resp, err

}

//method handles sending response
func handleResponse(resp http.ResponseWriter, writableResp []byte, err errors.Error, respHandler ResponseHandler) {
	if respHandler != nil {
		respHandler(resp, writableResp, err)
	} else {
		resp.Header().Set("Content-Type", "application/json")
		var jsonStream []byte
		if err != nil {
			jsonStream = errors.GetMessageBytes(err)
			resp.WriteHeader(err.StatusCode())
		} else {
			jsonStream = writableResp
			resp.WriteHeader(http.StatusOK)
		}

		resp.Write(jsonStream)
	}
}

func handleLog(reader logr.LogReader, logger *logr.Logger) {
	logger.LogWithReader(reader)
}
