package identify

type config struct {
	userAgent string
}

type Option func(*config)

func UserAgent(ua string) Option {
	return func(cfg *config) {
		cfg.userAgent = ua
	}
}
