package model

import "gomodel/internal/shared/model"

type ListResponse[R any] struct {
	Data []R `json:"data"`
}

func GetListResponse[R any](r []R) ListResponse[R] {
	if r == nil {
		r = make([]R, 0)
	}
	return ListResponse[R]{
		Data: r,
	}
}

type PaginatedListResponse[R any] struct {
	ListResponse[R]
	model.Pagination
}
