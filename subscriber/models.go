package subscriber

type Message struct {
	Message   []byte
	Topic     string
	Partition int
	Offset    int64
}
