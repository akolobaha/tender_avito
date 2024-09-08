package app

import (
	"tenders/internal/handler"
)

func Run() {

	//ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	handler.NewRouter()

}
