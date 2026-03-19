package Dungeon

import (
	"context"
	"dungeons/app/functions"
	"dungeons/app/models"
	"dungeons/app/mongodb"
	"dungeons/app/server"
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
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
func (d *Dungeon) Create(in *models.Dungeon) (*models.Dungeon, error) {
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

func (d *Dungeon) GetByID(id string) (models.Dungeon, error) {
	var dungeon models.Dungeon

	srv := server.GetServer()
	collection := srv.Database.Collection(dungeon.Collection())

	filter := bson.M{"customID": id} // ✅ simple et safe

	err := collection.FindOne(context.TODO(), filter).Decode(&dungeon)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Error().Err(err).Msg("No document found")
			return dungeon, err
		}
		return dungeon, err
	}

	return dungeon, nil
}

func (d *Dungeon) Update(id string, in *models.Dungeon) error {
	var (
		doc 	   interface{}
		result    *mongo.UpdateResult
		err	  error
		queryParams models.QueryParams
		Dungeon models.Dungeon
	)

	srv := server.GetServer()

	err = d.Validate.Struct(in)
	if err != nil {
		log.Error().Err(err).Msg("")
		return err
	}

	Dungeon, err = d.GetByID(id)
	if err != nil {
		log.Error().Err(err).Msg("")
		return err
	}

	err = functions.ConvertInputStructToDataStruct(in, &doc)
	if err != nil {
		log.Error().Err(err).Msg("")
		return err
	}

	Dungeon.UpdatedAt = time.Now()
	collection := srv.Database.Collection(Dungeon.Collection())

	queryParams.FilterClause = append(queryParams.FilterClause, "customID," + id)
	filter := mongodb.SelectConstructeur(queryParams)
	if doc, err = mongodb.ToDoc(Dungeon); err != nil {
		log.Error().Err(err).Msg("")
		return err
	}

	update := bson.M{"$set": doc}
	result, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Error().Err(err).Msg("")
		return err
	}

	if result.MatchedCount == 0 {
		err = errors.New("Dungeon to be modified was not found")
	}

	if err == nil && result.ModifiedCount == 0 {
		err = errors.New("Dungeon could not be updated")
	}
	if err != nil {
		log.Error().Err(err).Msg("")
	}
	return err
}

func (d *Dungeon) Publish(id string, in *models.Dungeon) error {
	var err error
	srv := server.GetServer()

	dungeon, err := d.GetByID(id)
	if err != nil {
		log.Error().Err(err).Msg("")
		return err
	}

	dungeon.Status = "published"

	collection := srv.Database.Collection(dungeon.Collection())

	filter := bson.M{"customID": id} // ✅ simple et safe
	update := bson.M{"$set": dungeon}

	_, err = collection.UpdateOne(context.TODO(), filter, update)
	return err
}
