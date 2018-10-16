package flotilla

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"unicode"

	"github.com/gorilla/mux"
	"github.com/stitchfix/flotilla-os/exceptions"
	"github.com/stitchfix/flotilla-os/services"
	"github.com/stitchfix/flotilla-os/state"
)

type endpoints struct {
	executionService  services.ExecutionService
	definitionService services.DefinitionService
	logService        services.LogService
}

type listRequest struct {
	limit      int
	offset     int
	sortBy     string
	order      string
	filters    map[string][]string
	envFilters map[string]string
	tagFilters map[string]string
}

type launchRequest struct {
	ClusterName string         `json:"cluster"`
	Env         *state.EnvList `json:"env"`
}

type launchRequestV2 struct {
	RunTags RunTags `json:"run_tags"`
	*launchRequest
}

type launchGenericReq struct {
	ClusterName     *string `json:"cluster"`
	RunTags         RunTags `json:"run_tags"`
	TaskID          *string `json:"task_id"`
	Command         *string
	CommandEncoding *string `json:"cmd_enconding"`
	Memory          *int64
	Image           *string
	Env             *state.EnvList
	UserTags        *state.UserTagMap `json:"user_tags"`
}

func (lr launchGenericReq) Validate() error {
	if lr.ClusterName == nil || lr.TaskID == nil || lr.Command == nil ||
		lr.Memory == nil || lr.Image == nil {
		msg := "Required fields: cluster, run_tags, task_id, command, memory, image"
		return exceptions.MalformedInput{msg}
	}

	return nil

}

//
// RunTags represents which user is responsible for a task run
//
type RunTags struct {
	OwnerEmail string `json:"owner_email"`
	TeamName   string `json:"team_name"`
	OwnerID    string `json:"owner_id"`
}

func (ep *endpoints) getURLParam(v url.Values, key string, defaultValue string) string {
	val, ok := v[key]
	if ok && len(val) > 0 {
		return val[0]
	}
	return defaultValue
}

func (ep *endpoints) getFilters(params url.Values, nonFilters map[string]bool) (map[string][]string, map[string]string, map[string]string) {
	filters := make(map[string][]string)
	envFilters := make(map[string]string)
	tagFilters := make(map[string]string)
	for k, v := range params {
		if !nonFilters[k] && len(v) > 0 {
			// Env filters have the "env" key and are "|" separated key-value pairs
			//
			// eg. env=FOO|BAR&env=CUPCAKE|SPRINKLES
			//
			if k == "env" {
				for _, kv := range v {
					split := strings.Split(kv, "|")
					if len(split) == 2 {
						envFilters[split[0]] = split[1]
					}
				}
			} else if k == "tag" {
				// Tag filters refer to user tags; syntax is same as env filters.
				for _, kv := range v {
					split := strings.Split(kv, "|")
					if len(split) == 2 {
						tagFilters[split[0]] = split[1]
					}
				}
			} else {
				filters[k] = v
			}
		}
	}
	return filters, envFilters, tagFilters
}

func (ep *endpoints) decodeListRequest(r *http.Request) listRequest {
	var lr listRequest
	params := r.URL.Query()

	lr.limit, _ = strconv.Atoi(ep.getURLParam(params, "limit", "1024"))
	lr.offset, _ = strconv.Atoi(ep.getURLParam(params, "offset", "0"))
	lr.sortBy = ep.getURLParam(params, "sort_by", "group_name")
	lr.order = ep.getURLParam(params, "order", "asc")
	lr.filters, lr.envFilters, lr.tagFilters = ep.getFilters(params, map[string]bool{
		"limit":   true,
		"offset":  true,
		"sort_by": true,
		"order":   true,
	})
	return lr
}

func (ep *endpoints) decodeRequest(r *http.Request, entity interface{}) error {
	return json.NewDecoder(r.Body).Decode(entity)
}

