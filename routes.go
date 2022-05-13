package main

func (s *server) Routes() {
	s.PublicRoutes()

}

func (s *server) PublicRoutes() {
	s.router.GET("/ping", s.HandlePing())
	s.router.GET("/new", s.HandleNewToken())
}
