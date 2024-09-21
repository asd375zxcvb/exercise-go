package main

import (
	"context"
	"ent-eics/ent"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func main() {
	client, err := ent.Open("mysql", "root:123456@tcp(localhost:3306)/eics-jeecg?parseTime=True", ent.Debug())
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	defer client.Close()
	ctx := context.Background()
	_, err = QueryTeamMember(ctx, client)
}

func QueryTeamMember(ctx context.Context, client *ent.Client) (*ent.EicsReceptionTeamMember, error) {
	u, err := client.EicsReceptionTeamMember.
		Query().
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed querying user: %w", err)
	}
	log.Println("user returned: ", u)
	return u, nil
}
