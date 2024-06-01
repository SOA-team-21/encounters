package repo

import (
	"context"
	"encounters/model"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type EncounterRepo struct {
	cli    *mongo.Client
	logger *log.Logger
}

func New(ctx context.Context, logger *log.Logger) (*EncounterRepo, error) {
	dburi := os.Getenv("MONGO_DB_URI")

	client, err := mongo.NewClient(options.Client().ApplyURI(dburi))
	if err != nil {
		return nil, err
	}

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return &EncounterRepo{
		cli:    client,
		logger: logger,
	}, nil
}

func (pr *EncounterRepo) Disconnect(ctx context.Context) error {
	err := pr.cli.Disconnect(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (pr *EncounterRepo) Ping() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check connection -> if no error, connection is established
	err := pr.cli.Ping(ctx, readpref.Primary())
	if err != nil {
		pr.logger.Println(err)
	}

	// Print available databases
	databases, err := pr.cli.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		pr.logger.Println(err)
	}
	fmt.Println(databases)
}

//hiddenEncounter

func (pr *EncounterRepo) GetAllHiddenEncounters() (model.HiddenEncounters, error) {
	// Initialise context (after 5 seconds timeout, abort operation)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	hiddenEncountersCollection := pr.getHiddenEncounterCollection()

	var hiddenEncounters model.HiddenEncounters
	encountersCursor, err := hiddenEncountersCollection.Find(ctx, bson.M{})
	if err != nil {
		pr.logger.Println(err)
		return nil, err
	}
	if err = encountersCursor.All(ctx, &hiddenEncounters); err != nil {
		pr.logger.Println(err)
		return nil, err
	}
	return hiddenEncounters, nil
}

func (pr *EncounterRepo) GetHiddenEncounterById(id string) (*model.HiddenEncounter, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	hiddenEncountersColllection := pr.getHiddenEncounterCollection()

	var hiddenEncounter model.HiddenEncounter
	objID, _ := primitive.ObjectIDFromHex(id)
	err := hiddenEncountersColllection.FindOne(ctx, bson.M{"_id": objID}).Decode(&hiddenEncounter)
	if err != nil {
		pr.logger.Println(err)
		return nil, err
	}
	return &hiddenEncounter, nil
}

func (pr *EncounterRepo) InsertHiddenEncounter(hiddenEncounter *model.HiddenEncounter) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	hiddenEncountersCollection := pr.getHiddenEncounterCollection()

	result, err := hiddenEncountersCollection.InsertOne(ctx, &hiddenEncounter)
	if err != nil {
		pr.logger.Println(err)
		return err
	}
	pr.logger.Printf("Documents ID: %v\n", result.InsertedID)
	return nil
}

func (pr *EncounterRepo) UpdateHiddenEncounter(id string, hiddenEncounter *model.HiddenEncounter) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	hiddenEncountersCollection := pr.getHiddenEncounterCollection()

	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objID}
	update := bson.M{"$set": bson.M{
		"name":        hiddenEncounter.Name,
		"description": hiddenEncounter.Description,
	}}
	result, err := hiddenEncountersCollection.UpdateOne(ctx, filter, update)
	pr.logger.Printf("Documents matched: %v\n", result.MatchedCount)
	pr.logger.Printf("Documents updated: %v\n", result.ModifiedCount)

	if err != nil {
		pr.logger.Println(err)
		return err
	}
	return nil
}

func (pr *EncounterRepo) ActivateHiddenEncounter(id string, participantLocation model.ParticipantLocation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	hiddenEncountersCollection := pr.getHiddenEncounterCollection()
	hiddenEncounter, _ := pr.GetHiddenEncounterById(id)

	result := hiddenEncounter.CanActivateEncounter(participantLocation.Username, participantLocation.Latitude, participantLocation.Longitude)

	if result {
		newParticipant := model.Participant{
			Username: participantLocation.Username,
		}

		objID, _ := primitive.ObjectIDFromHex(id)
		filter := bson.D{{Key: "_id", Value: objID}}
		update := bson.M{"$push": bson.M{
			"encounter.participants": newParticipant,
		}}

		result, err := hiddenEncountersCollection.UpdateOne(ctx, filter, update)
		pr.logger.Printf("Documents matched: %v\n", result.MatchedCount)
		pr.logger.Printf("Documents updated: %v\n", result.ModifiedCount)

		if err != nil {
			pr.logger.Println(err)
			return err
		}
	} else {
		return fmt.Errorf("activation failed")
	}
	return nil
}

