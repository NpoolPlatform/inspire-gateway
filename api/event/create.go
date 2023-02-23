//nolint:dupl
package event

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/event"
	alloccoupmgrpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/coupon/allocated"
	mgrpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/event"

	event1 "github.com/NpoolPlatform/inspire-gateway/pkg/event"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

//nolint
func (s *Server) CreateEvent(ctx context.Context, in *npool.CreateEventRequest) (*npool.CreateEventResponse, error) {
	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorw("CreateEvent", "AppID", in.GetAppID(), "Error", err)
		return &npool.CreateEventResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	switch in.GetEventType() {
	case basetypes.UsedFor_Signup:
	case basetypes.UsedFor_Signin:
	case basetypes.UsedFor_Update:
	case basetypes.UsedFor_Contact:
	case basetypes.UsedFor_SetWithdrawAddress:
	case basetypes.UsedFor_Withdraw:
	case basetypes.UsedFor_CreateInvitationCode:
	case basetypes.UsedFor_SetCommission:
	case basetypes.UsedFor_SetTransferTargetUser:
	case basetypes.UsedFor_WithdrawalRequest:
	case basetypes.UsedFor_WithdrawalCompleted:
	case basetypes.UsedFor_DepositReceived:
	case basetypes.UsedFor_KYCApproved:
	case basetypes.UsedFor_KYCRejected:
	case basetypes.UsedFor_Purchase:
		if _, err := uuid.Parse(in.GetGoodID()); err != nil {
			logger.Sugar().Errorw("ValidateCreate", "GoodID", in.GetGoodID(), "Error", err)
			return &npool.CreateEventResponse{}, status.Error(codes.InvalidArgument, err.Error())
		}
	default:
		logger.Sugar().Errorw("CreateEvent", "EventType", in.GetEventType())
		return &npool.CreateEventResponse{}, status.Error(codes.InvalidArgument, "EventType is invalid")
	}
	for _, coupon := range in.GetCoupons() {
		if _, err := uuid.Parse(coupon.GetID()); err != nil {
			logger.Sugar().Errorw("ValidateCreate", "Coupons", in.GetCoupons(), "Error", err)
			return &npool.CreateEventResponse{}, status.Error(codes.InvalidArgument, err.Error())
		}
		switch coupon.GetCouponType() {
		case alloccoupmgrpb.CouponType_FixAmount:
		case alloccoupmgrpb.CouponType_Discount:
		case alloccoupmgrpb.CouponType_SpecialOffer:
		case alloccoupmgrpb.CouponType_ThresholdFixAmount:
		case alloccoupmgrpb.CouponType_ThresholdDiscount:
		case alloccoupmgrpb.CouponType_GoodFixAmount:
		case alloccoupmgrpb.CouponType_GoodDiscount:
		case alloccoupmgrpb.CouponType_GoodThresholdFixAmount:
		case alloccoupmgrpb.CouponType_GoodThresholdDiscount:
		default:
			logger.Sugar().Errorw("ValidateCreate", "Coupons", in.GetCoupons())
			return &npool.CreateEventResponse{}, fmt.Errorf("coupontype is invalid")
		}
	}
	if _, err := decimal.NewFromString(in.GetCredits()); err != nil {
		logger.Sugar().Errorw("CreateEvent", "Credits", in.GetCredits(), "Error", err)
		return &npool.CreateEventResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if _, err := decimal.NewFromString(in.GetCreditsPerUSD()); err != nil {
		logger.Sugar().Errorw("CreateEvent", "CreditsPerUSD", in.GetCreditsPerUSD(), "Error", err)
		return &npool.CreateEventResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	req := &mgrpb.EventReq{
		AppID:         &in.AppID,
		EventType:     &in.EventType,
		Coupons:       in.Coupons,
		Credits:       in.Credits,
		CreditsPerUSD: in.CreditsPerUSD,
	}

	info, err := event1.CreateEvent(ctx, req)
	if err != nil {
		logger.Sugar().Errorw("CreateEvent", "Req", req, "Error", err)
		return &npool.CreateEventResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateEventResponse{
		Info: info,
	}, nil
}
