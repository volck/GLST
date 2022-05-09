package main

func (s *server) Routes() {
	s.PublicRoutes()

}

func (s *server) PublicRoutes() {
	s.router.GET("/ping", s.HandlePing())
	s.router.GET("/tokens", s.HandleGetTokens())
	s.router.GET("/usedtokens", s.HandleUsedTokens())
	s.router.GET("/newtoken", s.HandleGetFreshToken())
	s.router.POST("/retire", s.HandleRetireToken())
	s.router.GET("/insert", s.HandleInsertDB())

}
