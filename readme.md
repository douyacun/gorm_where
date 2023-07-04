# gorm-where

根绝struct tag自动拼接sql语句

## tag用法

```go
type AccountQuery struct {
	Keyword string       `json:"keyword" session:"name,op:like"`
	Id      string       `json:"id" session:"id,op:equal"`
	IdList  []string     `json:"id_list" session:"id,op:in"`
	LtId   int64        `json:"lt_uid" session:"id,op:lt"`
	GtId   int64        `json:"gt_uid" session:"id,op:gt"`

	Page int `es:"page" json:"page" session:"page"`
	PageSize int    `json:"page_size" session:"page_size"`
	Sort     string `json:"sort" session:"sort_by"`
	NoCount  bool   `json:"no_count" session:"no_count"`
}
```