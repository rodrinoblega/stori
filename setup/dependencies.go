package setup

import (
	"github.com/rodrinoblega/stori/adapters/email_sender"
	"github.com/rodrinoblega/stori/adapters/watchers"
	"github.com/rodrinoblega/stori/config"
	"github.com/rodrinoblega/stori/frameworks/database"
	"github.com/rodrinoblega/stori/uses_cases"
)

type AppDependencies struct {
	DB               uses_cases.Database
	FileWatcher      uses_cases.Watcher
	ProcessFile      *uses_cases.ProcessFileUseCase
	ProcessDirectory *uses_cases.ProcessDirectoryFilesUseCase
}

func InitializeDependencies(envConf *config.Config) *AppDependencies {
	db := database.New(envConf)
	watcher := watchers.NewWatcherPath("/path")
	processFile := initializeProcessFile(db, envConf)
	processDirectory := uses_cases.NewProcessDirectoryFilesUseCase(processFile)

	return &AppDependencies{
		DB:               db,
		FileWatcher:      watcher,
		ProcessFile:      processFile,
		ProcessDirectory: processDirectory,
	}
}

func initializeProcessFile(db uses_cases.Database, envConf *config.Config) *uses_cases.ProcessFileUseCase {
	emailSender := email_sender.NewSMTPEmailSender(
		envConf.EmailHost,
		envConf.EmailPort,
		envConf.EmailUsername,
		envConf.EmailPassword,
	)

	return uses_cases.NewProcessFileUseCase(
		uses_cases.NewFileReaderUseCase(),
		uses_cases.NewStoreTransactionsUseCase(db),
		uses_cases.NewEmailSummaryUseCase(emailSender, db),
	)
}
