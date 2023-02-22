// nolint
package integration_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	faker "github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/apabramov/hw-test/hw12_13_14_15_calendar/internal/server/pb"
)

type TestSuite struct {
	suite.Suite
	client pb.EventServiceClient
}

func TestIntegration(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) SetupSuite() {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	url := os.Getenv("GRPC_URL")
	if url == "" {
		url = "0.0.0.0:12000"
	}
	conn, err := grpc.Dial(url, opts...)
	s.Require().NoError(err)
	s.client = pb.NewEventServiceClient(conn)
}

func (s *TestSuite) TestAddEvent() {
	ctx := context.Background()
	id := faker.UUIDHyphenated()
	_, err := s.client.Add(ctx, &pb.EventRequest{
		Event: &pb.Event{
			ID:     id,
			UserId: faker.UUIDHyphenated(),
		},
	})
	s.Require().NoError(err)

	_, err = s.client.Del(ctx, &pb.IDRequest{
		ID: id,
	})
	s.Require().NoError(err)
}

func (s *TestSuite) TestUpdateEvent() {
	id := faker.UUIDHyphenated()
	userId := faker.UUIDHyphenated()
	ctx := context.Background()

	_, err := s.client.Add(ctx, &pb.EventRequest{
		Event: &pb.Event{
			ID:     id,
			UserId: userId,
			Title:  "Event add",
		},
	})
	s.Require().NoError(err)

	_, err = s.client.Update(ctx, &pb.EventRequest{
		Event: &pb.Event{
			ID:     id,
			UserId: userId,
			Title:  "Event update",
		},
	})
	s.Require().NoError(err)

	e, err := s.client.Get(ctx, &pb.IDRequest{ID: id})
	s.Require().NoError(err)
	s.Require().Equal("Event update", e.Event.GetTitle())

	_, err = s.client.Del(ctx, &pb.IDRequest{
		ID: id,
	})
	s.Require().NoError(err)
}

func (s *TestSuite) TestDeleteEvent() {
	id := faker.UUIDHyphenated()
	userId := faker.UUIDHyphenated()
	ctx := context.Background()

	_, err := s.client.Add(ctx, &pb.EventRequest{
		Event: &pb.Event{
			ID:     id,
			UserId: userId,
			Title:  "Event add",
		},
	})
	s.Require().NoError(err)

	_, err = s.client.Del(ctx, &pb.IDRequest{
		ID: id,
	})
	s.Require().NoError(err)

	_, err = s.client.Get(ctx, &pb.IDRequest{ID: id})
	s.Require().Error(err)
}

func (s *TestSuite) TestListDayEvent() {
	id := faker.UUIDHyphenated()
	userId := faker.UUIDHyphenated()
	ctx := context.Background()

	dt := time.Now()

	_, err := s.client.Add(ctx, &pb.EventRequest{
		Event: &pb.Event{
			ID:     id,
			UserId: userId,
			Title:  "Event add",
			Date:   timestamppb.New(dt.Add(time.Second * 10)),
		},
	})
	s.Require().NoError(err)

	l, err := s.client.ListByDay(ctx, &pb.ListRequest{Bg: timestamppb.New(dt), Fn: timestamppb.New(dt.AddDate(0, 0, 1))})

	s.Require().NoError(err)
	m := l.GetEvents()
	fmt.Println(m)

	s.Require().True(len(m) == 1)

	s.Require().Equal("Event add", m[0].Title)

	_, err = s.client.Del(ctx, &pb.IDRequest{
		ID: id,
	})
	s.Require().NoError(err)
}

func (s *TestSuite) TestListWeekEvent() {
	id := faker.UUIDHyphenated()
	userId := faker.UUIDHyphenated()
	ctx := context.Background()

	dt := time.Now()

	_, err := s.client.Add(ctx, &pb.EventRequest{
		Event: &pb.Event{
			ID:     id,
			UserId: userId,
			Title:  "Event add",
			Date:   timestamppb.New(dt.Add(time.Second * 10)),
		},
	})
	s.Require().NoError(err)

	l, err := s.client.ListByWeek(ctx, &pb.ListRequest{Bg: timestamppb.New(dt), Fn: timestamppb.New(dt.AddDate(0, 0, 1))})

	s.Require().NoError(err)
	m := l.GetEvents()
	s.Require().True(len(m) == 1)

	s.Require().Equal("Event add", m[0].Title)

	_, err = s.client.Del(ctx, &pb.IDRequest{
		ID: id,
	})
	s.Require().NoError(err)
}

func (s *TestSuite) TestListMonthEvent() {
	id := faker.UUIDHyphenated()
	userId := faker.UUIDHyphenated()
	ctx := context.Background()

	dt := time.Now()

	_, err := s.client.Add(ctx, &pb.EventRequest{
		Event: &pb.Event{
			ID:     id,
			UserId: userId,
			Title:  "Event add",
			Date:   timestamppb.New(dt.Add(time.Second * 10)),
		},
	})
	s.Require().NoError(err)

	l, err := s.client.ListByMonth(ctx, &pb.ListRequest{Bg: timestamppb.New(dt), Fn: timestamppb.New(dt.AddDate(0, 0, 1))})

	s.Require().NoError(err)
	m := l.GetEvents()
	s.Require().True(len(m) == 1)

	s.Require().Equal("Event add", m[0].Title)

	_, err = s.client.Del(ctx, &pb.IDRequest{
		ID: id,
	})
	s.Require().NoError(err)
}

//func (s *TestSuite) TestNotifyEvent() {
//	id := faker.UUIDHyphenated()
//	userId := faker.UUIDHyphenated()
//	ctx := context.Background()
//	dt := time.Now()
//
//	_, err := s.client.Add(ctx, &pb.EventRequest{
//		Event: &pb.Event{
//			ID:     id,
//			UserId: userId,
//			Title:  "Event add",
//			Date:   timestamppb.New(dt.Add(time.Second * 60)),
//			Notify: durationpb.New(time.Second * 55),
//		},
//	})
//	s.Require().NoError(err)
//	time.Sleep(time.Second * 120)
//	e, err := s.client.Get(ctx, &pb.IDRequest{ID: id})
//
//	s.Require().True(e.Event.GetSent())
//	s.Require().NoError(err)
//}
