package main

func (s *server) routes() {
	s.router.GET("/", s.appIndex)

	s.router.GET("/signup", s.signupGet)
	s.router.POST("/signup", s.signupPost)
	s.router.POST("/signin", s.signinPost)
	s.router.POST("/signout", s.signoutPost)
}
