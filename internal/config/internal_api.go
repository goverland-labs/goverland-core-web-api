package config

type InternalAPI struct {
	CoreStorageAddress string `env:"INTERNAL_API_CORE_STORAGE_ADDRESS" envDefault:"localhost:11100"`
	CoreFeedAddress    string `env:"INTERNAL_API_CORE_FEED_ADDRESS" envDefault:"localhost:11000"`
}
