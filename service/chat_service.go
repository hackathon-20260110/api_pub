package service

import (
	"context"
	"sort"
	"time"

	"github.com/hackathon-20260110/api/adapter"
	"github.com/hackathon-20260110/api/response"
	"github.com/hackathon-20260110/api/utils"
	"go.uber.org/dig"
)

type ChatService struct {
	container *dig.Container
}

func NewChatService(container *dig.Container) *ChatService {
	return &ChatService{container: container}
}

type chatListItem struct {
	ID              string
	PartnerID       string
	PartnerName     string
	PartnerImageURL string
	LastMessage     string
	LastMessageAt   time.Time
	MatchingScore   int
	IsMatched       bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (s *ChatService) GetChats(ctx context.Context, userID string) (*response.GetChatsResponse, error) {
	var userChatAdapter adapter.UserChatAdapter
	var avatarChatAdapter adapter.AvatarChatAdapter
	var avatarAdapter adapter.AvatarAdapter
	var userAdapter adapter.UserAdapter
	var matchingAdapter adapter.MatchingAdapter

	if err := s.container.Invoke(func(
		uca adapter.UserChatAdapter,
		aca adapter.AvatarChatAdapter,
		aa adapter.AvatarAdapter,
		ua adapter.UserAdapter,
		ma adapter.MatchingAdapter,
	) error {
		userChatAdapter = uca
		avatarChatAdapter = aca
		avatarAdapter = aa
		userAdapter = ua
		matchingAdapter = ma
		return nil
	}); err != nil {
		return nil, utils.WrapError(err)
	}

	var allChats []chatListItem

	// 1. マッチ済みユーザーとのチャットを取得
	matchedChats, err := s.getMatchedUserChats(ctx, userID, matchingAdapter, userChatAdapter, userAdapter)
	if err != nil {
		return nil, utils.WrapError(err)
	}
	allChats = append(allChats, matchedChats...)

	// マッチ済みユーザーIDのセットを作成（重複防止用）
	matchedUserIDs := make(map[string]bool)
	for _, chat := range matchedChats {
		matchedUserIDs[chat.PartnerID] = true
	}

	// 2. アバターとのチャットを取得（マッチ済みは除外）
	avatarChats, err := s.getAvatarChats(ctx, userID, avatarChatAdapter, avatarAdapter, userAdapter, matchedUserIDs)
	if err != nil {
		return nil, utils.WrapError(err)
	}
	allChats = append(allChats, avatarChats...)

	// 3. LastMessageAtで降順ソート
	sort.Slice(allChats, func(i, j int) bool {
		return allChats[i].LastMessageAt.After(allChats[j].LastMessageAt)
	})

	// 4. レスポンス形式に変換
	chats := make([]response.Chat, len(allChats))
	for i, chat := range allChats {
		chats[i] = response.Chat{
			ID:              chat.ID,
			UserID:          userID,
			PartnerID:       chat.PartnerID,
			PartnerName:     chat.PartnerName,
			PartnerImageURL: chat.PartnerImageURL,
			LastMessage:     chat.LastMessage,
			LastMessageAt:   chat.LastMessageAt.Format(time.RFC3339),
			UnreadCount:     0,
			MatchingScore:   chat.MatchingScore,
			IsMatched:       chat.IsMatched,
			CreatedAt:       chat.CreatedAt.Format(time.RFC3339),
			UpdatedAt:       chat.UpdatedAt.Format(time.RFC3339),
		}
	}

	return &response.GetChatsResponse{
		Chats: chats,
		Total: len(chats),
	}, nil
}

func (s *ChatService) getMatchedUserChats(
	ctx context.Context,
	userID string,
	matchingAdapter adapter.MatchingAdapter,
	userChatAdapter adapter.UserChatAdapter,
	userAdapter adapter.UserAdapter,
) ([]chatListItem, error) {
	matchings, err := matchingAdapter.GetMatchingsByUserID(userID)
	if err != nil {
		return nil, err
	}

	var chats []chatListItem
	for _, matching := range matchings {
		partnerID := matching.User1ID
		if matching.User1ID == userID {
			partnerID = matching.User2ID
		}

		partner, err := userAdapter.GetByID(partnerID)
		if err != nil {
			continue
		}

		messages, err := userChatAdapter.GetUserChatMessages(ctx, userID, partnerID)
		var lastMessage string
		var lastMessageAt time.Time
		if err == nil && len(messages) > 0 {
			lastMsg := messages[len(messages)-1]
			lastMessage = lastMsg.Message
			lastMessageAt = lastMsg.CreatedAt
		} else {
			lastMessageAt = matching.CreatedAt
		}

		chats = append(chats, chatListItem{
			ID:              matching.ID,
			PartnerID:       partnerID,
			PartnerName:     partner.DisplayName,
			PartnerImageURL: partner.ProfileImageURL,
			LastMessage:     lastMessage,
			LastMessageAt:   lastMessageAt,
			MatchingScore:   100,
			IsMatched:       true,
			CreatedAt:       matching.CreatedAt,
			UpdatedAt:       matching.UpdatedAt,
		})
	}

	return chats, nil
}

func (s *ChatService) getAvatarChats(
	ctx context.Context,
	userID string,
	avatarChatAdapter adapter.AvatarChatAdapter,
	avatarAdapter adapter.AvatarAdapter,
	userAdapter adapter.UserAdapter,
	matchedUserIDs map[string]bool,
) ([]chatListItem, error) {
	relations, err := avatarAdapter.GetUserAvatarRelationsByUserID(userID)
	if err != nil {
		return nil, err
	}

	var chats []chatListItem
	for _, relation := range relations {
		avatar, err := avatarAdapter.GetByID(relation.AvatarID)
		if err != nil {
			continue
		}

		// 既にマッチ済みのユーザーのアバターはスキップ
		if matchedUserIDs[avatar.UserID] {
			continue
		}

		owner, err := userAdapter.GetByID(avatar.UserID)
		if err != nil {
			continue
		}

		messages, err := avatarChatAdapter.GetAvatarChatMessages(ctx, userID, relation.AvatarID)
		var lastMessage string
		var lastMessageAt time.Time
		if err == nil && len(messages) > 0 {
			lastMsg := messages[len(messages)-1]
			lastMessage = lastMsg.Message
			lastMessageAt = lastMsg.CreatedAt
		} else {
			lastMessageAt = relation.CreatedAt
		}

		chats = append(chats, chatListItem{
			ID:              relation.ID,
			PartnerID:       relation.AvatarID,
			PartnerName:     owner.DisplayName,
			PartnerImageURL: avatar.AvatarIconURL,
			LastMessage:     lastMessage,
			LastMessageAt:   lastMessageAt,
			MatchingScore:   relation.MatchingPoint,
			IsMatched:       false,
			CreatedAt:       relation.CreatedAt,
			UpdatedAt:       relation.UpdatedAt,
		})
	}

	return chats, nil
}
