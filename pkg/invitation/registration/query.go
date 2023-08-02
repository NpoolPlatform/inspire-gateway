package registration

import (
	"context"
	"fmt"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	regmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/invitation/registration"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/invitation/registration"
	regmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/invitation/registration"
)

type queryHandler struct {
	*Handler
	registrations []*regmwpb.Registration
	users         map[string]*usermwpb.User
	infos         []*npool.Registration
}

func (h *queryHandler) getUsers(ctx context.Context) error {
	if h.AppID == nil {
		return fmt.Errorf("invalid appid")
	}

	userIDs := []string{}
	for _, registration := range h.registrations {
		userIDs = append(userIDs, registration.InviterID, registration.InviteeID)
	}
	users, _, err := usermwcli.GetUsers(ctx, &usermwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		IDs:   &basetypes.StringSliceVal{Op: cruder.IN, Value: userIDs},
	}, 0, int32(len(userIDs)))
	if err != nil {
		return err
	}
	for _, user := range users {
		h.users[user.ID] = user
	}
	return nil
}

func (h *queryHandler) formalize() {
	for _, registration := range h.registrations {
		inviter, ok := h.users[registration.InviterID]
		if !ok {
			continue
		}

		invitee, ok := h.users[registration.InviteeID]
		if !ok {
			continue
		}

		h.infos = append(h.infos, &npool.Registration{
			ID:                  registration.ID,
			AppID:               registration.AppID,
			InviterID:           registration.InviterID,
			InviterEmailAddress: inviter.EmailAddress,
			InviterPhoneNO:      inviter.PhoneNO,
			InviterUsername:     inviter.Username,
			InviteeID:           registration.InviteeID,
			InviteeEmailAddress: invitee.EmailAddress,
			InviteePhoneNO:      invitee.PhoneNO,
			InviteeUsername:     invitee.Username,
			CreatedAt:           registration.CreatedAt,
			UpdatedAt:           registration.UpdatedAt,
		})
	}
}

func (h *Handler) GetRegistration(ctx context.Context) (*npool.Registration, error) {
	if h.ID == nil {
		return nil, fmt.Errorf("invalid id")
	}

	info, err := regmwcli.GetRegistration(ctx, *h.ID)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, nil
	}

	handler := &queryHandler{
		Handler:       h,
		registrations: []*regmwpb.Registration{info},
	}
	if err := handler.getUsers(ctx); err != nil {
		return nil, err
	}
	handler.formalize()
	if len(handler.infos) == 0 {
		return nil, nil
	}
	return handler.infos[0], nil
}

func (h *Handler) GetRegistrations(ctx context.Context) ([]*npool.Registration, uint32, error) {
	if h.AppID == nil {
		return nil, 0, fmt.Errorf("invalid appid")
	}
	infos, total, err := regmwcli.GetRegistrations(ctx, &regmwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
	}, h.Offset, h.Limit)
	if err != nil {
		return nil, 0, err
	}
	if len(infos) == 0 {
		return nil, total, nil
	}

	handler := &queryHandler{
		Handler:       h,
		registrations: infos,
	}
	if err := handler.getUsers(ctx); err != nil {
		return nil, 0, err
	}
	handler.formalize()
	if len(handler.infos) == 0 {
		return nil, total, nil
	}
	return handler.infos, total, nil
}
