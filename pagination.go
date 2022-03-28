package pagination

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"html/template"
	"math"
	"regexp"
	"strconv"

)

type Paginate struct {
	Count int64 `form:"count" json:"count" `
	Page int `form:"page" json:"page"`
	Path string `form:"pathname" json:"pathname"`
	CurrentPageUrt string `json:"current_page_urt"`
	FirstPageUrl string `json:"first_page_url"`
	LastPageUrl string `json:"last_page_url"`
	PrevPageUrl string `json:"prev_page_url"`
	NextPageUrl string	`json:"next_page_url"`
	Size int `form:"size" json:"size"`
	Slot int `json:"slot"`
	PageCount int `json:"page_count"`
	Data interface{} `json:"data"`
}

//default Init
func Init(page int, size int, slot int, count int64, path string) *Paginate {
	var p Paginate
	p.Page = page
	p.Size = size
	p.Slot = slot
	p.Count = count
	p.Path = fmt.Sprintf("%s?",path)
	p.checkReq(size)
	return &p
}

//Init with gin
func GinInit(ctx *gin.Context, size int, slot int, count int64) *Paginate {
	var p Paginate
	p.Slot = slot
	ctx.ShouldBindQuery(&p)
	if p.Count == 0{
		p.Count = count
	}
	p.setPath(ctx)
	p.checkReq(size)

	return &p

}

//Init with gin and gorm
//scope just use with (tx.order(); tx.select(); tx.Omit() ...)
func GinOrmInit(ctx *gin.Context, tx *gorm.DB,scope func(tx *gorm.DB) *gorm.DB,data interface{}, size int, slot int) *Paginate {
	var p Paginate
	var count int64
	p.Slot = slot
	ctx.ShouldBindQuery(&p)
	if p.Count == 0{
		tx.Count(&count)
		p.Count = count
	}
	p.setPath(ctx)
	p.checkReq(size)
	offset := (p.Page-1) * p.Size
	tx.Scopes(scope).Offset(offset).Limit(p.Size).Find(data)
	p.Data = data
	return &p
}

func (p *Paginate) setPath(ctx *gin.Context)  {
	uri := ctx.Request.RequestURI
	path := ctx.Request.URL.Path
	if regexp.MustCompile(`\?[\w-]+=`).MatchString(uri){
		if regexp.MustCompile(`\?page=\d*`).MatchString(uri){
			p.Path = fmt.Sprintf(`%s?`, path)
		}else {
			s := regexp.MustCompile(`&(?:page|size|count)=\d*`).ReplaceAllString(uri, "")
			p.Path = fmt.Sprintf(`%s&`, s)
		}
	}else {
		p.Path = fmt.Sprintf(`%s?`, path)
	}
}

func (p *Paginate) checkReq(size int)  {
	if p.Page <= 0{
		p.Page = 1
	}
	switch  {
	case p.Size > 100:
		p.Size = 100
	case p.Size <= 0:
		p.Size = size
	}
}

