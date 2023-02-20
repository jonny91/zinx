package zdatabase

import (
	"context"
	"errors"
	"fmt"
	"github.com/jonny91/zinx/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	client *mongo.Client
}

func (db *Mongo) Init() {
}

func (db *Mongo) GetName() string {
	return "db"
}

func (db *Mongo) Connect(ctx context.Context) (bool, error) {
	var err error
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%d",
		utils.GlobalObject.Database.Username,
		utils.GlobalObject.Database.Password,
		utils.GlobalObject.Database.Host,
		utils.GlobalObject.Database.Port,
	)
	clientOptions := options.Client().ApplyURI(uri)
	db.client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		return false, err
	}

	err = db.client.Ping(ctx, nil)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (db *Mongo) Close(ctx context.Context) error {
	if db.client != nil {
		err := db.client.Disconnect(ctx)
		return err
	}
	return errors.New("database disconnected")
}
