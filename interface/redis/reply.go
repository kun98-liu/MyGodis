package redis

//Redis serialization protocal
type Reply interface {
	ToBytes() []byte
}
