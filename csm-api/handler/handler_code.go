package handler

import (
	"csm-api/entity"
	"csm-api/service"
	"net/http"
)

/**
 * @author 작성자: 김진우
 * @created 작성일: 2025-03-18
 * @modified 최종 수정일:
 * @modifiedBy 최종 수정자:
 * @modified description
 * -
 */

// struct, func: 코드 조회
type HandlerCode struct {
	Service service.ServiceCode
}

func (h *HandlerCode) ListByPCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	pCode := r.URL.Query().Get("p_code")
	if pCode == "" {
		BadRequestResponse(ctx, w)
		return
	}

	list, err := h.Service.GetCodeList(ctx, pCode)
	if err != nil {
		FailResponse(ctx, w, err)
		return
	}

	values := struct {
		List entity.Codes `json:"list"`
	}{List: *list}

	SuccessValuesResponse(ctx, w, values)
}

// struct, func: 코드 트리 조회
type HandlerCodeTree struct {
	Service service.ServiceCode
}

func (h *HandlerCodeTree) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	codeTrees, err := h.Service.GetCodeTree(ctx)
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
	}

	rsp := Response{
		Result: Success,
		Values: struct {
			CodeTrees entity.CodeTrees `json:"code_trees"`
		}{CodeTrees: *codeTrees},
	}

	RespondJSON(ctx, w, &rsp, http.StatusOK)

}
