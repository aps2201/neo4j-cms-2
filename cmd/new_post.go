package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"time"

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
		var post_created string
		var post_title string
		var post_content string
		post_created = time.Now().Format("20060102150405")
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
		).WithTheme(huh.ThemeFunc(huh.ThemeBase16))
		if err := form.Run(); err != nil {
			slog.Error("Form failed", ":", err)
		}
		post_id := writeNewPost(post_title, post_content, post_created)
		cmd.Println(post_id)

	},
}

func writeNewPost(post_title string, post_content string, post_created string) (post_id string) {
	var post_uuid uuid.UUID
	post_uuid, _ = uuid.NewRandom()
	post_id = post_uuid.String()
	d := GetDriver()
	ctx := context.Background()
	_, err := neo4j.ExecuteQuery(ctx, d, `
	CREATE (p:Post {post_id:$post_id, source:"cms"}) 
	SET p.title = $post_title,
		p.content 	= $post_content,
		p.created = $post_created
	`, map[string]any{"post_id": post_id,
		"post_title":   post_title,
		"post_content": post_content,
		"post_created": post_created,
	}, neo4j.EagerResultTransformer)
	if err != nil {
		slog.Error("cant execute query", "error:", err)
	}
	return
}

//TODO: write new post
