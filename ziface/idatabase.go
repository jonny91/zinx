package ziface

type IDatabase interface {
	Connect() error
}
