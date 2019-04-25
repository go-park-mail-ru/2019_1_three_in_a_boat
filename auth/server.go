package main

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/logger"
	"gopkg.in/square/go-jose.v2/jwt"

	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/db"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/formats"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/formats/pb"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/settings/auth"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/settings/shared"
)

type AuthService struct{}

func (*AuthService) Authorize(
	ctx context.Context, in *pb.AuthorizeRequest) (*pb.AuthorizeReply, error) {
	u, err := db.GetUserByUsernameOrEmail(
		settings.DB(),
		in.Username,
		in.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			// no such login
			return &pb.AuthorizeReply{
				Ok:      false,
				Claims:  nil,
				Message: formats.ErrInvalidCredentials,
			}, nil
		} else {
			return &pb.AuthorizeReply{
				Ok:      false,
				Claims:  nil,
				Message: formats.ErrSqlFailure,
			}, err
		}
	}

	ok, err := db.AccountComparePasswordToHash(in.Password, u.Account.Password)
	if err != nil {
		return &pb.AuthorizeReply{
			Ok:      false,
			Claims:  nil,
			Message: formats.ErrPasswordHashing,
		}, err
	}

	if !ok {
		// password mismatch
		return &pb.AuthorizeReply{
			Ok:      false,
			Claims:  nil,
			Message: formats.ErrInvalidCredentials,
		}, nil
	}

	if token, err := tokenizeUser(u); err != nil {
		return &pb.AuthorizeReply{
			Ok:      false,
			Claims:  nil,
			Message: formats.ErrJWTEncryptionFailure,
		}, err
	} else {
		return &pb.AuthorizeReply{
			Ok:      true,
			Claims:  nil,
			Message: token,
		}, nil
	}
}

func (*AuthService) CheckAuthorize(
	ctx context.Context, in *pb.CheckAuthorizeRequest) (
	*pb.CheckAuthorizeReply, error) {
	var parsedJWT *jwt.JSONWebToken
	var err error
	claims := &pb.Claims{}
	errMsg := ""

	parsedJWT, err = jwt.ParseSigned(in.Token)
	if err != nil {
		errMsg = formats.ErrJWTDecryptionFailure
	} else {
		err = parsedJWT.Claims(&auth_settings.GetSecretKey().PublicKey, claims)
		if err != nil {
			errMsg = formats.ErrJWTDecryptionFailure
		} else if claims.Uid == 0 {
			// special case: err == nil, but decryption failed
			// (can happen according to some stupid git issue i lost lol)
			errMsg = formats.ErrJWTDecryptionEmpty
			return &pb.CheckAuthorizeReply{
				Ok:      false,
				Claims:  claims,
				Message: errMsg,
			}, errors.New("failed to decrypt with no error: empty uid")
		}
	}

	if err != nil {
		ctx = formats.NewAuthContext(context.Background(), nil)
		//noinspection GoNilness
		logger.Errorf("Decryption failure: %v", err)
		return &pb.CheckAuthorizeReply{
			Ok:      false,
			Message: errMsg,
		}, err
	}

	return &pb.CheckAuthorizeReply{
		Ok:      true,
		Claims:  claims,
		Message: "",
	}, nil

}

func (*AuthService) Tokenize(
	ctx context.Context, in *pb.Claims) (
	*pb.Token, error) {
	token, err := tokenizeClaims(in)
	if err != nil {
		return nil, err
	}

	return &pb.Token{Token: token}, err
}
