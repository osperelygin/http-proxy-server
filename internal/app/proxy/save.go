package proxy

import (
	"context"
	"net/http"
)

const notSavedRequestID = -1

func (ps ProxyServer) saveRequest(r *http.Request) int {
	if ps.saver == nil {
		return notSavedRequestID
	}

	ps.logger.Infoln("start saving request")

	reqID, err := ps.saver.SaveRequest(context.Background(), r)
	if err != nil {
		ps.logger.Errorln("SaveRequest failed:", err.Error())
		return notSavedRequestID
	}

	ps.logger.Infoln("successful save request")

	return reqID
}

func (ps ProxyServer) saveResponse(resp *http.Response, reqID int) {
	if ps.saver == nil || reqID == notSavedRequestID {
		return
	}

	ps.logger.Infoln("start saving response")

	if err := ps.saver.SaveResponse(context.Background(), reqID, resp); err != nil {
		ps.logger.Errorln("SaveResponse failed:", err.Error())
		return
	}

	ps.logger.Infoln("successful save response")
}

// func (ps ProxyServer) save(r *http.Request, resp *http.Response) {
// 	if ps.saver == nil {
// 		return
// 	}

// 	ps.logger.Infoln("start saving request and response")

// 	reqID, err := ps.saver.SaveRequest(context.Background(), r)
// 	if err != nil {
// 		ps.logger.Errorln("SaveRequest failed:", err.Error())
// 		return
// 	}

// 	if err := ps.saver.SaveResponse(context.Background(), reqID, resp); err != nil {
// 		ps.logger.Errorln("SaveResponse failed:", err.Error())
// 	}

// 	ps.logger.Infoln("successful saved request and response")
// }
