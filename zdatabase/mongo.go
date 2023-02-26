package zdatabase

import (
	"context"
	"errors"
	"fmt"
	"github.com/jonny91/zinx/utils"
	"github.com/jonny91/zinx/zlog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Mongo struct {
	client *mongo.Client
}

var connectedChan chan bool

func (db *Mongo) Init(background context.Context) (bool, error) {
	ctx, cancel := context.WithTimeout(background, time.Second*5)
	defer cancel()

	connectedChan = make(chan bool)
	go db.Connect(ctx)
	select {
	case <-ctx.Done():
		zlog.Errorf("mongo connect timeout ...")
		return false, errors.New("mongo connect timeout")
	case ok := <-connectedChan:
		if !ok {
			zlog.Errorf("mongo connect failed ...")
			return false, errors.New("mongo connect failed")
		} else {
			return true, nil
		}
	}
}

func (db *Mongo) GetName() string {
	return "mongo"
}

func (db *Mongo) Connect(ctx context.Context) (bool, error) {
	var err error
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%d",
		utils.GlobalObject.Database.Username,
		utils.GlobalObject.Database.Password,
		utils.GlobalObject.Database.Host,
		utils.GlobalObject.Database.Port,
	)
	zlog.Infof("try to connect database %s\n", uri)
	clientOptions := options.Client().ApplyURI(uri)
	db.client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		connectedChan <- false
		return false, err
	}

	err = db.client.Ping(ctx, nil)
	if err != nil {
		connectedChan <- false
		return false, err
	}
	zlog.Infof("mongo connect success ...\n")
	connectedChan <- true
	return true, nil
}

func (db *Mongo) Close(ctx context.Context) error {
	if db.client != nil {
		err := db.client.Disconnect(ctx)
		return err
	}
	return errors.New("database disconnected")
}
