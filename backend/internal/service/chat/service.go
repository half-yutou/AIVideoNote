package chat

import (
	"context"

	"github.com/aivideonote/backend/internal/repository"
	"github.com/aivideonote/backend/internal/service/llm"
)

type Service struct {
	noteRepo     *repository.NoteRepository
	providerRepo *repository.ProviderRepository
	llmClient    *llm.Client
}

func NewService(noteRepo *repository.NoteRepository, providerRepo *repository.ProviderRepository, llmClient *llm.Client) *Service {
	return &Service{noteRepo: noteRepo, providerRepo: providerRepo, llmClient: llmClient}
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func (s *Service) IndexStatus(taskID string) bool {
	_, err := s.noteRepo.FindByTaskID(taskID)
	return err == nil
}

func (s *Service) BuildContext(taskID, question string) (markdownSnippets string, transcriptSnippets []string, err error) {
	note, err := s.noteRepo.FindByTaskID(taskID)
	if err != nil {
		return "", nil, err
	}

	if note.Transcript != "" {
		snippet := note.Transcript
		if len(snippet) > 2000 {
			snippet = snippet[:2000]
		}
		transcriptSnippets = append(transcriptSnippets, snippet)
	}

	markdownSnippets = note.Markdown
	if len(note.Markdown) > 4000 {
		markdownSnippets = note.Markdown[:4000]
	}

	return markdownSnippets, transcriptSnippets, nil
}

func (s *Service) Ask(ctx context.Context, providerID, modelName string, messages []llm.ChatMessage) (string, error) {
	provider, err := s.providerRepo.FindByID(providerID)
	if err != nil {
		return "", err
	}
	return s.llmClient.Chat(ctx, provider.BaseURL, provider.ApiKey, modelName, messages)
}

func (s *Service) AskStream(ctx context.Context, providerID, modelName string, messages []llm.ChatMessage) (<-chan string, <-chan error, error) {
	provider, err := s.providerRepo.FindByID(providerID)
	if err != nil {
		return nil, nil, err
	}
	ch, errCh := s.llmClient.ChatStream(ctx, provider.BaseURL, provider.ApiKey, modelName, messages)
	return ch, errCh, nil
}
