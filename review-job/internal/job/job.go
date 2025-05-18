package job

import "github.com/google/wire"

// ProviderSet is data providers. wire 提供者
var ProviderSet = wire.NewSet(NewJobWorker, NewKafkaReader, NewESClient)