func (ep endpoints) encodeError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err.(type) {
	case exceptions.MalformedInput:
		w.WriteHeader(http.StatusBadRequest)
	case exceptions.ConflictingResource:
		w.WriteHeader(http.StatusConflict)
	case exceptions.MissingResource:
		w.WriteHeader(http.StatusNotFound)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func (ep *endpoints) encodeResponse(w http.ResponseWriter, response interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(response)
}

func (ep *endpoints) ListDefinitions(w http.ResponseWriter, r *http.Request) {
	lr := ep.decodeListRequest(r)

	definitionList, err := ep.definitionService.List(
		lr.limit, lr.offset, lr.sortBy, lr.order, lr.filters, lr.envFilters)
	if definitionList.Definitions == nil {
		definitionList.Definitions = []state.Definition{}
	}
	if err != nil {
		ep.encodeError(w, err)
	} else {
		response := make(map[string]interface{})
		response["total"] = definitionList.Total
		response["definitions"] = definitionList.Definitions
		response["limit"] = lr.limit
		response["offset"] = lr.offset
		response["sort_by"] = lr.sortBy
		response["order"] = lr.order
		response["env_filters"] = lr.envFilters
		for k, v := range lr.filters {
			response[k] = v
		}
		ep.encodeResponse(w, response)
	}
}

func (ep *endpoints) GetDefinition(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	definition, err := ep.definitionService.Get(vars["definition_id"])
	if err != nil {
		ep.encodeError(w, err)
	} else {
		ep.encodeResponse(w, definition)
	}
}

func (ep *endpoints) GetDefinitionByAlias(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	definition, err := ep.definitionService.GetByAlias(vars["alias"])
	if err != nil {
		ep.encodeError(w, err)
	} else {
		ep.encodeResponse(w, definition)
	}
}

func (ep *endpoints) CreateDefinition(w http.ResponseWriter, r *http.Request) {
	var definition state.Definition
	err := ep.decodeRequest(r, &definition)
	if err != nil {
		ep.encodeError(w, err)
		return
	}

	created, err := ep.definitionService.Create(&definition)
	if err != nil {
		ep.encodeError(w, err)
	} else {
		ep.encodeResponse(w, created)
	}
}

func (ep *endpoints) UpdateDefinition(w http.ResponseWriter, r *http.Request) {
	var definition state.Definition
	err := ep.decodeRequest(r, &definition)
	if err != nil {
		ep.encodeError(w, err)
		return
	}

	vars := mux.Vars(r)
	updated, err := ep.definitionService.Update(vars["definition_id"], definition)

	if err != nil {
		ep.encodeError(w, err)
	} else {
		ep.encodeResponse(w, updated)
	}
}

func (ep *endpoints) DeleteDefinition(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	err := ep.definitionService.Delete(vars["definition_id"])
	if err != nil {
		ep.encodeError(w, err)
	} else {
		ep.encodeResponse(w, map[string]bool{"deleted": true})
	}
}

func (ep *endpoints) ListRuns(w http.ResponseWriter, r *http.Request) {
	lr := ep.decodeListRequest(r)

	vars := mux.Vars(r)
	definitionID, ok := vars["definition_id"]
	if ok {
		lr.filters["definition_id"] = []string{definitionID}
	}

	runList, err := ep.executionService.List(
		lr.limit, lr.offset, lr.order, lr.sortBy, lr.filters, lr.envFilters, lr.tagFilters)
	if err != nil {
		ep.encodeError(w, err)
	} else {
		response := make(map[string]interface{})
		response["total"] = runList.Total
		response["history"] = runList.Runs
		response["limit"] = lr.limit
		response["offset"] = lr.offset
		response["sort_by"] = lr.sortBy
		response["order"] = lr.order
		response["env_filters"] = lr.envFilters
		for k, v := range lr.filters {
			response[k] = v
		}
		ep.encodeResponse(w, response)
	}
}

func (ep *endpoints) GetRun(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	run, err := ep.executionService.Get(vars["run_id"])
	if err != nil {
		ep.encodeError(w, err)
	} else {
		ep.encodeResponse(w, run)
	}
}

func (ep *endpoints) CreateRun(w http.ResponseWriter, r *http.Request) {
	var lr launchRequest
	err := ep.decodeRequest(r, &lr)
	if err != nil {
		ep.encodeError(w, err)
		return
	}

	vars := mux.Vars(r)
	run, err := ep.executionService.Create(vars["definition_id"], lr.ClusterName, lr.Env, "v1-unknown")
	if err != nil {
		ep.encodeError(w, err)
	} else {
		ep.encodeResponse(w, run)
	}
}

func (ep *endpoints) CreateRunV2(w http.ResponseWriter, r *http.Request) {
	var lr launchRequestV2
	err := ep.decodeRequest(r, &lr)
	if err != nil {
		ep.encodeError(w, err)
		return
	}

	if len(lr.RunTags.OwnerEmail) == 0 || len(lr.RunTags.TeamName) == 0 {
		ep.encodeError(w, exceptions.MalformedInput{
			ErrorString: fmt.Sprintf("run_tags must exist in body and contain [owner_email] and [team_name]")})
		return
	}

	vars := mux.Vars(r)
	run, err := ep.executionService.Create(vars["definition_id"], lr.ClusterName, lr.Env, lr.RunTags.OwnerEmail)
	if err != nil {
		ep.encodeError(w, err)
	} else {
		ep.encodeResponse(w, run)
	}
}

func (ep *endpoints) CreateRunV4(w http.ResponseWriter, r *http.Request) {
	var lr launchRequestV2
	err := ep.decodeRequest(r, &lr)
	if err != nil {
		ep.encodeError(w, err)
		return
	}

	if len(lr.RunTags.OwnerID) == 0 {
		ep.encodeError(w, exceptions.MalformedInput{
			ErrorString: fmt.Sprintf("run_tags must exist in body and contain [owner_id]")})
		return
	}

	vars := mux.Vars(r)
	run, err := ep.executionService.Create(vars["definition_id"], lr.ClusterName, lr.Env, lr.RunTags.OwnerID)
	if err != nil {
		ep.encodeError(w, err)
	} else {
		ep.encodeResponse(w, run)
	}
}

func (ep *endpoints) decodeGenericReq(r *http.Request) (launchGenericReq, error) {
	var lr launchGenericReq
	err := ep.decodeRequest(r, &lr)
	if err != nil {
		return lr, exceptions.MalformedInput{err.Error()}
	}

	if len(lr.RunTags.OwnerID) == 0 {
		return lr, exceptions.MalformedInput{
			ErrorString: fmt.Sprintf("run_tags must exist in body and contain [owner_id]")}
	}

	if lr.Command == nil {
		return lr, exceptions.MalformedInput{"[command] is required"}
	}
	if lr.CommandEncoding != nil && *lr.CommandEncoding == "base64" {
		decoded, err := base64.StdEncoding.DecodeString(*lr.Command)
		if err != nil {
			return lr, err
		}
		decodedS := string(decoded)
		lr.Command = &decodedS
	}

	lr.Command, err = ep.wrapRunTimeCmd(*lr.Command)
	if err != nil {
		return lr, err
	}

	return lr, nil
}

func (ep *endpoints) wrapRunTimeCmd(cmd string) (*string, error) {
	var result bytes.Buffer
	if err := state.CommandTemplate.Execute(&result, map[string]string{
		"Command": cmd}); err != nil {
		return nil, err
	}
	resultS := result.String()
	return &resultS, nil
}

//
// Returns generic task def corresponding to image. If one does not exist,
// creates it and returns the newly created image.
//
func (ep *endpoints) findOrCreateGenericDef(image string) (state.Definition, error) {
	filter := map[string][]string{
		"group_name": []string{state.GenericDefPrefix},
		"image":      []string{image}}

	defList, err := ep.definitionService.List(100, 0, "image", "desc",
		filter, nil)
	if err != nil {
		return state.Definition{}, err
	}

	if defList.Total > 1 {
		// Inconsistent db state
		return state.Definition{}, fmt.Errorf("More than one generic task def found.")
	}

	if defList.Total == 1 {
		return defList.Definitions[0], nil
	}

	return ep.createGenericDef(image)
}

func (ep *endpoints) groupNameForImage(image string) string {
	splits := strings.Split(image, "/")
	var imgTag string
	if len(splits) == 2 {
		imgTag = splits[1]
	} else {
		imgTag = image
	}

	mapFunc := func(r rune) rune {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_':
			return r
		default:
			return '_'
		}
	}

	suffix := strings.Map(mapFunc, imgTag)

	if len(state.GenericDefPrefix)+len(suffix) > 255 {
		toRemove := len(state.GenericDefPrefix) + len(suffix) - 255
		return fmt.Sprintf("%s%s", state.GenericDefPrefix, suffix[toRemove:])
	}
	return fmt.Sprintf("%s%s", state.GenericDefPrefix, suffix)
}

