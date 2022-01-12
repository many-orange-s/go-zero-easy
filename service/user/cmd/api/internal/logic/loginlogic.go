package logic

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"go-zero-easy/commen/errorx"
	"go-zero-easy/commen/errorx/errconcrete"
	"go-zero-easy/service/user/cmd/api/internal/svc"
	"go-zero-easy/service/user/cmd/api/internal/types"
	"log"
	"time"

	"github.com/tal-tech/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) LoginLogic {
	return LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req types.UserLogin) (resp *types.JWT, err error) {
	account := req.Account
	password := req.Password

	usermsg, err := l.svcCtx.UserModel.FindOneByAccount(account)
	if err != nil {
		if ok := errors.Is(err, sqlx.ErrNotFound); ok {
			return nil, &errorx.CodeError{Code: errorx.InvalidParam, Msg: errconcrete.SqlNotFound}
		} else {
			log.Println("Login FindOneFound err :", err)
			return nil, &errorx.CodeError{Code: errorx.SystemBusy, Msg: errconcrete.InterErr}
		}
	}

	if usermsg.Password != password {
		return nil, errorx.NewCodeErr(errorx.InvalidParam, errconcrete.PasswordErr)
	}

	now := time.Now().Unix()
	secret := l.svcCtx.Config.Auth.AccessSecret
	expire := l.svcCtx.Config.Auth.AccessExpire
	jwtoken, err := l.getJwtToken(secret, now, expire, usermsg.Uid)
	if err != nil {
		log.Println("Login getJwtToken err:", err)
		return nil, errorx.NewCodeErr(errorx.SystemBusy, errconcrete.InterErr)
	}

	// 在返回的时候必须要初始化
	resp = new(types.JWT)
	resp.Token = jwtoken
	return
}

func (l *LoginLogic) getJwtToken(secretKey string, iat, seconds, userId int64) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["userId"] = userId
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}
