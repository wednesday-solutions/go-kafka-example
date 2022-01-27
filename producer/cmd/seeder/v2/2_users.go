package main

import (
	"context"
	"fmt"

	"github.com/volatiletech/sqlboiler/queries/qm"
	seeders "github.com/wednesday-solutions/go-template-producer/cmd/seeder"
	"github.com/wednesday-solutions/go-template-producer/internal/postgres"
	"github.com/wednesday-solutions/go-template-producer/models"
	"github.com/wednesday-solutions/go-template-producer/pkg/utl/secure"
)

func main() {

	sec := secure.New(1, nil)
	db, _ := postgres.Connect()
	// getting the latest location company and role id so that we can seed a new user

	role, _ := models.Roles(qm.OrderBy("id DESC")).One(context.Background(), db)
	var insertQuery = fmt.Sprintf("INSERT INTO public.users (first_name, last_name, username, password, "+
		"email, active, role_id) VALUES ('Admin', 'Admin', 'admin', '%s', 'johndoe@mail.com', true, %d);",
		sec.Hash("adminuser"), role.ID)
	_ = seeders.SeedData("users", insertQuery)
}
