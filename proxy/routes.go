package main

func (s *server) routes() {
	s.router.Get("/{dongle_id}", s.AuthenticateCookie(s.websocketHandler()))
	s.router.Post("/{dongle_id}", s.restHandler())
}
