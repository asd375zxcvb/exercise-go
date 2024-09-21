package main

import (
	"context"
	"entgo/ent"
	"entgo/ent/car"
	"entgo/ent/group"
	"entgo/ent/user"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"time"
)

func main() {
	client, err := ent.Open("postgres", "host=127.0.0.1 port=5432 user=postgres dbname=postgres password=123456 sslmode=disable", ent.Debug())
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	defer func(client *ent.Client) {
		err := client.Close()
		if err != nil {
			log.Fatalf("failed closing client: %v", err)
		}
	}(client)
	ctx := context.Background()
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	err = round1(ctx, client)
	if err != nil {
		log.Fatalf("failed round1 resources: %v", err)
	}
	err = round2(ctx, client)
	if err != nil {
		log.Fatalf("failed round2 resources: %v", err)
	}
	err = round3(ctx, client)
	if err != nil {
		log.Fatalf("failed round3 resources: %v", err)
	}
	err = client.Close()
	if err != nil {
		log.Fatalf("failed client Close resources: %v", err)
	}
}

func round1(ctx context.Context, client *ent.Client) error {
	log.Println("====================  round1  start =======================")
	_, err := CreateUser(ctx, client)
	if err != nil {
		return err
	}
	_, err = QueryUser(ctx, client)
	if err != nil {
		return err
	}
	client.User.Delete().Exec(ctx)
	log.Println("====================  round1  end   =======================")
	return nil
}

func round2(ctx context.Context, client *ent.Client) error {
	log.Println("====================  round2  start =======================")
	a8m, err := CreateCars(ctx, client)
	if err != nil {
		return err
	}
	err = QueryCars(ctx, a8m)
	if err != nil {
		return err
	}
	_, err = client.User.Delete().Exec(ctx)
	if err != nil {
		return err
	}
	_, err = client.Car.Delete().Exec(ctx)
	if err != nil {
		return err
	}
	log.Println("====================  round2  end   =======================")
	return nil
}

func round3(ctx context.Context, client *ent.Client) error {
	log.Println("====================  round3  start =======================")
	err := CreateGraph(ctx, client)
	if err != nil {
		return err
	}
	err = QueryGithub(ctx, client)
	if err != nil {
		return err
	}
	err = QueryArielCars(ctx, client)
	if err != nil {
		return err
	}
	err = QueryGroupWithUsers(ctx, client)
	if err != nil {
		return err
	}
	_, err = client.User.Delete().Exec(ctx)
	if err != nil {
		return err
	}
	_, err = client.Car.Delete().Exec(ctx)
	if err != nil {
		return err
	}
	_, err = client.Group.Delete().Exec(ctx)
	if err != nil {
		return err
	}
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
		return nil, fmt.Errorf("failed creating user: %w", err)
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
		return nil, fmt.Errorf("failed querying user: %w", err)
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
		return nil, fmt.Errorf("failed creating car: %w", err)
	}
	log.Println("car was created: ", tesla)
	ford, err := client.Car.
		Create().
		SetModel("Ford").
		SetRegisteredAt(time.Now()).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating car: %w", err)
	}
	log.Println("car was created: ", ford)
	a8m, err := client.User.
		Create().
		SetAge(30).
		SetName("a8m").
		AddCars(tesla, ford).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating user: %w", err)
	}
	log.Println("user was created: ", a8m)
	return a8m, nil
}

func QueryCars(ctx context.Context, a8m *ent.User) error {
	cars, err := a8m.QueryCars().All(ctx)
	if err != nil {
		return fmt.Errorf("failed querying user cars: %w", err)
	}
	log.Println("returned cars:", cars)
	ford, err := a8m.QueryCars().
		Where(car.Model("Ford")).
		Only(ctx)
	if err != nil {
		return fmt.Errorf("failed querying user cars: %w", err)
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
		return err
	}
	neta, err := client.User.
		Create().
		SetAge(28).
		SetName("Neta").
		Save(ctx)
	if err != nil {
		return err
	}
	err = client.Car.
		Create().
		SetModel("Tesla").
		SetRegisteredAt(time.Now()).
		SetOwner(a8m).
		Exec(ctx)
	if err != nil {
		return err
	}
	err = client.Car.
		Create().
		SetModel("Mazda").
		SetRegisteredAt(time.Now()).
		SetOwner(a8m).
		Exec(ctx)
	if err != nil {
		return err
	}
	err = client.Car.
		Create().
		SetModel("Ford").
		SetRegisteredAt(time.Now()).
		SetOwner(neta).
		Exec(ctx)
	if err != nil {
		return err
	}
	err = client.Group.
		Create().
		SetName("GitLab").
		AddUsers(neta, a8m).
		Exec(ctx)
	if err != nil {
		return err
	}
	err = client.Group.
		Create().
		SetName("GitHub").
		AddUsers(a8m).
		Exec(ctx)
	if err != nil {
		return err
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
		return fmt.Errorf("failed getting cars: %w", err)
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
		return fmt.Errorf("failed getting cars: %w", err)
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
		return fmt.Errorf("failed getting groups: %w", err)
	}
	log.Println("groups returned:", groups)
	return nil
}
