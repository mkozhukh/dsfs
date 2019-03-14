package main

import (
	"github.com/mkozhukh/go-remote"
)

func initAPI(s *remote.Server) {
	//s.RegisterWithGuard("admin", &SnippetAdminAPI{}, access.CheckRequest(CanAdminUser))
}
