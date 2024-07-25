package types

type PositionListSearchReq struct {
	OrganizationId uint `json:"organizationId"  binding:"required,number" form:"organizationId" uri:"organizationId" `
}
