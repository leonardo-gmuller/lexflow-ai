package dto

import "errors"

type ValidationPaginationFunc func(p Pagination) error

type Pagination struct {
	Page         int
	ItemsPerPage int
	SortBy       string
	SortType     string
}

func (p Pagination) Validate(fieldsSortBy map[string]bool) error {
	if p.Page < 1 {
		return errors.New("page must be greater than 0")
	}

	if p.ItemsPerPage < 1 || p.ItemsPerPage > 100 {
		return errors.New("itemsPerPage must be between 1 and 100")
	}

	if _, found := fieldsSortBy[p.SortBy]; !found {
		return errors.New("invalid sortBy param")
	}

	if p.SortType != "ASC" && p.SortType != "DESC" {
		return errors.New("sortType param must be 'ASC' or 'DESC'")
	}

	return nil
}

func (p Pagination) Offset() int {
	return (p.Page - 1) * p.ItemsPerPage
}
