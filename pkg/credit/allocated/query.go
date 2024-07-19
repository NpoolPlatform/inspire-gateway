package allocated

import (
	"context"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	"github.com/NpoolPlatform/go-service-framework/pkg/wlog"
	creditallocatedmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/credit/allocated"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	appusermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/credit/allocated"
	creditallocatedmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/credit/allocated"
)

type queryHandler struct {
	*Handler
	creditallocateds []*creditallocatedmwpb.CreditAllocated
	appuser          map[string]*appusermwpb.User
	infos            []*npool.CreditAllocated
}

func (h *queryHandler) getCreditAllocateds(ctx context.Context) error {
	conds := &creditallocatedmwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
	}
	if h.UserID != nil {
		conds.UserID = &basetypes.StringVal{Op: cruder.EQ, Value: *h.UserID}
	}
	infos, _, err := creditallocatedmwcli.GetCreditAllocateds(ctx, conds, h.Offset, h.Limit)
	if err != nil {
		return wlog.WrapError(err)
	}
	h.creditallocateds = infos
	return nil
}

func (h *queryHandler) getUsers(ctx context.Context) error {
	userIDs := []string{}
	for _, allocated := range h.creditallocateds {
		userIDs = append(userIDs, allocated.UserID)
	}
	users, _, err := usermwcli.GetUsers(ctx, &appusermwpb.Conds{
		EntIDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: userIDs},
	}, 0, int32(len(userIDs)))
	if err != nil {
		return wlog.WrapError(err)
	}

	for _, user := range users {
		h.appuser[user.EntID] = user
	}
	return nil
}

func (h *queryHandler) formalize() {
	for _, info := range h.creditallocateds {
		user, ok := h.appuser[info.UserID]
		if !ok {
			continue
		}

		h.infos = append(h.infos, &npool.CreditAllocated{
			ID:           info.ID,
			EntID:        info.EntID,
			AppID:        info.AppID,
			Credits:      info.Value,
			CreatedAt:    info.CreatedAt,
			UserID:       user.EntID,
			PhoneNO:      user.PhoneNO,
			EmailAddress: user.EmailAddress,
		})
	}
}

func (h *Handler) GetCreditAllocated(ctx context.Context) (*npool.CreditAllocated, error) {
	if h.EntID == nil {
		return nil, wlog.Errorf("invalid entid")
	}

	info, err := creditallocatedmwcli.GetCreditAllocated(ctx, *h.EntID)
	if err != nil {
		return nil, wlog.WrapError(err)
	}
	if info == nil {
		return nil, nil
	}

	handler := &queryHandler{
		Handler:          h,
		creditallocateds: []*creditallocatedmwpb.CreditAllocated{info},
		appuser:          map[string]*appusermwpb.User{},
	}

	if err := handler.getUsers(ctx); err != nil {
		return nil, wlog.WrapError(err)
	}

	handler.formalize()
	if len(handler.infos) == 0 {
		return nil, nil
	}

	return handler.infos[0], nil
}

func (h *Handler) GetCreditAllocateds(ctx context.Context) ([]*npool.CreditAllocated, uint32, error) {
	handler := &queryHandler{
		Handler:          h,
		creditallocateds: []*creditallocatedmwpb.CreditAllocated{},
		appuser:          map[string]*appusermwpb.User{},
	}

	if err := handler.getCreditAllocateds(ctx); err != nil {
		return nil, 0, wlog.WrapError(err)
	}

	if err := handler.getUsers(ctx); err != nil {
		return nil, 0, wlog.WrapError(err)
	}

	handler.formalize()
	return handler.infos, uint32(len(handler.infos)), nil
}
