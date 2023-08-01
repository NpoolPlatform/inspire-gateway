//nolint:dupl
package event

import (
	"context"
	"fmt"
	"time"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	timedef "github.com/NpoolPlatform/go-service-framework/pkg/const/time"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"

	commonpb "github.com/NpoolPlatform/message/npool"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	appgoodmgrpb "github.com/NpoolPlatform/message/npool/good/mgr/v1/appgood"
	npool "github.com/NpoolPlatform/message/npool/inspire/gw/v1/event"
	alloccoupmgrpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/coupon/allocated"
	mgrpb "github.com/NpoolPlatform/message/npool/inspire/mgr/v1/event"

	event1 "github.com/NpoolPlatform/inspire-gateway/pkg/event"

	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	appgoodmwcli "github.com/NpoolPlatform/good-middleware/pkg/client/appgood"
	mgrcli "github.com/NpoolPlatform/inspire-manager/pkg/client/event"
	coupmwcli "github.com/NpoolPlatform/inspire-middleware/pkg/client/coupon"

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

	app, err := appmwcli.GetApp(ctx, in.GetAppID())
	if err != nil {
		logger.Sugar().Errorw("CreateEvent", "AppID", in.GetAppID(), "Error", err)
		return &npool.CreateEventResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if app == nil {
		logger.Sugar().Errorw("CreateEvent", "AppID", in.GetAppID())
		return &npool.CreateEventResponse{}, status.Error(codes.InvalidArgument, "App not exist")
	}

	conds := &mgrpb.Conds{
		AppID: &basetypes.StringVal{Op: cruder.EQ, Value: in.GetAppID()},
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
		fallthrough //nolint
	case basetypes.UsedFor_AffiliatePurchase:
		if _, err := uuid.Parse(in.GetGoodID()); err != nil {
			logger.Sugar().Errorw("ValidateCreate", "GoodID", in.GetGoodID(), "Error", err)
			return &npool.CreateEventResponse{}, status.Error(codes.InvalidArgument, err.Error())
		}
		conds.GoodID = &basetypes.StringVal{Op: cruder.EQ, Value: in.GetGoodID()}

		good, err := appgoodmwcli.GetGoodOnly(ctx, &appgoodmgrpb.Conds{
			AppID:  &commonpb.StringVal{Op: cruder.EQ, Value: in.GetAppID()},
			GoodID: &commonpb.StringVal{Op: cruder.EQ, Value: in.GetGoodID()},
		})
		if err != nil {
			logger.Sugar().Errorw("ValidateCreate", "GoodID", in.GetGoodID(), "Error", err)
			return &npool.CreateEventResponse{}, status.Error(codes.InvalidArgument, err.Error())
		}
		if good == nil {
			logger.Sugar().Errorw("ValidateCreate", "GoodID", in.GetGoodID())
			return &npool.CreateEventResponse{}, status.Error(codes.InvalidArgument, "Good not exist")
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

		coup, err := coupmwcli.GetCoupon(ctx, coupon.ID, coupon.CouponType)
		if err != nil {
			logger.Sugar().Errorw("ValidateCreate", "Coupons", in.GetCoupons(), "Error", err)
			return &npool.CreateEventResponse{}, status.Error(codes.InvalidArgument, err.Error())
		}
		if coup == nil {
			logger.Sugar().Errorw("ValidateCreate", "Coupons", in.GetCoupons(), "Error", err)
			return &npool.CreateEventResponse{}, status.Error(codes.InvalidArgument, "Coupon not exist")
		}
		now := uint32(time.Now().Unix())
		if now < coup.StartAt || coup.StartAt+coup.DurationDays*timedef.SecondsPerDay < now {
			logger.Sugar().Errorw("ValidateCreate", "Coupons", in.GetCoupons(), "Error", err)
			return &npool.CreateEventResponse{}, status.Error(codes.InvalidArgument, "Coupon invalid")
		}
	}
	if in.Credits != nil {
		if _, err := decimal.NewFromString(in.GetCredits()); err != nil {
			logger.Sugar().Errorw("CreateEvent", "Credits", in.GetCredits(), "Error", err)
			return &npool.CreateEventResponse{}, status.Error(codes.InvalidArgument, err.Error())
		}
	}
	if in.CreditsPerUSD != nil {
		if _, err := decimal.NewFromString(in.GetCreditsPerUSD()); err != nil {
			logger.Sugar().Errorw("CreateEvent", "CreditsPerUSD", in.GetCreditsPerUSD(), "Error", err)
			return &npool.CreateEventResponse{}, status.Error(codes.InvalidArgument, err.Error())
		}
	}

	conds.EventType = &basetypes.Uint32Val{Op: cruder.EQ, Value: uint32(in.GetEventType())}
	exist, err := mgrcli.ExistEventConds(ctx, conds)
	if err != nil {
		logger.Sugar().Errorw("CreateEvent", "Conds", conds, "Error", err)
		return &npool.CreateEventResponse{}, status.Error(codes.Internal, err.Error())
	}
	if exist {
		logger.Sugar().Errorw("CreateEvent", "Conds", conds, "Exist", exist)
		return &npool.CreateEventResponse{}, status.Error(codes.AlreadyExists, "Event already exist")
	}

	req := &mgrpb.EventReq{
		AppID:          &in.AppID,
		EventType:      &in.EventType,
		Coupons:        in.Coupons,
		Credits:        in.Credits,
		CreditsPerUSD:  in.CreditsPerUSD,
		GoodID:         in.GoodID,
		MaxConsecutive: in.MaxConsecutive,
		InviterLayers:  in.InviterLayers,
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
