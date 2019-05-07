package swagger

import model "FirstProject/Model"

// type swaggerAPI struct {
// 	GetAllUsersResp []model.User
// }

// // Request containing array
// // swagger:parameters getAllUsersReq
// type swaggerGetAllUsersReq struct {
// 	// in:body
// 	Body model.User
// }

// Array with all users
// swagger:response getAllUsersResp
type swaggGetAllUsersResp struct {
	// in:body
	Body struct {
		// User model
		Users []model.User `json:"data"`
	}
}

// HTTP status code 400 response
// swagger:response queryErrResp
type swaggQueryErrResp struct {
	// in:body
	Body struct {
		// HTTP status code 400
		Code int `json:"code"`
	}
}

// HTTP status code 404 response
// swagger:response notFound
type swaggNotFound struct {
	// in:body
	Body struct {
		// HTTP status code 400
		Code int `json:"code"`
	}
}
