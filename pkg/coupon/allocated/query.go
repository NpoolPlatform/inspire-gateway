package allocated

import (
	"context"
	"fmt"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	allocatedmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/coupon/allocated"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	allocatedmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/coupon/allocated"
)

type queryHandler struct {
	*Handler
	infos []*allocatedmwpb.Coupon
	users map[string]*usermwpb.User
}

func (h *queryHandler) getUsers(ctx context.Context) error {
	userIDs := []string{}
	for _, info := range h.infos {
		userIDs = append(userIDs, info.UserID)
	}
	users, _, err := usermwcli.GetUsers(ctx, &usermwpb.Conds{
		AppID:  &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		EntIDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: userIDs},
	}, int32(0), int32(len(userIDs)))
	if err != nil {
		return err
	}
	for _, user := range users {
		h.users[user.EntID] = user
	}
	return nil
}

func (h *queryHandler) formalize() {
	for _, info := range h.infos {
		user, ok := h.users[info.UserID]
		if !ok {
			continue
		}
		info.Username = user.Username
		info.EmailAddress = user.EmailAddress
		info.PhoneNO = user.PhoneNO
	}
}

func (h *Handler) GetCoupon(ctx context.Context) (*allocatedmwpb.Coupon, error) {
	if h.EntID == nil {
		return nil, fmt.Errorf("invalid entid")
	}
	info, err := allocatedmwcli.GetCoupon(ctx, *h.EntID)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, nil
	}

	handler := &queryHandler{
		Handler: h,
		infos:   []*allocatedmwpb.Coupon{info},
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

func (h *Handler) GetCoupons(ctx context.Context) ([]*allocatedmwpb.Coupon, uint32, error) {
	conds := &allocatedmwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
	}
	if h.UserID != nil {
		conds.UserID = &basetypes.StringVal{Op: cruder.EQ, Value: *h.UserID}
	}

	infos, total, err := allocatedmwcli.GetCoupons(ctx, conds, h.Offset, h.Limit)
	if err != nil {
		return nil, 0, err
	}
	if len(infos) == 0 {
		return nil, total, nil
	}

	handler := &queryHandler{
		Handler: h,
		infos:   infos,
		users:   map[string]*usermwpb.User{},
	}
	if err := handler.getUsers(ctx); err != nil {
		return nil, 0, err
	}
	handler.formalize()
	return handler.infos, total, nil
}
