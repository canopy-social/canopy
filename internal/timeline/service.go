package timeline

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sumi-devs/canopy-social/canopy/internal/posts"
	"github.com/sumi-devs/canopy-social/canopy/pkg/config"
)

type Service struct {
	redis           *redis.Client
	cfg             *config.Config
	postRepo        PostRepository
	accountRepo     AccountRepository
	celebrityThresh int
}

func NewService(r *redis.Client, cfg *config.Config, postRepo PostRepository, accountRepo AccountRepository, celebrityThresh int) *Service {
	return &Service{
		redis:           r,
		cfg:             cfg,
		postRepo:        postRepo,
		accountRepo:     accountRepo,
		celebrityThresh: celebrityThresh,
	}
}

func (s *Service) FanoutPost(ctx context.Context, post *posts.Post) error {
	author, err := s.accountRepo.GetByID(ctx, post.AccountID)
	if err != nil {
		return err
	}

	score := float64(post.CreatedAt.UnixMilli())
	member := post.ID

	s.redis.ZAdd(ctx, "timeline:home:"+post.AccountID, redis.Z{
		Score:  score,
		Member: member,
	})
	s.redis.ZRemRangeByRank(ctx, "timeline:home:"+post.AccountID, 0, -801)

	if author.FollowersCount >= s.celebrityThresh {
		return nil
	}

	limit := 100
	offset := 0
	for {
		followers, err := s.accountRepo.ListFollowers(ctx, post.AccountID, limit, offset)
		if err != nil {
			return err
		}
		if len(followers) == 0 {
			break
		}

		pipe := s.redis.Pipeline()
		for _, f := range followers {
			key := "timeline:home:" + f.ID
			pipe.ZAdd(ctx, key, redis.Z{
				Score:  score,
				Member: member,
			})
			pipe.ZRemRangeByRank(ctx, key, 0, -801)
		}
		_, err = pipe.Exec(ctx)
		if err != nil {
			return err
		}

		if len(followers) < limit {
			break
		}
		offset += limit
	}

	return nil
}

func (s *Service) GetHomeTimeline(ctx context.Context, accountID string, limit int, maxID, sinceID string) (*TimelineResponse, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	key := "timeline:home:" + accountID

	var maxScore float64 = -1
	var sinceScore float64 = -1

	if maxID != "" {
		p, err := s.postRepo.GetByID(ctx, maxID)
		if err == nil && p != nil {
			maxScore = float64(p.CreatedAt.UnixMilli())
		}
	}
	if sinceID != "" {
		p, err := s.postRepo.GetByID(ctx, sinceID)
		if err == nil && p != nil {
			sinceScore = float64(p.CreatedAt.UnixMilli())
		}
	}

	var postIDs []string

	if maxScore > 0 && sinceScore > 0 {
		res, err := s.redis.ZRevRangeByScore(ctx, key, &redis.ZRangeBy{
			Min:    fmt.Sprintf("(%f", sinceScore),
			Max:    fmt.Sprintf("(%f", maxScore),
			Offset: 0,
			Count:  int64(limit),
		}).Result()
		if err == nil {
			postIDs = res
		}
	} else if maxScore > 0 {
		res, err := s.redis.ZRevRangeByScore(ctx, key, &redis.ZRangeBy{
			Min:    "-inf",
			Max:    fmt.Sprintf("(%f", maxScore),
			Offset: 0,
			Count:  int64(limit),
		}).Result()
		if err == nil {
			postIDs = res
		}
	} else if sinceScore > 0 {
		res, err := s.redis.ZRevRangeByScore(ctx, key, &redis.ZRangeBy{
			Min:    fmt.Sprintf("(%f", sinceScore),
			Max:    "+inf",
			Offset: 0,
			Count:  int64(limit),
		}).Result()
		if err == nil {
			postIDs = res
		}
	} else {
		res, err := s.redis.ZRevRangeByScore(ctx, key, &redis.ZRangeBy{
			Min:    "-inf",
			Max:    "+inf",
			Offset: 0,
			Count:  int64(limit),
		}).Result()
		if err == nil {
			postIDs = res
		}
	}

	var homePosts []*posts.Post
	for _, id := range postIDs {
		p, err := s.postRepo.GetByID(ctx, id)
		if err == nil && p != nil {
			homePosts = append(homePosts, p)
		}
	}

	var celebrityPosts []*posts.Post
	offset := 0
	batchSize := 100
	for {
		following, err := s.accountRepo.ListFollowing(ctx, accountID, batchSize, offset)
		if err != nil {
			break
		}
		if len(following) == 0 {
			break
		}

		for _, f := range following {
			if f.FollowersCount >= s.celebrityThresh {
				cPosts, err := s.postRepo.ListPostsByAccountWithBoosts(ctx, f.ID, limit, 0)
				if err == nil {
					for _, cp := range cPosts {
						t := cp.CreatedAt.UnixMilli()
						if maxScore > 0 && float64(t) >= maxScore {
							continue
						}
						if sinceScore > 0 && float64(t) <= sinceScore {
							continue
						}
						celebrityPosts = append(celebrityPosts, cp)
					}
				}
			}
		}

		if len(following) < batchSize {
			break
		}
		offset += batchSize
	}

	merged := append(homePosts, celebrityPosts...)
	uniqueMap := make(map[string]*posts.Post)
	for _, p := range merged {
		uniqueMap[p.ID] = p
	}

	var uniquePosts []*posts.Post
	for _, p := range uniqueMap {
		uniquePosts = append(uniquePosts, p)
	}

	sort.Slice(uniquePosts, func(i, j int) bool {
		return uniquePosts[i].CreatedAt.After(uniquePosts[j].CreatedAt)
	})

	if len(uniquePosts) > limit {
		uniquePosts = uniquePosts[:limit]
	}

	nextCursor := ""
	prevCursor := ""
	if len(uniquePosts) > 0 {
		nextCursor = uniquePosts[len(uniquePosts)-1].ID
		prevCursor = uniquePosts[0].ID
	}

	return &TimelineResponse{
		Data:       uniquePosts,
		NextCursor: nextCursor,
		PrevCursor: prevCursor,
	}, nil
}

func (s *Service) GetPublicTimeline(ctx context.Context, local bool, limit, offset int) ([]*posts.Post, error) {
	cacheKey := fmt.Sprintf("cache:public_timeline:%t:%d:%d", local, limit, offset)

	cachedData, err := s.redis.Get(ctx, cacheKey).Bytes()
	if err == nil {
		var cachedPosts []*posts.Post
		if err := json.Unmarshal(cachedData, &cachedPosts); err == nil {
			return cachedPosts, nil
		}
	}

	resPosts, err := s.postRepo.ListPublicTimeline(ctx, local, limit, offset)
	if err != nil {
		return nil, err
	}

	if data, err := json.Marshal(resPosts); err == nil {
		s.redis.Set(ctx, cacheKey, data, 30*time.Second)
	}

	return resPosts, nil
}
