package invitationcode

import (
	"context"
	"fmt"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	invitationcodemwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/invitation/invitationcode"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/invitation/invitationcode"
	invitationcodemwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/invitation/invitationcode"
)

type queryHandler struct {
	*Handler
	codes []*invitationcodemwpb.InvitationCode
	users map[string]*usermwpb.User
	infos []*npool.InvitationCode
}

func (h *queryHandler) getUsers(ctx context.Context) error {
	ids := []string{}
	for _, code := range h.codes {
		ids = append(ids, code.UserID)
	}
	users, _, err := usermwcli.GetUsers(ctx, &usermwpb.Conds{
		AppID:  &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		EntIDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: ids},
	}, 0, int32(len(ids)))
	if err != nil {
		return err
	}
	for _, user := range users {
		h.users[user.EntID] = user
	}
	return nil
}

func (h *queryHandler) formalize() {
	for _, code := range h.codes {
		user, ok := h.users[code.UserID]
		if !ok {
			continue
		}
		h.infos = append(h.infos, &npool.InvitationCode{
			ID:             code.ID,
			EntID:          code.EntID,
			AppID:          code.AppID,
			UserID:         code.UserID,
			EmailAddress:   user.EmailAddress,
			PhoneNO:        user.PhoneNO,
			Username:       user.Username,
			InvitationCode: code.InvitationCode,
			Disabled:       code.Disabled,
			CreatedAt:      code.CreatedAt,
			UpdatedAt:      code.UpdatedAt,
		})
	}
}

func (h *Handler) GetInvitationCode(ctx context.Context) (*npool.InvitationCode, error) {
	if h.EntID == nil {
		return nil, fmt.Errorf("invalid entid")
	}
	info, err := invitationcodemwcli.GetInvitationCode(ctx, *h.EntID)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, nil
	}

	handler := &queryHandler{
		Handler: h,
		codes:   []*invitationcodemwpb.InvitationCode{info},
		users:   map[string]*usermwpb.User{},
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

func (h *Handler) GetInvitationCodes(ctx context.Context) ([]*npool.InvitationCode, uint32, error) {
	infos, total, err := invitationcodemwcli.GetInvitationCodes(ctx, &invitationcodemwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
	}, h.Offset, h.Limit)
	if err != nil {
		return nil, 0, err
	}
	if len(infos) == 0 {
		return nil, total, nil
	}

	handler := &queryHandler{
		Handler: h,
		codes:   infos,
		users:   map[string]*usermwpb.User{},
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