//
// Create a generic task def for a particular image
//
func (ep *endpoints) createGenericDef(image string) (state.Definition, error) {
	groupName := ep.groupNameForImage(image)
	cmd := ` echo "ERROR This is a generic task def. It should be run with run time params, see docs."
			 exit 1
		   `
	mem := int64(100)

	def := state.Definition{
		Alias:     groupName,
		GroupName: groupName,
		Image:     image,
		Memory:    &mem,
		Command:   cmd,
		TaskType:  state.TaskTypeGeneric}

	return ep.definitionService.Create(&def)
}

func (ep *endpoints) CreateGenericRun(w http.ResponseWriter, r *http.Request) {
	// decode
	lr, err := ep.decodeGenericReq(r)
	if err != nil {
		ep.encodeError(w, err)
		return
	}

	err = lr.Validate()
	if err != nil {
		ep.encodeError(w, err)
		return
	}

	// find generic def; create if not exists
	genericDef, err := ep.findOrCreateGenericDef(*lr.Image)
	if err != nil {
		ep.encodeError(w, err)
		return
	}

	// create run
	run, err := ep.executionService.CreateFromDefinition(genericDef, *lr.ClusterName, lr.Env, lr.RunTags.OwnerID, lr.TaskID, lr.Command, lr.Memory, lr.UserTags)
	if err != nil {
		ep.encodeError(w, err)
	} else {
		ep.encodeResponse(w, run)
	}
}

