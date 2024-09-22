package main

import (
	"omsms/cmd"
	"omsms/db"
)

func main() {
	db.InitDB()
	db.RegisterModels(db.DB)
	cmd.Execute()
}
