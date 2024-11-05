package handlers

func (s *Server) InitRoutes() {
	s.Router.HandleFunc("/", s.Ping).Methods("GET")
	v1Router := s.Router.PathPrefix("/api/v1").Subrouter()

	threadRouter := v1Router.PathPrefix("/thread").Subrouter()
	threadRouter.HandleFunc("", s.ListThreads).Methods("GET")
	threadRouter.HandleFunc("", s.CreateThread).Methods("POST")
	threadRouter.HandleFunc("/{id}", s.GetThread).Methods("GET")
	threadRouter.HandleFunc("/{id}", s.UpdateThread).Methods("PUT")
	threadRouter.HandleFunc("/{id}", s.DeleteThread).Methods("DELETE")

	messageRouter := v1Router.PathPrefix("/message").Subrouter()
	messageRouter.HandleFunc("/{id}", s.GetMessage).Methods("GET")
	messageRouter.HandleFunc("/{id}", s.UpdateMessage).Methods("PUT")
	messageRouter.HandleFunc("/{id}", s.DeleteMessage).Methods("DELETE")

	messageThreadIDRouter := messageRouter.PathPrefix("/thread/{thread_id}").Subrouter()
	messageThreadIDRouter.HandleFunc("", s.CreateMessage).Methods("POST")
	messageThreadIDRouter.HandleFunc("", s.ListMessages).Methods("GET")
}
