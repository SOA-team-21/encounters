package handler

import (
	"context"
	//"encoding/json"
	"encounters/model"
	"encounters/proto/encounters"
	"encounters/repo"
	"fmt"

	"google.golang.org/protobuf/types/known/timestamppb"
	//"github.com/gorilla/mux"
)

type KeyProduct struct{}

type EncounterHandler struct {
	Repo *repo.EncounterRepo
	encounters.UnimplementedEncountersServiceServer
}

func (handler *EncounterHandler) GetAllSocialEncounters(ctx context.Context, request *encounters.EmptyRequest) (*encounters.SocialEncountersResponse, error) {

	var fromDb, err = handler.Repo.GetAllSocialEncounters()
	if err != nil {
		return &encounters.SocialEncountersResponse{}, err
	}
	return SocialEncountersToRpc(fromDb), nil
}

func (handler *EncounterHandler) GetSocialEncounterById(ctx context.Context, request *encounters.SocialEncounterIdRequest) (*encounters.SocialEncounterResponse, error) {

	socialEncounterId := fmt.Sprint(request.Id)

	var fromDb, err = handler.Repo.GetSocialEncounterById(socialEncounterId)
	if err != nil {
		return &encounters.SocialEncounterResponse{}, err
	}
	return SocialEncounterToRpc(fromDb), nil
}

func (handler *EncounterHandler) PostSocialEncounter(ctx context.Context, request *encounters.SocialEncounterResponse) (*encounters.SocialEncounterResponse, error) {
	socialEncounter := RpcToSocialEncounter(request)

	var fromDb, err = handler.Repo.InsertSocialEncounter(socialEncounter)
	if err != nil {
		return &encounters.SocialEncounterResponse{}, err
	}
	return SocialEncounterToRpc(fromDb), nil
}

func (handler *EncounterHandler) ActivateSocialEncounter(ctx context.Context, request *encounters.ActivateSocialEncounterRequest) (*encounters.EmptyResponse, error) {
	id := request.Id
	participantLocation := RpcToParticipantLocation(request.ParticipantLocation)

	err := handler.Repo.ActivateSocialEncounter(id, *participantLocation)
	if err != nil {
		return &encounters.EmptyResponse{}, err
	}
	return &encounters.EmptyResponse{}, nil
}

/*
func (p *EncounterHandler) MiddlewareContentTypeSet(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		p.logger.Println("Method [", h.Method, "] - Hit path :", h.URL.Path)

		rw.Header().Add("Content-Type", "application/json")

		next.ServeHTTP(rw, h)
	})
}

func (p *EncounterHandler) MiddlewareSocialEncounterDeserialization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		socialEncounter := &model.SocialEncounter{}
		err := socialEncounter.FromJSON(h.Body)
		if err != nil {
			http.Error(rw, "Unable to decode json", http.StatusBadRequest)
			p.logger.Fatal(err)
			return
		}

		ctx := context.WithValue(h.Context(), KeyProduct{}, socialEncounter)
		h = h.WithContext(ctx)

		next.ServeHTTP(rw, h)
	})
}*/

// modelToRpc
func SocialEncounterToRpc(socialEncounter *model.SocialEncounter) *encounters.SocialEncounterResponse {
	return &encounters.SocialEncounterResponse{
		Id:                   socialEncounter.Id,
		Name:                 socialEncounter.Name,
		Description:          socialEncounter.Description,
		Location:             LocationToRpc(socialEncounter.Location),
		Experience:           int64(socialEncounter.Experience),
		Status:               encounters.EncounterStatus(socialEncounter.Status),
		Type:                 encounters.EncounterType(socialEncounter.Type),
		Radius:               int64(socialEncounter.Radius),
		Participants:         ParticipantsToRpc(socialEncounter.Participants),
		Completers:           CompletersToRpc(socialEncounter.Completers),
		RequiredParticipants: int64(socialEncounter.RequiredParticipants),
		CurrentlyInRange:     ParticipantsToRpc(socialEncounter.CurrentlyInRange),
	}
}

func SocialEncountersToRpc(socialEncounters []model.SocialEncounter) *encounters.SocialEncountersResponse {
	result := make([]*encounters.SocialEncounterResponse, len(socialEncounters))
	for i, e := range socialEncounters {
		result[i] = SocialEncounterToRpc(&e)
	}
	return &encounters.SocialEncountersResponse{SocialEncounters: result}
}

