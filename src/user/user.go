package user

// User struct
type User struct {
	Username string `db:"username"`
	Nickname string `db:"nickname"`
	Profile  string `db:"profile"`
}

// GetByNameReq : get req
type GetByNameReq struct {
	Username string
}

// GetByNameRsp get rsp
type GetByNameRsp struct {
	UserItem User
}

// UpdateByNameReq : get req
type UpdateByNameReq struct {
	UserItem User
}

// UpdateByNameRsp get rsp
type UpdateByNameRsp struct {
}

// IUser user op interface, for client & svr
type IUser interface {
	GetByName(req *GetByNameReq, rsp *GetByNameRsp) int
	UpdateByName(req *UpdateByNameReq, rsp *UpdateByNameRsp) int
}
