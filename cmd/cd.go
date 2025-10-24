package cmd

//func init() {
//
//}
//
//var cdCmd = &cobra.Command{
//	Use:   "cd [worktree]",
//	Short: "Change directory to the specified worktree",
//	Args:  cobra.MaximumNArgs(1),
//	RunE: func(cmd *cobra.Command, args []string) error {
//		var worktree string
//		if len(args) == 1 {
//			worktree = args[0]
//		} else {
//			// If no argument is provided, select via fzf
//
//		}
//
//		g, err := git.New()
//		if err != nil {
//			return err
//		}
//
//		path, err := g.GetWorktreePath(worktree)
//		if err != nil {
//			return err
//		}
//
//		fmt.Println(path)
//		return nil
//	}
//}
