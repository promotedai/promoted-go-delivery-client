package delivery

import "github.com/promotedai/schema/generated/go/proto/delivery"

// NewPaging creats a paging instance with start and offset since proto is overly complicated.
func NewPaging(size, offset int32) *delivery.Paging {
	return &delivery.Paging{
		Starting: &delivery.Paging_Offset{
			Offset: offset,
		},
		Size: size,
	}
}
