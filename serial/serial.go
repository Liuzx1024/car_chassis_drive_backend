package serial

func OpenPort(config *Config) (*Port, error) {
	return openPort(config)
}
