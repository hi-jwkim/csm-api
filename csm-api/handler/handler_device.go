package handler

import (
	"csm-api/entity"
	"csm-api/service"
	"csm-api/utils"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

/**
 * @author 작성자: 김진우
 * @created 작성일: 2025-02-12
 * @modified 최종 수정일: 2025-02-21
 * @modifiedBy 최종 수정자: 정지영
 * @modified description
 * - 검색 및 정렬 조건 추가, url의 query parameter로 받음
 */

// struct: 근태인식기 조회
type DeviceHandler struct {
	Service service.DeviceService
}

// func: 근태인식기 조회
// @param
// - response: http get paramter
func (d *DeviceHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// http get paramter를 저장할 구조체 생성
	page := entity.Page{}
	search := entity.Device{}

	pageNum := r.URL.Query().Get(entity.PageNumKey)
	rowSize := r.URL.Query().Get(entity.RowSizeKey)
	order := r.URL.Query().Get(entity.OrderKey)

	if pageNum == "" || rowSize == "" {
		BadRequestResponse(ctx, w)
		return
	}
	page.PageNum, _ = strconv.Atoi(pageNum)
	page.RowSize, _ = strconv.Atoi(rowSize)
	page.Order = order

	search.DeviceNm = utils.ParseNullString(r.URL.Query().Get("device_nm"))
	search.DeviceSn = utils.ParseNullString(r.URL.Query().Get("device_sn"))
	search.SiteNm = utils.ParseNullString(r.URL.Query().Get("site_nm"))
	search.Etc = utils.ParseNullString(r.URL.Query().Get("etc"))
	search.IsUse = utils.ParseNullString(r.URL.Query().Get("is_use"))
	retrySearchText := r.URL.Query().Get("retry_search_text")

	// 근태인식기 목록
	list, err := d.Service.GetDeviceList(ctx, page, search, retrySearchText)
	if err != nil {
		FailResponse(ctx, w, err)
		return
	}

	// 근태인식기 전체 개수
	count, err := d.Service.GetDeviceListCount(ctx, search, retrySearchText)
	if err != nil {
		FailResponse(ctx, w, err)
		return
	}

	values := struct {
		List  entity.Devices `json:"list"`
		Count int            `json:"count"`
	}{List: *list, Count: count}

	SuccessValuesResponse(ctx, w, values)
}

// func: 근태인식기 추가
// @param
// - response: entity.Device - json(raw)
func (d *DeviceHandler) Add(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	//response 데이터 파싱
	device := entity.Device{}
	if err := json.NewDecoder(r.Body).Decode(&device); err != nil {
		FailResponse(ctx, w, err)
		return
	}

	// 근태인식기 추가
	err := d.Service.AddDevice(ctx, device)
	if err != nil {
		FailResponse(ctx, w, err)
		return
	}

	SuccessResponse(ctx, w)
}

// func: 근태인식기 수정
// @param
// - response: entity.Device - json(raw)
func (d *DeviceHandler) Modify(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	//response 데이터 파싱
	device := entity.Device{}
	if err := json.NewDecoder(r.Body).Decode(&device); err != nil {
		FailResponse(ctx, w, err)
		return
	}

	// 근태인식기 수정
	err := d.Service.ModifyDevice(ctx, device)
	if err != nil {
		FailResponse(ctx, w, err)
		return
	}

	SuccessResponse(ctx, w)
}

// func: 근태인식기 삭제
// @param
// - response: url parameter
func (d *DeviceHandler) Remove(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// URL에서 {id} 값을 가져오기
	id := chi.URLParam(r, "id")
	int64ID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		FailResponse(ctx, w, err)
		return
	}

	// 서비스 호출하여 삭제 처리
	err = d.Service.RemoveDevice(ctx, int64ID)
	if err != nil {
		FailResponse(ctx, w, err)
		return
	}

	// 성공 응답
	SuccessResponse(ctx, w)
}
