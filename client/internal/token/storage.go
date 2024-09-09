package token

type Storage interface {
	Load() (accessToken string, err error)
	Save(accessToken string) error
}