func LocationToRpc(location model.Location) *encounters.Location {
	return &encounters.Location{
		Id:        location.Id,
		Latitude:  float32(location.Latitude),
		Longitude: float32(location.Longitude),
	}
}

func ParticipantToRpc(participant *model.Participant) *encounters.Participant {
	return &encounters.Participant{
		Id:          participant.Id,
		Username:    participant.Username,
		EncounterId: int64(participant.EncounterId),
	}
}

func ParticipantsToRpc(participants []model.Participant) []*encounters.Participant {
	result := make([]*encounters.Participant, len(participants))
	for i, e := range participants {
		result[i] = ParticipantToRpc(&e)
	}
	return result
}

func CompleterToRpc(completer *model.Completer) *encounters.Completer {
	return &encounters.Completer{
		Id:             completer.Id,
		Username:       completer.Username,
		CompletionDate: timestamppb.New(completer.CompletionDate),
		EncounterId:    int64(completer.EncounterId),
	}
}

func CompletersToRpc(completers []model.Completer) []*encounters.Completer {
	result := make([]*encounters.Completer, len(completers))
	for i, e := range completers {
		result[i] = CompleterToRpc(&e)
	}
	return result
}

//rpcToModel

func RpcToSocialEncounter(rpcSocialEncounter *encounters.SocialEncounterResponse) *model.SocialEncounter {
	return &model.SocialEncounter{
		Encounter: model.Encounter{
			Id:           rpcSocialEncounter.Id,
			Name:         rpcSocialEncounter.Name,
			Description:  rpcSocialEncounter.Description,
			Location:     *RpcToLocation(rpcSocialEncounter.Location),
			Experience:   rpcSocialEncounter.Experience,
			Status:       model.EncounterStatus(rpcSocialEncounter.Status),
			Type:         model.EncounterType(rpcSocialEncounter.Type),
			Radius:       rpcSocialEncounter.Radius,
			Participants: RpcToParticipants(rpcSocialEncounter.Participants),
			Completers:   RpcToCompleters(rpcSocialEncounter.Completers),
		},
		RequiredParticipants: int64(rpcSocialEncounter.RequiredParticipants),
		CurrentlyInRange:     RpcToParticipants(rpcSocialEncounter.CurrentlyInRange),
	}
}

func RpcsToSocialEncounters(rpcSocialEncounters *encounters.SocialEncountersResponse) []model.SocialEncounter {
	result := make([]model.SocialEncounter, len(rpcSocialEncounters.SocialEncounters))
	for i, e := range rpcSocialEncounters.SocialEncounters {
		result[i] = *RpcToSocialEncounter(e)
	}
	return result
}

func RpcToLocation(rpcLocation *encounters.Location) *model.Location {
	return &model.Location{
		Id:        rpcLocation.Id,
		Latitude:  float64(rpcLocation.Latitude),
		Longitude: float64(rpcLocation.Longitude),
	}
}

func RpcToParticipant(rpcParticipant *encounters.Participant) *model.Participant {
	return &model.Participant{
		Id:          rpcParticipant.Id,
		Username:    rpcParticipant.Username,
		EncounterId: int64(rpcParticipant.EncounterId),
	}
}

func RpcToParticipants(rpcParticipants []*encounters.Participant) []model.Participant {
	result := make([]model.Participant, len(rpcParticipants))
	for i, e := range rpcParticipants {
		result[i] = *RpcToParticipant(e)
	}
	return result
}

func RpcToCompleter(rpcCompleter *encounters.Completer) *model.Completer {
	return &model.Completer{
		Id:             rpcCompleter.Id,
		Username:       rpcCompleter.Username,
		CompletionDate: rpcCompleter.CompletionDate.AsTime(),
		EncounterId:    int64(rpcCompleter.EncounterId),
	}
}

func RpcToCompleters(rpcCompleters []*encounters.Completer) []model.Completer {
	result := make([]model.Completer, len(rpcCompleters))
	for i, e := range rpcCompleters {
		result[i] = *RpcToCompleter(e)
	}
	return result
}

func RpcToParticipantLocation(rpcParticipantLocation *encounters.ParticipantLocation) *model.ParticipantLocation {
	return &model.ParticipantLocation{
		Username:  rpcParticipantLocation.Username,
		Latitude:  float64(rpcParticipantLocation.Latitude),
		Longitude: float64(rpcParticipantLocation.Longitude),
	}
}
