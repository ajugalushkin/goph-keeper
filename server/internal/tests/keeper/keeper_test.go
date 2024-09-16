package keeper

import (
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ajugalushkin/goph-keeper/client/secret"
	"github.com/ajugalushkin/goph-keeper/client/vaulttypes"
	keeperv1 "github.com/ajugalushkin/goph-keeper/gen/keeper/v1"
	"github.com/ajugalushkin/goph-keeper/server/internal/tests/keeper/suite"
)

type Item struct {
	Name     string
	Email    string
	Password string
}

// Type возвращает тип хранимой информации
func (i Item) Type() vaulttypes.VaultType {
	var data vaulttypes.VaultType
	data = "item"
	return data
}

// String функция отображения приватной информации
func (i Item) String() string {
	return "ITEM DATA"
}

func TestCRUDItem_CRUDItem_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)
	defer st.Closer()

	nameExpected := gofakeit.Name()

	data := Item{
		Name:     nameExpected,
		Email:    gofakeit.Email(),
		Password: suite.RandomFakePassword(),
	}

	content, err := secret.EncryptSecret(data)
	require.NoError(t, err)

	resp, err := st.KeeperClient.CreateItemV1(ctx, &keeperv1.CreateItemRequestV1{
		Name:    nameExpected,
		Content: content,
	})
	require.NoError(t, err)
	assert.Equal(t, nameExpected, resp.GetName())

	respGet, err := st.KeeperClient.GetItemV1(ctx, &keeperv1.GetItemRequestV1{
		Name: nameExpected,
	})
	require.NoError(t, err)
	assert.Equal(t, nameExpected, respGet.GetName())
	assert.Equal(t, resp.GetVersion(), respGet.GetVersion())

	dataUpd := Item{
		Name:     nameExpected,
		Email:    gofakeit.Email(),
		Password: suite.RandomFakePassword(),
	}

	contentUpd, err := secret.EncryptSecret(dataUpd)
	require.NoError(t, err)

	respUpd, err := st.KeeperClient.UpdateItemV1(ctx, &keeperv1.UpdateItemRequestV1{
		Name:    nameExpected,
		Content: contentUpd,
	})
	require.NoError(t, err)
	assert.Equal(t, nameExpected, respUpd.GetName())

	respGet, err = st.KeeperClient.GetItemV1(ctx, &keeperv1.GetItemRequestV1{
		Name: nameExpected,
	})
	require.NoError(t, err)
	assert.Equal(t, nameExpected, respGet.GetName())
	assert.Equal(t, respUpd.GetVersion(), respGet.GetVersion())
}

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
