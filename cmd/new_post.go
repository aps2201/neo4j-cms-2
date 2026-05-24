package cmd

import (
	"context"
	"fmt"
	"log/slog"

	"charm.land/huh/v2"
	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v6/neo4j"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(write_post)
}

var write_post = &cobra.Command{
	Use:   "write-post",
	Short: "Write a new post",
	Run: func(cmd *cobra.Command, args []string) {
		var post_title string
		var post_content string
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().Title("Title").Value(&post_title).Validate(func(s string) error {
					if len(s) < 5 {
						return fmt.Errorf("something went wrong, title too short")
					}
					return nil
				}),
				huh.NewText().Title("Content").Lines(20).Value(&post_content).Validate(func(s string) error {
					if post_content == "" {
						return fmt.Errorf("something went wrong")
					}
					return nil
				}),
			),
		)
		if err := form.Run(); err != nil {
			slog.Error("Form failed", ":", err)
		}
		post_id := writeNewPost(post_title, post_content)
		cmd.Println(post_id)

	},
}

func writeNewPost(post_title string, post_content string) (post_id string) {
	var post_uuid uuid.UUID
	post_uuid, _ = uuid.NewRandom()
	post_id = post_uuid.String()
	d := GetDriver()
	ctx := context.Background()
	_, err := neo4j.ExecuteQuery(ctx, d, `
	CREATE (p:Post {post_id:$post_id, source:"cms"}) 
	SET p.title = $post_title,
		p.content 	= $post_content
	`, map[string]any{"post_id": post_id,
		"post_title":   post_title,
		"post_content": post_content,
	}, neo4j.EagerResultTransformer)
	if err != nil {
		slog.Error("cant execute query", "error:", err)
	}
	return
}

//TODO: write new post
