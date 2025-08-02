package page

type PageUtils[T any] struct {
	PageNum  int32
	PageSize int32
	Src      []T
}

func NewPageUtils[T any](src []T, pageNum, pageSize int32) *PageUtils[T] {
	return &PageUtils[T]{
		Src:      src,
		PageNum:  pageNum,
		PageSize: pageSize,
	}
}

func (p *PageUtils[T]) GetRecords() []T {
	total := p.GetTotal()
	var startIndex = (p.PageNum - 1) * p.PageSize
	var endIndex = startIndex + p.PageSize
	if startIndex > total {
		return nil
	}
	if endIndex > total {
		endIndex = total
	}
	return p.Src[startIndex:endIndex]
}

func (p *PageUtils[T]) GetTotal() int32 {
	return int32(len(p.Src))
}

func (p *PageUtils[T]) SetPageNum(pageNum int32) {
	p.PageNum = pageNum
}

func (p *PageUtils[T]) SetPageSize(pageSize int32) {
	p.PageSize = pageSize
}
