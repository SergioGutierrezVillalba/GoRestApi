// Copyright 2019 Sergio Gutiérrez. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Golang RestAPI
//
// This documentation describes a RESTful API based in Golang
//
//     Schemes: https, ws, wss
//     BasePath: /doc
//     Version: 1.1.0
//     License: MIT http://opensource.org/licenses/MIT
//     Contact: Sergio Gutierrez <gv.sergio@gmail.com>
//
//     Consumes:
//     - application/json
//     - multipart/form-data
//
//     Produces:
//     - application/json
//
//     Security:
//     - bearer
//
//     SecurityDefinitions:
//     bearer:
//          type: apiKey
//          name: Authorization
//          in: header
//
// swagger:meta

package main

import (
	"FirstProject/Model/api"
	"FirstProject/Model/database"
)

var (
	Api api.Api
	Db  database.Db
)

func main() {
	Db.StartConnection()
	Api.Start(Db.Session)
}
