package token_cache

type Storage interface {
	Load() (accessToken string, err error)
	Save(accessToken string) error
}
