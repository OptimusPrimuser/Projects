package globals

import (
	"backend/mysql_client"

	"backend/redis_store"
)

// var RMQ rabbit_mq.RMQ
var Redis redis_store.Redis
var MYSQL mysql_client.Client
