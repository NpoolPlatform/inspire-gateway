//nolint:dupl
package event

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"

	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/event"
	alloccoupmgrpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/coupon/allocated"
	mgrpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/event"

	event1 "github.com/NpoolPlatform/inspire-gateway/pkg/event"

	mgrcli "github.com/NpoolPlatform/inspire-manager/pkg/client/event"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func (s *Server) UpdateEvent(ctx context.Context, in *npool.UpdateEventRequest) (*npool.UpdateEventResponse, error) { //nolint
	if _, err := uuid.Parse(in.GetID()); err != nil {
		logger.Sugar().Errorw("UpdateEvent", "ID", in.GetID(), "Error", err)
		return &npool.UpdateEventResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorw("UpdateEvent", "AppID", in.GetAppID(), "Error", err)
		return &npool.UpdateEventResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	for _, coupon := range in.GetCoupons() {
		if _, err := uuid.Parse(coupon.GetID()); err != nil {
			logger.Sugar().Errorw("ValidateUpdate", "Coupons", in.GetCoupons(), "Error", err)
			return &npool.UpdateEventResponse{}, status.Error(codes.InvalidArgument, err.Error())
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
			logger.Sugar().Errorw("ValidateUpdate", "Coupons", in.GetCoupons())
			return &npool.UpdateEventResponse{}, fmt.Errorf("coupontype is invalid")
		}
	}
	if _, err := decimal.NewFromString(in.GetCredits()); err != nil {
		logger.Sugar().Errorw("UpdateEvent", "Credits", in.GetCredits(), "Error", err)
		return &npool.UpdateEventResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if in.Credits != nil {
		if _, err := decimal.NewFromString(in.GetCredits()); err != nil {
			logger.Sugar().Errorw("UpdateEvent", "Credits", in.GetCredits(), "Error", err)
			return &npool.UpdateEventResponse{}, status.Error(codes.InvalidArgument, err.Error())
		}
	}
	if in.CreditsPerUSD != nil {
		if _, err := decimal.NewFromString(in.GetCreditsPerUSD()); err != nil {
			logger.Sugar().Errorw("UpdateEvent", "CreditsPerUSD", in.GetCreditsPerUSD(), "Error", err)
			return &npool.UpdateEventResponse{}, status.Error(codes.InvalidArgument, err.Error())
		}
	}

	exist, err := mgrcli.ExistEventConds(ctx, &mgrpb.Conds{
		ID:    &basetypes.StringVal{Op: cruder.EQ, Value: in.GetID()},
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: in.GetAppID()},
	})
	if err != nil {
		logger.Sugar().Errorw("UpdateEvent", "Error", err)
		return &npool.UpdateEventResponse{}, status.Error(codes.Internal, err.Error())
	}
	if !exist {
		logger.Sugar().Errorw("UpdateEvent", "ID", in.GetID(), "AppID", in.GetAppID())
		return &npool.UpdateEventResponse{}, status.Error(codes.Internal, "Event is invalid")
	}

	req := &mgrpb.EventReq{
		ID:            &in.ID,
		AppID:         &in.AppID,
		Coupons:       in.Coupons,
		Credits:       in.Credits,
		CreditsPerUSD: in.CreditsPerUSD,
	}

	info, err := event1.UpdateEvent(ctx, req)
	if err != nil {
		logger.Sugar().Errorw("UpdateEvent", "Req", req, "Error", err)
		return &npool.UpdateEventResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UpdateEventResponse{
		Info: info,
	}, nil
}
