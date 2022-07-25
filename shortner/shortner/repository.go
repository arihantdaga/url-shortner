package shortner

type DbRepository struct {
}

type Repository interface {
	findOne(id string) (ShortUrl, error)
	findUrl(shorturl string) (string, error)
}

func (r *DbRepository) findOne(id string) (ShortUrl, error) {
	return ShortUrl{}, nil
}
func (r *DbRepository) findUrl(shorturl string) (string, error) {
	return "", nil
}

func NewDbRepository() Repository {
	return &DbRepository{}
}
