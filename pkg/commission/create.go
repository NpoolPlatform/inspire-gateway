package commission

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/go-service-framework/pkg/pubsub"
	"github.com/google/uuid"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	appgoodmwcli "github.com/NpoolPlatform/good-middleware/pkg/client/app/good"
	constant "github.com/NpoolPlatform/inspire-gateway/pkg/const"
	commissionmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/commission"
	registrationmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/invitation/registration"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
	types "github.com/NpoolPlatform/message/npool/basetypes/inspire/v1"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	appgoodmwpb "github.com/NpoolPlatform/message/npool/good/mw/v1/app/good"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/commission"
	commissionmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/commission"
	registrationmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/invitation/registration"

	"github.com/shopspring/decimal"
)

type createHandler struct {
	*Handler
	user       *usermwpb.User
	inviter    *registrationmwpb.Registration
	targetUser *usermwpb.User
	req        *commissionmwpb.CommissionReq
	goodID     *string
}

func (h *createHandler) _getUser(ctx context.Context, userID string) (*usermwpb.User, error) {
	if h.AppID == nil {
		return nil, fmt.Errorf("invalid appid")
	}
	user, err := usermwcli.GetUser(ctx, *h.AppID, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("invalid user")
	}
	return user, nil
}

func (h *createHandler) getUser(ctx context.Context) error {
	if !h.CheckAffiliate {
		return nil
	}
	if h.UserID == nil {
		return fmt.Errorf("invalid userid")
	}
	user, err := h._getUser(ctx, *h.UserID)
	if err != nil {
		return err
	}
	if !user.Kol {
		return fmt.Errorf("permission denied")
	}
	h.user = user
	return nil
}

func (h *createHandler) getTargetUser(ctx context.Context) error {
	if h.TargetUserID == nil {
		return fmt.Errorf("invalid targetuserid")
	}
	user, err := h._getUser(ctx, *h.TargetUserID)
	if err != nil {
		return err
	}
	h.targetUser = user
	return nil
}

func (h *createHandler) validateRegistration(ctx context.Context) error {
	info, err := registrationmwcli.GetRegistrationOnly(ctx, &registrationmwpb.Conds{
		AppID:     &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		InviteeID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.TargetUserID},
	})
	if err != nil {
		return err
	}
	h.inviter = info
	if !h.CheckAffiliate {
		if info != nil {
			h.UserID = &info.InviterID
		}
		return nil
	}
	if h.UserID == nil {
		return fmt.Errorf("permission denied")
	}
	if info == nil {
		return fmt.Errorf("permission denied")
	}
	if info.InviterID != *h.UserID {
		return fmt.Errorf("permission denied")
	}
	return nil
}

func (h *createHandler) validateInviter(ctx context.Context) error {
	if h.AmountOrPercent == nil {
		return fmt.Errorf("invalid amountorpercent")
	}
	commission, err := decimal.NewFromString(*h.AmountOrPercent)
	if err != nil {
		return err
	}
	if commission.Cmp(decimal.NewFromInt(0)) < 0 {
		return fmt.Errorf("invalid amountorpercent")
	}
	if h.inviter == nil {
		if h.CheckAffiliate {
			return fmt.Errorf("permission denied")
		}
		// That means we don't have a inviter
		return nil
	}

	info, err := commissionmwcli.GetCommissionOnly(ctx, &commissionmwpb.Conds{
		AppID:      &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		UserID:     &basetypes.StringVal{Op: cruder.EQ, Value: *h.UserID},
		AppGoodID:  &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppGoodID},
		EndAt:      &basetypes.Uint32Val{Op: cruder.EQ, Value: 0},
		SettleType: &basetypes.Uint32Val{Op: cruder.EQ, Value: uint32(types.SettleType_GoodOrderPayment)},
	})
	if err != nil {
		return nil
	}
	if info == nil {
		return fmt.Errorf("invalid inviter commission")
	}

	commission1, err := decimal.NewFromString(info.AmountOrPercent)
	if err != nil {
		return err
	}
	if commission.Cmp(commission1) > 0 {
		return fmt.Errorf("invalid invitee commission")
	}

	return nil
}

func (h *createHandler) validateInvitees(ctx context.Context) error {
	invitees := []*registrationmwpb.Registration{}
	offset := int32(0)
	limit := constant.DefaultRowLimit

	for {
		_invitees, _, err := registrationmwcli.GetRegistrations(ctx, &registrationmwpb.Conds{
			AppID:     &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
			InviterID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.TargetUserID},
		}, offset, limit)
		if err != nil {
			return err
		}
		if len(_invitees) == 0 {
			break
		}
		invitees = append(invitees, _invitees...)
		offset += limit
	}
	userIDs := []string{}
	for _, invitee := range invitees {
		userIDs = append(userIDs, invitee.InviteeID)
	}
	commissions, _, err := commissionmwcli.GetCommissions(ctx, &commissionmwpb.Conds{
		AppID:      &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		UserIDs:    &basetypes.StringSliceVal{Op: cruder.IN, Value: userIDs},
		AppGoodID:  &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppGoodID},
		EndAt:      &basetypes.Uint32Val{Op: cruder.EQ, Value: 0},
		SettleType: &basetypes.Uint32Val{Op: cruder.EQ, Value: uint32(types.SettleType_GoodOrderPayment)},
	}, 0, int32(len(userIDs)))
	if err != nil {
		return err
	}
	commission, err := decimal.NewFromString(*h.AmountOrPercent)
	if err != nil {
		return err
	}
	for _, _commission := range commissions {
		value, err := decimal.NewFromString(_commission.AmountOrPercent)
		if err != nil {
			return err
		}
		if commission.Cmp(value) < 0 {
			return fmt.Errorf("invalid inviter commission")
		}
	}
	return nil
}

