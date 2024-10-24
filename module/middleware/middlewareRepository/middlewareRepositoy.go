package middlewareRepository

type MiddlewareRepository interface{
	AccessTokenSearch(grpcUrl string, accessToken string) error
	BookShelfSearch(grpcUrl string, userID uint, bookID uint) error 
}