func (ep *endpoints) CreateRunByAlias(w http.ResponseWriter, r *http.Request) {
	var lr launchRequestV2
	err := ep.decodeRequest(r, &lr)
	if err != nil {
		ep.encodeError(w, err)
		return
	}

	if len(lr.RunTags.OwnerID) == 0 {
		ep.encodeError(w, exceptions.MalformedInput{
			ErrorString: fmt.Sprintf("run_tags must exist in body and contain [owner_id]")})
		return
	}

	vars := mux.Vars(r)
	run, err := ep.executionService.CreateByAlias(vars["alias"], lr.ClusterName, lr.Env, lr.RunTags.OwnerID)
	if err != nil {
		ep.encodeError(w, err)
	} else {
		ep.encodeResponse(w, run)
	}
}

func (ep *endpoints) StopRun(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	err := ep.executionService.Terminate(vars["run_id"])
	if err != nil {
		ep.encodeError(w, err)
	} else {
		ep.encodeResponse(w, map[string]bool{"terminated": true})
	}
}

func (ep *endpoints) UpdateRun(w http.ResponseWriter, r *http.Request) {
	var run state.Run
	err := ep.decodeRequest(r, &run)
	if err != nil {
		ep.encodeError(w, err)
		return
	}

	vars := mux.Vars(r)
	err = ep.executionService.UpdateStatus(vars["run_id"], run.Status, run.ExitCode)
	if err != nil {
		ep.encodeError(w, err)
	} else {
		ep.encodeResponse(w, map[string]bool{"updated": true})
	}
}

func (ep *endpoints) GetLogs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	params := r.URL.Query()

	lastSeen := ep.getURLParam(params, "last_seen", "")
	logs, newLastSeen, err := ep.logService.Logs(vars["run_id"], &lastSeen)
	if err != nil {
		ep.encodeError(w, err)
		return
	}

	res := map[string]string{
		"log": logs,
	}
	if newLastSeen != nil {
		res["last_seen"] = *newLastSeen
	}
	ep.encodeResponse(w, res)
}

func (ep *endpoints) GetGroups(w http.ResponseWriter, r *http.Request) {
	lr := ep.decodeListRequest(r)

	var name string
	if len(lr.filters["name"]) > 0 {
		name = lr.filters["name"][0]
	}

	groups, err := ep.definitionService.ListGroups(lr.limit, lr.offset, &name)
	if err != nil {
		ep.encodeError(w, err)
	} else {
		response := make(map[string]interface{})
		response["total"] = groups.Total
		response["groups"] = groups.Groups
		response["limit"] = lr.limit
		response["offset"] = lr.offset
		ep.encodeResponse(w, response)
	}
}

func (ep *endpoints) GetTags(w http.ResponseWriter, r *http.Request) {
	lr := ep.decodeListRequest(r)

	var name string
	if len(lr.filters["name"]) > 0 {
		name = lr.filters["name"][0]
	}

	tags, err := ep.definitionService.ListTags(lr.limit, lr.offset, &name)
	if err != nil {
		ep.encodeError(w, err)
	} else {
		response := make(map[string]interface{})
		response["total"] = tags.Total
		response["tags"] = tags.Tags
		response["limit"] = lr.limit
		response["offset"] = lr.offset
		ep.encodeResponse(w, response)
	}
}

func (ep *endpoints) ListClusters(w http.ResponseWriter, r *http.Request) {
	clusters, err := ep.executionService.ListClusters()
	if err != nil {
		ep.encodeError(w, err)
	} else {
		response := make(map[string]interface{})
		response["clusters"] = clusters
		ep.encodeResponse(w, response)
	}
}
