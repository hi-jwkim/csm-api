package handler

import (
	"csm-api/entity"
	"csm-api/service"
	"encoding/json"
	"net/http"
	"strconv"
)

/**
 * @author 작성자: 정지영
 * @created 작성일: 2025-02-17
 * @modified 최종 수정일:
 * @modifiedBy 최종 수정자:
 * @modified description
 * -
 */

// struct : 공지사항 조회
type NoticeListHandler struct {
	Service service.NoticeService
}

// func : 공지사항 전체조회
// @param
// - response: hhtp get parameter
func (n *NoticeListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	page := entity.Page{}

	pageNum := r.URL.Query().Get(entity.PageNumKey)
	rowSize := r.URL.Query().Get(entity.RowSizeKey)

	if pageNum == "" || rowSize == "" {
		RespondJSON(
			ctx,
			w,
			&ErrResponse{
				Result:         Failure,
				Message:        "get parameter is missing",
				Details:        NotFoundParam,
				HttpStatusCode: http.StatusBadRequest,
			},
			http.StatusOK)
		return
	}

	page.PageNum, _ = strconv.Atoi(pageNum)
	page.RowSize, _ = strconv.Atoi(rowSize)

	notices, err := n.Service.GetNoticeList(ctx, page)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{Message: err.Error()}, http.StatusInternalServerError)
		return
	}

	count, err := n.Service.GetNoticeListCount(ctx)
	if err != nil {
		RespondJSON(
			ctx,
			w,
			&ErrResponse{
				Result:         Failure,
				Message:        err.Error(),
				HttpStatusCode: http.StatusInternalServerError,
			},
			http.StatusOK)
		return
	}

	rsp := Response{
		Result: Success,
		Values: struct {
			Notices entity.Notices `json:"notices"`
			Count   int            `json:"count"`
		}{Notices: *notices, Count: count},
	}

	RespondJSON(ctx, w, &rsp, http.StatusOK)

}

// struct: 공지사항 추가
type NoticeAddHandler struct {
	Service service.NoticeService
}

// func: 공지사항 추가
// @param
// - response: entity.Notice - json(raw)
func (n *NoticeAddHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	notice := entity.Notice{}
	if err := json.NewDecoder(r.Body).Decode(&notice); err != nil {
		RespondJSON(
			ctx,
			w,
			&ErrResponse{
				Result:         Failure,
				Message:        err.Error(),
				Details:        BodyDataParseError,
				HttpStatusCode: http.StatusInternalServerError,
			},
			http.StatusOK)
		return
	}

	err := n.Service.AddNotice(ctx, notice)
	if err != nil {
		RespondJSON(
			ctx,
			w,
			&ErrResponse{
				Result:         Failure,
				Message:        err.Error(),
				Details:        DataAddFailed,
				HttpStatusCode: http.StatusInternalServerError,
			},
			http.StatusOK)
		return
	}
}
