package handler

import (
	"context"
	"errors"
	"fmt"
	pbtask "task-service/gen/go/task-service/task"
	apptask "task-service/internal/application/task"
	domaintask "task-service/internal/domain/task"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TaskHandler struct {
	pbtask.UnimplementedTaskServiceServer

	getHandler         apptask.GetHandler
	searchHandler      apptask.SearchHandler
	completeHandler    apptask.CompleteHandler
	createHandler      apptask.CreateHandler
	updateHandler      apptask.UpdateHandler
	deleteHandler      apptask.DeleteHandler
	restoreHandler     apptask.RestoreHandler
	listDeletedHandler apptask.ListDeletedHandler
}

func extractRequesterIDFromCtx(ctx context.Context) (uuid.UUID, error) {
	unparsedID, ok := ctx.Value("requester-id").(string)
	if !ok {
		return uuid.Nil, errors.New("missing requester id")
	}
	id, err := uuid.Parse(unparsedID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("requester id: %w", err)
	}
	return id, nil
}

func fromDomainNotifyAt(notifyAt []time.Time) []*timestamppb.Timestamp {
	pbNotifyAt := make([]*timestamppb.Timestamp, len(notifyAt))
	for i, notifyTime := range notifyAt {
		pbNotifyAt[i] = timestamppb.New(notifyTime)
	}
	return pbNotifyAt
}

func fromDomainGeoPoint(point domaintask.GeoPoint) *pbtask.GeoPoint {
	return &pbtask.GeoPoint{
		Latitude:  new(point.Latitude),
		Longitude: new(point.Longitude),
	}
}

func toProtoGetTaskResponse(task *domaintask.Task) *pbtask.GetTaskResponse {
	var groupID *string
	if value, ok := task.GroupID.Get(); ok {
		groupID = new(value.String())
	}

	var location *pbtask.GeoPoint
	if value, ok := task.Location.Get(); ok {
		location = fromDomainGeoPoint(value)
	}

	var deadline *timestamppb.Timestamp
	if value, ok := task.Deadline.Get(); ok {
		deadline = timestamppb.New(value)
	}

	var purgeAt *timestamppb.Timestamp
	if value, ok := task.PurgeAt.Get(); ok {
		purgeAt = timestamppb.New(value)
	}

	notifyAt := fromDomainNotifyAt(task.NotifyAt)

	var completedAt *timestamppb.Timestamp
	if value, ok := task.CompletedAt.Get(); ok {
		completedAt = timestamppb.New(value)
	}

	return &pbtask.GetTaskResponse{
		Id:          new(task.ID.String()),
		OwnerId:     new(task.OwnerID.String()),
		GroupId:     groupID,
		Title:       new(task.Title),
		Description: new(task.Description),
		Location:    location,
		IsFavorite:  new(task.IsFavorite),
		Priority:    new(pbtask.TaskPriority(task.Priority)),
		Icon:        new(pbtask.TaskIcon(task.Icon)),
		Status:      new(pbtask.TaskStatus(task.Status)),
		Deadline:    deadline,
		PurgeAt:     purgeAt,
		NotifyAt:    notifyAt,
		CompletedAt: completedAt,
	}
}

func fromProtoGeoPoint(point *pbtask.GeoPoint) domaintask.GeoPoint {
	res, _ := domaintask.NewGeoPoint(*point.Latitude, *point.Longitude)
	return res
}

func fromProtoNotifyAt(notifyAt []*timestamppb.Timestamp) []time.Time {
	res := make([]time.Time, len(notifyAt))
	for i, notifyTime := range notifyAt {
		res[i] = notifyTime.AsTime()
	}
	return res
}

func fromDomainUUIDs(ids []uuid.UUID) []string {
	strIDs := make([]string, len(ids))
	for i, id := range ids {
		strIDs[i] = id.String()
	}
	return strIDs
}

func NewSearchTasksResponse(ids []uuid.UUID, key time.Time) *pbtask.SearchTasksResponse {
	return &pbtask.SearchTasksResponse{
		Ids:         fromDomainUUIDs(ids),
		NextSyncKey: timestamppb.New(key),
	}
}

func toGeoFilterQuery(geoFilter *pbtask.GeoFilter) (apptask.GeoFilterQuery, error) {
	return apptask.NewGeoFilterQuery(
		geoFilter.GetLocation().GetLatitude(),
		geoFilter.GetLocation().GetLongitude(),
		geoFilter.GetRadius(),
	)
}

func (h *TaskHandler) GetTask(ctx context.Context, req *pbtask.GetTaskRequest) (*pbtask.GetTaskResponse, error) {
	requesterID, err := extractRequesterIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	taskID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, err
	}

	query, err := apptask.NewFindTaskQuery(taskID, requesterID)
	if err != nil {
		return nil, err
	}

	res, err := h.getHandler.Handle(ctx, query)
	if err != nil {
		return nil, err
	}

	return toProtoGetTaskResponse(res), nil
}

