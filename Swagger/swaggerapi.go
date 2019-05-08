package swagger

import model "FirstProject/Model"

// Users requested
// swagger:model getAllUsersResp
type GetAllUsersResp struct {
	// in:body
	Users []model.User `json:"users"`
}

// User requested
// swagger:model getUserResp
type GetUserResp struct {
	// in:body
	User model.User `json:"user"`
}

// User's JWT updated
// swagger:model updateUserResp
type UpdateUserResp struct {
	// in:body
	Jwt string `json:"jwt"`
}

// User's JWT
// swagger:model loginResp
type LoginResp struct {
	// in:body
	Jwt string `json:"jwt"`
}

// HTTP status code 200 message response
// swagger:model genericSuccessResp
type GenericSuccessResp struct {
	// in:body
	Response string `json:"response"`
}

// HTTP status code 400 response message
// swagger:model queryErrResp
type GenericQueryErrResp struct {
	// in:body
	Error error `json:"error"`
}

// HTTP status code 404 response
// swagger:model notFound
type GenericNotFound struct {
	// in:body
	Error error `json:"error"`
}

// HTTP status code 500 response
// swagger:model internalErr
type GenericInternalErr struct {
	// in:body
	Error error `json:"error"`
}
