package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"tasks17-server/internal/auth"
	"tasks17-server/internal/tasks"
)

func ProjectAddStageHandler(ts tasks.TaskStorage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		type createProjectStage struct {
			Name        string
			Description string
		}

		var cps createProjectStage
		decoder := json.NewDecoder(r.Body)
		_ = decoder.Decode(&cps)

		projectId, _ := strconv.Atoi(mux.Vars(r)["project"])

		id := ts.CreateProjectStage(tasks.ProjectStage{
			ProjectId:   projectId,
			Name:        cps.Name,
			Description: cps.Description,
			Status:      0,
		})

		jsonResponse(map[string]int{"id": id}, w)
	}
}

func ProjectInfoHandler(ts tasks.TaskStorage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		type response struct {
			Project  tasks.Project        `json:"project"`
			Comments []tasks.Comment      `json:"comments"`
			Stages   []tasks.ProjectStage `json:"stages"`
			Tasks    []tasks.Task         `json:"tasks"`
		}
		var rs response
		pid, _ := strconv.Atoi(mux.Vars(r)["project"])
		rs.Project = ts.LoadProject(pid)
		rs.Comments = ts.LoadComments(rs.Project.DiscussionId)
		rs.Stages = ts.LoadProjectStages(pid)

		var stageIds []int
		for _, v := range rs.Stages {
			stageIds = append(stageIds, v.Id)
		}

		rs.Tasks = ts.LoadTasksByProjectStages(stageIds)

		jsonResponse(rs, w)
	}
}

func ProjectListHandler(ts tasks.TaskStorage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		credentials, _ := auth.FromRequest(r)

		jsonResponse(ts.LoadProjects(credentials.Oid), w)
	}
}

func AddProjectHandler(ts tasks.TaskStorage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		type createProject struct {
			Name        string
			Description string
		}

		var cp createProject
		credentials, _ := auth.FromRequest(r)
		decoder := json.NewDecoder(r.Body)
		_ = decoder.Decode(&cp)

		discussionId := uuid.New().String()
		ts.CreateDiscussion(discussionId)

		projectId := ts.CreateProject(tasks.Project{
			OrgId:        credentials.Oid,
			Name:         cp.Name,
			Description:  cp.Description,
			Status:       0,
			DiscussionId: discussionId,
		})

		jsonResponse(map[string]int{"id": projectId}, w)
	}
}

func AddTaskForProjectStageHandler(ts tasks.TaskStorage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var t tasks.Task
		credentials, _ := auth.FromRequest(r)

		decoder := json.NewDecoder(r.Body)
		_ = decoder.Decode(&t)

		t.AuthorId = credentials.Uid
		t.ProjectStageId, _ = strconv.Atoi(mux.Vars(r)["stage"])

		id, err := tasks.CreateTask(ts, &t)

		if err != nil {
			log.Printf("Error while task creation: %v", err)
			http.Error(w, "Something went wrong", 500)
			return
		}

		_, _ = fmt.Fprintf(w, "{\"id\": \"%s\"}", id)
	}
}
