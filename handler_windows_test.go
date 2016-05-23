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

// +build windows

package windowsauthtoken

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"unicode/utf16"

	"golang.org/x/sys/windows"
)

func TestCallbackIsNotInvokedIfHeaderIsNotPresent(t *testing.T) {
	calls := 0
	callbackSpy := func(_ string) error {
		calls++
		return nil
	}
	dummyHandler := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {

	})
	responseRecorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/", nil)

	handler := Handler(dummyHandler, callbackSpy, TokenUsername)
	handler.ServeHTTP(responseRecorder, request)

	if 0 != calls {
		t.Error("The callback was invoked but should not have been")
	}
}

func TestErrorIsReportedIfTokenHandleIsInvalid(t *testing.T) {
	dummyCallback := func(_ string) error {
		return nil
	}
	dummyHandler := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {

	})
	responseRecorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/", nil)
	request.Header.Set("X-IIS-WindowsAuthToken", "THIS-WILL-CAUSE-AN-ERROR")

	handler := Handler(dummyHandler, dummyCallback, TokenUsername)
	handler.ServeHTTP(responseRecorder, request)

	if http.StatusBadRequest != responseRecorder.Code {
		t.Error("Expected StatusBadRequest for the HTTP status code")
	}
}

func getWindowsUsernameForCurrentUser(t *testing.T, token windows.Token) string {
	user, err := token.GetTokenUser()
	if nil != err {
		t.Error(err)
	}

	var usernameLength uint32 = 50
	var domainNameLength uint32 = 50
	usernameChars := make([]uint16, usernameLength)
	domainNameChars := make([]uint16, domainNameLength)
	var use uint32
	err = windows.LookupAccountSid(
		nil,
		user.User.Sid,
		&usernameChars[0],
		&usernameLength,
		&domainNameChars[0],
		&domainNameLength,
		&use)
	if nil != err {
		t.Error(err)
	}

	username := string(utf16.Decode(usernameChars[0:usernameLength]))
	domainName := string(utf16.Decode(domainNameChars[0:domainNameLength]))
	return fmt.Sprintf("%s\\%s", domainName, username)
}

func TestHandlerInvokesCallbackWithWindowsUsername(t *testing.T) {
	token, err := windows.OpenCurrentProcessToken()
	if nil != err {
		t.Error(err)
	}

	expectedUsername := getWindowsUsernameForCurrentUser(t, token)
	var actualUsername string
	mockCallback := func(username string) error {
		actualUsername = username
		return nil
	}
	dummyHandler := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {

	})
	responseRecorder := httptest.NewRecorder()
	request, err := http.NewRequest("GET", "/", nil)
	if nil != err {
		t.Error(err)
	}

	request.Header.Set("X-IIS-WindowsAuthToken", strconv.FormatUint(uint64(token), 16))

	handler := Handler(dummyHandler, mockCallback, TokenUsername)
	handler.ServeHTTP(responseRecorder, request)

	if expectedUsername != actualUsername {
		t.Errorf("Expected \"%s\", but was \"%s\"", expectedUsername, actualUsername)
	}
}

func getSidForCurrentUser(t *testing.T, token windows.Token) string {
	user, err := token.GetTokenUser()
	if nil != err {
		t.Error(err)
	}

	expectedSid, err := user.User.Sid.String()
	if nil != err {
		t.Error(err)
	}

	return expectedSid
}

func TestHandlerInvokesCallbackWithUserSid(t *testing.T) {
	token, err := windows.OpenCurrentProcessToken()
	if nil != err {
		t.Error(err)
	}

	expectedSid := getSidForCurrentUser(t, token)
	var actualSid string
	mockCallback := func(sid string) error {
		actualSid = sid
		return nil
	}
	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
	responseRecorder := httptest.NewRecorder()
	request, err := http.NewRequest("GET", "/", nil)
	if nil != err {
		t.Error(err)
	}

	request.Header.Set("X-IIS-WindowsAuthToken", strconv.FormatUint(uint64(token), 16))

	handler := Handler(dummyHandler, mockCallback, TokenSid)
	handler.ServeHTTP(responseRecorder, request)

	if expectedSid != actualSid {
		t.Errorf("Expected \"%s\", but was \"%s\"", expectedSid, actualSid)
	}
}
