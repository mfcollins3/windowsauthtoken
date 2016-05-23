// Copyright 2016 Michael F. Collins, III
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this softwar and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
// DEALINGS IN THE SOFTWARE.

package windowsauthtoken

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"golang.org/x/sys/windows"
)

func getWindowsAuthToken(r *http.Request) (windows.Token, error) {
	header := r.Header.Get("X-IIS-WindowsAuthToken")
	if "" == header {
		return 0, nil
	}

	handle, err := strconv.ParseUint(header, 16, 0)
	if nil != err {
		return 0, err
	}

	return windows.Token(handle), nil
}

func getWindowsUsername(token windows.Token) (string, error) {
	user, err := token.GetTokenUser()
	if nil != err {
		return "", err
	}

	username, domainName, _, err := user.User.Sid.LookupAccount("")
	if nil != err {
		return "", err
	}

	return fmt.Sprintf("%s\\%s", domainName, username), nil
}

func getSid(token windows.Token) (string, error) {
	user, err := token.GetTokenUser()
	if nil != err {
		return "", err
	}

	return user.User.Sid.String()
}

type handler struct {
	next       http.Handler
	callback   Callback
	tokenValue TokenValue
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	token, err := getWindowsAuthToken(r)
	if nil != err {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if 0 != token {
		defer token.Close()

		var tokenValue string
		switch h.tokenValue {
		case TokenUsername:
			tokenValue, err = getWindowsUsername(token)
		case TokenSid:
			tokenValue, err = getSid(token)
		default:
			err = errors.New("Unrecognized TokenValue")
		}

		if nil != err {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = h.callback(tokenValue)
		if nil != err {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	h.next.ServeHTTP(w, r)
}

// Handler returns an HTTP handler that will process the token for the
// authenticated user that is passed to the web application in the
// X-IIS-WindowsAuthToken HTTP header. Handler will obtain the Windows
// username for the authenticated user and will pass the username to
// the web application using the callback parameter.
func Handler(next http.Handler, callback Callback, tokenValue TokenValue) http.Handler {
	return &handler{next, callback, tokenValue}
}
