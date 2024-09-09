package auth

//const (
//	passDefaultLen = 10
//)

//func TestLogin_Login_HappyPath(t *testing.T) {
//	ctx, st := suite.New(t)
//
//	email := gofakeit.Email()
//	pass := randomFakePassword()
//
//	respReq, err := st.AuthClient.RegisterV1(ctx, &authv1.RegisterRequestV1{
//		Email:    email,
//		Password: pass,
//	})
//	require.NoError(t, err)
//	assert.NotEmpty(t, respReq.GetUserId())
//
//	respLog, err := st.AuthClient.LoginV1(ctx, &authv1.LoginRequestV1{
//		Email:    email,
//		Password: pass,
//	})
//	require.NoError(t, err)
//
//	loginTime := time.Now()
//
//	token := respLog.GetToken()
//	require.NotEmpty(t, token)
//
//	tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
//		return []byte(st.Cfg.Token.Secret), nil
//	})
//	require.NoError(t, err)
//
//	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
//	assert.True(t, ok)
//
//	assert.Equal(t, respReq.GetUserId(), int64(claims["uid"].(float64)))
//	assert.Equal(t, email, claims["email"].(string))
//
//	const deltaSeconds = 1
//	assert.InDelta(t, loginTime.Add(st.Cfg.Token.TTL).Unix(), claims["exp"].(float64), deltaSeconds)
//}
//
//func TestRegisterLogin_DuplicatedRegistration(t *testing.T) {
//	ctx, st := suite.New(t)
//
//	email := gofakeit.Email()
//	pass := randomFakePassword()
//
//	respReg, err := st.AuthClient.RegisterV1(ctx, &authv1.RegisterRequestV1{
//		Email:    email,
//		Password: pass,
//	})
//	require.NoError(t, err)
//	require.NotEmpty(t, respReg.GetUserId())
//
//	respReg, err = st.AuthClient.RegisterV1(ctx, &authv1.RegisterRequestV1{
//		Email:    email,
//		Password: pass,
//	})
//	require.Error(t, err)
//	assert.Empty(t, respReg.GetUserId())
//	assert.ErrorContains(t, err, "user already exists")
//}
//
//func TestRegister_FailCases(t *testing.T) {
//	ctx, st := suite.New(t)
//
//	tests := []struct {
//		name        string
//		email       string
//		password    string
//		expectedErr string
//	}{
//		{
//			name:        "Register with Empty Password",
//			email:       gofakeit.Email(),
//			password:    "",
//			expectedErr: "validation error:\n - password: value length must be at least 8 characters [string.min_len]",
//		},
//		{
//			name:        "Register with Empty Email",
//			email:       "",
//			password:    randomFakePassword(),
//			expectedErr: "validation error:\n - email: value is empty, which is not a valid email address [string.email_empty]",
//		},
//		{
//			name:        "Register with Both Empty",
//			email:       "",
//			password:    "",
//			expectedErr: "validation error:\n - email: value is empty, which is not a valid email address [string.email_empty]",
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			_, err := st.AuthClient.RegisterV1(ctx, &authv1.RegisterRequestV1{
//				Email:    tt.email,
//				Password: tt.password,
//			})
//			require.Error(t, err)
//			require.Contains(t, err.Error(), tt.expectedErr)
//
//		})
//	}
//}
//
//func TestLogin_FailCases(t *testing.T) {
//	ctx, st := suite.New(t)
//
//	tests := []struct {
//		name        string
//		email       string
//		password    string
//		expectedErr string
//	}{
//		{
//			name:        "Login with Empty Password",
//			email:       gofakeit.Email(),
//			password:    "",
//			expectedErr: "validation error:\n - password: value length must be at least 8 characters [string.min_len]",
//		},
//		{
//			name:        "Login with Empty Email",
//			email:       "",
//			password:    randomFakePassword(),
//			expectedErr: "validation error:\n - email: value is empty, which is not a valid email address [string.email_empty]",
//		},
//		{
//			name:        "Login with Both Empty Email and Password",
//			email:       "",
//			password:    "",
//			expectedErr: "validation error:\n - email: value is empty, which is not a valid email address [string.email_empty]",
//		},
//		{
//			name:        "Login with Non-Matching Password",
//			email:       gofakeit.Email(),
//			password:    randomFakePassword(),
//			expectedErr: "invalid email or password",
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			_, err := st.AuthClient.RegisterV1(ctx, &authv1.RegisterRequestV1{
//				Email:    gofakeit.Email(),
//				Password: randomFakePassword(),
//			})
//			require.NoError(t, err)
//
//			_, err = st.AuthClient.LoginV1(ctx, &authv1.LoginRequestV1{
//				Email:    tt.email,
//				Password: tt.password,
//			})
//			require.Error(t, err)
//			require.Contains(t, err.Error(), tt.expectedErr)
//		})
//	}
//}
//
//func randomFakePassword() string {
//	return gofakeit.Password(true, true, true, true, false, passDefaultLen)
//}
