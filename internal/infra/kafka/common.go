package kafka

func getKafkaBroker(kafkaBroker string, isDev bool) string {
	if broker := kafkaBroker; broker != "" {
		return broker
	}

	if isDev {
		return "localhost:9094"
	}

	return "kafka:9092"
	// fallback по окружению (вдруг переменной нет)
	//if _, err := os.Stat("/.dockerenv"); err == nil {
	// внутри Docker
	//return "kafka:9092"
	//}
	// локально
	//return "localhost:9094"
}