func (h *TaskHandler) SearchTasks(ctx context.Context, req *pbtask.TaskFilter) (*pbtask.SearchTasksResponse, error) {
	requesterID, err := extractRequesterIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	builder := apptask.NewSearchTasksQueryBuilder()

	builder.WithSearchQuery(*req.SearchQuery)
	builder.WithRequesterID(requesterID)

	if req.GroupId != nil {
		groupID, err := uuid.Parse(req.GetGroupId())
		if err != nil {
			return nil, fmt.Errorf("parse group id: %w", err)
		}

		builder.WithGroupID(groupID)
	}

	if req.GeoFilter != nil {
		query, err := toGeoFilterQuery(req.GetGeoFilter())
		if err != nil {
			return nil, err
		}

		builder.WithGeoQuery(query)
	}

	if req.IsFavorite != nil {
		builder.WithIsFavorite(req.GetIsFavorite())
	}

	if req.Priority != nil {
		builder.WithPriority(int32(req.GetPriority()))
	}

	if req.Status != nil {
		builder.WithStatus(int32(req.GetStatus()))
	}

	builder.WithCursor(
		apptask.NewSearchCursor(
			req.GetCursor().GetLimit(),
			req.GetCursor().GetKey().AsTime(),
		),
	)

	query, err := builder.Build()
	if err != nil {
		return nil, err
	}

	res, key, err := h.searchHandler.Handle(ctx, query)
	if err != nil {
		return nil, err
	}
	return NewSearchTasksResponse(res, key), nil
}

