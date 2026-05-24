package cmd

import (
	"context"
	"fmt"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
	"github.com/neo4j/neo4j-go-driver/v6/neo4j"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(list_posts)
}

var list_posts = &cobra.Command{
	Use:   "list-posts",
	Short: "List all posts.",
	Long:  "List all posts in neo4j db.",
	Run: func(cmd *cobra.Command, args []string) {
		listPosts()
	},
}

func listPosts() {
	d := GetDriver()
	ctx := context.Background()
	result, err := neo4j.ExecuteQuery(ctx, d, `
	MATCH (p:Post) 
	//WHERE p.source="cms"
    RETURN p.title as title, p.post_id as post_id
    ORDER BY p.created;
	`, nil, neo4j.EagerResultTransformer)
	if err != nil {
		panic(err)
	}

	t := table.New().
		Border(lipgloss.NormalBorder()).
		Headers("Title", "Post ID")
	for _, record := range result.Records {
		title, _ := record.Get("title")
		post_id, _ := record.Get("post_id")

		t.Row(title.(string), post_id.(string))
	}
	fmt.Println(t)
}

// w := tabwriter.NewWriter(os.Stdout, 100, 1, 1, ' ', 0)
// _, err = fmt.Fprintln(w, "Title\t|", "Post ID")
// if err != nil {
// 	slog.Error("Something wrong", "error", err)
// }
// _, err = fmt.Fprintln(w, "--------------------------------------------------")
// if err != nil {
// 	slog.Error("Something wrong", "error", err)
// }

// for _, record := range result.Records {
// 	title, _ := record.Get("title")
// 	post_id, _ := record.Get("post_id")

// 	_, err = fmt.Fprintln(w, title, "\t|", post_id)
// 	if err != nil {
// 		slog.Error("Something wrong", "error", err)
// 	}
// 	err = w.Flush()
// 	if err != nil {
// 		slog.Error("Something wrong", "error", err)
// 	}
// }
