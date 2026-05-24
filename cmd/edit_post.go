package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"charm.land/huh/v2"
	"github.com/neo4j/neo4j-go-driver/v6/neo4j"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(edit_post)
}

var edit_post = &cobra.Command{
	Use:   "edit-post [post-id]",
	Short: "Edit a post",
	Long:  "Get a post in neo4j db and modify it with forms. Post ID can be partial ",
	Args:  cobra.ExactArgs(1), // post-id argument
	Run: func(cmd *cobra.Command, args []string) {
		var post_title string
		var post_content string
		var post_id string
		var post_modified string
		var confirm bool

		if len(args[0]) < 8 {
			slog.Error("post_id too short")
		}
		post := getPost(args[0])

		post_title = post.Props["title"].(string)
		post_content = post.Props["content"].(string)
		post_id = post.Props["post_id"].(string)
		post_modified = time.Now().Format("20060102150405")

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().Title("Title").Value(&post_title).Validate(func(s string) error {
					if len(s) < 5 {
						return fmt.Errorf("something went wrong")
					}
					return nil
				}),
				huh.NewText().Title("Content").Lines(20).Value(&post_content).Validate(func(s string) error {
					if post_content == "" {
						return fmt.Errorf("something went wrong")
					}
					return nil
				}),
				huh.NewConfirm().
					Title(fmt.Sprintf("Are you sure you want to edit %s?", post_title)).
					Value(&confirm),
			),
		).WithTheme(huh.ThemeFunc(huh.ThemeBase16))
		if err := form.Run(); err != nil {
			slog.Error("Form failed", ":", err)
		}
		if !confirm {
			cmd.Println("Edit cancelled.")
			return
		}
		writePost(post_id, post_title, post_content, post_modified)

	},
}

func getPost(post_id string) (post neo4j.Node) {
	d := GetDriver()
	ctx := context.Background()
	result, err := neo4j.ExecuteQuery(ctx, d, `
	MATCH (p:Post) 
	WHERE p.post_id STARTS WITH $post_id
    RETURN p;
	`, map[string]any{"post_id": post_id}, neo4j.EagerResultTransformer)
	if err != nil {
		slog.Error("cant execute query", "error:", err)
	}

	record := result.Records[0]
	post_, ok := record.Get("p")
	if !ok {
		slog.Error("post not found")
	}

	post, ok = post_.(neo4j.Node)
	if !ok {
		slog.Error("not quite the type, got %T", post_)
	}

	return
}

func writePost(post_id string, post_title string, post_content string, post_modified string) {
	d := GetDriver()
	ctx := context.Background()
	_, err := neo4j.ExecuteQuery(ctx, d, `
	MATCH (p:Post {post_id:$post_id}) 
	SET p.title = $post_title,
		p.content = $post_content,
		p.source = "cms",
		p.modified = $post_modified
	`, map[string]any{"post_id": post_id,
		"post_title":    post_title,
		"post_content":  post_content,
		"post_modified": post_modified,
	}, neo4j.EagerResultTransformer)
	if err != nil {
		slog.Error("cant execute query", "error:", err)
	}
}
