package cmd

import (
	"context"
	"fmt"
	"log/slog"

	"charm.land/huh/v2"
	"github.com/neo4j/neo4j-go-driver/v6/neo4j"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(delete_post)
}

var delete_post = &cobra.Command{
	Use:   "delete-post [post-id]",
	Short: "Delete a post",
	Long:  "Delete a post from the neo4j database using its post_id.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		postID := args[0]

		var confirm bool
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title(fmt.Sprintf("Are you sure you want to delete post %s?", postID)).
					Value(&confirm),
			),
		).WithTheme(huh.ThemeFunc(huh.ThemeBase16))

		if err := form.Run(); err != nil {
			slog.Error("Form failed", "error", err)
			return
		}

		if !confirm {
			cmd.Println("Delete cancelled.")
			return
		}

		deletePost(postID)
		cmd.Printf("Post %s deleted successfully\n", postID)
	},
}

func deletePost(post_id string) {
	d := GetDriver()
	ctx := context.Background()
	_, err := neo4j.ExecuteQuery(ctx, d, `
	MATCH (p:Post) 
	WHERE p.post_id = $post_id
	DETACH DELETE p
	`, map[string]any{"post_id": post_id}, neo4j.EagerResultTransformer)
	if err != nil {
		slog.Error("failed to delete post", "post_id", post_id, "error", err)
	}
}
