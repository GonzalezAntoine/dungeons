package dungeons

import (
	"context"
	"dungeons/app/functions"
	"dungeons/app/models"
	"dungeons/app/mongodb"
	"dungeons/app/server"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Dungeon struct {
	Validate *validator.Validate
}

func New() *Dungeon {
	return &Dungeon{
		Validate: validator.New(),
	}
}

func (d *Dungeon) Get(queryParams models.QueryParams) ([]models.Dungeon, error) {
	var(
		err error
		dungeons []models.Dungeon
		dungeon models.Dungeon
		cursor *mongo.Cursor
	)

	srv := server.GetServer()
	collection := srv.Database.Collection(dungeon.Collection())

	filter := mongodb.SelectConstructeur(queryParams)
	cursor, err = collection.Find(context.TODO(), filter)
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		// A new result variable should be declared for each document.
		var dungeon models.Dungeon
		err = cursor.Decode(&dungeon)
		if err != nil {
			log.Error().Err(err).Msg("")
			return nil, err
		}
		dungeons = append(dungeons, dungeon)
	}

	err = cursor.Err()
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}
	
	return dungeons, err
}

// create new dungeon in db
func (d *Dungeon) Create(in models.Dungeon) (*models.Dungeon, error) {
	var dungeon models.Dungeon

	srv := server.GetServer()
	collection := srv.Database.Collection(dungeon.Collection())

	// check input fields
	err := d.Validate.Struct(in)
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	err = functions.ConvertInputStructToDataStruct(in, &dungeon)
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	dungeon.CustomID = functions.NewUUID()
	dungeon.CreatedAt = time.Now()

	_, err = collection.InsertOne(context.TODO(), dungeon)
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	return &dungeon, nil
}