func (h *TaskHandler) CompleteTask(ctx context.Context, req *pbtask.CompleteTaskRequest) (*emptypb.Empty, error) {
	requesterID, err := extractRequesterIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	parsedTaskID, err := uuid.Parse(*req.Id)
	if err != nil {
		return nil, err
	}

	cmd, err := apptask.NewCompleteTaskCommand(parsedTaskID, requesterID)
	if err != nil {
		return nil, err
	}

	err = h.completeHandler.Handle(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (h *TaskHandler) CreateTask(ctx context.Context, req *pbtask.CreateTaskRequest) (*pbtask.CreateTaskResponse, error) {
	requesterID, err := extractRequesterIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	builder := apptask.NewCreateTaskCommandBuilder().
		WithOwnerID(requesterID).
		WithTitle(req.GetTitle()).
		WithDescription(req.GetDescription()).
		WithIsFavorite(req.GetIsFavorite()).
		WithPriority(int32(req.GetPriority())).
		WithIcon(int32(req.GetIcon()))

	if req.GroupId != nil {
		groupID, err := uuid.Parse(req.GetGroupId())
		if err != nil {
			return nil, err
		}

		builder.WithGroupID(groupID)
	}

	if req.Location != nil {
		builder.WithLocation(fromProtoGeoPoint(req.Location))
	}

	if req.Deadline != nil {
		builder.WithDeadline(req.Deadline.AsTime())
	}

	if req.NotifyAt != nil {
		builder.WithNotifyAt(fromProtoNotifyAt(req.GetNotifyAt()))
	}

	cmd := builder.Build()
	id, err := h.createHandler.Handle(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return &pbtask.CreateTaskResponse{
		Id: new(id.String()),
	}, nil
}

func (h *TaskHandler) UpdateTask(ctx context.Context, req *pbtask.UpdateTaskRequest) (*emptypb.Empty, error) {
	requesterID, err := extractRequesterIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	builder := apptask.NewUpdateTaskCommandBuilder()
	builder.WithRequesterID(requesterID)

	taskID, err := uuid.Parse(*req.Id)
	if err != nil {
		return nil, err
	}

	builder.WithID(taskID)

	if req.GroupId != nil {
		return nil, err
	}

	if req.Title != nil {
		builder.WithTitle(req.GetTitle())
	}

	if req.Description != nil {
		builder.WithDescription(req.GetDescription())
	}

	if req.Location != nil {
		builder.WithLocation(fromProtoGeoPoint(req.Location))
	}

	if req.IsFavorite != nil {
		builder.WithIsFavorite(req.GetIsFavorite())
	}

	if req.Priority != nil {
		priority, err := domaintask.ParsePriority(req.GetPriority().String())
		if err != nil {
			return nil, err
		}

		builder.WithPriority(priority)
	}

	if req.Deadline != nil {
		builder.WithDeadline(req.GetDeadline().AsTime())
	}

	if req.NotifyAt != nil {
		builder.WithNotifyAt(fromProtoNotifyAt(req.NotifyAt))
	}

	cmd, err := builder.Build()
	if err != nil {
		return nil, err
	}

	if err = h.updateHandler.Handle(ctx, cmd); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (h *TaskHandler) DeleteTask(ctx context.Context, req *pbtask.DeleteTaskRequest) (*pbtask.DeleteTaskResponse, error) {
	requesterID, err := extractRequesterIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	parsedTaskID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, err
	}

	cmd, err := apptask.NewDeleteTaskCommand(parsedTaskID, requesterID)
	if err != nil {
		return nil, err
	}

	res, err := h.deleteHandler.Handle(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return &pbtask.DeleteTaskResponse{
		Id:      new(res.ID.String()),
		PurgeAt: timestamppb.New(res.PurgeAt),
	}, nil
}

func (h *TaskHandler) RestoreTask(ctx context.Context, req *pbtask.RestoreTaskRequest) (*emptypb.Empty, error) {
	requesterID, err := extractRequesterIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	taskID, err := uuid.Parse(*req.Id)
	if err != nil {
		return nil, err
	}

	cmd, err := apptask.NewRestoreTaskCommand(taskID, requesterID)
	if err = h.restoreHandler.Handle(ctx, cmd); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (h *TaskHandler) ListDeleted(ctx context.Context, req *pbtask.ListDeletedRequest) (*pbtask.SearchTasksResponse, error) {
	requesterID, err := extractRequesterIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	b := apptask.NewListDeletedTasksQueryBuilder()
	b.WithRequesterID(requesterID)

	if req.GroupId != nil {
		groupID, err := uuid.Parse(req.GetGroupId())
		if err != nil {
			return nil, err
		}

		b.WithGroupID(groupID)
	}

	b.WithCursor(
		apptask.NewSearchCursor(
			req.GetCursor().GetLimit(),
			req.GetCursor().GetKey().AsTime(),
		),
	)

	query, err := b.Build()
	if err != nil {
		return nil, fmt.Errorf("query build: %w", err)
	}

	ids, key, err := h.listDeletedHandler.Handle(ctx, query)
	if err != nil {
		return nil, err
	}

	return NewSearchTasksResponse(ids, key), nil
}

type TaskHandlerBuilder struct {
	getHandler         apptask.GetHandler
	searchHandler      apptask.SearchHandler
	completeHandler    apptask.CompleteHandler
	createHandler      apptask.CreateHandler
	updateHandler      apptask.UpdateHandler
	deleteHandler      apptask.DeleteHandler
	restoreHandler     apptask.RestoreHandler
	listDeletedHandler apptask.ListDeletedHandler
}

func (b *TaskHandlerBuilder) WithGetHandler(h apptask.GetHandler) *TaskHandlerBuilder {
	b.getHandler = h
	return b
}

func (b *TaskHandlerBuilder) WithSearchHandler(h apptask.SearchHandler) *TaskHandlerBuilder {
	b.searchHandler = h
	return b
}

func (b *TaskHandlerBuilder) WithCompleteHandler(h apptask.CompleteHandler) *TaskHandlerBuilder {
	b.completeHandler = h
	return b
}

func (b *TaskHandlerBuilder) WithCreateHandler(h apptask.CreateHandler) *TaskHandlerBuilder {
	b.createHandler = h
	return b
}

func (b *TaskHandlerBuilder) WithUpdateHandler(h apptask.UpdateHandler) *TaskHandlerBuilder {
	b.updateHandler = h
	return b
}

func (b *TaskHandlerBuilder) WithDeleteHandler(h apptask.DeleteHandler) *TaskHandlerBuilder {
	b.deleteHandler = h
	return b
}

func (b *TaskHandlerBuilder) WithRestoreHandler(h apptask.RestoreHandler) *TaskHandlerBuilder {
	b.restoreHandler = h
	return b
}

func (b *TaskHandlerBuilder) WithListDeletedHandler(h apptask.ListDeletedHandler) *TaskHandlerBuilder {
	b.listDeletedHandler = h
	return b
}

func (b *TaskHandlerBuilder) Build() (*TaskHandler, error) {
	if b.getHandler == nil {
		return nil, errors.New("get handler is missing")
	}

	if b.searchHandler == nil {
		return nil, errors.New("search handler is missing")
	}

	if b.completeHandler == nil {
		return nil, errors.New("complete handler is missing")
	}

	if b.createHandler == nil {
		return nil, errors.New("create handler is missing")
	}

	if b.updateHandler == nil {
		return nil, errors.New("update handler is missing")
	}

	if b.deleteHandler == nil {
		return nil, errors.New("delete handler is missing")
	}

	if b.restoreHandler == nil {
		return nil, errors.New("restore handler is missing")
	}

	if b.listDeletedHandler == nil {
		return nil, errors.New("list deleted handler is missing")
	}

	return &TaskHandler{
		getHandler:         b.getHandler,
		searchHandler:      b.searchHandler,
		completeHandler:    b.completeHandler,
		createHandler:      b.createHandler,
		updateHandler:      b.updateHandler,
		deleteHandler:      b.deleteHandler,
		restoreHandler:     b.restoreHandler,
		listDeletedHandler: b.listDeletedHandler,
	}, nil
}

func NewTaskHandlerBuilder() *TaskHandlerBuilder {
	return &TaskHandlerBuilder{}
}
