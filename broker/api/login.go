// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package api

import (
	"net/http"

	"github.com/lindb/lindb/broker/middleware"
	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/pkg/logger"
)

var log = logger.GetLogger("broker", "api")

// LoginAPI represents login param
type LoginAPI struct {
	user config.User
	auth middleware.Authentication
}

// NewLoginAPI creates login api instance
func NewLoginAPI(user config.User, auth middleware.Authentication) *LoginAPI {
	return &LoginAPI{
		user: user,
		auth: auth,
	}
}

// Login responses unique token
// if use name or password is empty will responses error msg
// if use name or password is error also will responses error msg
func (l *LoginAPI) Login(w http.ResponseWriter, r *http.Request) {
	user := config.User{}
	err := GetJSONBodyFromRequest(r, &user)
	// login request is error
	if err != nil {
		log.Error("cannot get user info from request")
		OK(w, "")
		return
	}
	// user name is empty
	if len(user.UserName) == 0 {
		log.Error("username is empty")
		OK(w, "")
		return
	}
	// password is empty
	if len(user.Password) == 0 {
		log.Error("password is empty")
		OK(w, "")
		return
	}
	// user name is error
	if l.user.UserName != user.UserName {
		log.Error("username is invalid")
		OK(w, "")
		return
	}
	// password is error
	if l.user.Password != user.Password {
		log.Error("password is invalid")
		OK(w, "")
		return
	}
	token, err := l.auth.CreateToken(user)
	if err != nil {
		OK(w, "")
		return
	}
	OK(w, token)
}

// Check responses use msg
// this method use for test
func (l *LoginAPI) Check(w http.ResponseWriter, r *http.Request) {
	OK(w, l.user)
}
