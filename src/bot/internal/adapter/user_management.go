package adapter

import (
	"avtor.ru/bot/client"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type UMAdapter struct {
	client *client.Client
}

func NewUMServiceAdapter(address string) (*UMAdapter, error) {
	umClient, err := client.NewClient(address)
	if err != nil {
		return nil, fmt.Errorf("failed to create um client: %s", err)
	}

	return &UMAdapter{client: umClient}, nil
}

func (um *UMAdapter) GetUserRole(ctx context.Context, username string) (string, error) {
	resp, err := um.client.GetUserAuthUserID(ctx, username)
	if err != nil {
		return "", fmt.Errorf("failed to get user role: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get user role: invalid status code: %d", resp.StatusCode)
	}

	var role string
	if err := json.NewDecoder(resp.Body).Decode(&role); err != nil {
		return "", fmt.Errorf("failed to decode user role: %w", err)
	}

	log.Printf("user role: %s", role)

	return role, nil
}

func (um *UMAdapter) CreateUser(ctx context.Context, username, role string) error {
	resp, err := um.client.PostUserCreateUserID(ctx, username, client.PostUserCreateUserIDJSONRequestBody{
		Role:     role,
		UserName: username,
	})
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to create user: invalid status code: %d", resp.StatusCode)
	}

	return nil
}

func (um *UMAdapter) GetUsers(ctx context.Context) (*client.Users, error) {
	resp, err := um.client.GetUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get users: invalid status code: %d", resp.StatusCode)
	}

	users := &client.Users{}
	if err := json.NewDecoder(resp.Body).Decode(users); err != nil {
		return nil, fmt.Errorf("failed to decode users: %w", err)
	}

	return users, nil
}

func (um *UMAdapter) DeleteUser(ctx context.Context, username string) error {
	resp, err := um.client.DeleteUsersUserID(ctx, username)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("failed to delete user: invalid status code: %d", resp.StatusCode)
	}

	return nil
}
