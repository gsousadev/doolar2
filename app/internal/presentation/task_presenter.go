package presentation

import "net/http"

type TaskPresenter struct {
	Presenter
}

func (tp *TaskPresenter) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Task handler is working"))
}
