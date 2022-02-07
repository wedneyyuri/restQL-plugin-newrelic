package restqlnewrelic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/b2wdigital/restQL-golang/v6/pkg/restql"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/pkg/errors"
)

const (
	newRelicPluginName = "NewRelicPlugin"

	newRelicResponseWriter  = "__newRelicResponseWriter__"
	newRelicTransaction     = "__newRelicTransaction__"
	newRelicExternalSegment = "__newRelicExternalSegment__"
	newRelicExternalRequest = "__newRelicExternalRequest__"
)

func init() {
	restql.RegisterPlugin(restql.PluginInfo{
		Name: newRelicPluginName,
		Type: restql.LifecyclePluginType,
		New: func(logger restql.Logger) (restql.Plugin, error) {
			return MakeNewRelicPlugin(logger)
		},
	})
}

func MakeNewRelicPlugin(logger restql.Logger) (restql.LifecyclePlugin, error) {
	app, err := newrelic.NewApplication(
		newrelic.ConfigFromEnvironment(),
		ExtraConfigFromEnvironment(),
	)

	if err != nil {
		logger.Error("failed to initialize new relic", err)
		return nil, err
	}

	return &NewRelicPlugin{log: logger, application: app}, nil
}

type NewRelicPlugin struct {
	log         restql.Logger
	application *newrelic.Application
}

func (n *NewRelicPlugin) Name() string {
	return newRelicPluginName
}

func (n *NewRelicPlugin) BeforeTransaction(ctx context.Context, tr restql.TransactionRequest) context.Context {
	txn := n.application.StartTransaction(txnName(tr.Method, tr.Url))
	txn.SetWebRequest(newrelic.WebRequest{
		Header: tr.Header,
		URL:    tr.Url,
		Method: tr.Method,
	})
	txn.AddAttribute("query", tr.Url.Query().Encode())
	w := txn.SetWebResponse(noOpResponseWriter{})

	txnCtx := context.WithValue(ctx, newRelicTransaction, txn)
	writerCtx := context.WithValue(txnCtx, newRelicResponseWriter, w)

	return writerCtx
}

func (n *NewRelicPlugin) AfterTransaction(ctx context.Context, tr restql.TransactionResponse) context.Context {
	txn, err := n.getTransactionFromContext(ctx)
	if err != nil {
		return ctx
	}
	defer txn.End()

	w, ok := ctx.Value(newRelicResponseWriter).(http.ResponseWriter)
	if !ok {
		n.log.Error(
			"failed to retrieve new relic writer from context",
			errors.Errorf("incorrect writer type : %T", w),
		)
		return ctx
	}

	segment := txn.StartSegment("Flush")
	defer segment.End()
	//todo: set headers
	w.WriteHeader(tr.Status)
	_, err = w.Write(tr.Body)
	if err != nil {
		n.log.Error("failed to write response for new relic", err)
	}

	return ctx
}

func (n *NewRelicPlugin) BeforeQuery(ctx context.Context, query string, queryCtx restql.QueryContext) context.Context {
	return ctx
}
func (n *NewRelicPlugin) AfterQuery(ctx context.Context, query string, result map[string]interface{}) context.Context {
	return ctx
}

func (n *NewRelicPlugin) BeforeRequest(ctx context.Context, request restql.HTTPRequest) context.Context {
	txn, err := n.getTransactionFromContext(ctx)
	if err != nil {
		return ctx
	}

	externalRequest, err := makeExternalRequest(request)
	if err != nil {
		n.log.Error("failed to make request to report for new relic", err)
		return ctx
	}

	reqTxn := txn.NewGoroutine()
	segment := newrelic.StartExternalSegment(reqTxn, externalRequest)

	segmentCtx := context.WithValue(ctx, newRelicExternalSegment, segment)
	extReqCtx := context.WithValue(segmentCtx, newRelicExternalRequest, externalRequest)

	return extReqCtx
}

func (n *NewRelicPlugin) AfterRequest(ctx context.Context, request restql.HTTPRequest, response restql.HTTPResponse, errordetail error) context.Context {
	seg := ctx.Value(newRelicExternalSegment)
	if seg == nil {
		return ctx
	}

	segment, ok := seg.(*newrelic.ExternalSegment)
	if !ok {
		n.log.Error(
			"failed to retrieve new relic segment from context",
			errors.Errorf("incorrect new relic segment type : %T", seg),
		)
		segment.End()
		return ctx
	}

	extReq, ok := ctx.Value(newRelicExternalRequest).(*http.Request)
	if !ok {
		n.log.Error(
			"failed to retrieve new relic external request from context",
			errors.Errorf("incorrect new relic external request type : %T", seg),
		)
		segment.End()
		return ctx
	}

	externalResponse, err := makeExternalResponse(extReq, response)
	if err != nil {
		n.log.Error("failed to make response to report for new relic", err)
		segment.End()
		return ctx
	}

	segment.AddAttribute("errordetail", errordetail)
	segment.Response = externalResponse
	segment.End()

	// here you could send custom metrics to New Relic Insights using n.sendCustomEvent

	return ctx
}

func (n *NewRelicPlugin) sendCustomEvent(name string, request restql.HTTPRequest, response restql.HTTPResponse, errs ...error) {
	n.log.Debug("preparing custom event")

	event := map[string]interface{}{
		"url":    request.Host,
		"status": response.StatusCode,
		"time":   response.Duration.Milliseconds(),
		"method": request.Method,
	}

	if len(errs) > 0 {
		var errsStr string

		b, err := json.Marshal(errs)
		if err != nil {
			errsStr = fmt.Sprintf("%v", errs)
		} else {
			errsStr = string(b)
		}

		event["errors"] = errsStr
	}

	n.application.RecordCustomEvent(name, event)

	n.log.Debug("custom event send")
}

func (n *NewRelicPlugin) getTransactionFromContext(ctx context.Context) (*newrelic.Transaction, error) {
	txn, ok := ctx.Value(newRelicTransaction).(*newrelic.Transaction)
	if !ok {
		err := errors.Errorf("incorrect transaction type : %T", txn)
		n.log.Error(
			"failed to retrieve new relic transaction from context",
			err,
		)
		return nil, err
	}
	return txn, nil
}

func txnName(method string, reqUrl *url.URL) string {
	return method + " " + reqUrl.Path
}

type noOpResponseWriter struct {
	headers http.Header
	body    []byte
	status  int
}

func (n noOpResponseWriter) Header() http.Header {
	return n.headers
}

func (n noOpResponseWriter) Write(body []byte) (int, error) {
	n.body = body
	return len(body), nil
}

func (n noOpResponseWriter) WriteHeader(statusCode int) {
	n.status = statusCode
}
