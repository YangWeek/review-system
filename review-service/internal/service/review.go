package service

import (
	"context"
	"fmt"
	"review-service/internal/biz"
	"review-service/internal/data/model"

	pb "review-service/api/review/v1"
)

type ReviewService struct {
	pb.UnimplementedReviewServer

	uc *biz.ReviewUsecase
}

func NewReviewService(uc *biz.ReviewUsecase) *ReviewService {
	return &ReviewService{uc: uc}
}

// 创建评论
func (s *ReviewService) CreateReview(ctx context.Context, req *pb.CreateReviewRequest) (*pb.CreateReviewReply, error) {
	fmt.Printf("[service] CreateReview, req:%#v\n", req)
	// 参数转换
	var anonymous int32
	// 是否匿名 1 为表示匿名
	if req.Anonymous {
		anonymous = 1
	}
	// 调用biz层
	review, err := s.uc.CreateReview(ctx, &model.ReviewInfo{
		UserID:       req.UserID,
		OrderID:      req.OrderID,
		Score:        req.Score,
		ServiceScore: req.Score,
		ExpressScore: req.ExpressScore,
		Content:      req.Content,
		PicInfo:      req.PicInfo,
		VideoInfo:    req.VideoInfo,
		Anonymous:    anonymous,
		Status:       0,
	})
	// 拼装返回结果
	if err != nil {
		return nil, err
	}

	return &pb.CreateReviewReply{ReviewID: review.ReviewID}, err
}

// GetReview C端获取评价详情
func (s *ReviewService) GetReview(ctx context.Context, req *pb.GetReviewRequest) (*pb.GetReviewReply, error) {
	fmt.Printf("GetReview req:%#v\n", req)
	review, err := s.uc.GetReview(ctx, req.ReviewID)
	return &pb.GetReviewReply{
		Data: &pb.ReviewInfo{
			ReviewID:     review.ReviewID,
			UserID:       review.UserID,
			OrderID:      review.OrderID,
			Score:        review.Score,
			ServiceScore: review.ServiceScore,
			ExpressScore: review.ExpressScore,
			Content:      review.Content,
			PicInfo:      review.PicInfo,
			VideoInfo:    review.VideoInfo,
			Status:       review.Status,
		},
	}, err
}

// AuditReview 审计评论
func (s *ReviewService) AuditReview(ctx context.Context, req *pb.AuditReviewRequest) (*pb.AuditReviewReply, error) {
	fmt.Printf("AuditReview req:%#v\n", req)
	err := s.uc.AuditReview(ctx, &biz.AuditParam{
		ReviewID:  req.GetReviewID(),
		OpUser:    req.GetOpUser(),
		OpReason:  req.GetOpReason(),
		OpRemarks: req.GetOpRemarks(),
		Status:    req.GetStatus(),
	})
	if err != nil {
		return nil, err
	}
	return &pb.AuditReviewReply{
		ReviewID: req.ReviewID,
		Status:   req.Status,
	}, nil
}

// ReplyReview 回复评价
func (s *ReviewService) ReplyReview(ctx context.Context, req *pb.ReplyReviewRequest) (*pb.ReplyReviewReply, error) {
	fmt.Printf("[service] ReplyReview req:%#v\n", req)
	// 调用biz层
	reply, err := s.uc.CreateReply(ctx, &biz.ReplyParam{
		ReviewID:  req.GetReviewID(),
		StoreID:   req.GetStoreID(),
		Content:   req.GetContent(),
		PicInfo:   req.GetPicInfo(),
		VideoInfo: req.GetVideoInfo(),
	})
	if err != nil {
		return nil, err
	}
	return &pb.ReplyReviewReply{ReplyID: reply.ReplyID}, nil
}

// AppealReview 上诉审查
func (s *ReviewService) AppealReview(ctx context.Context, req *pb.AppealReviewRequest) (*pb.AppealReviewReply, error) {
	fmt.Printf("[service] AppealReview req:%#v\n", req)
	ret, err := s.uc.AppealReview(ctx, &biz.AppealParam{
		ReviewID:  req.GetReviewID(),
		StoreID:   req.GetStoreID(),
		Reason:    req.GetReason(),
		Content:   req.GetContent(),
		PicInfo:   req.GetPicInfo(),
		VideoInfo: req.GetVideoInfo(),
	})
	if err != nil {
		return nil, err
	}
	fmt.Printf("[service] AppealReview ret:%v err:%v\n", ret, err)
	return &pb.AppealReviewReply{AppealID: ret.AppealID}, nil
}

func (s *ReviewService) AuditAppeal(ctx context.Context, req *pb.AuditAppealRequest) (*pb.AuditAppealReply, error) {
	fmt.Printf("[service] AuditAppeal req:%#v\n", req)
	err := s.uc.AuditAppeal(ctx, &biz.AuditAppealParam{
		ReviewID: req.GetReviewID(),
		AppealID: req.GetAppealID(),
		OpUser:   req.GetOpUser(),
		Status:   req.GetStatus(),
	})
	if err != nil {
		return nil, err
	}
	return &pb.AuditAppealReply{}, nil
}

func (s *ReviewService) ListReviewByUserID(ctx context.Context, req *pb.ListReviewByUserIDRequest) (*pb.ListReviewByUserIDReply, error) {
	fmt.Printf("[service] ListReviewByUserID req:%#v\n", req)
	dataList, err := s.uc.ListReviewByUserID(ctx, req.GetUserID(), int(req.GetPage()), int(req.GetSize()))
	if err != nil {
		return nil, err
	}
	list := make([]*pb.ReviewInfo, 0, len(dataList))
	for _, review := range dataList {
		list = append(list, &pb.ReviewInfo{
			ReviewID:     review.ReviewID,
			UserID:       review.UserID,
			OrderID:      review.OrderID,
			Score:        review.Score,
			ServiceScore: review.ServiceScore,
			ExpressScore: review.ExpressScore,
			Content:      review.Content,
			PicInfo:      review.PicInfo,
			VideoInfo:    review.VideoInfo,
			Status:       review.Status,
		})
	}
	return &pb.ListReviewByUserIDReply{List: list}, nil
}

func (s *ReviewService) ListReviewByStoreID(ctx context.Context, req *pb.ListReviewByStoreIDRequest) (*pb.ListReviewByStoreIDReply, error) {
	fmt.Printf("[service] ListReviewByStoreID req:%#v\n", req)
	reviewList, err := s.uc.ListReviewByStoreID(ctx, req.StoreID, int(req.Page), int(req.Size))
	if err != nil {
		return nil, err
	}
	// format
	list := make([]*pb.ReviewInfo, 0, len(reviewList))
	for _, r := range reviewList {
		list = append(list, &pb.ReviewInfo{
			ReviewID:     r.ReviewID,
			UserID:       r.UserID,
			OrderID:      r.OrderID,
			Score:        r.Score,
			ServiceScore: r.ServiceScore,
			ExpressScore: r.ExpressScore,
			Content:      r.Content,
			PicInfo:      r.PicInfo,
			VideoInfo:    r.VideoInfo,
			Status:       r.Status,
		})
	}

	return &pb.ListReviewByStoreIDReply{List: list}, nil
}
