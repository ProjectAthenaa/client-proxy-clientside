package main

import (
	"github.com/ProjectAthenaa/sonic-core/fasttls"
	client_proxy "github.com/ProjectAthenaa/sonic-core/protos/clientProxy"
	"time"
)

func convertToRequest(reqIn *client_proxy.Request) *fasttls.Request {
	//pr := "localhost:8080"

	return &fasttls.Request{
		URL:             reqIn.URL,
		Method:          fasttls.Method(reqIn.Method),
		Headers:         headerConvertTo(reqIn.Headers),
		Data:            reqIn.Data,
		FollowRedirects: reqIn.FollowRedirects,
		Timeout:         (*time.Duration)(reqIn.Timeout),
		UseHttp2:        true,
		ServerName:      reqIn.ServerName,
		UseMobile:       reqIn.UseMobile,
	}
}

func convertToResponse(resIn *fasttls.Response, taskID string) *client_proxy.Response {
	return &client_proxy.Response{
		StatusCode:    int32(resIn.StatusCode),
		Body:          resIn.Body,
		Headers:       headerConvertFrom(resIn.Headers),
		TimeTaken:     int64(resIn.TimeTaken),
		IsHttp2:       resIn.IsHttp2,
		ContentLength: resIn.ContentLength,
		TaskID:        taskID,
	}
}

func headerConvertTo(headersIn map[string]string) fasttls.Headers {
	headersOut := make(fasttls.Headers)

	for k, v := range headersIn {
		headersOut[k] = []string{v}
	}

	return headersOut
}

func headerConvertFrom(headersIn fasttls.Headers) map[string]string {
	headersOut := make(map[string]string)

	for k, v := range headersIn {
		headersOut[k] = v[0]
	}
	return headersOut
}
