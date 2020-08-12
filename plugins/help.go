package plugins

func ShowHelp() string {
	helpMessage := `
		Commands available:
			!help: Show message
			!image: Return random image link by search string
			!gif: Return random gif image link
			!video: return random youtube link
			!gs: Common google search
			!commit: return random commit message
			!adme: Return link from adme.ru by search string

		Usage for Skype group chat: @Cho !image pussy cat
		Usage for Skype direct message: !image pussy cat
	`
	return helpMessage
}
