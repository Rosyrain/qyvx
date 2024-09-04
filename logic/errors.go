package logic

import "errors"

var (
	ErrorGithubAccessToken  = errors.New("github access token error")
	ErrorAddressAccessToken = errors.New("address access token error")
	ErrorRefreshToken       = errors.New("refresh token error")
	ErrorInvite             = errors.New("invite failed")
	ErrorGithubUser         = errors.New("invalid github account")
	ErrorGetGithubID        = errors.New("get githubID failed")
	ErrorUpdateAlias        = errors.New("update alias failed")
)
