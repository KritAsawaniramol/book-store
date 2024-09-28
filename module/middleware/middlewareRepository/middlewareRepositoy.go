package middlewareRepository

type MiddlewareRepository interface{
	AccessTokenSearch(grpcUrl string, accessToken string) error
}