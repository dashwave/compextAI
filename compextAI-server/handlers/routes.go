package handlers

import "github.com/burnerlee/compextAI/middlewares"

func (s *Server) InitRoutes() {
	s.Router.HandleFunc("/", s.Ping).Methods("GET")
	v1Router := s.Router.PathPrefix("/api/v1").Subrouter()

	threadRouter := v1Router.PathPrefix("/thread").Subrouter()

	threadRouter.HandleFunc("", middlewares.AuthMiddleware(s.ListThreads, s.DB)).Methods("GET")
	threadRouter.HandleFunc("", middlewares.AuthMiddleware(s.CreateThread, s.DB)).Methods("POST")
	threadRouter.HandleFunc("/{id}", middlewares.AuthMiddleware(s.GetThread, s.DB)).Methods("GET")
	threadRouter.HandleFunc("/{id}", middlewares.AuthMiddleware(s.UpdateThread, s.DB)).Methods("PUT")
	threadRouter.HandleFunc("/{id}", middlewares.AuthMiddleware(s.DeleteThread, s.DB)).Methods("DELETE")
	threadRouter.HandleFunc("/{id}/execute", middlewares.AuthMiddleware(s.ExecuteThread, s.DB)).Methods("POST")

	threadExecRouter := v1Router.PathPrefix("/threadexec").Subrouter()
	threadExecRouter.HandleFunc("/{id}/status", middlewares.AuthMiddleware(s.GetThreadExecutionStatus, s.DB)).Methods("GET")
	threadExecRouter.HandleFunc("/{id}/response", middlewares.AuthMiddleware(s.GetThreadExecutionResponse, s.DB)).Methods("GET")

	messageRouter := v1Router.PathPrefix("/message").Subrouter()

	messageRouter.HandleFunc("/{id}", middlewares.AuthMiddleware(s.GetMessage, s.DB)).Methods("GET")
	messageRouter.HandleFunc("/{id}", middlewares.AuthMiddleware(s.UpdateMessage, s.DB)).Methods("PUT")
	messageRouter.HandleFunc("/{id}", middlewares.AuthMiddleware(s.DeleteMessage, s.DB)).Methods("DELETE")

	messageThreadIDRouter := messageRouter.PathPrefix("/thread/{thread_id}").Subrouter()

	messageThreadIDRouter.HandleFunc("", middlewares.AuthMiddleware(s.CreateMessage, s.DB)).Methods("POST")
	messageThreadIDRouter.HandleFunc("", middlewares.AuthMiddleware(s.ListMessages, s.DB)).Methods("GET")

	userRouter := v1Router.PathPrefix("/user").Subrouter()
	userRouter.HandleFunc("/signup", s.CreateUser).Methods("POST")
	userRouter.HandleFunc("/login", s.Login).Methods("POST")

	threadExecutionParamsRouter := v1Router.PathPrefix("/execparams").Subrouter()
	threadExecutionParamsRouter.HandleFunc("/fetchall", middlewares.AuthMiddleware(s.ListThreadExecutionParams, s.DB)).Methods("GET")
	threadExecutionParamsRouter.HandleFunc("/create", middlewares.AuthMiddleware(s.CreateThreadExecutionParams, s.DB)).Methods("POST")
	threadExecutionParamsRouter.HandleFunc("/fetch", middlewares.AuthMiddleware(s.GetThreadExecutionParamsByNameAndEnv, s.DB)).Methods("POST")
	threadExecutionParamsRouter.HandleFunc("/delete", middlewares.AuthMiddleware(s.DeleteThreadExecutionParams, s.DB)).Methods("DELETE")

	threadExecutionParamsTemplateRouter := v1Router.PathPrefix("/execparamstemplate").Subrouter()
	threadExecutionParamsTemplateRouter.HandleFunc("", middlewares.AuthMiddleware(s.ListThreadExecutionParamsTemplates, s.DB)).Methods("GET")
	threadExecutionParamsTemplateRouter.HandleFunc("", middlewares.AuthMiddleware(s.CreateThreadExecutionParamsTemplate, s.DB)).Methods("POST")
	threadExecutionParamsTemplateRouter.HandleFunc("/{id}", middlewares.AuthMiddleware(s.GetThreadExecutionParamsTemplateByID, s.DB)).Methods("GET")
	threadExecutionParamsTemplateRouter.HandleFunc("/{id}", middlewares.AuthMiddleware(s.DeleteThreadExecutionParamsTemplate, s.DB)).Methods("DELETE")
	threadExecutionParamsTemplateRouter.HandleFunc("/{id}", middlewares.AuthMiddleware(s.UpdateThreadExecutionParamsTemplate, s.DB)).Methods("PUT")

}
