package handler

import (
	"csm-api/auth"
	"csm-api/entity"
	"csm-api/service"
	"encoding/json"
	"net/http"
	"time"
)

/**
 * @author 작성자: 김진우
 * @created 작성일: 2025-02-12
 * @modified 최종 수정일:
 * @modifiedBy 최종 수정자:
 * @modified description
 * -
 */

type HandlerSite struct {
	Service     service.SiteService
	CodeService service.CodeService
	Jwt         *auth.JWTUtils
}

// func: 현장 관리 리스트
// @param
// - response: targetDate(현재날짜), pCode(부모코드) - url parameter
func (s *HandlerSite) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// GET 요청에서 파라미터 값 읽기
	targetDateString := r.URL.Query().Get("targetDate")
	if targetDateString == "" {
		BadRequestResponse(ctx, w)
		return
	}
	targetDate, err := time.Parse("2006-01-02", targetDateString)
	if err != nil {
		FailResponse(ctx, w, err)
		return
	}

	pCode := r.URL.Query().Get("pCode")
	if pCode == "" {
		BadRequestResponse(ctx, w)
		return
	}

	// 현장 관리 리스트 조회
	sites, err := s.Service.GetSiteList(ctx, targetDate)
	if err != nil {
		FailResponse(ctx, w, err)
		return
	}

	// 현장 관리 코드 조회
	codes, err := s.CodeService.GetCodeList(ctx, pCode)
	if err != nil {
		FailResponse(ctx, w, err)
		return
	}

	values := entity.SiteRes{
		Site: *sites,
		Code: *codes,
	}
	SuccessValuesResponse(ctx, w, values)
}

// func: 현장 관리 수정
// @param
// -
func (s *HandlerSite) Modify(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	site := entity.Site{}
	if err := json.NewDecoder(r.Body).Decode(&site); err != nil {
		BadRequestResponse(ctx, w)
		return
	}

	if err := s.Service.ModifySite(ctx, site); err != nil {
		FailResponse(ctx, w, err)
		return
	}

	SuccessResponse(ctx, w)
}

// func: 현장명 조회
// @param
// -
func (s *HandlerSite) SiteNameList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	list, err := s.Service.GetSiteNmList(ctx)
	if err != nil {
		FailResponse(ctx, w, err)
		return
	}

	values := struct {
		List entity.Sites `json:"list"`
	}{List: *list}
	SuccessValuesResponse(ctx, w, values)
}

// func: 현장 상태 조회
// @param
// -
func (s *HandlerSite) StatsList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// GET 요청에서 파라미터 값 읽기
	targetDateString := r.URL.Query().Get("targetDate")
	if targetDateString == "" {
		BadRequestResponse(ctx, w)
		return
	}
	targetDate, err := time.Parse("2006-01-02", targetDateString)
	if err != nil {
		FailResponse(ctx, w, err)
		return
	}

	list, err := s.Service.GetSiteStatsList(ctx, targetDate)
	if err != nil {
		FailResponse(ctx, w, err)
		return
	}

	values := struct {
		List entity.Sites `json:"list"`
	}{List: *list}
	SuccessValuesResponse(ctx, w, values)
}

// func: 현장 생성(추가)
// @param
// -
func (h *HandlerSite) Add(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	request := struct {
		Jno      int64  `json:"jno"`
		Uno      int64  `json:"uno"`
		UserName string `json:"user_name"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		FailResponse(ctx, w, err)
		return
	}
	user := entity.User{}
	user = user.SetUser(request.Uno, request.UserName)

	err := h.Service.AddSite(ctx, request.Jno, user)
	if err != nil {
		FailResponse(ctx, w, err)
		return
	}

	SuccessResponse(ctx, w)
}

// func: 현장 사용안함 변경
// @param
// -
func (h *HandlerSite) ModifyNonUse(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	temp := struct {
		Sno int64 `json:"sno"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&temp); err != nil {
	}
	if temp.Sno == 0 {
		BadRequestResponse(ctx, w)
	}

	if err := h.Service.ModifySiteIsNonUse(ctx, temp.Sno); err != nil {
		FailResponse(ctx, w, err)
	}

	SuccessResponse(ctx, w)
}