func (p *Paginate) GetList() []string {
	var lists []string
	p.PageCount = int(math.Ceil(float64(p.Count)/float64(p.Size)))
	p.FirstPageUrl = fmt.Sprintf("%spage=1&size=%d&count=%d", p.Path, p.Size, p.Count)
	p.LastPageUrl = fmt.Sprintf("%spage=%d&size=%d&count=%d", p.Path, p.PageCount, p.Size, p.Count)
	if p.Page != 1{
		p.PrevPageUrl = fmt.Sprintf("%spage=%d&size=%d&count=%d", p.Path, p.Page-1, p.Size, p.Count)
	}
	if p.Page != p.PageCount{
		p.NextPageUrl = fmt.Sprintf("%spage=%d&size=%d&count=%d", p.Path, p.Page+1, p.Size, p.Count)
	}
	if p.PageCount <= p.Slot + 2{
		for i := 1; i <= p.PageCount; i++ {
			lists = append(lists, strconv.Itoa(i))
		}
	}else {
		switch {
		case p.Page < p.Slot:
			for i := 1; i <= p.Slot ; i++ {
				lists = append(lists, strconv.Itoa(i))
			}
			lists = append(lists,"...", strconv.Itoa(p.PageCount))
		case p.Page > p.PageCount - p.Slot + 1:
			lists = append(lists, "1", "...")
			for i := p.PageCount - p.Slot + 1; i <= p.PageCount ; i++ {
				lists = append(lists, strconv.Itoa(i))
			}
		default:
			lists = append(lists, "1", "...")
			for i := p.Page - (p.Slot - 1)/2; i <= p.Page + (p.Slot - 1)/2 ; i++ {
				lists = append(lists, strconv.Itoa(i))
			}
			lists = append(lists,"...", strconv.Itoa(p.PageCount))
		}
	}
	return lists
}
//pagination with bootstrap5
func (p *Paginate) BsPage() template.HTML {
	var html, navH, navF, prev, next, items string
	lists := p.GetList()
	navH = `<nav aria-label="pagination"> <ul class="pagination my-3">`
	navF = `</ul> </nav>`
	if p.PrevPageUrl == ""{
		prev = `<li class="page-item disabled"> <a class="page-link disabled" href="#" aria-label="Previous"> <span aria-hidden="true">&laquo;</span> </a> </li>`
	}else {
		prev = fmt.Sprintf(`<li class="page-item"> <a class="page-link" href="%s" aria-label="Previous"> <span aria-hidden="true">&laquo;</span> </a> </li>`, p.PrevPageUrl)
	}
	if p.NextPageUrl == ""{
		next = `<li class="page-item disabled"> <a class="page-link" href="#" aria-label="Next"> <span aria-hidden="true">&raquo;</span> </a> </li>`
	}else {
		next = fmt.Sprintf(`<li class="page-item"> <a class="page-link" href="%s" aria-label="Next"> <span aria-hidden="true">&raquo;</span> </a> </li>`, p.NextPageUrl)
	}
	for _, list := range lists {
		var item, linkUrl string
		linkUrl = fmt.Sprintf("%spage=%s&size=%d&count=%d", p.Path, list, p.Size, p.Count)

		switch{
		case list == "...":
			item = fmt.Sprintf(`<li class="page-item disabled"><a class="page-link" href="#">%s</a></li>`,list)
		case list == strconv.Itoa(p.Page):
			item = fmt.Sprintf(`<li class="page-item active" aria-current="page"><a class="page-link" href="%s">%s</a></li>`,linkUrl, list)
		default:
			item = fmt.Sprintf(`<li class="page-item"><a class="page-link" href="%s">%s</a></li>`,linkUrl, list)
		}
		items += item
	}
	if p.Count != 0{
		html = navH + prev + items + next +navF
	}
	return template.HTML(html)
}

//simple pagination with bootstrap5
func (p *Paginate) SimpleBsPage() template.HTML {
	var html, navH, navF, prev, next, current string
	navH = `<nav aria-label="pagination"> <ul class="pagination my-3">`
	navF = `</ul> </nav>`
	if p.PrevPageUrl == ""{
		prev = `<li class="page-item disabled"> <a class="page-link disabled" href="#" aria-label="Previous"> <span aria-hidden="true">&laquo;</span> </a> </li>`
	}else {
		prev = fmt.Sprintf(`<li class="page-item"> <a class="page-link" href="%s" aria-label="Previous"> <span aria-hidden="true">&laquo;</span> </a> </li>`, p.PrevPageUrl)
	}
	if p.NextPageUrl == ""{
		next = `<li class="page-item disabled"> <a class="page-link" href="#" aria-label="Next"> <span aria-hidden="true">&raquo;</span> </a> </li>`
	}else {
		next = fmt.Sprintf(`<li class="page-item"> <a class="page-link" href="%s" aria-label="Next"> <span aria-hidden="true">&raquo;</span> </a> </li>`, p.NextPageUrl)
	}
	linkUrl := fmt.Sprintf("%spage=%d&size=%d&count=%d", p.Path, p.Page, p.Size, p.Count)
	current = fmt.Sprintf(`<li class="page-item active" aria-current="page"><a class="page-link" href="%s">%d</a></li>`,linkUrl, p.Page)
	if p.Count != 0{
		html = navH + prev + current + next +navF
	}
	return template.HTML(html)
}

