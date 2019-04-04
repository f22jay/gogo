package auth

// UserAuth :
type UserAuth struct {
	Username string `db:"username"`
	Passwd   string `db:"passwd"`
}

// LoginByPassReq :
type LoginByPassReq struct {
	UserAuth
}

// LoginByPassRsp :
type LoginByPassRsp struct {
	Token string
}

// LoginByTokenReq check token req
type LoginByTokenReq struct {
	Token string
}

// LoginByTokenRsp check token rsp
type LoginByTokenRsp struct {
	Username string
}

// IAuth interface
type IAuth interface {
	LoginByPass(req *LoginByPassReq, rsp *LoginByPassRsp) int    // check login username passwd req
	LoginByToken(req *LoginByTokenReq, rsp *LoginByTokenRsp) int // check login by token
}
