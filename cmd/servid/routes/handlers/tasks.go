package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"tasks17-server/internal/auth"
	"tasks17-server/internal/platform"
	"tasks17-server/internal/tasks"
)

func AddTeamHandler(ts tasks.TaskStorage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		type createTeam struct {
			Name string
		}

		var ct createTeam

		credentials, _ := auth.FromRequest(r)
		decoder := json.NewDecoder(r.Body)
		_ = decoder.Decode(&ct)

		id, _ := ts.SaveTeam(tasks.Team{
			OrgId: credentials.Oid,
			Name:  ct.Name,
		})

		jsonResponse(map[string]int{"id": id}, w)
	}
}

func TeamInfoHandler(ts tasks.TaskStorage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		teamId, _ := strconv.Atoi(mux.Vars(r)["team"])
		team := ts.LoadTeam(teamId)

		jsonResponse(team, w)
	}
}

func TeamListHandler(ts tasks.TaskStorage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		credentials, _ := auth.FromRequest(r)

		jsonResponse(ts.LoadTeams(credentials.Oid), w)
	}
}

func TaskListHandler(ts tasks.TaskStorage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		teamId, _ := strconv.Atoi(mux.Vars(r)["team"])

		jsonCollectionResponse("tasks", ts.LoadTasks(teamId), w)
	}
}

func TaskStateListHandler(ts tasks.TaskStorage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		credentials, _ := auth.FromRequest(r)
		teamId, _ := strconv.Atoi(mux.Vars(r)["team"])
		tList := tasks.LoadTaskStateList(ts, credentials.Uid, teamId)

		jsonCollectionResponse("tasks", tList, w)
	}
}

func TaskCreateHandler(ts tasks.TaskStorage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var t tasks.Task
		credentials, _ := auth.FromRequest(r)

		decoder := json.NewDecoder(r.Body)
		_ = decoder.Decode(&t)

		t.AuthorId = credentials.Uid
		t.TeamId, _ = strconv.Atoi(mux.Vars(r)["team"])

		id, err := tasks.CreateTask(ts, &t)

		if err != nil {
			log.Printf("Error while task creation: %v", err)
			http.Error(w, "Something went wrong", 500)
			return
		}

		_, _ = fmt.Fprintf(w, "{\"id\": \"%s\"}", id)
	}
}

func SubTaskCloseHandler(ts tasks.TaskStorage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		subTaskIdString := mux.Vars(r)["subTask"]

		stId, _ := strconv.Atoi(subTaskIdString)

		tasks.CompleteSubTask(ts, stId)
	}
}

func SubTaskCreateHandler(h *platform.Hub, ts tasks.TaskStorage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var st tasks.SubTask
		decoder.Decode(&st)

		taskId := mux.Vars(r)["task"]
		credentials, _ := auth.FromRequest(r)

		st.AuthorId = credentials.Uid
		st.TaskId = taskId

		id, err := tasks.AddSubTask(h, ts, st)

		if err != nil {
			http.Error(w, "Something went wrong", 500)
			return
		}

		fmt.Fprintf(w, "{\"id\": %d}", id)
	}
}

func StagesLoadHandler(ts tasks.TaskStorage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		teamId, err := strconv.Atoi(vars["team"])

		if err != nil {
			http.Error(w, "Bad team in path", http.StatusBadRequest)
		}

		jsonCollectionResponse("stages", ts.LoadStages(teamId), w)
	}
}

func TaskLoadHandler(ts tasks.TaskStorage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		type response struct {
			Task     tasks.Task         `json:"task"`
			Comments []tasks.Comment    `json:"comments"`
			Progress tasks.FullProgress `json:"progress"`
		}

		var rs response

		taskId := mux.Vars(r)["task"]
		rs.Task = ts.LoadTask(taskId)
		rs.Comments = ts.LoadComments(rs.Task.DiscussionId)
		rs.Progress = tasks.LoadProgress(ts, taskId)

		jsonResponse(rs, w)
	}
}

func TaskUpdateDescription(ts tasks.TaskStorage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		type newDescription struct {
			Description string `json:"description"`
		}

		var nd newDescription
		d := json.NewDecoder(r.Body)
		d.Decode(&nd)

		taskId := mux.Vars(r)["task"]

		err := ts.UpdateDescription(taskId, nd.Description)

		if err != nil {
			log.Println(err)
		}
	}
}

func TaskClose(ts tasks.TaskStorage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		taskId := mux.Vars(r)["task"]

		err := ts.CloseTask(taskId)

		if err != nil {
			log.Println(err)
		}
	}
}
