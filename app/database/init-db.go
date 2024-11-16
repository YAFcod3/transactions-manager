package database

func Init() {
	InitRedis()
	InitMongo()
}

func Close() {
	CloseRedis()
	CloseMongo()
}