func (h *createHandler) checkGood(ctx context.Context) error {
	if h.AppGoodID == nil {
		return nil
	}

	appgood, err := appgoodmwcli.GetGoodOnly(ctx, &appgoodmwpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		EntID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppGoodID},
	})
	if err != nil {
		return err
	}
	h.goodID = &appgood.GoodID
	return nil
}

func (h *createHandler) createCommission(ctx context.Context) error {
	h.req = &commissionmwpb.CommissionReq{
		EntID:            h.EntID,
		AppID:            h.AppID,
		UserID:           h.TargetUserID,
		GoodID:           h.goodID,
		AppGoodID:        h.AppGoodID,
		SettleType:       h.SettleType,
		SettleMode:       h.SettleMode,
		SettleAmountType: h.SettleAmountType,
		SettleInterval:   h.SettleInterval,
		StartAt:          h.StartAt,
		AmountOrPercent:  h.AmountOrPercent,
		Threshold:        h.Threshold,
	}
	if _, err := commissionmwcli.CreateCommission(ctx, h.req); err != nil {
		return err
	}
	return nil
}

func (h *createHandler) notifyCreateCommission() {
	if err := pubsub.WithPublisher(func(publisher *pubsub.Publisher) error {
		comm := &commissionmwpb.Commission{
			AppID:            *h.AppID,
			EntID:            *h.EntID,
			UserID:           *h.TargetUserID,
			AppGoodID:        *h.AppGoodID,
			SettleType:       *h.SettleType,
			SettleMode:       *h.SettleMode,
			SettleAmountType: *h.SettleAmountType,
			SettleInterval:   *h.SettleInterval,
			StartAt:          *h.StartAt,
			AmountOrPercent:  *h.AmountOrPercent,
		}
		if h.goodID != nil {
			comm.GoodID = *h.goodID
		}
		return publisher.Update(
			basetypes.MsgID_CreateCommissionReq.String(),
			nil,
			nil,
			nil,
			comm,
		)
	}); err != nil {
		logger.Sugar().Errorw(
			"rewardSignup",
			"AppID", h.AppID,
			"UserID", h.TargetUserID,
			"Error", err,
		)
	}
}

func (h *createHandler) validateCommissions(ctx context.Context) error {
	if h.StartAt == nil {
		return nil
	}

	commissions := []*commissionmwpb.Commission{}
	offset := int32(0)
	limit := constant.DefaultRowLimit

	for {
		_commissions, _, err := commissionmwcli.GetCommissions(ctx, &commissionmwpb.Conds{
			AppID:      &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
			UserID:     &basetypes.StringVal{Op: cruder.EQ, Value: *h.TargetUserID},
			GoodID:     &basetypes.StringVal{Op: cruder.EQ, Value: *h.goodID},
			EndAt:      &basetypes.Uint32Val{Op: cruder.NEQ, Value: 0},
			SettleType: &basetypes.Uint32Val{Op: cruder.EQ, Value: uint32(*h.SettleType)},
		}, offset, limit)
		if err != nil {
			return err
		}
		if len(_commissions) == 0 {
			break
		}
		commissions = append(commissions, _commissions...)
		offset += limit
	}
	for _, commission := range commissions {
		if commission.EndAt > *h.StartAt {
			return fmt.Errorf("invalid startat")
		}
	}
	return nil
}

func (h *Handler) CreateCommission(ctx context.Context) (*npool.Commission, error) {
	id := uuid.NewString()
	if h.EntID == nil {
		h.EntID = &id
	}

	handler := &createHandler{
		Handler: h,
	}
	if err := handler.getUser(ctx); err != nil {
		return nil, err
	}
	if err := handler.getTargetUser(ctx); err != nil {
		return nil, err
	}
	if err := handler.validateRegistration(ctx); err != nil {
		return nil, err
	}
	if err := handler.checkGood(ctx); err != nil {
		return nil, err
	}
	if err := handler.validateInviter(ctx); err != nil {
		return nil, err
	}
	if err := handler.validateInvitees(ctx); err != nil {
		return nil, err
	}
	if err := handler.validateCommissions(ctx); err != nil {
		return nil, err
	}
	if err := handler.createCommission(ctx); err != nil {
		return nil, err
	}
	handler.notifyCreateCommission()

	return h.GetCommission(ctx)
}
