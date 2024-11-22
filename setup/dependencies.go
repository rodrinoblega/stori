package setup

import (
	"github.com/rodrinoblega/stori/adapters/email_sender"
	"github.com/rodrinoblega/stori/adapters/watchers"
	"github.com/rodrinoblega/stori/config"
	"github.com/rodrinoblega/stori/frameworks/database"
	"github.com/rodrinoblega/stori/frameworks/email_service"
	"github.com/rodrinoblega/stori/uses_cases"
)

type AppDependencies struct {
	DB               uses_cases.Database
	FileWatcher      uses_cases.Watcher
	ProcessFile      *uses_cases.ProcessFileUseCase
	ProcessDirectory *uses_cases.ProcessDirectoryFilesUseCase
}

func InitializeProdDependencies(envConf *config.Config) *AppDependencies {
	db := database.New(envConf)
	watcher := watchers.NewWatcherPath("/path")
	processFile := initializeProdProcessFile(db)
	processDirectory := uses_cases.NewProcessDirectoryFilesUseCase(processFile)

	return &AppDependencies{
		DB:               db,
		FileWatcher:      watcher,
		ProcessFile:      processFile,
		ProcessDirectory: processDirectory,
	}
}

func InitializeLocalDependencies(envConf *config.Config) *AppDependencies {
	db := database.New(envConf)
	watcher := watchers.NewWatcherPath("/path")
	processFile := initializeLocalProcessFile(db, envConf)
	processDirectory := uses_cases.NewProcessDirectoryFilesUseCase(processFile)

	return &AppDependencies{
		DB:               db,
		FileWatcher:      watcher,
		ProcessFile:      processFile,
		ProcessDirectory: processDirectory,
	}
}

func initializeLocalProcessFile(db uses_cases.Database, envConf *config.Config) *uses_cases.ProcessFileUseCase {
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

func initializeProdProcessFile(db uses_cases.Database) *uses_cases.ProcessFileUseCase {
	emailSender := email_sender.NewSESEmailSender(
		email_service.CreateSESSession(),
		"stori.summary@gmail.com",
	)

	return uses_cases.NewProcessFileUseCase(
		uses_cases.NewFileReaderUseCase(),
		uses_cases.NewStoreTransactionsUseCase(db),
		uses_cases.NewEmailSummaryUseCase(emailSender, db),
	)
}
