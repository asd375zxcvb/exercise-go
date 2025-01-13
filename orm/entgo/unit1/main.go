package main

import (
	"context"
	"entgo/ent"
	"entgo/ent/car"
	"entgo/ent/group"
	"entgo/ent/user"
	"github.com/cockroachdb/errors"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	// _ "github.com/mattn/go-sqlite3"
	_ "github.com/xiaoqidun/entps"
	"log"
	"time"
)

func main() {
	//client := PanicErr(ent.Open("postgres", "host=127.0.0.1 port=5432 user=postgres dbname=postgres password=123456 sslmode=disable", ent.Debug()))
	//client := PanicErr(ent.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/entgo?parseTime=True", ent.Debug()))
	//client := PanicErr(ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1", ent.Debug()))
	client := PanicErr(ent.Open("sqlite3", "file:./ent.db?_fk=1", ent.Debug()))
	defer func(client *ent.Client) {
		err := client.Close()
		if err != nil {
			log.Fatalf("failed closing the client: %v", err)
		} else {
			log.Println("closed the client")
		}
	}(client)
	ctx := context.Background()
	PanicErrExec(client.Schema.Create(ctx))
	PanicErrExec(round1(ctx, client))
	PanicErrExec(round2(ctx, client))
	PanicErrExec(round3(ctx, client))
}

func round1(ctx context.Context, client *ent.Client) error {
	log.Println("====================  round1  start =======================")
	PanicErr(CreateUser(ctx, client))
	PanicErr(QueryUser(ctx, client))
	PanicErr(client.User.Delete().Exec(ctx))
	log.Println("====================  round1  end   =======================")
	return nil
}

func round2(ctx context.Context, client *ent.Client) error {
	log.Println("====================  round2  start =======================")
	a8m := PanicErr(CreateCars(ctx, client))
	PanicErrExec(QueryCars(ctx, a8m))
	PanicErr(client.User.Delete().Exec(ctx))
	PanicErr(client.Car.Delete().Exec(ctx))
	log.Println("====================  round2  end   =======================")
	return nil
}

func round3(ctx context.Context, client *ent.Client) error {
	log.Println("====================  round3  start =======================")
	PanicErrExec(CreateGraph(ctx, client))
	PanicErrExec(QueryGithub(ctx, client))
	PanicErrExec(QueryArielCars(ctx, client))
	PanicErrExec(QueryGroupWithUsers(ctx, client))
	PanicErr(client.User.Delete().Exec(ctx))
	PanicErr(client.Car.Delete().Exec(ctx))
	PanicErr(client.Group.Delete().Exec(ctx))
	log.Println("====================  round3  end   =======================")
	return nil
}

func CreateUser(ctx context.Context, client *ent.Client) (*ent.User, error) {
	u, err := client.User.
		Create().
		SetAge(30).
		SetName("a8m").
		Save(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	log.Println("user was created: ", u)
	return u, nil
}

func QueryUser(ctx context.Context, client *ent.Client) (*ent.User, error) {
	u, err := client.User.
		Query().
		Where(user.Name("a8m")).
		Only(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	log.Println("user returned: ", u)
	return u, nil
}

func CreateCars(ctx context.Context, client *ent.Client) (*ent.User, error) {
	tesla, err := client.Car.
		Create().
		SetModel("Tesla").
		SetRegisteredAt(time.Now()).
		Save(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	log.Println("car was created: ", tesla)
	ford, err := client.Car.
		Create().
		SetModel("Ford").
		SetRegisteredAt(time.Now()).
		Save(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	log.Println("car was created: ", ford)
	a8m, err := client.User.
		Create().
		SetAge(30).
		SetName("a8m").
		AddCars(tesla, ford).
		Save(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	log.Println("user was created: ", a8m)
	return a8m, nil
}

func QueryCars(ctx context.Context, a8m *ent.User) error {
	cars, err := a8m.QueryCars().All(ctx)
	if err != nil {
		return errors.WithStack(err)
	}
	log.Println("returned cars:", cars)
	ford, err := a8m.QueryCars().
		Where(car.Model("Ford")).
		Only(ctx)
	if err != nil {
		return errors.WithStack(err)
	}
	log.Println(ford)
	return nil
}

func CreateGraph(ctx context.Context, client *ent.Client) error {
	a8m, err := client.User.
		Create().
		SetAge(30).
		SetName("Ariel").
		Save(ctx)
	if err != nil {
		return errors.WithStack(err)
	}
	neta, err := client.User.
		Create().
		SetAge(28).
		SetName("Neta").
		Save(ctx)
	if err != nil {
		return errors.WithStack(err)
	}
	err = client.Car.
		Create().
		SetModel("Tesla").
		SetRegisteredAt(time.Now()).
		SetOwner(a8m).
		Exec(ctx)
	if err != nil {
		return errors.WithStack(err)
	}
	err = client.Car.
		Create().
		SetModel("Mazda").
		SetRegisteredAt(time.Now()).
		SetOwner(a8m).
		Exec(ctx)
	if err != nil {
		return errors.WithStack(err)
	}
	err = client.Car.
		Create().
		SetModel("Ford").
		SetRegisteredAt(time.Now()).
		SetOwner(neta).
		Exec(ctx)
	if err != nil {
		return errors.WithStack(err)
	}
	err = client.Group.
		Create().
		SetName("GitLab").
		AddUsers(neta, a8m).
		Exec(ctx)
	if err != nil {
		return errors.WithStack(err)
	}
	err = client.Group.
		Create().
		SetName("GitHub").
		AddUsers(a8m).
		Exec(ctx)
	if err != nil {
		return errors.WithStack(err)
	}
	log.Println("The graph was created successfully")
	return nil
}

func QueryGithub(ctx context.Context, client *ent.Client) error {
	cars, err := client.Group.
		Query().
		Where(group.Name("GitHub")).
		QueryUsers().
		QueryCars().
		All(ctx)
	if err != nil {
		return errors.WithStack(err)
	}
	log.Println("cars returned:", cars)
	return nil
}

func QueryArielCars(ctx context.Context, client *ent.Client) error {
	a8m := client.User.
		Query().
		Where(
			user.HasCars(),
			user.Name("Ariel"),
		).
		OnlyX(ctx)
	cars, err := a8m.
		QueryGroups().
		QueryUsers().
		QueryCars().
		Where(
			car.Not(
				car.Model("Mazda"),
			),
		).
		All(ctx)
	if err != nil {
		return errors.WithStack(err)
	}
	log.Println("cars returned:", cars)
	return nil
}

func QueryGroupWithUsers(ctx context.Context, client *ent.Client) error {
	groups, err := client.Group.
		Query().
		Where(group.HasUsers()).
		All(ctx)
	if err != nil {
		return errors.WithStack(err)
	}
	log.Println("groups returned:", groups)
	return nil
}

func PanicErr[T any](v T, err error) T {
	if err != nil {
		log.Printf("unexpected error: %+v", err)
		panic(err)
	}
	return v
}

func PanicErrExec(err error) {
	if err != nil {
		log.Printf("unexpected error: %+v", err)
		panic(err)
	}
}