func (pr *EncounterRepo) SolveHiddenEncounter(id string, participantLocation model.ParticipantLocation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	hiddenEncountersCollection := pr.getHiddenEncounterCollection()
	hiddenEncounter, _ := pr.GetHiddenEncounterById(id)

	result := hiddenEncounter.CanSolveEncounter(participantLocation.Username, participantLocation.Longitude, participantLocation.Latitude)

	if result {
		now := time.Now()
		completer := model.Completer{
			Username:       participantLocation.Username,
			CompletionDate: now,
		}

		objID, _ := primitive.ObjectIDFromHex(id)

		filter := bson.D{{Key: "_id", Value: objID}}

		updateCompleters := bson.M{"$push": bson.M{
			"encounter.completers": completer,
		}}
		resultCompleters, err := hiddenEncountersCollection.UpdateOne(ctx, filter, updateCompleters)
		if err != nil {
			pr.logger.Println(err)
			return err
		}
		pr.logger.Printf("Documents matched: %v\n", resultCompleters.MatchedCount)
		pr.logger.Printf("Documents updated: %v\n", resultCompleters.ModifiedCount)

		updateParticipants := bson.M{"$pull": bson.M{
			"encounter.participants": bson.M{"username": completer.Username},
		}}
		resultParticipants, err := hiddenEncountersCollection.UpdateOne(ctx, filter, updateParticipants)
		if err != nil {
			pr.logger.Println(err)
			return err
		}
		pr.logger.Printf("Documents matched: %v\n", resultParticipants.MatchedCount)
		pr.logger.Printf("Documents updated: %v\n", resultParticipants.ModifiedCount)

	} else {
		return fmt.Errorf("activation failed")
	}
	return nil
}

func (hr *EncounterRepo) getHiddenEncounterCollection() *mongo.Collection {
	hiddenEncounterDatabase := hr.cli.Database("mongoDemo")
	hiddenEncounterCollection := hiddenEncounterDatabase.Collection("hiddenEncounters")
	return hiddenEncounterCollection
}

//socialEncounter

func (pr *EncounterRepo) GetAllSocialEncounters() ([]model.SocialEncounter, error) {
	// Initialise context (after 5 seconds timeout, abort operation)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	socialEncountersCollection := pr.getSocialEncounterCollection()

	var socialEncounters []model.SocialEncounter
	encountersCursor, err := socialEncountersCollection.Find(ctx, bson.M{})
	if err != nil {
		pr.logger.Println(err)
		return nil, err
	}
	if err = encountersCursor.All(ctx, &socialEncounters); err != nil {
		pr.logger.Println(err)
		return nil, err
	}
	return socialEncounters, nil
}

func (pr *EncounterRepo) GetSocialEncounterById(id string) (*model.SocialEncounter, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	socialEncountersColllection := pr.getSocialEncounterCollection()

	var socialEncounter model.SocialEncounter
	objID, _ := primitive.ObjectIDFromHex(id)
	err := socialEncountersColllection.FindOne(ctx, bson.M{"_id": objID}).Decode(&socialEncounter)
	if err != nil {
		pr.logger.Println(err)
		return nil, err
	}
	return &socialEncounter, nil
}

func (pr *EncounterRepo) InsertSocialEncounter(socialEncounter *model.SocialEncounter) (*model.SocialEncounter, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	socialEncountersCollection := pr.getSocialEncounterCollection()

	result, err := socialEncountersCollection.InsertOne(ctx, &socialEncounter)
	if err != nil {
		pr.logger.Println(err)
		return nil, err
	}
	pr.logger.Printf("Documents ID: %v\n", result.InsertedID)
	return socialEncounter, nil
}

func (pr *EncounterRepo) UpdateSocialEncounter(id string, socialEncounter *model.SocialEncounter) (*model.SocialEncounter, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	socialEncountersCollection := pr.getSocialEncounterCollection()

	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objID}
	update := bson.M{"$set": bson.M{
		"name":        socialEncounter.Name,
		"description": socialEncounter.Description,
	}}
	result, err := socialEncountersCollection.UpdateOne(ctx, filter, update)
	pr.logger.Printf("Documents matched: %v\n", result.MatchedCount)
	pr.logger.Printf("Documents updated: %v\n", result.ModifiedCount)

	if err != nil {
		pr.logger.Println(err)
		return nil, err
	}
	return socialEncounter, nil
}

func (pr *EncounterRepo) ActivateSocialEncounter(id string, participantLocation model.ParticipantLocation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	socialEncountersCollection := pr.getSocialEncounterCollection()
	socialEncounter, _ := pr.GetSocialEncounterById(id)

	result := socialEncounter.CanActivateEncounter(participantLocation.Username, participantLocation.Latitude, participantLocation.Longitude)

	if result {
		newParticipant := model.Participant{
			Username: participantLocation.Username,
		}

		objID, _ := primitive.ObjectIDFromHex(id)
		filter := bson.D{{Key: "_id", Value: objID}}
		update := bson.M{"$push": bson.M{
			"encounter.participants": newParticipant,
		}}

		result, err := socialEncountersCollection.UpdateOne(ctx, filter, update)
		pr.logger.Printf("Documents matched: %v\n", result.MatchedCount)
		pr.logger.Printf("Documents updated: %v\n", result.ModifiedCount)

		if err != nil {
			pr.logger.Println(err)
			return err
		}
	} else {
		return fmt.Errorf("activation failed")
	}
	return nil
}

func (hr *EncounterRepo) getSocialEncounterCollection() *mongo.Collection {
	socialEncounterDatabase := hr.cli.Database("mongoDemo")
	socialEncounterCollection := socialEncounterDatabase.Collection("socialEncounters")
	return socialEncounterCollection
}
