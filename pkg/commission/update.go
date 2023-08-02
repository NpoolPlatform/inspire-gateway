package commission

import (
	"context"
	"fmt"

	commmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/commission"
	registrationmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/invitation/registration"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/commission"
	commmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/commission"
	registrationmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/invitation/registration"
)

type updateHandler struct {
	*Handler
	info *commmwpb.Commission
}

func (h *updateHandler) validateInviter(ctx context.Context) error {
	if h.UserID == nil {
		return fmt.Errorf("invalid userid")
	}

	exist, err := registrationmwcli.ExistRegistrationConds(
		ctx,
		&registrationmwpb.Conds{
			AppID:     &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
			InviterID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.UserID},
			InviteeID: &basetypes.StringVal{Op: cruder.EQ, Value: h.info.UserID},
		},
	)
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("permission denied")
	}
	return nil
}

func (h *Handler) UpdateCommission(ctx context.Context) (*npool.Commission, error) {
	if h.ID == nil {
		return nil, fmt.Errorf("invalid id")
	}
	if h.AppID == nil {
		return nil, fmt.Errorf("invalid appid")
	}

	info, err := commmwcli.GetCommission(ctx, *h.ID)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, fmt.Errorf("invalid commission")
	}
	if info.AppID != *h.AppID {
		return nil, fmt.Errorf("permission denied")
	}

	handler := &updateHandler{
		Handler: h,
		info:    info,
	}
	if err := handler.validateInviter(ctx); err != nil {
		return nil, err
	}

	_, err = commmwcli.UpdateCommission(ctx, &commmwpb.CommissionReq{
		ID:              h.ID,
		SettleType:      h.SettleType,
		StartAt:         h.StartAt,
		AmountOrPercent: h.AmountOrPercent,
		Threshold:       h.Threshold,
	})
	if err != nil {
		return nil, err
	}

	return h.GetCommission(ctx)
}
