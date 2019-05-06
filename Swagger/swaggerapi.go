package swagger

import model "FirstProject/Model"

// type swaggerAPI struct {
// 	GetAllUsersResp []model.User
// }

// // Request containing array
// // swagger:parameters getAllUsersReq
// type swaggerGetAllUsersReq struct {
// 	// in:body
// 	Body swaggerAPI
// }

// HTTP status code 200 and array of users in data
// swagger:response getAllUsersResp
type swaggGetAllUsersResp struct {
	// in:body
	Body struct {
		// HTTP status code 200/201
		Code int `json:"code"`
		// Repository model
		Users []model.User `json:"data"`
	}
}
