package types

type OrganizationGetSearchReq struct {
	Name string `json:"name" form:"name" uri:"name" `
}

type OrganizationPutReq struct {
	Name string `json:"name" form:"name" uri:"name" binding:"required,max=32"`
}
